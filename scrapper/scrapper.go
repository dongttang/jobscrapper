package scrapper

import (
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// JobCard includes information about a job
type JobCard struct {
	ID       string
	Title    string
	Location string
	Salary   string
	Summary  string
}

// GetPageNum returns last number of pages
func GetPageNum(baseURL string) (pageNum int) {

	// dummyBigNumber is used just for getting last page from target URL
	const dummyBigNumber int = 99999

	lastPageURL := urlBuilder(baseURL, dummyBigNumber)

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

// RequestJobInfoArray sends JobInfo type array to chanel
func RequestJobInfoArray(baseURL string, targetPageNum int, mainChannel chan<- []JobCard) {

	channel := make(chan JobCard)

	// After waiting response, jobCardArray is sent to channel
	jobCardArray := new([]JobCard)

	targetURL := urlBuilder(baseURL, targetPageNum)

	doc := getPageDocObject(targetURL)

	log.Println("Target url is : ", targetURL)

	jobCardNumCounter := 0

	doc.Find(".jobsearch-SerpJobCard").Each(func(i int, s *goquery.Selection) {

		go scrapJob(s, channel)

		jobCardNumCounter++
	})

	for i := 0; i < jobCardNumCounter; i++ {

		extractedJobCard := <-channel

		*jobCardArray = append(*jobCardArray, extractedJobCard)
	}

	mainChannel <- *jobCardArray
}

func scrapJob(s *goquery.Selection, channel chan<- JobCard) {

	jobCard := new(JobCard)

	jobCard.Title = s.Find(".title").Find("a").Text()

	jobCard.Location = s.Find(".sjcl").Find(".accessible-contrast-color-location").Text()

	jobCard.Summary = s.Find(".summary").Text()

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

	targetURL = baseURL + "&start=" + strconv.Itoa(targetPage*10)

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
