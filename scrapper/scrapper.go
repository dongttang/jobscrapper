package scrapper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// JobInfo includes each information of content
type JobInfo struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

type pageButtonInfo struct {
	url       string
	buttonNum int
}

// GetPageNum returns last number of pages
func GetPageNum(baseURL string) (pageNum int) {

	lastPageURL := urlBuilder(baseURL, 9999)

	fmt.Println("GetPageNum lastPageUrl is ", lastPageURL)

	doc := getPageDocObject(lastPageURL)

	doc.Find(".pagination-list").Each(func(i int, s *goquery.Selection) {

		pageNum, _ = strconv.Atoi(s.Find("b").Text())

	})

	if pageNum == 0 {
		log.Fatal("Total size of page is 0... Failed to load pages.")
	} else {
		log.Printf("There are %d pages.", pageNum)
	}

	return
}

// RequestJobInfoArray sends JabInfo type array to chanel
func RequestJobInfoArray(baseURL string, targetPage int, mainChannel chan<- []JobInfo) {

	channel := make(chan JobInfo)

	// After waiting response, jobCardArray is sent to channel
	jobCardArray := []JobInfo{}

	targetURL := urlBuilder(baseURL, targetPage)

	doc := getPageDocObject(targetURL)

	jobCardNumCounter := 0
	doc.Find(".jobsearch-SerpJobCard unifiedRow row result clickcard").Each(func(i int, s *goquery.Selection) {

		go scrapJob(s, channel)
		jobCardNumCounter++
		fmt.Println("jobCardNumCounter:", jobCardNumCounter)

	})

	for i := 0; i < jobCardNumCounter; i++ {

		jobCardArray = append(jobCardArray, <-channel)

	}

	mainChannel <- jobCardArray
}

func scrapJob(s *goquery.Selection, channel chan<- JobInfo) {

	jobCard := new(JobInfo)
	jobCard.title = s.Find("title").Text()
	jobCard.location = s.Find("sjcl").Text()
	jobCard.summary = s.Find("summary").Text()

	channel <- *jobCard

}

// getPageDocObject returns *goquery.Document object using goquery framework
func getPageDocObject(targetURL string) (doc *goquery.Document) {

	res, err := http.Get(targetURL)

	errorCheck(err)

	statusCodeCheck(res)

	defer res.Body.Close()

	doc, err = goquery.NewDocumentFromReader(res.Body)

	errorCheck(err)

	return
}

func urlBuilder(baseURL string, targetPage int) (targetURL string) {
	targetURL = baseURL + "&start=" + strconv.Itoa(targetPage)
	return
}

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func statusCodeCheck(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
}
