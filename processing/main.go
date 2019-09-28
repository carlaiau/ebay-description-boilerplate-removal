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
	"encoding/json"
	"regexp"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"bytes"
	"math"
)

type Set struct {
	XMLName xml.Name `xml:"Superset"` // Remember about the Capitalisatio of this
	Items   []Item   `xml:"Item"`
}

type Count struct {
	total       int
	titles      int
	description int
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

type IndexLine struct{
	ID string `json:"_id"`
}
type Jsonl struct{
	indexLine IndexLine
	doc ConvertedDoc
}
type ConvertedDoc struct {
	Docno		string   
	Title	    string
	Description string
}

/*
 *
 * Load a set given a filePath
 *
 */
func loadSet(filePath string) *Set {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer xmlFile.Close()

	data, _ := ioutil.ReadAll(xmlFile)
	set := &Set{}
	err = xml.Unmarshal(data, &set)
	if err != nil {
		fmt.Printf("%s", err)
	}
	xmlFile.Close()
	return set
}

/*
 *
 * Given a folder of XML dumps for each file,
 * remove the items that did not have a title
 * defined, and save the remaining items in an
 * identically named file in the output folder
 *
 */
func removeEmpties(in string, out string) {
	files, _ := ioutil.ReadDir(in)
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		inPath := "./" + in + "/" + file.Name()
		outPath := "./" + out + "/" + file.Name()

		set := loadSet(inPath)
		filtered := &Set{}

		for _, i := range set.Items {
			if i.Title != "" {
				filtered.Items = append(filtered.Items, i) // Add another Item to the array
			}
		}
		filteredFile, _ := xml.MarshalIndent(filtered, "", " ")
		_ = ioutil.WriteFile(outPath, filteredFile, 0644)

	}

}

/*
 *
 * Create basic analysis of a folder of XMLs.
 * Prints to std.out the number of items, titles and descriptions
 *
 */
func countItemsInFolder(in string) {
	files, _ := ioutil.ReadDir(in)
	totalCounts := Count{0, 0, 0}
	fmt.Printf("File Counts\nFilename\tFiles\tTitle\tDescription\n")
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		filePath := "./" + in + "/" + file.Name()
		set := loadSet(filePath)

		fileCounts := Count{len(set.Items), 0, 0}

		for _, i := range set.Items {
			if i.Title != "" {
				fileCounts.titles++
			}
			if i.Description != "" {
				fileCounts.description++
			}
		}

		totalCounts.total += fileCounts.total
		totalCounts.titles += fileCounts.titles
		totalCounts.description += fileCounts.description
		fmt.Printf("%s\t\t%d\t%d\t%d\n", file.Name(), fileCounts.total, fileCounts.titles, fileCounts.description)
	}
	fmt.Printf("\tTotals\nFiles\tTitle\tDescription\n")
	fmt.Printf("%d\t%d\t%d\n", totalCounts.total, totalCounts.titles, totalCounts.description)
}

/*
 *
 * outputs to std.out every docID found in the folder.
 * Used to generate a one dimensional array of the document IDs
 * that were found in our scrape, so that we can determine the
 * missing documents that need reintroduced via documents.tsv
 *
 */
func docIDsInFolder(in string) {
	files, _ := ioutil.ReadDir(in)
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		filePath := in + "/" + file.Name()
		set := loadSet(filePath)
		for _, i := range set.Items {
			fmt.Println(i.OrigID)
		}
	}
}

/*
 *
 * Creates files that are ready to be indexed by ATIRE
 * The initial Files only contain docID, title, category and description
 * Adding other fields from the original struct is straight forward
 */
