OpenSSH can resolve a host alias into a hostname, user, port, and other settings. Use the `Host` block to scope an alias. `ssh -G` prints effective values locally without starting a connection. `-F` selects a config file for this exercise.

## Commands and operations

An SSH alias is the word matched by a `Host` block. `HostName` supplies the actual destination, `User` supplies the account, and `Port` supplies the service port. Change the block for `build-box` rather than editing your normal SSH configuration.

- `ssh -G -F lab-ssh.conf build-box` resolves settings for that alias and prints the effective configuration. The `-F` flag selects this lab file. The `-G` flag prints settings before making a network connection.
- In the output, use `awk '$1=="hostname" {print $2}'` to select the hostname value. Matching the first field directly avoids selecting keys whose names merely contain the same text.
- The other report needs three complete settings in a specified order. Print the effective lines by their exact keys and verify their spelling and order before saving the report.

If the effective host or user still looks like a default, check which alias you passed and which block matches it. The names and ports are fixture values, so you can finish while offline.

## Connect

The remote commands in these reference sections illustrate SSH syntax. The local exercise checks its configuration with `ssh -G` and makes no remote connection.

```bash
ssh user@server.example.com
ssh -p 2222 user@server.example.com
```

The server identity is checked through the host-key mechanism. Do not blindly accept unexpected host-key changes.

## Key authentication

Create a modern key pair:

```bash
ssh-keygen -t ed25519
```

Copy the public key using a trusted mechanism:

```bash
ssh-copy-id user@server.example.com
```

Keep the private key private.

```bash
chmod 700 ~/.ssh
chmod 600 ~/.ssh/id_ed25519
```

## SSH config

Use `~/.ssh/config` to create named hosts:

```sshconfig
Host staging
    HostName staging.example.com
    User deploy
    IdentityFile ~/.ssh/id_ed25519
```

Then:

```bash
ssh staging
```

## Debugging

```bash
ssh -v user@server.example.com
```

Verbose output helps distinguish DNS, TCP, authentication, key, and configuration failures.

## File transfer

```bash
scp report.txt user@server:/tmp/
rsync -av ./project/ user@server:/opt/project/
```

For large or repeated transfers, `rsync` can avoid retransmitting unchanged data.

## Tunneling

Local forwarding:

```bash
ssh -L 8080:internal-service:80 user@bastion
```

This exposes the remote destination through a local port. Understand the security implications before forwarding sensitive services.

## SSH agent

An agent can hold keys so you do not repeatedly enter passphrases. Be careful when forwarding agents across untrusted hosts.

## Production hygiene

- Prefer key authentication.
- Protect private keys with passphrases.
- Restrict administrative access.
- Keep a second session open while changing SSH server configuration.
- Test configuration before disconnecting the only working session.

## Before you practice

**Outcome:** Verify effective SSH alias settings without opening a network connection.

Understand files, arguments, and host/port names. Install the OpenSSH client for ssh -G.

## Common mistakes and recovery

ssh -G prints effective configuration. It does not log in. Use -F lab-ssh.conf so you inspect the exercise file.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Review a deployment host alias before attempting remote access.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
