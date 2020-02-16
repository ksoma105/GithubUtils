# GitHubCommitRatio
## What is this
You can get the commit ratio for each company.
(Calculated from the top 15 Contributors.)

Like this
```json
{
  "owner": "grafana",
  "name": "grafana",
  "companycommits": {
    "": 1697,
    "@grafana ": 3229,
    "@grafana and @raintank": 667,
    "Atomler": 371,
    "Grafana Labs": 6924
  }
}
```

## How to use
```bash
export GITHUB_TOKEN = YOUR_GITHUB_API_TOKEN
go run main.go
```

Update ossList.csv files when adding products.
```
https://github.com/vuejs/vue
https://github.com/grafana/grafana
```