func createRaw(in string, out string, fields string) {
	files, _ := ioutil.ReadDir(in)
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		inPath := in + "/" + file.Name()
		outPath := out + "/" + file.Name()
		set := loadSet(inPath)

		var b strings.Builder
		for _, i := range set.Items {
			fmt.Fprintf(&b, "<DOC>\n")
			fmt.Fprintf(&b, "<DOCNO>")
			fmt.Fprintf(&b, i.OrigID)
			fmt.Fprintf(&b, "</DOCNO>\n")

			if strings.Contains(fields, "t") {
				fmt.Fprintf(&b, "<ORIGTITLE>%s</ORIGTITLE>\n", i.OrigTitle)
			}

			if strings.Contains(fields, "c") {
				fmt.Fprintf(&b, "<CATEGORY>")
				fmt.Fprintf(&b, i.OrigCategoryBreadcrumb)
				fmt.Fprintf(&b, "</CATEGORY>\n")
			}

			if strings.Contains(fields, "d") {
				fmt.Fprintf(&b, "<CSDESCRIPTION>\n")
				fmt.Fprintf(&b, i.Description)
				fmt.Fprintf(&b, "</CSDESCRIPTION>\n")
			}
			
			fmt.Fprintf(&b, "\n</DOC>\n")
		}
		_ = ioutil.WriteFile(outPath, []byte(b.String()), 0644)

	}
}


func individualFilesDescriptionsOnly(in string, out string){
	files, _ := ioutil.ReadDir(in)
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		inPath := in + "/" + file.Name()
		set := loadSet(inPath)
		for _, i := range set.Items {
			
			dirPath := out + "/" + strings.Split(file.Name(),".")[0]
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				err = os.MkdirAll(dirPath, 0755)
				if err != nil {
					panic(err)
				}
			}
			outPath :=  dirPath + "/" + i.OrigID
			err := ioutil.WriteFile(outPath, []byte(i.Description), 0644)
			if err != nil{
				fmt.Printf("%s", err)
			}
		}
	}
}


/*
 *
 * Used for converting the raw Atire-indexable files back into JSON so that
 * we can get them into elasticsearch. Because once they have been generated raw
 * they are no longer valid XML, so XML parsing/unmarshalling is not usable.
 * Simple string split, regex, and then JSON encode.
 * Creates multiple files to assist in indexing into Elastic Search as well
 * as reducing RAM requirements 
 *
 */
func convertToJSONL(in string, out string){
	b, err := ioutil.ReadFile(in) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	docs := strings.Split(string(b), "<DOC>")[1:]

	num_buckets := 42
	bucket_size := len(docs)/num_buckets
	
	for i := 0;  i < num_buckets; i++ {
		//convertedDocs := []ConvertedDoc{}
		elementsToWrite := []Jsonl{}

		for index, doc := range docs[i * bucket_size:(i + 1) * bucket_size]{	
			docnoRegex := regexp.MustCompile(`<DOCNO>(\S+)</DOCNO>`)
			titleRegex := regexp.MustCompile(`<ORIGTITLE>([\S\s]+)</ORIGTITLE`)
			descRegex := regexp.MustCompile(`<CSDESCRIPTION>([\S\s]+)</CSDESCRIPTION>`)

			docno := docnoRegex.FindSubmatch([]byte(doc))[1]
			title := titleRegex.FindSubmatch([]byte(doc))[1]
			descriptionMatches := descRegex.FindSubmatch([]byte(doc))
			description := []byte(``)
			if(len(descriptionMatches) > 0){
				description = descriptionMatches[1] 
			} 
			
			doc := ConvertedDoc{
				Docno: string(docno),
				Title: string(title),
				Description: string(description),
			}
			indexLine:= IndexLine{
				ID: string(docno),
			}
			element := Jsonl{
				indexLine: indexLine,
				doc: doc,
			}
			elementsToWrite = append(elementsToWrite, element)
			//convertedDocs = append(convertedDocs, doc)

			if index % 5000 == 0 && index != 0{
				fmt.Printf("%d\n", index)
			}
		}
		
		var buffer bytes.Buffer
		for _, element := range elementsToWrite{
			indexLine, err := json.Marshal(element.indexLine)
			docLine, _ := json.Marshal(element.doc)

			if err != nil{
				panic(err)
			}

			// Need to do the hack to get it to index
			buffer.WriteString("{\"index\": " + string(indexLine) + "}\n" + string(docLine) + "\n")
		}

		err := ioutil.WriteFile(fmt.Sprintf("%s-%d.out", out, i), []byte(buffer.String()), 0644)
		if err != nil{
			fmt.Printf("Err!")
			panic(err)
		} else{
			fmt.Printf("\nFile writen\n\n")
		}
		
	}
}

