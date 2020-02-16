package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
)

// StarHistory Github StarredAt json struct
type StarHistory struct {
	StarredAt string `json:"starred_at"`
}

// CommitTotalCount Github StarredAt json struct
type CommitTotalCount struct {
	Data struct {
		Repository struct {
			DefaultBranchRef struct {
				Name   string `json:"name"`
				Target struct {
					History struct {
						TotalCount int `json:"totalCount"`
					} `json:"history"`
				} `json:"target"`
			} `json:"defaultBranchRef"`
		} `json:"repository"`
	} `json:"data"`
}

// CommitTotalCounts CommitTotalCount slice
type CommitTotalCounts []CommitTotalCount

// StartHistories StarHistory slice
type StartHistories []StarHistory

// Contributor Github contributors json struct
type Contributor struct {
	Author Author `json:"author"`
	Total  int    `json:"total"`
}

// Author Github contributors URL json struct
type Author struct {
	URL string `json:"url"`
}

// User Github user struct
type User struct {
	Company string `json:"company"`
}

// Contributors contributos slice
type Contributors []Contributor

func (a Contributors) Len() int           { return len(a) }
func (a Contributors) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Contributors) Less(i, j int) bool { return a[i].Total < a[j].Total }

func main() {
	ossList, _ := getOssList()
	for _, j := range ossList {
		respCommit, err := getCommitHistory(j[0], j[1], 2015, 5)
		respStar, err := getStarHistory(j[0], j[1], 2015, 5)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(j[0], j[1])
		fmt.Println(respCommit)
		fmt.Println(respStar)
	}
}

func getStarHistory(owner, name string, year, period int) (map[int]int, error) {
	starCount := make(map[int]int)
	for i := 0; i < period; i++ {
		starCount[year+i] = 0
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+owner+"/"+name+"/stargazers?per_page=100", nil)
	apiToken := "bearer " + os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", apiToken)
	req.Header.Set("Accept", "application/vnd.github.v3.star+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	header, _ := httputil.DumpResponse(resp, false)
	slice := strings.Split(string(header), "page=")
	slice = strings.Split(slice[4], "rel=")
	pageNum, _ := strconv.Atoi(strings.TrimRight(slice[0], ">; "))
	_, err = ioutil.ReadAll(resp.Body)
	for i := 1; i <= pageNum; i++ {
		req, err = http.NewRequest("GET", "https://api.github.com/repos/"+owner+"/"+name+"/stargazers?per_page=100&page="+strconv.Itoa(i), nil)
		req.Header.Set("Authorization", apiToken)
		req.Header.Set("Accept", "application/vnd.github.v3.star+json")
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		StarHistories := new(StartHistories)
		err = json.Unmarshal(body, StarHistories)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < period; i++ {
			for _, j := range *StarHistories {
				if strings.Contains(j.StarredAt, strconv.Itoa(year+i)) {
					starCount[year+i]++
				}
			}
		}
		fmt.Println(starCount)
	}
	return starCount, err
}

func getCommitHistory(owner, name string, year, period int) (commitCount map[int]int, err error) {
	commitCount = make(map[int]int)
	for i := 0; i < period; i++ {
		commitTotalCount := new(CommitTotalCount)
		client := http.Client{}
		jsonQuery := `{repository(owner:\"` + owner + `\",name:\"` + name + `\"){defaultBranchRef{ name, target{ ... on Commit{history(since:\"` + strconv.Itoa(year+i) + `-01-01T00:00:00\",until:\"` + strconv.Itoa(year+i) + `-12-31T00:00:00\"){totalCount}}}}}}`
		query := strings.NewReader(` { "query": "query` + jsonQuery + `"}`)
		req, err := http.NewRequest("POST", "https://api.github.com/graphql", query)
		apiToken := "bearer " + os.Getenv("GITHUB_TOKEN")
		req.Header.Set("Authorization", apiToken)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, commitTotalCount)
		fmt.Println(commitTotalCount.Data.Repository.DefaultBranchRef.Target.History.TotalCount)
		commitCount[year+i] = commitTotalCount.Data.Repository.DefaultBranchRef.Target.History.TotalCount
	}
	return commitCount, err
}

func getOssList() (ossList [][]string, err error) {
	file, err := os.Open("./ossList.csv")
	if err != nil {
		log.Fatalf("CSV file reading error.")
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var line []string
	for {
		line, err = reader.Read()
		if err != nil {
			break
		}
		slice := strings.Split(line[0], "/")
		ossList = append(ossList, slice[3:5])
	}
	return ossList, err
}
