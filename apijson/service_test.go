package apijson

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestReadRecords(t *testing.T) {
	type ReadWant struct {
		ErrText error
		Resp    Response
	}

	type ReadInput struct {
		Url         string
		NrOfRecords int
	}

	type errorTestCases struct {
		name  string
		input ReadInput
		want  ReadWant
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		record1 := Record{First: "John", Last: "Doe", Email: "johndoe@example.com", Address: "1 Main St", Created: "2022-01-01", Balance: "$100.00"}
		record2 := Record{First: "Jane", Last: "Doe", Email: "janedoe@example.com", Address: "2 Main St", Created: "2022-01-02", Balance: "$200.00"}
		record3 := Record{First: "Jim", Last: "Smith", Email: "jimsmith@example.com", Address: "3 Main St", Created: "2022-01-03", Balance: "$300.00"}
		response := Response{Results: []Record{record1, record2, record3}}
		responseData, _ := json.Marshal(response)
		w.Write(responseData)
	}))

	for _, scenario := range []errorTestCases{
		{
			name:  "invalid url",
			input: ReadInput{"", 2},
			want:  ReadWant{errors.New("the http request failed"), Response{}},
		},
		{
			name:  "invalid number of records",
			input: ReadInput{"https://randomapi.com/api/6de6abfedb24f889e0b5f675edc50deb?fmt=prettyjson&sole", 0},
			want:  ReadWant{errors.New("invalid number of records"), Response{}},
		},
		{
			name:  "check the number of records returned is equal to the expected number",
			input: ReadInput{ts.URL, 3},
			want:  ReadWant{nil, Response{}},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			gotResponse, got := ReadRecords(scenario.input.Url, scenario.input.NrOfRecords)

			if scenario.want.ErrText != nil {
				if got.Error() != scenario.want.ErrText.Error() {
					t.Errorf("got %v, wanted %v", got.Error(), scenario.want)
				}
			} else {
				nrOfRecordsGot := len(gotResponse.Results)
				if nrOfRecordsGot != scenario.input.NrOfRecords {
					t.Errorf("got %v, wanted %v", nrOfRecordsGot, scenario.input.NrOfRecords)
				}
			}
		})
	}

	ts.Close()
}

func TestEliminateDuplicates(t *testing.T) {
	type testCases struct {
		name  string
		input Response
		want  Response
	}

	for _, scenario := range []testCases{
		{
			name: "Eliminate duplicates from result",
			input: Response{
				Results: []Record{
					{"John", "Doe", "john@example.com", "123 Main St", "01/01/2021", "$100"},
					{"John", "Doe", "john@example.com", "123 Main St", "01/01/2021", "$100"},
					{"Maria", "Maria", "maria@example.com", "456 Main St", "01/02/2021", "$200"},
					{"Ion", "Jon", "ion@example.com", "456 Main St", "01/02/2021", "$200"},
				},
			},
			want: Response{
				Results: []Record{
					{"John", "Doe", "john@example.com", "123 Main St", "01/01/2021", "$100"},
					{"Maria", "Maria", "maria@example.com", "456 Main St", "01/02/2021", "$200"},
					{"Ion", "Jon", "ion@example.com", "456 Main St", "01/02/2021", "$200"},
				},
			},
		},
		{
			name: "Result with no duplicates",
			input: Response{
				Results: []Record{
					{"John", "Doe", "john@example.com", "123 Main St", "01/01/2021", "$100"},
					{"Ioana", "Doe", "ioana@example.com", "456 Main St", "01/02/2021", "$200"},
					{"Alex", "Ion", "alex@example.com", "456 Main St", "01/02/2021", "$300"},
				},
			},
			want: Response{
				Results: []Record{
					{"John", "Doe", "john@example.com", "123 Main St", "01/01/2021", "$100"},
					{"Ioana", "Doe", "ioana@example.com", "456 Main St", "01/02/2021", "$200"},
					{"Alex", "Ion", "alex@example.com", "456 Main St", "01/02/2021", "$300"},
				},
			},
		},
		{
			name: "Only one field is changed",
			input: Response{
				Results: []Record{
					{"John", "Doe", "john@example.com", "123 Main St", "01/01/2021", "$100"},
					{"John", "Doe", "john@example.com", "123 Main St", "01/10/2021", "$100"},
				},
			},
			want: Response{
				Results: []Record{
					{"John", "Doe", "john@example.com", "123 Main St", "01/01/2021", "$100"},
					{"John", "Doe", "john@example.com", "123 Main St", "01/10/2021", "$100"},
				},
			},
		},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			scenario.input.EliminateDuplicates()
			//if !reflect.DeepEqual(scenario.input, scenario.want) {
			if !(fmt.Sprint(scenario.input) == fmt.Sprint(scenario.want)) {
				t.Errorf("got %v, wanted %v", scenario.input, scenario.want)
			}
		})
	}
}

