package main

// Filter the XML for documents that have descriptions defined.
// Build a CLI that allows the user to define what field that want spat out of the XML.
// After all XML documents are parsed and determined, then spit out a list of all document ID's.
// Create a one-dimensional array of these document_ids
// Traverse the documents.tsv file, for all documents that do not exist in this one-dimensional array, then create the XML with the corresponding structure.

/*
i.xml contains 34057 items
f.xml contains 34580 items
h.xml contains 34067 items
e.xml contains 34679 items
b.xml contains 35373 items
g.xml contains 34367 items
a.xml contains 35900 items
d.xml contains 34755 items
c.xml contains 35302 items
Total Items: 313080

j.xml contains 33789 items
p.xml contains 33135 items
r.xml contains 32710 items
q.xml contains 32922 items
n.xml contains 33265 items
o.xml contains 33184 items
m.xml contains 33590 items
l.xml contains 33578 items
k.xml contains 33653 items
_missing.xml contains 65800 items
Total Items: 365626

Total total: 678706



*/
import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
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

func filterEmpties(in string, out string) {
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

func individual(wg *sync.WaitGroup, filePath string, docid bool) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)
	set := &Set{}
	_ = xml.Unmarshal(data, &set)

	if docid {
		for _, i := range set.Items {
			fmt.Println(i.OrigID)
		}
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
		go individual(&wg, filePath, false)
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
		go individual(&wg, filePath, true)

	}
	wg.Wait()

}
func main() {
	opName := os.Args[1]
	switch opName {
	case "count":
		inputFolder := os.Args[2]
		countItemsInFolder(inputFolder)

	case "docid":
		inputFolder := os.Args[2]
		docIDsInFolder(inputFolder)

	case "filter":
		inputFile := os.Args[2]
		outputFile := os.Args[3]
		filterEmpties(inputFile, outputFile)
	default:
		fmt.Println("Incorrect command line args")
	}
}
