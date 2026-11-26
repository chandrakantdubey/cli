This incident has two sources. The log shows which workers returned server errors, and the JSON inventory shows who owns each worker.

First reduce the log to the worker names that actually failed. Inspect the intermediate list before counting. Then use the inventory to find the team that owns the worker with the most failures.

```bash
awk '$3 >= 500 && $3 < 600 { print $2 }' sample.log
jq -r '.workers[] | select(.name == "worker-example") | .team' inventory.json
```

## Commands and operations

This challenge applies the log-report workflow to a different data set. Use the status field to find real failures, and preserve the sequence until you have written the first report.

- `awk '$3 >= 500 && $3 < 600 { print $2 }' sample.log` tests the third field numerically and returns the worker in field two. A worker name containing a three-digit number can still have returned success.
- `sort workers.txt | uniq -c` groups the workers and counts their failures. An additional `awk '{print $1, $2}'` removes count padding. Keep the counts report sorted by worker, even if you separately sort by count to identify the highest one.
- `jq -r '.workers[] | select(.name == "worker-example") | .team' inventory.json` visits each inventory object, selects exactly the requested name, and prints its team. Use the worker you found from the incident data rather than an example name.

Compare the three outputs with their source files. The raw matching worker list should retain duplicates, the counts report should group them, and the owner report should contain the team alone. The review score is stored independently of the original lab.
