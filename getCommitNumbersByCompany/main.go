package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

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
		companyList, err := getComanyList(j[0], j[1])
		if err != nil {
			log.Fatal(err)
		}
		//resp, err := json.Marshal(companyList)
		if err != nil {
			log.Fatal(err)
		}
		a := List{}
		for k, v := range companyList {
			e := Entry{k, v}
			a = append(a, e)
		}

		sort.Sort(sort.Reverse(a))
		output, err := os.OpenFile("../../results/AIML/CommitRatio_"+j[0]+"_"+j[1]+".csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer output.Close()
		output.Write([]byte(fmt.Sprintf("%v,%v \n", j[0], j[1])))
		for _, t := range a {
			if t.name == "" {
				output.Write([]byte(fmt.Sprintf(`"記載なし" ,%v`, t.value)))
				output.Write([]byte("\n"))
			} else {
				output.Write([]byte(fmt.Sprintf(`"%v" ,%v`, t.name, t.value)))
				output.Write([]byte("\n"))
			}
		}
		time.Sleep(time.Second)
	}
}

func getComanyList(owner, name string) (results map[string]int, err error) {
	log.Printf("%s %s Start getComany List \n", owner, name)
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+owner+"/"+name+"/stats/contributors", nil)
	apiToken := "bearer " + os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Printf("Status Code is %d. Auto Retry after 20 second.", resp.StatusCode)
		time.Sleep(20 * time.Second)
		// Fix: Recursive processing. golang is not good at recursive process.
		return getComanyList(owner, name)
	}
	contributors := new(Contributors)
	err = json.Unmarshal(body, contributors)
	if err != nil {
		log.Fatal(err)
	}
	sort.Sort(sort.Reverse(contributors))
	result := make(map[string]int)
	for i, j := range *contributors {
		log.Printf("%d ", i)
		if i == 100 {
			break
		}
		companyName, err := getComanyName(j.Author.URL)
		if err != nil {
			log.Fatal(err)
		}
		result[companyName] += (*contributors)[i].Total
	}
	log.Printf("%s %s Finished getComany List \n", owner, name)
	return result, err
}

func getComanyName(url string) (comanyName string, err error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	apiToken := "bearer " + os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", apiToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	user := new(User)
	err = json.Unmarshal(body, user)
	if len(user.Company) > 0 {
		if string(user.Company[0]) == "@" {
			user.Company = user.Company[1:]
		}
	}

	return strings.ToLower(strings.TrimSpace(user.Company)), err
}

func getOssList() (ossList [][]string, err error) {
	file, err := os.Open("../inputData/AIML2.csv")
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

type Entry struct {
	name  string
	value int
}
type List []Entry

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return (l[i].name < l[j].name)
	} else {
		return (l[i].value < l[j].value)
	}
}
