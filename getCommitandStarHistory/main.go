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
	"time"
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
		fmt.Println("###########################")
		fmt.Println(j[0], j[1])
		fmt.Println("Commit")
		fmt.Println(respCommit)
		fmt.Println("Star")
		fmt.Println(respStar)
		fmt.Println("###########################")
		time.Sleep(time.Second)
	}
}

//TODO: use year and period
func getStarHistory(owner, name string, year, period int) (map[int]int, error) {
	log.Printf("Start getting Stars \n")
	starCount := make(map[int]int)
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
		//TODO: use slice
		const layout = "2006-01-02"
		first2017, _ := time.Parse(layout, "2017-01-01")
		second2017, _ := time.Parse(layout, "2017-07-01")
		first2018, _ := time.Parse(layout, "2018-01-01")
		second2018, _ := time.Parse(layout, "2018-07-01")
		first2019, _ := time.Parse(layout, "2019-01-01")
		second2019, _ := time.Parse(layout, "2019-07-01")
		first2020, _ := time.Parse(layout, "2020-01-01")

		for _, j := range *StarHistories {
			starredAt, _ := time.Parse(layout, (j.StarredAt)[:10])
			switch {
			case starredAt.After(first2017) && starredAt.Before(second2017):
				starCount[201701]++
			case starredAt.After(second2017) && starredAt.Before(first2018):
				starCount[201707]++
			case starredAt.After(first2018) && starredAt.Before(second2018):
				starCount[201801]++
			case starredAt.After(second2018) && starredAt.Before(first2019):
				starCount[201807]++
			case starredAt.After(first2019) && starredAt.Before(second2019):
				starCount[201901]++
			case starredAt.After(second2019) && starredAt.Before(first2020):
				starCount[201907]++
			case starredAt.After(first2020) && starredAt.Before(time.Now()):
				starCount[20201]++
			default:
				starCount[99999]++
			}
		}
		log.Printf("%d %d All:%d finished:%d \n", owner, name, pageNum, i)
	}
	return starCount, err
}

//TODO: use year and period
func getCommitHistory(owner, name string, year, period int) (commitCount map[int]int, err error) {

	log.Printf("Start getting Commits \n")
	test := []string{"2017-01-01", "2017-07-01", "2018-01-01", "2018-07-01", "2019-01-01", "2019-07-01", "2020-01-01"}
	commitCount = make(map[int]int)
	for i := range test {
		commitTotalCount := new(CommitTotalCount)
		client := http.Client{}
		jsonQuery := `{repository(owner:\"` + owner + `\",name:\"` + name + `\"){defaultBranchRef{ name, target{ ... on Commit{history(since:\"` + test[i] + `T00:00:00\",until:\"` + test[i+1] + `T00:00:00\"){totalCount}}}}}}`
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
		commitCount[i] = commitTotalCount.Data.Repository.DefaultBranchRef.Target.History.TotalCount
		if i == 5 {
			break
		}
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
