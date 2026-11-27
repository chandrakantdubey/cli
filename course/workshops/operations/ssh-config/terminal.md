## Try it in the terminal

### Verify the resolved settings locally

Read `lab-ssh.conf` and identify the alias and its setting block. Run `ssh -G -F lab-ssh.conf build-box` to print the effective configuration without opening a connection. Find the hostname, user, and port among the output. An SSH alias is a local label and may be different from the hostname it expands to.

Make one change, read the effective configuration again, and confirm only the setting you intended moved. Compare the format of the one-value host report with the three full setting lines requested by the rubric. If they have the same content shape, one command is probably selecting too many or too few fields.
