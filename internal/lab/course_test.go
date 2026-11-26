package lab

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cli-silly/course"
)

type labScenario struct {
	first  string
	finish string
}

var scenarios = map[string]labScenario{
	"review/worker-incident": {
		`awk '$3 >= 500 && $3 < 600 { print $2 }' events.log > failed-workers.txt`,
		`sort failed-workers.txt | uniq -c | awk '{print $1, $2}' > failures.txt
worker=$(sort -nr failures.txt | head -n 1 | awk '{print $2}')
jq -r --arg worker "$worker" '.workers[] | select(.name == $worker) | .team' workers.json > owner.txt`,
	},
	"shell/make-it-runnable": {
		`chmod u+x bin/health-check`,
		`./bin/health-check > status.txt`,
	},
	"shell/paths-and-quoting": {
		`cp "Draft Files/server config.txt" config/server.conf`,
		`mv "Draft Files/server config.txt" config/server.conf`,
	},
	"streams/find-config": {
		`grep -rl '^debug=true$' config > debug-files.txt`,
		`grep '^debug=true$' "$(cat debug-files.txt)" > debug-setting.txt`,
	},
	"streams/log-report": {
		`awk '$3 >= 500 && $3 < 600 {print $2}' access.log > errors.txt`,
		`sort errors.txt | uniq -c | awk '{print $1, $2}' > counts.txt`,
	},
	"streams/json-report": {
		`jq -r '.services[] | select(.enabled) | .name' deploy.json > enabled.txt`,
		`jq -r '.services[] | select(.enabled) | [.name, .port] | @tsv' deploy.json > ports.tsv`,
	},
	"streams/split-streams": {
		`bash sync-check.sh > data.txt 2>/dev/null || true`,
		`status=0; bash sync-check.sh > data.txt 2>errors.txt || status=$?; printf '%s\n' "$status" > status.txt`,
	},
	"git/selective-staging": {
		`git add app.cfg`,
		`git commit -m 'Enable ready mode'`,
	},
	"git/merge-feature": {
		`git switch feature`,
		`git switch main; git merge feature`,
	},
	"git/resolve-conflict": {
		`git merge feature || test -f .git/MERGE_HEAD`,
		`printf 'theme=auto\n' > settings.txt; git add settings.txt; git commit -m 'Resolve theme'`,
	},
	"git/local-remote": {
		`printf 'release=ready\n' > release.txt; git add release.txt; git commit -m 'Prepare release'`,
		`git push -u origin main`,
	},
	"bash/arguments-and-exit": {
		`cat > check-config.sh <<'SCRIPT'
if [ -f "$1" ]; then cat -- "$1"; else printf 'missing config\n' >&2; exit 2; fi
SCRIPT
bash check-config.sh "fixtures/demo config" > good.out`,
		`status=0; bash check-config.sh missing 2> bad.out || status=$?; printf '0\n%s\n' "$status" > exits.txt`,
	},
	"bash/release-report": {
		`cat > report.sh <<'SCRIPT'
#!/usr/bin/env bash
while IFS= read -r service; do
    if [ -f "config/$service.conf" ]; then
        printf 'present %s\n' "$service"
    else
        printf 'missing %s\n' "$service"
    fi
done < services.txt > report.txt
SCRIPT
chmod u+x report.sh`,
		`./report.sh; ./report.sh`,
	},
	"operations/service-startup": {
		`sed -n 's/^active_release=//p' deployment.txt > diagnosis.txt`,
		`printf 'WorkingDirectory=%s\n' "$(cat diagnosis.txt)" > service-fix.conf`,
	},
	"operations/disk-and-archives": {
		`tar -czf incident.tar.gz logs/incident.txt; tar -tzf incident.tar.gz > members.txt`,
		`tar -xOzf incident.tar.gz logs/incident.txt > evidence.txt`,
	},
	"operations/ssh-config": {
		`printf 'Host build-box\n    HostName build.internal.example\n    User deploy\n    Port 2222\n' > lab-ssh.conf; ssh -G -F lab-ssh.conf build-box | awk '$1=="hostname"{print $2}' > host.txt`,
		`ssh -G -F lab-ssh.conf build-box | awk '$1=="hostname"{h=$0} $1=="user"{u=$0} $1=="port"{p=$0} END{print h; print u; print p}' > effective.txt`,
	},
	"docker/compose-health": {
		`docker compose up --build -d --wait --wait-timeout 60`,
		`curl -fsS "http://$(docker compose port api 8080)/health" > health.txt`,
	},
	"docker/persistent-data": {
		`docker compose up --build -d --wait --wait-timeout 60`,
		`docker compose exec -T api sh -c 'printf "released\n" > /data/status.txt'
docker compose ps -q api > previous-container.txt
docker compose up -d --force-recreate --wait --wait-timeout 60
mapping=$(docker compose port api 8080)
printf '%s\n' "${mapping##*:}" > published.port
curl -fsS "http://$mapping/health" > saved.txt`,
	},
	"docker/service-network": {
		`sed -i 's|localhost:8081|worker:8081|' compose.yaml; docker compose up --build -d --wait --wait-timeout 60`,
		`curl -fsS "http://$(docker compose port api 8080)/health" > health.txt`,
	},
}

