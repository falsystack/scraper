package main

import (
	"log"
	"net/http"
	"os"
	"scraper/scrapper"
	"scraper/utils"
	"strings"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./resources")))
	http.HandleFunc("/scrape", scrapeJobs)
	http.ListenAndServe(":8080", nil)

}

func scrapeJobs(w http.ResponseWriter, req *http.Request) {
	defer os.Remove("jobs.csv")
	err := req.ParseForm()
	utils.HasErr(err)

	term := req.FormValue("term")
	log.Println(term)
	term = strings.ToLower(utils.CleanString(term))
	scrapper.Scrape(term)

	http.ServeFile(w, req, "jobs.csv")
}
