## Try it in the terminal

### Read an unfamiliar command before running it

Use `type head` to see how Bash resolves the command, then use `head --help` to find the option that accepts a line count. Read one row from `inventory.csv` and mark the option, its value, and the file operand. Changing the count from one to two should leave the contents of each row intact.

Try `cut -d ',' -f 1 inventory.csv` and identify the field you selected. Repeat with field 3. Read the values before composing the requested port report. Finally compare `wc -l inventory.csv` and `wc -l ports.txt`. They should have the same number of rows, even though their columns differ.
