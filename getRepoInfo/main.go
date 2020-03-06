package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	ossList, _ := getOssList()
	for _, j := range ossList {
		respRepoInfo, err := getRepoInfo(j[0], j[1])
		if err != nil {
			log.Fatalf("Error: Cannot get repo info.")
		}
		respNumCommits, err := getNumCommitsForHalfYear(j[0], j[1])
		if err != nil {
			log.Fatalf("Error: Cannot get number of commits.")
		}
		numContributors, err := getNumContributors(j[0], j[1])
		if err != nil {
			log.Fatalf("Error: Cannot get number of contributors.")
		}
		result := "{\"repoInfo\":" + respRepoInfo + ",\"commitsForHalfYear\":" + respNumCommits + "," + "\"contributors\":" + strconv.Itoa(numContributors) + "}"
		fmt.Println(result)
		time.Sleep(time.Second)
	}
}

func getRepoInfo(owner, name string) (string, error) {
	client := http.Client{}
	json := `{repository(owner:\"` + owner + `\",name:\"` + name + `\"){name, url, licenseInfo{name}, createdAt, primaryLanguage{name}, defaultBranchRef{ name, target{ ... on Commit{history{totalCount}}}} releases(last:5,orderBy:{field:CREATED_AT, direction:ASC}){nodes{tagName, createdAt}}, stargazers{totalCount}}}`
	query := strings.NewReader(` { "query": "query` + json + `"}`)
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", query)
	apiToken := "bearer " + os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func getNumCommitsForHalfYear(owner, name string) (string, error) {
	client := http.Client{}
	json := `{repository(owner:\"` + owner + `\",name:\"` + name + `\"){defaultBranchRef{ name, target{ ... on Commit{history(since:\"2019-01-01T00:00:00\",until:\"2020-01-01T00:00:00\"){totalCount}}}}}}`
	query := strings.NewReader(` { "query": "query` + json + `"}`)
	req, err := http.NewRequest("POST", "https://api.github.com/graphql", query)
	apiToken := "bearer " + os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
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

func getNumContributors(owner, name string) (int, error) {
	resp, err := http.Get("http://github.com/" + owner + "/" + name + "/contributors_size")
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	slice := strings.Split(string(body), "\n")
	i, _ := strconv.Atoi(strings.TrimSpace(strings.Replace(slice[3], ",", "", -1)))
	return i, err
}
