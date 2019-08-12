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
		if item.Title != "" {
			foundCount++
			if strings.Trim(strings.ToLower(item.Title), " ") != strings.Trim(strings.ToLower(item.OrigTitle), " ") {
				inCorrectCount++
				fmt.Printf("Potential Bug #%d\n", inCorrectCount)
				fmt.Printf("%s\n", item.Title)
				fmt.Printf("%s\n\n", item.OrigTitle)
			}
			if item.Description != "" {
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

// 1760267	4.76	Adventures Of Cyclops And Phoenix Lot 4 Bks Set 1 2 3 4 VF NM X-Men Apocalypse  	 Collectibles > Comics > Modern Age (1992-Now) > Superhero > Other Modern Age Superheroes	http://i.ebayimg.com/00/s/MTA2N1gxNTk5/z/jw8AAOSwzvlW9Iu0/$_57.JPG?set_id=8800005007

// <OrigDocId>1000006</OrigDocId>
// <OrigPrice>12.74</OrigPrice>
// <OrigTitle>SCOUTS OF THAILAND - SENIOR GIRL SCOUTS (GUIDES) Epaulettes Patch (PAIR)</OrigTitle>
// <OrigCategoryBreadcrumb>Collectibles &gt; Historical Memorabilia &gt; Fraternal Organizations &gt; Boy Scouts &gt; Badges &amp; Patches &gt; Jamboree Patches</OrigCategoryBreadcrumb>
// <OrigItemIDImageURL>http://i.ebayimg.com/00/s/MTYwMFgxNDQ5/z/tn0AAOSwmfhX2sxd/$_1.JPG?set_id=880000500F</OrigItemIDImageURL>
