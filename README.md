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
```csv
##### COMMITS #####
concourse,concourse,11746,13842,15505,16776,18292,19812,20928,21464 
argoproj,argo-cd,0,0,0,385,869,1387,1892,1983 
jenkins-x,jx,0,0,0,2047,5290,7730,9066,9275 
##### STARS #####
argoproj,argo-cd,0,0,0,40,503,1254,2217,2564 
jenkins-x,jx,0,0,0,1229,1831,2722,3252,3376 
fluxcd,flux,15,134,291,575,1282,2253,3603,4048 
```

## Number of Commit group by Companies
You can get nuber of commits group by Companies.
```csv
argoproj,argo-cd 
"intuit" ,1430
"5-:\$J$7" ,183
"floqast" ,67
"argoproj, uc berkeley grad" ,32
"yieldlab" ,10
"quipper" ,7
"nvidia" ,5
"tinkoff" ,4
"tesla" ,4
"riskified" ,4
"nomitor @distcloud" ,4
"mesosphere" ,4
"mambu" ,4
"majorleaguebaseball" ,4
"davidkarlsen.com" ,4
"tower-research" ,3
"sendgrid @sendgrid-ops @sendgrid-dev" ,3
"red hat" ,3
"mirantis" ,3
"ibm" ,3
"yieldlab ag" ,2
"vmware" ,2
"viaduct" ,2
"peloton interactive" ,2
"codility" ,2
"bad ass devops" ,2
"australiansynchrotron" ,2
"appdirect" ,2
"apalia" ,2
"akunca capital" ,2
"yros" ,1
"wongnai" ,1
"wisersolutions" ,1
"threefoldsys" ,1
"thousandeyes" ,1
"syncier" ,1
"swissquote" ,1
"moodev" ,1
"maxkelsen" ,1
"kasa-network" ,1
"kakao" ,1
"invisionapp" ,1
"hipages" ,1
"engineering lead @celonis" ,1
"edx" ,1
"edf-re" ,1
"commbank" ,1
"cloudphysics" ,1
"bookingcom" ,1
"biobox-analytics" ,1
"banzaicloud" ,1
"baloise" ,1
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
```json
{
	"repoInfo": {
		"data": {
			"repository": {
				"name": "spinnaker",
				"url": "https://github.com/spinnaker/spinnaker",
				"licenseInfo": {
					"name": "Apache License 2.0"
				},
				"createdAt": "2014-07-02T19:02:36Z",
				"primaryLanguage": {
					"name": "Shell"
				},
				"defaultBranchRef": {
					"name": "master",
					"target": {
						"history": {
							"totalCount": 2043
						}
					}
				},
				"releases": {
					"nodes": [
						{
							"tagName": "v0.79.0",
							"createdAt": "2017-04-05T15:46:53Z"
						},
						{
							"tagName": "v0.80.0",
							"createdAt": "2017-04-24T20:54:23Z"
						},
						{
							"tagName": "v0.81.0",
							"createdAt": "2017-05-23T14:03:40Z"
						},
						{
							"tagName": "v0.82.0",
							"createdAt": "2017-06-12T20:04:17Z"
						},
						{
							"tagName": "v0.83.0",
							"createdAt": "2017-12-19T18:01:30Z"
						}
					]
				},
				"stargazers": {
					"totalCount": 6927
				}
			}
		}
	},
	"commitsForHalfYear": {
		"data": {
			"repository": {
				"defaultBranchRef": {
					"name": "master",
					"target": {
						"history": {
							"totalCount": 156
						}
					}
				}
			}
		}
	},
	"contributors": 112
}
```