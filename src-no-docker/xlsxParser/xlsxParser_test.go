package xlsxParser_test

import (
	"SWIFT/src-no-docker/xlsxParser"
	"SWIFT/src-no-docker/structs"
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
)

// Define the test cases
var testCases_Parse = []structs.Testcase[structs.Input_x_parse]{
	{
		Name: "Valid file with multiple entries",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: structs.Input_x_parse{
			FileName: "../data/test.xlsx", 
			SheetName: "Valid_sheet",
		},
	},
	{
		Name: "File with only headers",
		ExpectedOutcome: false,
		ExpectedError: fmt.Errorf("the file contains one or less lines, so is invalid"),
		Input: structs.Input_x_parse{
			FileName: "../data/test.xlsx", 
			SheetName: "One-row-sheet",
		},
	},
	{
		Name: "Non-existent file",
		ExpectedOutcome: false,
		ExpectedError: fmt.Errorf("failed to open the '../data/non_existent.xlsx' file"),
		Input: structs.Input_x_parse{
			FileName: "../data/non_existent.xlsx", 
			SheetName: "Sheet1",
		},
	},
	{
		Name: "Non-existent sheet",
		ExpectedOutcome: false,
		ExpectedError: fmt.Errorf("failed to open the 'invalid_sheet' sheet in a file"),
		Input: structs.Input_x_parse{
			FileName: "../data/test.xlsx", 
			SheetName: "invalid_sheet",
		},
	},
	{
		Name: "File with empty rows",
		ExpectedOutcome: false,
		ExpectedError: fmt.Errorf("the file contains one or less lines, so is invalid"),
		Input: structs.Input_x_parse{
			FileName: "../data/test.xlsx", 
			SheetName: "Empty-sheet",
		},
	},
}

func TestParse(t *testing.T) {
	for _, tc := range testCases_Parse {
		t.Run(tc.Name, func(t *testing.T) {
			_, out, err := xlsxParser.Parse(tc.Input.FileName, tc.Input.SheetName)

			assert.Equal(t, tc.ExpectedOutcome, out)
			assert.Equal(t, tc.ExpectedError, err)
		})
	}
}