/* 
	Used for getting Statistics of each index.
	Mean and standard deviation of the character lengths
*/
func getStats(in string){
	b, err := ioutil.ReadFile(in) // just pass the file name
	fileNamePath := strings.Split(in, "/")
	if err != nil {
		fmt.Print(err)
	}
	docs := strings.Split(string(b), "<DOC>")[1:]

	lengths := []uint32{}

	for _, doc := range docs{	
		descRegex := regexp.MustCompile(`<CSDESCRIPTION>([\S\s]+)</CSDESCRIPTION>`)
		descriptionMatches := descRegex.FindSubmatch([]byte(doc))
		if(len(descriptionMatches) > 0){
			lengths = append(lengths, uint32(len(descriptionMatches[1])))
		} else{
			lengths = append(lengths, 0)
		}
	}

	num_docs := uint64(len(lengths));
	var num_chars uint64
	for _, i := range lengths{
		num_chars += uint64(i)
	}

	mean := float64(num_chars) / float64(num_docs)

	var variance float64
	for _, i := range lengths{
		variance += math.Pow(float64(i) - mean, 2)
	}
	variance /= float64(num_docs)
	std_dev := math.Pow(variance, 0.5)

	fmt.Printf("\nFile:\t\t%s\nDocuments:\t%d\nCharacters:\t%d\nAverage:\t%.2f\nVariance:\t%.2f\nStd_Dev:\t%.2f\n", 
		fileNamePath[len(fileNamePath) - 1], num_docs, num_chars, mean, variance, std_dev)


	
}

/*
 *
 *
 *
 */
func createDocESDeleteJSONL(in string){
	docidBytes, _ := ioutil.ReadFile(in)
	IDStringArray := strings.Split(string(docidBytes), "\n")

	for _, id := range IDStringArray{
		fmt.Printf("{ \"delete\" : {\"_id\" : \"%s\" } }\n", id)
	}
}

func main() {
	opName := os.Args[1]
	switch opName {

	case "countItems":
		inputFolder := os.Args[2]
		countItemsInFolder(inputFolder)

	// Needs Piped
	case "getDocIDs":
		inputFolder := os.Args[2]
		docIDsInFolder(inputFolder)

	case "removeEmpties":
		inputFolder := os.Args[2]
		outputFolder := os.Args[3]
		removeEmpties(inputFolder, outputFolder)

	// Needs piped
	case "createMissingXML":
		alreadyFoundDocIDsFile := os.Args[2]
		originalDocuments := os.Args[3]
		outputXMLpath := os.Args[4]

		foundDocIDs := loadDocIDArray(alreadyFoundDocIDsFile)
		fmt.Printf("Documents Found %d\n", len(foundDocIDs))

		createXMLFromMissingDocs(foundDocIDs, originalDocuments, outputXMLpath)

	// Loads a whole folder in
	case "createRaw":
		inputFolder := os.Args[2]
		outputFolder := os.Args[3]
		fields 	:= os.Args[4]
		createRaw(inputFolder, outputFolder, fields)

	// Needs piped
	case "convertJudgements":
		inFile := os.Args[2]
		convertJudgementsFromTSV(inFile)

	case "individualFiles":
		inFolder := os.Args[2]
		outFolder := os.Args[3]
		individualFilesDescriptionsOnly(inFolder, outFolder)

	case "convertToJSONL":
		inFile := os.Args[2]
		outFile := os.Args[3]
		convertToJSONL(inFile, outFile)

	case "getStats":
		inFile := os.Args[2]
		getStats(inFile)


	case "outputEmptyDescriptions":	
		indexFile := os.Args[2]
		outputEmptyDescriptions(indexFile)
	
	case "outputFilteredJudgements":
		judgements := os.Args[2]
		ids := os.Args[3]
		filterJudgements(judgements, ids)

	case "outputDeletes":
		idsFile := os.Args[2]
		createDocESDeleteJSONL(idsFile)
	}


}
