package lab

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"cli-silly/course"
)

func archiveMember(root *os.Root, archive, member string) string {
	contents, err := readRegular(root, archive)
	if err != nil {
		return ""
	}
	compressed, err := gzip.NewReader(bytes.NewReader(contents))
	if err != nil {
		return ""
	}
	defer compressed.Close()
	reader := tar.NewReader(compressed)
	header, err := reader.Next()
	if err != nil || header.Name != member || header.Typeflag != tar.TypeReg || header.Size > 64<<10 {
		return ""
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return ""
	}
	if _, err := reader.Next(); err != io.EOF {
		return ""
	}
	return string(data)
}

func exerciseScript(root *os.Root, check course.Criterion) (string, error) {
	script, err := readRegular(root, check.Path)
	if err != nil {
		return "", nil
	}
	dir, err := os.MkdirTemp("", "cli-silly-script-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	name := filepath.Base(check.Path)
	if err := os.WriteFile(filepath.Join(dir, name), script, 0600); err != nil {
		return "", err
	}
	if check.Kind == "script-config" {
		if err := os.WriteFile(filepath.Join(dir, "another config"), []byte("sample=value\n"), 0600); err != nil {
			return "", err
		}
		stdout, stderr, status := runScript(dir, name, "another config")
		if stdout != "sample=value\n" || stderr != "" || status != 0 {
			return "", nil
		}
		stdout, stderr, status = runScript(dir, name, "another missing config")
		if stdout == "" && stderr == "missing config\n" && status == 2 {
			return "passed", nil
		}
		return "", nil
	}

	if err := os.Mkdir(filepath.Join(dir, "config"), 0700); err != nil {
		return "", err
	}
	for path, contents := range map[string]string{
		"services.txt": "search worker\nqueue\n", "config/queue.conf": "workers=3\n",
	} {
		if err := os.WriteFile(filepath.Join(dir, path), []byte(contents), 0600); err != nil {
			return "", err
		}
	}
	for range 2 {
		_, _, status := runScript(dir, name)
		if status != 0 {
			return "", nil
		}
	}
	scratch, err := os.OpenRoot(dir)
	if err != nil {
		return "", err
	}
	defer scratch.Close()
	result, err := readRegular(scratch, "report.txt")
	if err == nil && string(result) == "missing search worker\npresent queue\n" {
		return "passed", nil
	}
	return "", nil
}

type boundedOutput struct{ bytes.Buffer }

func (b *boundedOutput) Write(p []byte) (int, error) {
	if b.Len()+len(p) > 64<<10 {
		return 0, errors.New("tool output exceeded 64 KiB")
	}
	return b.Buffer.Write(p)
}

func runScript(dir, name string, args ...string) (string, string, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var stdout, stderr boundedOutput
	cmd := exec.CommandContext(ctx, "bash", append([]string{"--noprofile", "--norc", "./" + name}, args...)...)
	cmd.Dir, cmd.Stdout, cmd.Stderr = dir, &stdout, &stderr
	cmd.WaitDelay = time.Second
	if err := cmd.Run(); err != nil {
		var result *exec.ExitError
		if errors.As(err, &result) {
			return stdout.String(), stderr.String(), result.ExitCode()
		}
		return "", "", -1
	}
	return stdout.String(), stderr.String(), 0
}
