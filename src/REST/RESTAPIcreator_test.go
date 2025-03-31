package REST

import (
	"SWIFT/src/databaseControl"
	"SWIFT/src/structs"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setUpBeforeRun() bool{
	out := databaseControl.SetUpBeforeAdd()
	if !out {
		return false
	}

	out, err := databaseControl.AddBranch(databaseControl.Dsn, structs.ReqBranch{
		Address: "ssssssss",
		BankName: "ssss ss",
		CountryISO2: "ss",
		CountryName: "ssss",
		SwiftCode: "SSSSSSSSSSS",
	})

	if !out || err != nil {
		return false
	}

	out, err = databaseControl.AddHeadquarter(databaseControl.Dsn, structs.ReqBranch{
		Address: "ssssssss",
		BankName: "ssss ss",
		CountryISO2: "ss",
		CountryName: "ssss",
		SwiftCode: "SSSSSSSSXXX",
	})

	if !out || err != nil {
		return false
	}

	return true
}

// checking if the previous tests 'xlsxParser_test' and 'databaseControl_test' were performed succesfully
func TestIsSafe(t* testing.T) {
	t.Run("Checking if the previous tests succeeded", func(t *testing.T) {
		if os.Getenv("db_cont") == "false" {
			t.Fatal("Previous test for 'databaseControl' package failed")
		}
		out := setUpBeforeRun()
		if !out {
			t.Fatal("Set up ofthe enviromnent for the test failed")
		}
	})
}

var testCases_GetAll = []structs.Testcase[struct{}]{
	{
		Name: "Success - Get all SWIFT codes",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: struct{}{},
	},
}

func TestGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/v1/swift-codes", getAll)

	for _, tc := range testCases_GetAll {
		t.Run(tc.Name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/swift-codes", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Code)
		})
	}
}

var testCases_GetBySWIFT = []structs.Testcase[string]{
	{
		Name: "Success - Get branch by SWIFT code",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: "SSSSSSSSSSS",
	},
	{
		Name: "Success - Get headquarter by SWIFT code",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: "SSSSSSSSXXX",
	},
	{
		Name: "Fail - Branch not found",
		ExpectedOutcome: false,
		ExpectedError: assert.AnError,
		Input: "TTTTTTTTTTT",
	},
	{
		Name: "Fail - Headquarter not found",
		ExpectedOutcome: false,
		ExpectedError: assert.AnError,
		Input: "TTTTTTTTXXX",
	},
}

func TestGetBySWIFT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/v1/swift-codes/:swift-code", getBySWIFTcode)

	for _, tc := range testCases_GetBySWIFT {
		t.Run(tc.Name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/v1/swift-codes/"+tc.Input, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if tc.ExpectedOutcome {
				assert.Equal(t, http.StatusOK, resp.Code)
			} else {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			}
		})
	}
}

var testCases_DeleteEntry = []structs.Testcase[string]{
	{
		Name: "Success - Delete a branch",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: "SSSSSSSSSSS",
	},
	{
		Name: "Success - Delete a headquarter",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: "SSSSSSSSXXX",
	},
	{
		Name: "Fail - Headquarter does not exist",
		ExpectedOutcome: false,
		ExpectedError: assert.AnError,
		Input: "UUUUUUUUXXX",
	},
	{
		Name: "Fail - Branch does not exist",
		ExpectedOutcome: false,
		ExpectedError: assert.AnError,
		Input: "UUUUUUUUUUU",
	},
}

func TestDeleteEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/v1/swift-codes/:swift-code", deleteEntry)

	for _, tc := range testCases_DeleteEntry {
		t.Run(tc.Name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/v1/swift-codes/"+tc.Input, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if tc.ExpectedOutcome {
				assert.Equal(t, http.StatusOK, resp.Code)
			} else {
				assert.Equal(t, http.StatusNotFound, resp.Code)
			}
		})
	}
}