package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"github.com/valyala/tsvreader"
	"strings"
	"strconv"
	"regexp"
)

type Query struct {
	ID       int
	Relevant int
}

type Judgement struct {
	DocID   string
	Queries []Query
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
			// Only add relevancy if it equals 1.
			if relevancy == 1 {
				query := Query{}
				query.ID = i
				query.Relevant = relevancy
				judgement.Queries = append(judgement.Queries, query)
			}
		}

		judgements = append(judgements, judgement)
	}

	for _, judgement := range judgements {
		for _, query := range judgement.Queries {
			fmt.Printf("%d 0 %s %d\n", query.ID, judgement.DocID, query.Relevant)
		}
	}

}


func outputEmptyDescriptions(indexFile string){
	b, err := ioutil.ReadFile(indexFile)
	if err != nil {
		fmt.Print(err)
	}
	
	//noDescriptionDocID:= []uint32{}
	docs := strings.Split(string(b), "<DOC>")[1:]
	found := 0
	for _, doc := range docs{
		descRegex := regexp.MustCompile(`<CSDESCRIPTION>([\S\s]+)</CSDESCRIPTION>`)
		descriptionMatches := descRegex.FindSubmatch([]byte(doc))
		if(len(descriptionMatches[1]) > 2){
			found += 1
		} else{
			docnoRegex := regexp.MustCompile(`<DOCNO>(\S+)</DOCNO>`)
			docno := docnoRegex.FindSubmatch([]byte(doc))[1]
			fmt.Printf("%s\n", string(docno))
		}
	}
}

// For judgements get each line.
// Split each line by spaces
// If the third element converted to int is found in array of non description ids,
// do not spit out the line, otherwise do
func filterJudgements(judgementsFile string, missingFile string){
	judgements, _ := ioutil.ReadFile(judgementsFile)
	missing, _ := ioutil.ReadFile(missingFile)
	
	IDStringArray := strings.Split(string(missing), "\n")
	
	var ids []int

	for _, id := range IDStringArray{
		idInt, _ := strconv.Atoi(id)
		ids = append(ids, idInt)
	}


	judgementLines := strings.Split(string(judgements), "\n")
	
	for _, j := range judgementLines[10:]{
		judgementIDString := strings.Split(j, " ")[2]
		judgementID, _ := strconv.Atoi(judgementIDString)

		found := false
		for i := range ids{
			if ids[i] == judgementID{
				found = true
				break
			}
		}

		if(!found){
			fmt.Printf("%s\n", j)
		}
	}
}
