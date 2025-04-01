package REST

import (
	"SWIFT/src/databaseControl"
	"SWIFT/src/structs"
	
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)
var restDsn = databaseControl.Dsn

func setUpBeforeRun() bool{
	// checking if the app is runing in Docker
	value, exists := os.LookupEnv("DOCKERIZED")
	if exists && value == "yes"{
		fmt.Println("entered")
		restDsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
	    	os.Getenv("DB_USER"),
	    	os.Getenv("DB_PASSWORD"),
	    	os.Getenv("DB_HOST"),
	    	os.Getenv("DB_PORT"),
	    	os.Getenv("DB_TESTNAME"),
	    )
	}
	out := databaseControl.SetUpBeforeAdd(restDsn)
	if !out {
		return false
	}

	out, err := databaseControl.AddBranch(restDsn, structs.ReqBranch{
		Address: "ssssssss",
		BankName: "ssss ss",
		CountryISO2: "ss",
		CountryName: "ssss",
		SwiftCode: "SSSSSSSSSSS",
	})

	if !out || err != nil {
		fmt.Println(err)
		return false
	}

	out, err = databaseControl.AddHeadquarter(restDsn, structs.ReqBranch{
		Address: "ssssssss",
		BankName: "ssss ss",
		CountryISO2: "ss",
		CountryName: "ssss",
		SwiftCode: "SSSSSSSSXXX",
	})

	if !out || err != nil {
		fmt.Println(err)
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
			t.Fatal("Set up of the environment for the test failed")
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

	if !t.Failed() {
		os.Setenv("rest_getall", "true") 
	} else {
		os.Setenv("rest_getall", "false")
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

	if !t.Failed() {
		os.Setenv("rest_getswift", "true") 
	} else {
		os.Setenv("rest_getswift", "false")
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

	if !t.Failed() {
		os.Setenv("rest_delete", "true") 
	} else {
		os.Setenv("rest_delete", "false")
	}
}

// EVALUATION - DID ALL THE TESTS PASS?
func TestAllPassed(t * testing.T) {
	t.Run("Package 'REST' - successfull testing", func(t *testing.T) {
		if os.Getenv("rest_getall") == "false" || os.Getenv("rest_getswift") == "false" || os.Getenv("rest_delete") == "false" {
			fmt.Println("TESTS_PASSED=yes")
			t.Fatalf("The package 'REST' didn't pass all the tests")
		}
		fmt.Println("TESTS_PASSED=no")
	})
}