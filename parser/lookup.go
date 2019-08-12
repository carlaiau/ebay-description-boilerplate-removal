package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type lookupResponseXML struct {
	XMLName xml.Name    `xml:"GetMultipleItemsResponse"`
	Items   []SuperItem `xml:"Item"`
}

type SuperSet struct {
	Items []SuperItem `xml:"Items"`
}

// Care of https://www.onlinetool.io/xmltogo/
type SuperItem struct {
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

func formBulkIDString(docs []Document) string {
	var bulkQuery strings.Builder
	numItems := 0
	for _, doc := range docs {
		if doc.EbayID != 0 {
			if numItems != 0 {
				bulkQuery.WriteString(",")
			}
			bulkQuery.WriteString(fmt.Sprintf("%d", doc.EbayID))
			numItems++
		}
	}
	return bulkQuery.String()
}

func getMultipleLookupURL(docs []Document) string {
	endpoint, _ := url.Parse("http://open.api.ebay.com/shopping")
	params := url.Values{}
	params.Add("callname", "GetMultipleItems")
	params.Add("responseencoding", "XML")
	params.Add("appid", appID)
	params.Add("siteid", "0")
	params.Add("version", "967")
	params.Add("ItemID", formBulkIDString(docs))
	params.Add("IncludeSelector", "Description,Details")

	endpoint.RawQuery = params.Encode()
	return endpoint.String()
}

func getLookupResponse(url string) lookupResponseXML {
	r, _ := http.Get(url)
	response, _ := ioutil.ReadAll(r.Body)
	v := lookupResponseXML{}
	err := xml.Unmarshal([]byte(response), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	return v
}

func process(docs []Document) {
	lookupURL := getMultipleLookupURL(docs)
	v := getLookupResponse(lookupURL)
	fmt.Printf("Length of returned lookups: %d\n", len(v.Items))
	itemIndex := 0
	for _, doc := range docs {
		if doc.EbayID != 0 && itemIndex < len(v.Items) {
			superItem := v.Items[itemIndex]
			superItem.OrigID = doc.ID
			superItem.OrigPrice = doc.Price
			superItem.OrigTitle = doc.Title
			superItem.OrigCategoryBreadcrumb = doc.CategoryTree
			superItem.OrigItemIDImageURL = doc.PictureURL
			itemIndex++
			processedDocuments = append(processedDocuments, superItem) // Add another Item to the array
		} else {
			// Product Not Found on Ebay
			superItem := SuperItem{
				OrigID:                 doc.ID,
				OrigPrice:              doc.Price,
				OrigTitle:              doc.Title,
				OrigCategoryBreadcrumb: doc.CategoryTree,
				OrigItemIDImageURL:     doc.PictureURL,
			}
			processedDocuments = append(processedDocuments, superItem) // Add another Item to the array
		}
	}

}
