package main

import (
	"encoding/csv"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var baseURL = "https://www.wantedly.com/projects?type=mixed&occupation_types%5B%5D=jp__engineering&keywords%5B%5D=golang&page=1"

type Job struct {
	id      string
	title   string
	summary string
}

func main() {
	var jobs []Job
	pages := getPages()
	for i := 1; i <= pages; i++ {
		extractJobs := getPage(i)
		jobs = append(jobs, extractJobs...)
	}
	writeJobs(jobs)
}

func writeJobs(jobs []Job) {
	file, err := os.Create("jobs.csv")
	hasErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Id", "Title", "Summary"}
	err = w.Write(headers)
	hasErr(err)

	for _, job := range jobs {
		err := w.Write([]string{
			"https://www.wantedly.com/projects/" + job.id,
			job.title,
			job.summary,
		})
		hasErr(err)
	}
}

func getPage(page int) []Job {
	var jobs []Job
	pageURL := baseURL + "&page=" + strconv.Itoa(page)
	fmt.Println("リクエストURL：", pageURL)

	resp, err := http.Get(pageURL)
	hasErr(err)
	hasErrCodes(resp)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	hasErr(err)

	doc.Find(".projects-index-single").Each(func(i int, s *goquery.Selection) {
		jobs = append(jobs, extractJob(s))
	})
	return jobs
}

func extractJob(s *goquery.Selection) Job {
	id, _ := s.Attr("data-project-id")
	title := cleanString(s.Find(".project-title").Text())
	summary := cleanString(s.Find(".project-excerpt").Text())

	return Job{
		id:      id,
		title:   title,
		summary: summary,
	}
}

func getPages() int {
	pages := 0
	resp, err := http.Get(baseURL)
	hasErr(err)
	hasErrCodes(resp)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	hasErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		// paginationクラス中の全てのaタグを探し出す。
		pages = s.Find("a").Length()
	})
	return pages
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func hasErrCodes(resp *http.Response) {
	if resp.StatusCode != 200 {
		log.Fatalln("request failed with status : ", resp.StatusCode)
	}
}

func hasErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
