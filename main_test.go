package main

import (
	"secondProjectGO/apijson"
	"testing"
)

func TestHandlerMain(t *testing.T) {
	type errorTestCases struct {
		name          string
		readFunction  func(string, int) (apijson.Response, error)
		writeFunction func(map[string][]apijson.Record) ([][]byte, error)
		want          error
	}
	//for _, scenario := range []errorTestCases{
	//	{
	//		name: "error while reading from file",
	//
	//	},
	//} {
	//
	//}
}