func TestCatalogLabs(t *testing.T) {
	entries, err := course.Entries()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		lesson := entry.Lesson
		t.Run(lesson.ID, func(t *testing.T) {
			for _, tool := range lesson.Tools {
				if _, err := exec.LookPath(tool); err != nil {
					t.Skipf("%s is unavailable", tool)
				}
			}
			store := Store{DataDir: t.TempDir(), StateDir: t.TempDir()}
			path, err := store.Open(lesson)
			if err != nil {
				t.Fatal(err)
			}

			if lesson.Setup == "compose" {
				runTask(t, path, "docker compose config --quiet")
				if os.Getenv("CLI_SILLY_TEST_DOCKER") != "1" {
					t.Skip("Compose syntax verified; set CLI_SILLY_TEST_DOCKER=1 for live checks")
				}
				t.Cleanup(func() { runTask(t, path, "docker compose down --volumes --remove-orphans") })
			}
			if strings.HasSuffix(lesson.ID, "inspect-local-api") || lesson.ID == "capstone/release-rescue" {
				testLocalService(t, store, lesson, path)
				return
			}
			scenario, ok := scenarios[lesson.ID]
			if !ok {
				t.Fatal("add untouched, partial, and completed coverage for this lesson")
			}
			checkScore(t, store, lesson, path, "starter", false)
			runTask(t, path, scenario.first)
			checkScore(t, store, lesson, path, "partial", false)
			runTask(t, path, scenario.finish)
			checkScore(t, store, lesson, path, "finished", true)
		})
	}
}

func runTask(t *testing.T, path, script string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-e", "-c", script)
	cmd.Dir = path
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("run lab solution: %v\n%s\n%s", err, script, output)
	}
}

func checkScore(t *testing.T, store Store, lesson course.Lesson, path, state string, complete bool) {
	t.Helper()
	result, err := store.Check(lesson, path)
	if err != nil || result.Complete != complete || (state == "partial" && (result.Score == 0 || result.Score == 100)) {
		t.Fatalf("%s: score=%d complete=%v err=%v items=%+v", state, result.Score, result.Complete, err, result.Items)
	}
}

func testLocalService(t *testing.T, store Store, lesson course.Lesson, path string) {
	t.Helper()
	checkScore(t, store, lesson, path, "starter", false)
	isCapstone := lesson.ID == "capstone/release-rescue"
	expected := "service: ready\n"
	if isCapstone {
		expected = "release: recovered\n"
	}
	if err := os.WriteFile(filepath.Join(path, "health.txt"), []byte(expected), 0600); err != nil {
		t.Fatal(err)
	}
	checkScore(t, store, lesson, path, "partial", false)
	if isCapstone {
		if err := os.WriteFile(filepath.Join(path, "release.env"), []byte("APP_MODE=ready\n"), 0600); err != nil {
			t.Fatal(err)
		}
	}

	cmd := exec.Command("python3", "app.py")
	cmd.Dir, cmd.Env = path, append(os.Environ(), "APP_MODE=ready")
	log, err := os.Create(filepath.Join(path, "test-server.log"))
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	cmd.Stdout, cmd.Stderr = log, log
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()
	portFile := filepath.Join(path, "server.port")
	for i := 0; i < 100; i++ {
		if port, err := os.ReadFile(portFile); err == nil && len(port) != 0 {
			runTask(t, path, fmt.Sprintf("curl -fsS http://127.0.0.1:%s/health > health.txt", strings.TrimSpace(string(port))))
			checkScore(t, store, lesson, path, "finished", true)
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	logText, _ := os.ReadFile(log.Name())
	t.Fatalf("local server failed to start: %s", logText)
}
