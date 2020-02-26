# GithubUtils
## How to use
```bash
export GITHUB_TOKEN="YOUR_GITHUB_API_TOKEN"
go run main.go
```

Update ossList.csv files when adding products.
```
https://github.com/vuejs/vue
https://github.com/grafana/grafana
```

## Commit and Star Histories
You can get Commits and Stars histories.
```
spinnaker spinnaker
Commit: map[2015:637 2016:514 2017:397 2018:300 2019:156]
Star: map[2015:1398 2016:1097 2017:1430 2018:1544 2019:1301]
argoproj argo-cd
Commit: map[2019:1023 2015:0 2016:0 2017:0 2018:869]
Star: map[2015:0 2016:0 2017:0 2018:504 2019:1717
```

## Number of Commit group by Companies
You can get nuber of commits group by Companies.
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

## RepoInfo
You can get Repository information.
- Stars 
- License
- Primary Language
- Created Date
- Commits(total)
- Commits(last Year)
- Relases date(last:5)
- Versions(last:5)
- Contributors
