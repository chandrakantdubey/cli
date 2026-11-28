Learn to interpret identities before assigning permissions. The workspace contains captured account and group records for a fictitious team. You can inspect them as a regular user and compare their format with safe, read-only commands on your own system.

## Commands and operations

An account record has colon-separated fields. Field 1 is the account name, field 3 is the UID, field 4 is the primary GID, field 6 the home directory, and field 7 the login shell. A group record uses field 1 for its name, field 3 for GID, and field 4 for listed members.

On your own system, `id -u` and `id -g` show your effective user and group numbers. `getent passwd "$(id -un)"` asks the system's configured identity services for the account record. The scored lab uses fixed captures so every learner gets the same result.

With `awk -F:`, the colon is the field separator. Compare `$1` to an exact name to avoid matching similarly named accounts. To find the primary group, use the account's fourth field and match it against group field 3. A primary membership can be absent from field 4. Keep system identities unchanged while doing the exercise.

## Identity

```bash
whoami
id
id "$USER"
groups
```

Important concepts:

- UID: numeric user identity
- GID: primary group identity
- supplementary groups: additional group memberships

## Important files

Traditional local identity information is represented through:

```text
/etc/passwd
/etc/group
/etc/shadow
```

Do not edit `/etc/shadow` manually. Modern systems may also obtain identities from LDAP, Active Directory, or other NSS providers.

Inspect safely:

```bash
getent passwd "$USER"
getent group
```

## Create users

These administrative examples require a disposable VM with an account you control. The scored exercise uses local captures.

```bash
sudo useradd --create-home --shell /bin/bash alice
sudo passwd alice
```

Distribution-specific tools such as `adduser` may provide friendlier workflows.

## Groups

```bash
sudo groupadd developers
sudo usermod -aG developers alice
id alice
```

The `-a` in `-aG` matters: without append, you can replace supplementary group membership.

A new login session is usually required before group membership is reflected in that user's session.

## Switching identity

```bash
su - alice
sudo -iu alice
```

`sudo -iu` starts a login shell as the target user when policy permits.

## Ownership

```bash
ls -l file.txt
sudo chown alice:developers file.txt
```

Ownership and permissions are separate concepts.

## Principle of least privilege

Services should run with only the permissions they need. Avoid running application processes as root unless there is a specific, reviewed reason.

## Before you practice

**Outcome:** Join account and group records using the numeric primary group ID.

Know how to read files. The command section introduces awk field selection.

## Common mistakes and recovery

A primary group can have a different name from the user. Match the GID rather than guessing from names or the supplemental-member column.

Inspect the relevant input and output before repeating a command. Use `cli-silly hint` for a starting point and `cli-silly check` for criterion-by-criterion feedback. If the starter is damaged, `cli-silly reset` archives this attempt before recreating it.

## Use it in development

Interpret service account ownership when a process cannot access an application file.

After passing, explain which files or state changed and how you verified them. Open related practice from **Do** to apply the idea to different data.
