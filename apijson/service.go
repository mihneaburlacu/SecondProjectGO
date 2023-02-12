package apijson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Record struct {
	First   string `json:"first"`
	Last    string `json:"last"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Created string `json:"created"`
	Balance string `json:"balance"`
}

type Response struct {
	Results []Record `json:"results"`
}

type Group struct {
	Index        string   `json:"index"`
	Records      []Record `json:"records"`
	TotalRecords int      `json:"totalRecords"`
}

func ReadRecords(url string, nrOfRecords int) (Response, error) {
	var allResults []Record
	var resp Response

	// Limit number of results to `numRecords` if specified
	if nrOfRecords > 0 {
		for len(allResults) < nrOfRecords {

			// Get data from API
			response, err := http.Get(url)
			if err != nil {
				fmt.Printf("The HTTP request failed with error %v\n", err)
				return resp, errors.New("the http request failed")
			}

			// Read response body
			var data []byte
			data, err = io.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("Reading the response body fail with error %v\n", err)
				return resp, err
			}

			// Unmarshal data into struct
			err = json.Unmarshal(data, &resp)
			if err != nil {
				fmt.Println("Error while unmarshal data into struct")
				return resp, err
			}

			// Add results to allResults slice
			allResults = append(allResults, resp.Results...)

			// Truncate allResults slice to numRecords if necessary
			if len(allResults) > nrOfRecords {
				allResults = allResults[:nrOfRecords]
			}
		}
		resp.Results = allResults
	} else {
		fmt.Printf("Invalid number of records")
		return resp, errors.New("invalid number of records")
	}

	return resp, nil
}

func (resp *Response) EliminateDuplicates() {
	// Eliminate duplicates using maps (keys can not be duplicate)
	m := make(map[Record]string)
	for _, oneResult := range resp.Results {
		m[oneResult] = oneResult.First
	}

	var uniqueResults []Record
	for k, _ := range m {
		uniqueResults = append(uniqueResults, k)
	}

	resp.Results = uniqueResults
}

func GroupByFirstLetter(resp Response) map[string][]Record {
	// Group results by first letter of First name (first field of the record)
	groups := make(map[string][]Record)
	for _, r := range resp.Results {
		firstLetter := string(r.First[0])
		groups[firstLetter] = append(groups[firstLetter], r)
	}

	return groups
}

func WriteGroups(groups map[string][]Record) ([][]byte, error) {
	// Write groups to separate .json files
	var allData [][]byte

	for firstLetter, records := range groups {
		var oneGroup Group
		oneGroup.Index = firstLetter
		oneGroup.Records = records
		oneGroup.TotalRecords = len(groups[firstLetter])

		//groupData, err := json.MarshalIndent(records, "", "    ")
		groupData, err := json.MarshalIndent(oneGroup, "", "    ")
		if err != nil {
			fmt.Println("Error while marshal data")
			return allData, err
		}

		err = os.WriteFile(firstLetter+".json", groupData, 0644)
		if err != nil {
			fmt.Printf("An error occurred while writing to file %s.json: %s\n", firstLetter, err)
			return allData, err
		}

		allData = append(allData, groupData)
	}

	return allData, nil
}
