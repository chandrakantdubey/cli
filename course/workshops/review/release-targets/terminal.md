## Try it in the terminal

Try filtering just the enabled records and count what remains. Add the environment condition and verify that a staging entry disappears. Inspect the resulting fields with `cat -T` to make tabs visible.

Compare `cut -f2 targets.tsv | sort -u` with numeric sort. Look for the duplicate shared port, then verify both reports before running the check.
