package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type SuperItemsLocal struct {
	//Set []SuperItem `xml:"SuperSet"`
	Items []SuperItem `xml:"Item"`
}

func loadSupersetFromFile(filePath string) SuperItemsLocal {
	data, openErr := os.Open(filePath)
	if openErr != nil {
		fmt.Println("Open Error")
		panic(openErr)
	}
	response, _ := ioutil.ReadAll(data)
	v := SuperItemsLocal{}
	err := xml.Unmarshal([]byte(response), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return v

}

func getMetrics(filePath string) {
	var itemCount = 0
	var foundCount = 0
	var inCorrectCount = 0
	var descriptionCount = 0
	xml := loadSupersetFromFile(filePath)
	for _, item := range xml.Items {
		itemCount++
		if item.Title != "" { // The Title is defined
			foundCount++
			if strings.Trim(strings.ToLower(item.Title), " ") != strings.Trim(strings.ToLower(item.OrigTitle), " ") {
				inCorrectCount++
				fmt.Printf("Potential Bug #%d\n", inCorrectCount)
				fmt.Printf("%s\n", item.Title)
				fmt.Printf("%s\n\n", item.OrigTitle)
			}
			if item.Description != "" { // The Description is defined
				descriptionCount++
			}
		}
	}

	fmt.Printf("Counts\n Total: %d\tFound: %d\t Descriptions: %d\tTitle Not Identical:%d\n", itemCount, foundCount, descriptionCount, inCorrectCount)
}

func getMissing(filePath string) {
	xml := loadSupersetFromFile(filePath)
	for _, item := range xml.Items {
		if item.Title == "" {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\n", item.OrigID, item.OrigPrice, item.OrigTitle, item.OrigCategoryBreadcrumb, item.OrigItemIDImageURL)
			//if strings.Trim(strings.ToLower(item.Title), " ") != strings.Trim(strings.ToLower(item.OrigTitle), " ") {
			//	fmt.Printf("%s\t%s\t%s\t%s\t%s\n", item.OrigID, item.OrigPrice, item.OrigTitle, item.OrigCategoryBreadcrumb, item.OrigItemIDImageURL)
			//}
		}
	}
}
