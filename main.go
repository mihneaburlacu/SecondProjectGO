package main

import (
	"fmt"
	"secondProjectGO/apijson"
)

func main() {
	numRecords := 100
	url := "https://randomapi.com/api/6de6abfedb24f889e0b5f675edc50deb?fmt=prettyjson&sole"

	// Read from API
	resp, err := apijson.ReadRecords(url, numRecords)
	if err != nil {
		fmt.Printf("Error while reading records from API: %v", err)
	}

	// Eliminate duplicates records
	apijson.EliminateDuplicates(&resp)

	// Group results by first letter of First name (first field of the record)
	groups := apijson.GroupByFirstLetter(resp)

	// Write groups to separate .json files
	err = apijson.WriteGroups(groups)
	if err != nil {
		fmt.Printf("Error while writing .json groups to files: %v", err)
	}
}
