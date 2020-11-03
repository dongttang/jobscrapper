package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/dongttang/jobscrapper/scrapper"
)

// BaseURL is base url
const BaseURL string = "https://kr.indeed.com/jobs?q=python&l="

func main() {

	// for measuring execution time
	start := time.Now()

	var jobList []scrapper.JobCard

	mainChannel := make(chan []scrapper.JobCard)

	totalPageNum := scrapper.GetPageNum(BaseURL)

	for i := 0; i < totalPageNum; i++ {

		go scrapper.RequestJobInfoArray(BaseURL, i, mainChannel)
	}

	for i := 0; i < totalPageNum; i++ {

		extractedJob := <-mainChannel

		jobList = append(jobList, extractedJob...)
	}

	file, err := os.Create("./jobList.csv")

	if err != nil {

		panic(err)
	}

	wr := csv.NewWriter(bufio.NewWriter(file))

	wr.Write([]string{"Title", "Location", "Summary"})

	for _, job := range jobList {

		wr.Write([]string{job.Title, job.Location, job.Summary})
	}

	wr.Flush()

	// for measuring execution time
	elapsed := time.Since(start)

	log.Printf("%d items are scrawled in %s sec.", len(jobList), elapsed)

	log.Println("Saved path: ./jobList.csv")

}
