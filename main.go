package main

/* Done
Filter the XML for documents that have descriptions defined.
After all XML documents are parsed and determined, then spit out a list of all document ID's.
Create a one-dimensional array of these document_ids
Traverse the documents.tsv file, for all documents that do not exist in this one-dimensional array, then create the XML with the corresponding structure.

Todo
Build a CLI that allows the user to define what field that want spat out of the XML. <- Edit raw and individual
*/

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/valyala/tsvreader"
)

var count = 0

type Set struct {
	XMLName xml.Name `xml:"Superset"`
	Items   []Item   `xml:"Item"`
}

// Care of https://www.onlinetool.io/xmltogo/
type Item struct {
	XMLName                xml.Name `xml:"Item"`
	OrigID                 string   `xml:"OrigDocId"`
	OrigPrice              string   `xml:"OrigPrice"`
	OrigTitle              string   `xml:"OrigTitle"`
	OrigCategoryBreadcrumb string   `xml:"OrigCategoryBreadcrumb"`
	OrigItemIDImageURL     string   `xml:"OrigItemIDImageURL"`
	BestOfferEnabled       string   `xml:"BestOfferEnabled"`
	BuyItNowPrice          struct {
		Text       string `xml:",chardata"`
		CurrencyID string `xml:"currencyID,attr"`
	} `xml:"BuyItNowPrice"`
	Description            string `xml:"Description"`
	ItemID                 int    `xml:"ItemID"`
	BuyItNowAvailable      string `xml:"BuyItNowAvailable"`
	ConvertedBuyItNowPrice struct {
		Text       string `xml:",chardata"`
		CurrencyID string `xml:"currencyID,attr"`
	} `xml:"ConvertedBuyItNowPrice"`
	EndTime                     string   `xml:"EndTime"`
	StartTime                   string   `xml:"StartTime"`
	ViewItemURLForNaturalSearch string   `xml:"ViewItemURLForNaturalSearch"`
	ListingType                 string   `xml:"ListingType"`
	Location                    string   `xml:"Location"`
	PaymentMethods              []string `xml:"PaymentMethods"`
	GalleryURL                  string   `xml:"GalleryURL"`
	PictureURL                  []string `xml:"PictureURL"`
	PostalCode                  string   `xml:"PostalCode"`
	PrimaryCategoryID           string   `xml:"PrimaryCategoryID"`
	PrimaryCategoryName         string   `xml:"PrimaryCategoryName"`
	Quantity                    string   `xml:"Quantity"`
	Seller                      struct {
		UserID                  string `xml:"UserID"`
		FeedbackRatingStar      string `xml:"FeedbackRatingStar"`
		FeedbackScore           string `xml:"FeedbackScore"`
		PositiveFeedbackPercent string `xml:"PositiveFeedbackPercent"`
		TopRatedSeller          string `xml:"TopRatedSeller"`
	} `xml:"Seller"`
	BidCount              string `xml:"BidCount"`
	ConvertedCurrentPrice struct {
		Text       string `xml:",chardata"`
		CurrencyID string `xml:"currencyID,attr"`
	} `xml:"ConvertedCurrentPrice"`
	CurrentPrice struct {
		Text       string `xml:",chardata"`
		CurrencyID string `xml:"currencyID,attr"`
	} `xml:"CurrentPrice"`
	ListingStatus         string   `xml:"ListingStatus"`
	QuantitySold          string   `xml:"QuantitySold"`
	ShipToLocations       []string `xml:"ShipToLocations"`
	Site                  string   `xml:"Site"`
	TimeLeft              string   `xml:"TimeLeft"`
	Title                 string   `xml:"Title"`
	HitCount              string   `xml:"HitCount"`
	PrimaryCategoryIDPath string   `xml:"PrimaryCategoryIDPath"`
	Country               string   `xml:"Country"`
	ReturnPolicy          struct {
		ReturnsAccepted                 string `xml:"ReturnsAccepted"`
		Refund                          string `xml:"Refund"`
		ReturnsWithin                   string `xml:"ReturnsWithin"`
		ShippingCostPaidBy              string `xml:"ShippingCostPaidBy"`
		InternationalReturnsAccepted    string `xml:"InternationalReturnsAccepted"`
		InternationalRefund             string `xml:"InternationalRefund"`
		InternationalReturnsWithin      string `xml:"InternationalReturnsWithin"`
		InternationalShippingCostPaidBy string `xml:"InternationalShippingCostPaidBy"`
	} `xml:"ReturnPolicy"`
	MinimumToBid struct {
		Text       string `xml:",chardata"`
		CurrencyID string `xml:"currencyID,attr"`
	} `xml:"MinimumToBid"`
	AutoPay                             string `xml:"AutoPay"`
	IntegratedMerchantCreditCardEnabled string `xml:"IntegratedMerchantCreditCardEnabled"`
	HandlingTime                        string `xml:"HandlingTime"`
	ConditionID                         string `xml:"ConditionID"`
	ConditionDisplayName                string `xml:"ConditionDisplayName"`
	GlobalShipping                      string `xml:"GlobalShipping"`
	QuantitySoldByPickupInStore         string `xml:"QuantitySoldByPickupInStore"`
	NewBestOffer                        string `xml:"NewBestOffer"`
	HighBidder                          struct {
		UserID             string `xml:"UserID"`
		FeedbackPrivate    string `xml:"FeedbackPrivate"`
		FeedbackRatingStar string `xml:"FeedbackRatingStar"`
		FeedbackScore      string `xml:"FeedbackScore"`
	} `xml:"HighBidder"`
	Storefront struct {
		StoreURL  string `xml:"StoreURL"`
		StoreName string `xml:"StoreName"`
	} `xml:"Storefront"`
	ExcludeShipToLocation   []string `xml:"ExcludeShipToLocation"`
	QuantityAvailableHint   string   `xml:"QuantityAvailableHint"`
	QuantityThreshold       string   `xml:"QuantityThreshold"`
	SKU                     string   `xml:"SKU"`
	PaymentAllowedSite      []string `xml:"PaymentAllowedSite"`
	TopRatedListing         string   `xml:"TopRatedListing"`
	ConditionDescription    string   `xml:"ConditionDescription"`
	Subtitle                string   `xml:"Subtitle"`
	SecondaryCategoryID     string   `xml:"SecondaryCategoryID"`
	SecondaryCategoryName   string   `xml:"SecondaryCategoryName"`
	SecondaryCategoryIDPath string   `xml:"SecondaryCategoryIDPath"`
	DiscountPriceInfo       struct {
		PricingTreatment    string `xml:"PricingTreatment"`
		OriginalRetailPrice struct {
			Text       string `xml:",chardata"`
			CurrencyID string `xml:"currencyID,attr"`
		} `xml:"OriginalRetailPrice"`
		SoldOneBay  string `xml:"SoldOneBay"`
		SoldOffeBay string `xml:"SoldOffeBay"`
	} `xml:"DiscountPriceInfo"`
	ReserveMet string `xml:"ReserveMet"`
	ProductID  struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
	} `xml:"ProductID"`
}

