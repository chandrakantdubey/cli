OpenSSH can resolve a host alias into a hostname, user, port, and other settings. Use the `Host` block to scope an alias. `ssh -G` prints effective values locally without starting a connection. `-F` selects a config file for this exercise.

## Commands and operations

An SSH alias is the word matched by a `Host` block. `HostName` supplies the actual destination, `User` supplies the account, and `Port` supplies the service port. Change the block for `build-box` rather than editing your normal SSH configuration.

- `ssh -G -F lab-ssh.conf build-box` resolves settings for that alias and prints the effective configuration. The `-F` flag selects this lab file. The `-G` flag prints settings before making a network connection.
- In the output, use `awk '$1=="hostname" {print $2}'` to select the hostname value. Matching the first field directly avoids selecting keys whose names merely contain the same text.
- The other report needs three complete settings in a specified order. Print the effective lines by their exact keys and verify their spelling and order before saving the report.

If the effective host or user still looks like a default, check which alias you passed and which block matches it. The names and ports are fixture values, so you can finish while offline.
