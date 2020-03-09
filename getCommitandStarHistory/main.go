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
	"sync"
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
	var wg sync.WaitGroup
	ossList, _ := getOssList()
	fmt.Println("owner,name,~2017/01,2017/7,2018/1,2018/7,2019/1,2019/7,2020/1,2020/3")
	fmt.Println("##### COMMITS #####")
	for _, j := range ossList {
		wg.Add(1)
		go func(i []string) {
			respCommit, err := getCommitHistory(i[0], i[1], 2015, 5)
			//		respStar, err := getStarHistory(j[0], j[1], 2015, 5)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v \n", i[0], i[1], respCommit[0], respCommit[1], respCommit[2], respCommit[3], respCommit[4], respCommit[5], respCommit[6], respCommit[7])
			wg.Done()
		}(j)
	}
	wg.Wait()
	fmt.Println("##### STARS #####")
	for _, j := range ossList {
		wg.Add(1)
		go func(i []string) {
			respStar, err := getStarHistory(i[0], i[1], 2015, 5)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v \n", i[0], i[1], respStar[201701], respStar[201707], respStar[201801], respStar[201807], respStar[201901], respStar[201907], respStar[202001], respStar[202003])
			wg.Done()
		}(j)
	}
	wg.Wait()
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
		if resp.StatusCode != 200 {
			log.Printf("%v,%v,Status Code is %d. Auto retry after 30 sec.", owner, name, resp.StatusCode)
			time.Sleep(30 * time.Second)
			// Fix: Recursive processing. golang is not good at recursive process.
			return getStarHistory(owner, name, year, period)
		}
		StarHistories := new(StartHistories)
		err = json.Unmarshal(body, StarHistories)
		if err != nil {
			log.Fatal(err)
		}
		//DEBUG
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
			if starredAt.Before(first2017) {
				starCount[201701]++
			}
			if starredAt.Before(second2017) {
				starCount[201707]++
			}
			if starredAt.Before(first2018) {
				starCount[201801]++
			}
			if starredAt.Before(second2018) {
				starCount[201807]++
			}
			if starredAt.Before(first2019) {
				starCount[201901]++
			}
			if starredAt.Before(second2019) {
				starCount[201907]++
			}
			if starredAt.Before(first2020) {
				starCount[202001]++
			}
			if starredAt.Before(time.Now()) {
				starCount[202003]++
			}
		}
		log.Printf("%d %d All:%d finished:%d \n", owner, name, pageNum, i)
		time.Sleep(time.Second)
	}
	if starCount[202003] == 40000 {
		log.Printf("%v,%v STARS LIMIT !!! \n", owner, name)
	}
	return starCount, err
}

//TODO: use year and period
func getCommitHistory(owner, name string, year, period int) (commitCount map[int]int, err error) {

	log.Printf("Start getting Commits \n")
	test := []string{"2017-01-01", "2017-01-01", "2017-07-01", "2018-01-01", "2018-07-01", "2019-01-01", "2019-07-01", "2020-01-01", "2020-03-01"}
	commitCount = make(map[int]int)
	for i := range test {
		commitTotalCount := new(CommitTotalCount)
		client := http.Client{}
		jsonQuery := ""
		if i == 0 {
			jsonQuery = `{repository(owner:\"` + owner + `\",name:\"` + name + `\"){defaultBranchRef{ name, target{ ... on Commit{history(until:\"` + test[i] + `T00:00:00\"){totalCount}}}}}}`
		} else {
			jsonQuery = `{repository(owner:\"` + owner + `\",name:\"` + name + `\"){defaultBranchRef{ name, target{ ... on Commit{history(since:\"` + test[i] + `T00:00:00\",until:\"` + test[i+1] + `T00:00:00\"){totalCount}}}}}}`
		}
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
		commitCount[i] = commitTotalCount.Data.Repository.DefaultBranchRef.Target.History.TotalCount + commitCount[i-1]
		if i == 7 {
			break
		}
	}
	return commitCount, err
}

func getOssList() (ossList [][]string, err error) {
	file, err := os.Open("../inputData/Obserbility.csv")
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