type Judgement struct {
	DocID   string
	Queries []Query
}
type Query struct {
	ID       int
	Relevant int
}

// Give this function the original XML from Ebay, and create an output XML only containing documents that
// had defined descriptions
func removeEmpties(in string, out string) {
	xmlFile, err := os.Open(in)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)
	set := &Set{}

	_ = xml.Unmarshal(data, &set)

	filteredData := &Set{}
	for _, i := range set.Items {
		if i.Description != "" {
			filteredData.Items = append(filteredData.Items, i) // Add another Item to the array
		}
	}
	filteredFile, _ := xml.MarshalIndent(filteredData, "", " ")
	_ = ioutil.WriteFile(out, filteredFile, 0644)
}

func individual(wg *sync.WaitGroup, filePath string, operation string) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)
	set := &Set{}
	_ = xml.Unmarshal(data, &set)

	switch operation {
	case "count":
		count += len(set.Items)

	case "docids":
		for _, i := range set.Items {
			fmt.Println(i.OrigID)
		}
	case "raw":
		for _, i := range set.Items {
			// All of the files within the collection have been tested to
			// ensure that these tags do not occur anywhere
			// This section needs to be CLIed
			fmt.Println("<DOC>")
			fmt.Printf("<DOCID>")
			fmt.Printf(i.OrigID)
			fmt.Printf("</DOCID>\n")
			fmt.Printf("<ORIGTITLE>%s</ORIGTITLE>\n", i.OrigTitle)
			fmt.Printf("<CATEGORY>")
			fmt.Printf(i.OrigCategoryBreadcrumb)
			fmt.Printf("</CATEGORY>\n")
			/*


				fmt.Println(i.Description)
			*/
			fmt.Println("</DOC>")
		}
	default:
		fmt.Println("We need an operation!")
	}

	xmlFile.Close()
	data = nil
	set = nil
	wg.Done()
}
func countItemsInFolder(in string) {
	files, _ := ioutil.ReadDir(in)
	var wg sync.WaitGroup
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		wg.Add(1)
		filePath := in + "/" + file.Name()
		go individual(&wg, filePath, "count")
	}
	wg.Wait()
	fmt.Printf("Total Items: %d\n", count)
}

