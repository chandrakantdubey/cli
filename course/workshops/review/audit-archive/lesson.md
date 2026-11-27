Package one file for a teammate and inspect the archive before handing it over. The file you include has a space in its name. The member stored inside the archive should be separate from the directory used to find the original.

## Commands and operations

In `tar -czf sample.tar.gz -C reports "api log.txt"`, `-c` creates the archive, `-z` compresses it with gzip, and `-f` sets the archive filename. With `-C reports`, the following input name is looked up inside reports. Quoting keeps the entire filename together.

Use `tar -tzf sample.tar.gz` to list members without extracting them. To inspect one member, use `tar -xOzf sample.tar.gz "api log.txt"`. The capital `-O` sends its contents to stdout, so you can capture evidence separately from the source.

Check the stored member name first. If you instead archive a path beginning with reports/, an extraction command looking for just the basename will not find it.
