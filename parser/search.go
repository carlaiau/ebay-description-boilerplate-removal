package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type ResponseXML struct {
	Items []struct {
		ItemID              int    `xml:"itemId"`
		Title               string `xml:"title"`
		PrimaryCategoryName string `xml:"primaryCategory>categoryName"`
	} `xml:"searchResult>item"`
	TimeStamp string `xml:"timestamp"`
}

func getDocumentEbayID(d Document) int {
	ebayID := search(d, true)
	searches++
	return ebayID
	// Removal of backup findCompletedItems Query
}

// Trim
func cleanWord(s string) string {
	reg, err := regexp.Compile("[^a-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(strings.ToLower(s), "")
}

func search(d Document, expired bool) int {
	searchURL := createSearch(d.Title, expired)
	v := getResponse(searchURL)
	if len(v.Items) > 0 {
		i := v.Items[0]
		// Only Return item if title and cateogories equal
		if cleanWord(d.Title) == cleanWord(i.Title) && cleanWord(d.Category) == cleanWord(i.PrimaryCategoryName) {
			return i.ItemID
		}

	}
	return 0
}

func createSearch(keywords string, expired bool) string {
	endpoint, _ := url.Parse("http://svcs.ebay.com/services/search/FindingService/v1?")

	params := url.Values{}
	params.Add("SECURITY-APPNAME", appID)
	params.Add("keywords", keywords)
	params.Add("GLOBAL-ID", "EBAY-US")
	params.Add("RESPONSE-DATA-FORMAT", "XML")
	params.Add("REST-PAYLOAD", "")
	if expired {
		params.Add("OPERATION-NAME", "findCompletedItems")
		params.Add("SERVICE-VERSION", "1.7.0")
	} else {
		params.Add("OPERATION-NAME", "findItemsByKeywords")
		params.Add("SERVICE-VERSION", "1.0.0")
	}

	endpoint.RawQuery = params.Encode()

	return endpoint.String()

}

func getResponse(url string) ResponseXML {
	r, _ := http.Get(url)
	response, _ := ioutil.ReadAll(r.Body)

	//	2fmt.Printf("%s", response)
	v := ResponseXML{}
	err := xml.Unmarshal([]byte(response), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	return v
}
