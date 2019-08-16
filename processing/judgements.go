package main

import (
	"fmt"
	"os"

	"github.com/valyala/tsvreader"
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