func TestGroupByFirstLetter(t *testing.T) {
	type testCases struct {
		name  string
		input Response
		want  map[string][]Record
	}
	for _, scenario := range []testCases{
		{
			name: "Group only by one letter",
			input: Response{
				Results: []Record{
					{"John", "Doe", "johndoe@example.com", "123 Main St", "2022-01-01", "100"},
					{"Jane", "Doe", "janedoe@example.com", "456 Main St", "2022-01-02", "200"},
					{"Jim", "Smith", "jimsmith@example.com", "789 Main St", "2022-01-03", "300"},
				},
			},
			want: map[string][]Record{
				"J": []Record{
					{"John", "Doe", "johndoe@example.com", "123 Main St", "2022-01-01", "100"},
					{"Jane", "Doe", "janedoe@example.com", "456 Main St", "2022-01-02", "200"},
					{"Jim", "Smith", "jimsmith@example.com", "789 Main St", "2022-01-03", "300"},
				},
			},
		},
		{
			name: "Group by more letters",
			input: Response{
				Results: []Record{
					{"John", "Doe", "johndoe@example.com", "123 Main St", "2022-01-01", "$100"},
					{"Anne", "Maria", "maria@example.com", "5 Main St", "2020-05-07", "$200"},
				},
			},
			want: map[string][]Record{
				"J": []Record{
					{"John", "Doe", "johndoe@example.com", "123 Main St", "2022-01-01", "$100"},
				},
				"A": []Record{
					{"Anne", "Maria", "maria@example.com", "5 Main St", "2020-05-07", "$200"},
				},
			},
		},
	} {
		got := GroupByFirstLetter(scenario.input)
		if !reflect.DeepEqual(got, scenario.want) {
			t.Errorf("got %v, wanted: %v", got, scenario.want)
		}
	}
}

func TestWriteGroups(t *testing.T) {
	groups := map[string][]Record{
		"A": {{First: "Adam", Last: "Smith", Email: "adam@example.com", Address: "123 Main St", Created: "2022-01-01", Balance: "100"}},
		"B": {{First: "Bob", Last: "Johnson", Email: "bob@example.com", Address: "456 Main St", Created: "2022-01-01", Balance: "100"}},
	}

	data, err := WriteGroups(groups)
	if err != nil {
		t.Errorf("Error while writing groups to separate .json files: %v", err)
	}

	for firstLetter, groupData := range data {
		if len(groupData) == 0 {
			t.Errorf("Group data for letter %v should not be empty", firstLetter)
		}

		var group Group
		err = json.Unmarshal(groupData, &group)
		if err != nil {
			t.Errorf("Error while unmarshalling group data: %v", err)
		}

		if group.TotalRecords != len(groups[group.Index]) {
			t.Errorf("Unexpected number of records in group data: %d", group.TotalRecords)
		}

		for i, record := range group.Records {
			if record != groups[group.Index][i] {
				t.Errorf("Unexpected record in group data: %v", record)
			}
		}
	}
}
