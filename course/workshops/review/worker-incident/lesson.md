This incident has two sources. The log shows which workers returned server errors, and the JSON inventory shows who owns each worker.

First reduce the log to the worker names that actually failed. Inspect the intermediate list before counting. Then use the inventory to find the team that owns the worker with the most failures.

```bash
awk '$3 >= 500 && $3 < 600 { print $2 }' sample.log
jq -r '.workers[] | select(.name == "worker-example") | .team' inventory.json
```
