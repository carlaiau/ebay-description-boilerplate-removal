package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"
)

var appID string
var processedDocuments = []SuperItem{}
var queue []Document
var searches = 0
var expiredSearches = 0
var expiredFinds = 0
var totalQueued = 0
var queuedFoundDocuments = 0
var totalProcessedDocuments = 0
var totalFoundDocuments = 0

func addToQueue(doc Document) {
	ebayID := getDocumentEbayID(doc)
	if ebayID != 0 {
		queuedFoundDocuments++
		totalFoundDocuments++
		doc.EbayID = ebayID
		fmt.Printf("%s Y\n", doc.ID)
	} else {
		fmt.Printf("%s N\n", doc.ID)
	}
	queue = append(queue, doc)
}

func addToQueueConcurrently(wg *sync.WaitGroup, doc Document) {
	rand.Seed(time.Now().UnixNano())
	sleepTime := rand.Intn(200) + 100
	time.Sleep(time.Duration(sleepTime) * time.Millisecond) // Random sleep of 100ms - 300ms, to prevent API ban

	ebayID := getDocumentEbayID(doc)
	totalProcessedDocuments++
	if ebayID != 0 {
		queuedFoundDocuments++
		totalFoundDocuments++
		doc.EbayID = ebayID
		fmt.Printf("Y %d\n", totalProcessedDocuments)
	} else {
		fmt.Printf("N %d\n", totalProcessedDocuments)
	}
	queue = append(queue, doc)
	wg.Done()
}

// Need to comment
func bigKahuna(unProcessedDocuments []Document) {
	for totalQueued < len(unProcessedDocuments)-20 { // This is a while loop
		var wg sync.WaitGroup // Make a Wait Group
		fmt.Printf("Concurrency\n[\n")
		for i := 1; i <= 20; i++ { // Add 20 goroutines to the Waitgroup
			doc := unProcessedDocuments[totalQueued]
			totalQueued++
			wg.Add(1)
			go addToQueueConcurrently(&wg, doc)
		}
		wg.Wait()
		fmt.Printf("]\n")
		// After the waitgroup is completed, sequentially add more documents to queue
		// until the queu contains 20 correctly identified ebay Items or we've reached
		// the end of the document array
		for queuedFoundDocuments < 20 && totalQueued < len(unProcessedDocuments) {
			doc := unProcessedDocuments[totalQueued]
			addToQueue(doc)
			totalQueued++
		}
		fmt.Printf("Processing Queue...\n")
		process(queue)
		fmt.Printf("Complete\n")
		queue = []Document{} // Empty the Queue
		queuedFoundDocuments = 0
	}

	// If we have a small number of unprocessedDocuments remaining
	for totalQueued < len(unProcessedDocuments) {
		doc := unProcessedDocuments[totalQueued]
		fmt.Printf("Processing Queue...\n")
		addToQueue(doc)
		fmt.Printf("Complete\n")
		totalQueued++

	}
	fmt.Printf("Process Queue!\n")
	process(queue)
}

func main() {
	opName := os.Args[1]
	switch opName {
	case "scrape": // Scrape data from Ebay
		appID = os.Args[2] // AppKey must be fed in via command line
		inputFilePath := os.Args[3]
		outputFilePath := os.Args[4]

		unProcessedDocuments := getDocs(inputFilePath) // Load Documents
		bigKahuna(unProcessedDocuments)

		fmt.Printf("Finished! %d of %d found!\t%d Searchs\t%d Expired Searchs\t%d Expired Finds\n", totalFoundDocuments, len(unProcessedDocuments), searches, expiredSearches, expiredFinds)
		file, _ := xml.MarshalIndent(SuperSet{Items: processedDocuments}, "", "    ")
		_ = ioutil.WriteFile(outputFilePath, file, 0644)

	case "test":
		testPathName := os.Args[2]
		getMetrics(testPathName)

	case "missing":
		testPathName := os.Args[2]
		getMissing(testPathName)
	default:
		fmt.Println("Incorrect command line args")
	}

}
