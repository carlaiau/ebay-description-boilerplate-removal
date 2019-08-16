package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/valyala/tsvreader"
)

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

/*
 *
 * Given an array of alreadySeendDocIDs, and the original
 *  documents TSV file, create documents in the correct
 * structure for index
 */
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
