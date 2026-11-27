#!/usr/bin/env bash
trap 'rm -f -- ready.txt; printf "handled TERM\n" > stopped.txt; exit 0' TERM
printf 'ready\n' > ready.txt
for ((tick=0; tick<100; tick++)); do
    sleep 0.2
done
rm -f -- ready.txt
