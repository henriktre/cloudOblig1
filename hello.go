package hello

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

//Data contains all the dataz
/*
 * Struct Data
 * @value Project [string] Something about string
 * @value Owner [string] Something about string
 * @value Committer [string] Something about string
 * @value Commits [int] Something about int
 * @value Language [string] Something about string
 */
type Data struct {
	Project   string   `json:"project"`
	Owner     string   `json:"owner"`
	Committer string   `json:"committer"`
	Commits   int      `json:"commits"`
	Language  []string `json:"language"`
}

// User containing data about the user
type User struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}

// Contributors data about contributions
type Contributors struct {
	Users []User
}

// Owner containing data about the owner
type Owner struct {
	Login string `json:"login"`
}

// Repo struct containing data about repo
type Repo struct {
	Project      string `json:"name"`
	Name         string `json:"full_name"`
	Owner        Owner  `json:"owner"`
	Contributors string `json:"contributors_url"`
	Languages    string `json:"languages_url"`
}

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler)
	r.HandleFunc("/projectinfo/v1/github.com/{username}/{reponame}", handlerRepo)
	http.Handle("/", r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Usage: write /projectinfo/v1/github.com/USERNAME/REPONAME in your get request")
}

func handlerRepo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.Get("https://api.github.com/repos/" + vars["username"] + "/" + vars["reponame"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Fprintf(w, "HTTP GET returned status %v", resp.Status)

	repo1 := Repo{}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	jsonError := json.Unmarshal(body, &repo1)

	if jsonError != nil {
		log.Fatal(jsonError)
	}

	cons := getContributors(w, r, repo1)
	langs := getLanguages(w, r, repo1)

	// arr, err := json.Marshal(langs)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	data := Data{}

	for k := range langs {
		data.Language = append(data.Language, k)
	}
	data.Project = string(repo1.Project)
	data.Owner = repo1.Owner.Login
	data.Committer = cons.Users[0].Login
	data.Commits = cons.Users[0].Contributions

	output, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, string(output))
}

func getContributors(w http.ResponseWriter, r *http.Request, repo1 Repo) Contributors {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.Get(repo1.Contributors)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// fmt.Fprintf(w, "HTTP GET returned status %v", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		log.Fatal(readErr)
	}

	cons := Contributors{}

	jsonError := json.Unmarshal(body, &cons.Users)

	if jsonError != nil {
		log.Fatal(jsonError)
	}

	return cons
}

func getLanguages(w http.ResponseWriter, r *http.Request, repo1 Repo) map[string]interface{} {
	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)
	resp, err := client.Get(repo1.Languages)

	if err != nil {
		var error interface{}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		error = err.Error()
		log.Fatal(error)
	}

	// fmt.Fprintf(w, "HTTP GET returned status %v", resp.Status)

	body, readErr := ioutil.ReadAll(resp.Body)

	var langs map[string]interface{}

	jsonError := json.Unmarshal(body, &langs)

	if readErr != nil || jsonError != nil {
		log.Fatal(readErr)
	}

	return langs
}