func docIDsInFolder(in string) {
	files, _ := ioutil.ReadDir(in)
	var wg sync.WaitGroup
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		wg.Add(1)
		filePath := in + "/" + file.Name()
		go individual(&wg, filePath, "docids")
	}
	wg.Wait()
}

func raw(in string) {
	files, _ := ioutil.ReadDir(in)
	var wg sync.WaitGroup
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		wg.Add(1)
		filePath := in + "/" + file.Name()
		go individual(&wg, filePath, "raw")
	}
	wg.Wait()
}

// Need to load in an array of docids
// Need to sort this
// Need to then create a function to see if ID is within this array

func loadDocIDArray(in string) (numbers []int) {
	fd, err := os.Open(in)
	defer fd.Close()
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", in, err))
	}
	var line int
	for {
		_, err := fmt.Fscanf(fd, "%d\n", &line)

		if err != nil {
			fmt.Println(err)
			if err == io.EOF {
				sort.Ints(numbers[:]) // Sort these
				return
			}
			panic(fmt.Sprintf("Scan Failed %s: %v", in, err))

		}
		numbers = append(numbers, line)
	}
}

func intInArray(haystack []int, needle int) (found bool) {
	i := sort.Search(len(haystack), func(i int) bool { return haystack[i] >= needle })
	if i < len(haystack) && haystack[i] == needle {
		return true
	}
	return false
}

func createXMLFromMissingDocs(alreadySavedDocIDs []int, inputTSVFile string, out string) {
	var missingItems []Item

	data, _ := os.Open(inputTSVFile)
	r := tsvreader.New(data)

	for r.Next() {
		id := r.String()
		price := strings.Trim(r.String(), " ")
		title := strings.Trim(r.String(), " ")
		cat := strings.Trim(r.String(), " ")
		image := strings.Trim(r.String(), " ")

		item := Item{
			OrigID:                 id,
			OrigPrice:              price,
			OrigTitle:              title,
			OrigCategoryBreadcrumb: cat,
			OrigItemIDImageURL:     image,
		}
		needle, _ := strconv.Atoi(id)

		if !intInArray(alreadySavedDocIDs, needle) {
			fmt.Printf("%d not found\n", needle)
			missingItems = append(missingItems, item)
		}
	}
	// Create a Set struct, to create correct XML
	set := Set{
		Items: missingItems,
	}
	filteredFile, _ := xml.MarshalIndent(set, "", " ")
	_ = ioutil.WriteFile(out, filteredFile, 0644)

}

func convertJudgementsFromTSV(in string) {
	data, _ := os.Open(in)
	r := tsvreader.New(data)

	judgements := []Judgement{}

	for r.Next() {
		docID := r.String()
		judgement := Judgement{
			DocID:   docID,
			Queries: []Query{},
		}

		// Each of the following columns represent the binary relevancy of that columns query ID
		for i := 1; r.HasCols(); i++ {
			relevancy := r.Int()
			if relevancy != 1 {
				relevancy = 0 // Set -1's to 0.
			}
			query := Query{}
			query.ID = i
			query.Relevant = relevancy
			judgement.Queries = append(judgement.Queries, query)
		}

		judgements = append(judgements, judgement)
	}

	for _, judgement := range judgements {
		for _, query := range judgement.Queries {
			fmt.Printf("%d 0 %s %d\n", query.ID, judgement.DocID, query.Relevant)
		}
	}

	/* filteredFile, _ := xml.MarshalIndent(set, "", " ")
	_ = ioutil.WriteFile(out, filteredFile, 0644) */

}
func main() {
	opName := os.Args[1]
	switch opName {

	case "count":
		inputFolder := os.Args[2]
		countItemsInFolder(inputFolder)

	case "docids":
		inputFolder := os.Args[2]
		docIDsInFolder(inputFolder)

	case "filter":
		inputFile := os.Args[2]
		outputFile := os.Args[3]
		removeEmpties(inputFile, outputFile)

	case "raw":
		inputFolder := os.Args[2]
		raw(inputFolder)

	case "createmissingxml":
		alreadyFoundDocIDsFile := os.Args[2]
		allDocIDs := os.Args[3]
		outputXMLpath := os.Args[4]

		foundDocIds := loadDocIDArray(alreadyFoundDocIDsFile)
		fmt.Printf("Documents Found %d\n", len(foundDocIds))

		createXMLFromMissingDocs(foundDocIds, allDocIDs, outputXMLpath)

	case "convertjudgements":
		inFile := os.Args[2]
		convertJudgementsFromTSV(inFile)

	default:
		fmt.Println("Incorrect command line args")
	}
}
