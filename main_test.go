package main

import (
	"errors"
	"secondProjectGO/apijson"
	"testing"
)

func TestHandlerMain(t *testing.T) {
	type errorTestCases struct {
		name          string
		readFunction  func(string, int) (apijson.Response, error)
		writeFunction func(map[string][]apijson.Record) ([][]byte, error)
		wantError     bool
	}
	for _, scenario := range []errorTestCases{
		{
			name: "error while reading records",
			readFunction: func(s string, i int) (apijson.Response, error) {
				return apijson.Response{}, errors.New("error while reading records from the API")
			},
			writeFunction: func(m map[string][]apijson.Record) ([][]byte, error) {
				return nil, nil
			},
			wantError: true,
		},
		{
			name: "error while writing groups to .json files",
			readFunction: func(s string, i int) (apijson.Response, error) {
				return apijson.Response{}, nil
			},
			writeFunction: func(m map[string][]apijson.Record) ([][]byte, error) {
				return nil, errors.New("error while writing groups")
			},
			wantError: true,
		},
		{
			name: "success, no errors",
			readFunction: func(s string, i int) (apijson.Response, error) {
				return apijson.Response{Results: []apijson.Record{
					{"John", "Doe", "johndoe@example.com", "123 Main St", "2022-01-01", "$100"},
				}}, nil
			},
			writeFunction: func(m map[string][]apijson.Record) ([][]byte, error) {
				return [][]byte{}, nil
			},
		},
	} {
		var got bool

		errGot := HandlerMain(scenario.readFunction, scenario.writeFunction)
		if errGot != nil {
			got = true
		} else {
			got = false
		}

		if got != scenario.wantError {
			t.Errorf("wanted the error to be: %v, got: %v", scenario.wantError, got)
		}
	}
}
