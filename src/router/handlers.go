package router

import (
	"encoding/json"
	"fmt"
	"git/processor"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Title      string
	Buttondata []Buttondata
}
type Buttondata struct {
	Name       string
	Id         string
	DateString string
	Repository string
}

const (
	TITLE = "CHM"
)

var templates = template.Must(template.ParseFiles("commits.html", "headAndNavbar.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	//	threeMonthAgo := time.Now().AddDate(0, -3, 0)
	//	query := processor.Query{Since: threeMonthAgo}
	vars := mux.Vars(r)

	query := getQueryFromVars(vars)

	queryResult := processor.GetCommits(query)

	commitData := []Buttondata{}
	for _, com := range queryResult {
		formatedDate := com.Time.Format(time.RFC822)[:10]
		commitData = append(commitData, Buttondata{com.Comment, com.Sha, formatedDate, com.Repo})
	}
	templates.ExecuteTemplate(w, "commits.html", Page{Title: TITLE, Buttondata: commitData}) //Page{Title: "Home"})
}

func shutdownCHM(w http.ResponseWriter, r *http.Request) {
	defer os.Exit(0)
}

func Log(w http.ResponseWriter, r *http.Request) {
	w.Write(LogBuffer.Bytes())
}

func AuthorsShowJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(processor.GetCachedAuthors()); err != nil {
		panic(err)
	}
}

func AuthorsShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing url %v", err), 500)
	}
	authorButtons := []Buttondata{}
	for _, author := range processor.GetCachedAuthors() {
		authorButtons = append(authorButtons, Buttondata{author, author, "", ""})
	}
	templates.ExecuteTemplate(w, "commits.html", Page{Title: TITLE, Buttondata: authorButtons})
}
func ReposShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(processor.GetCachedRepos()); err != nil {
		panic(err)
	}
}

func ShowSingleCommit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	vars := mux.Vars(r)
	sha, ok := vars["sha"]
	if !ok {
		// TODO
	}
	if err := json.NewEncoder(w).Encode(processor.GetSingleCommit(sha)); err != nil {
		panic(err)
	}
}

func CommitShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	vars := mux.Vars(r)

	query := getQueryFromVars(vars)

	queryResult := processor.GetCommits(query)

	if err := json.NewEncoder(w).Encode(queryResult); err != nil {
		panic(err)
	}
}

func SetConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jsonString := vars["jsonString"]
	jsonString = strings.Replace(jsonString, "+", "", -1)
	fmt.Fprintln(w, jsonString)
	if err := json.NewDecoder(strings.NewReader(jsonString)).Decode(&[]string{}); err != nil {
		fmt.Fprintln(w, err)
		return
	}
}
func GetConfig(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode("Test"); err != nil {
		panic(err)
	}
}

func getQueryFromVars(vars map[string]string) processor.Query {

	threeMonthAgo := time.Now().AddDate(0, -3, 0)
	query := processor.Query{Since: threeMonthAgo}

	authors, ok := vars["author"]
	if ok {
		if form, err := url.QueryUnescape(authors); err == nil {
			authors = form
		}
		query.Authors = strings.Split(authors, ";")
	}

	repos, ok := vars["repo"]
	if ok {
		if form, err := url.QueryUnescape(repos); err == nil {
			repos = form
		}
		query.Repos = strings.Split(repos, ";")
	}

	since, ok := vars["date"]
	if ok {
		if d, err := time.Parse(time.RFC3339, since); err == nil {
			query.Since = d
		}

	}
	return query
}
