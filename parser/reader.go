package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/valyala/tsvreader"
)

/*
Struct representing all documents in the documents collection.
Each Document is a Row.
TSV reader parses each tab seperated column into a variable on the struct.
*/
type Document struct {
	ID           string
	Price        string
	Title        string
	CategoryTree string
	Category     string
	PictureURL   string
	EbayID       int
}

func getCategoryFromTree(t string) string {
	categories := strings.Split(t, "> ")
	return categories[len(categories)-1]
}
func getDocs(filePath string) []Document {
	var documents []Document

	data, openErr := os.Open(filePath)
	if openErr != nil {
		fmt.Println("Open Error")
		panic(openErr)
	}

	r := tsvreader.New(data)

	for r.Next() {
		id := r.String()
		price := strings.Trim(r.String(), " ")
		title := strings.Trim(r.String(), " ")
		categoryTree := strings.Trim(r.String(), " ")
		pictureURL := strings.Trim(r.String(), " ")

		doc := Document{
			ID:           id,
			Price:        price,
			Title:        title,
			CategoryTree: categoryTree,
			Category:     getCategoryFromTree(categoryTree),
			PictureURL:   pictureURL,
		}
		documents = append(documents, doc)
	}

	if parseErr := r.Error(); parseErr != nil {
		fmt.Println("Parse Error")
		panic(parseErr)

	}

	return documents

}
