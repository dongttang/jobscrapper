package main

import (
	"fmt"

	"github.com/dongttang/jobscrapper/scrapper"
)

// BaseURL is base url
const BaseURL string = "https://kr.indeed.com/jobs?q=python&l="

func main() {

	var jobList []scrapper.JobInfo

	mainChannel := make(chan []scrapper.JobInfo)

	totalPageNum := scrapper.GetPageNum(BaseURL)

	fmt.Println("main function", totalPageNum)

	for i := 0; i < totalPageNum; i++ {

		go scrapper.RequestJobInfoArray(BaseURL, i, mainChannel)

	}

	for i := 0; i < totalPageNum; i++ {
		extractedJob := <-mainChannel
		jobList = append(jobList, extractedJob...)
	}

	fmt.Println("Completed. size of job list is:", len(jobList))

}
