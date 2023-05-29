package scrapper

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	. "scraper/utils"
	"strconv"
)

type job struct {
	id      string
	title   string
	summary string
}

func Scrape(term string) {
	baseURL := "https://www.wantedly.com/projects?type=mixed&occupation_types%5B%5D=jp__engineering&keywords%5B%5D=" + term + "&page=1"
	var jobs []job
	c := make(chan []job)
	pages := getPages(baseURL)
	for i := 1; i <= pages; i++ {
		go getPage(i, c, baseURL)
	}

	for i := 0; i < pages; i++ {
		jobs = append(jobs, <-c...)
	}
	writeJobs(jobs)
}

func writeJobs(jobs []job) {
	file, err := os.Create("jobs.csv")
	HasErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Id", "Title", "Summary"}
	err = w.Write(headers)
	HasErr(err)

	for _, job := range jobs {
		err := w.Write([]string{
			"https://www.wantedly.com/projects/" + job.id,
			job.title,
			job.summary,
		})
		HasErr(err)
	}
}

func getPage(page int, mChan chan<- []job, url string) {
	var jobs []job
	c := make(chan job)
	pageURL := url + "&page=" + strconv.Itoa(page)
	fmt.Println("リクエストURL：", pageURL)

	resp, err := http.Get(pageURL)
	HasErr(err)
	hasErrCodes(resp)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	HasErr(err)

	jobSection := doc.Find(".projects-index-single")
	jobSection.Each(func(i int, s *goquery.Selection) {
		go extractJob(s, c)
	})

	for i := 0; i < jobSection.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	mChan <- jobs
}

func extractJob(s *goquery.Selection, c chan<- job) {
	id, _ := s.Attr("data-project-id")
	title := CleanString(s.Find(".project-title").Text())
	summary := CleanString(s.Find(".project-excerpt").Text())

	c <- job{
		id:      id,
		title:   title,
		summary: summary,
	}
}

func getPages(url string) int {
	pages := 0
	resp, err := http.Get(url)
	HasErr(err)
	hasErrCodes(resp)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	HasErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// paginationクラス中の全てのaタグを探し出す。
		pages = s.Find("a").Length()
	})
	return pages
}

func hasErrCodes(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatalln("request failed with status : ", resp.StatusCode)
	}
}
