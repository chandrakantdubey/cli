A process normally receives three streams. File descriptor 0 is stdin, 1 is stdout, and 2 is stderr. A pipe sends stdout forward by default. Use `2>` to record diagnostic messages without mixing them into data.

The shell opens redirections left to right. Use two clear filenames when a caller needs to keep the streams separate.
