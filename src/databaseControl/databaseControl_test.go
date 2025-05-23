package databaseControl

import (
	"fmt"
	"testing"
	"errors"
	"os"

	"SWIFT/src/structs"

	"github.com/stretchr/testify/assert"
)

// ------------------- FILE: fundamentals -------------------- //

var valid_dsn = Dsn_test
func TestConnectToDb(t *testing.T) {
	// checking if the app is runing in Docker
	value := os.Getenv("DOCKERIZED")
	if value == "yes"{
		valid_dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
	    	os.Getenv("DB_USER"),
	    	os.Getenv("DB_PASSWORD"),
	    	os.Getenv("DB_HOST"),
	    	os.Getenv("DB_PORT"),
	    	os.Getenv("DB_TESTNAME"),
	    )
	}

	var testCases_connectToDb = []structs.Testcase[string]{
		{
			Name: "Valid dsn", 
			ExpectedOutcome: true, 
			ExpectedError: nil,
			Input: valid_dsn,
		},
		{
			Name: "Invalid dsn - invalid password", 
			ExpectedOutcome: false, 
			ExpectedError: fmt.Errorf("error verifying connection"),
			Input: "SWIFTuser:SWIFT@tcp(localhost:3306)/testswiftdb",
		},
		{
			Name: "Invalid dsn - non-existent database", 
			ExpectedOutcome: false, 
			ExpectedError: fmt.Errorf("error verifying connection"),
			Input: "SWIFTuser:SWIFTpass@tcp(localhost:3306)/thisdoesntexist",
		},
		{
			Name: "Invalid dsn - wrong host", 
			ExpectedOutcome: false, 
			ExpectedError: fmt.Errorf("error verifying connection"),
			Input: "SWIFTuser:SWIFTpass@tcp(2000:3306)/testswiftdb",
		},
	}

	for _, testCase := range testCases_connectToDb {
		t.Run(testCase.Name, func(t *testing.T){
			assert := assert.New(t)

			db, outcome, err := connectToDb(testCase.Input)

			assert.Equal(testCase.ExpectedOutcome, outcome)
			assert.Equal(testCase.ExpectedError, err)

			if testCase.Name == "Valid dsn" {
				defer db.Close()
			}
		})
	}
	if !t.Failed() {
		os.Setenv("db_conn", "true") 
	} else {
		os.Setenv("db_conn", "false")
	}
}

/* 
   performance of AddTheInitialData() is fully dependent on the 
   performance of the function used inside therefore its test  
   would be redundant                                          */

// ------------------- FILE: createTable --------------------- //

var testCases_CreateTable = []structs.Testcase[structs.Input_cr_tbl]{
    // Valid cases
    {
        Name: "Valid simple table",
        ExpectedOutcome: true,
        ExpectedError: nil,
        Input: structs.Input_cr_tbl{
            TableName: "users",
            TableRows: []structs.TableRow{
                {Name: "id", Data_type: "INT", Addition: " PRIMARY KEY"},
				{Name: "name", Data_type: "VARCHAR(30)", Addition: ""},
            },
        },
    },
    {
        Name: "Valid table with underscore",
        ExpectedOutcome: true,
        ExpectedError: nil,
        Input: structs.Input_cr_tbl{
            TableName: "new_orders",
            TableRows: []structs.TableRow{
                {Name: "id", Data_type: "INT", Addition: " PRIMARY KEY"},
            },
        },
    },
    {
        Name: "Valid table with maximum length (64 chars)",
        ExpectedOutcome: true,
        ExpectedError: nil,
        Input: structs.Input_cr_tbl{
            TableName: "a_very_long_table_name_that_is_exactly_64_characters_long_123456",
            TableRows: []structs.TableRow{
                {Name: "id", Data_type: "INT", Addition: " PRIMARY KEY"},
            },
        },
    },

    // Invalid table name cases
    {
        Name: "Empty table name",
        ExpectedOutcome: false,
        ExpectedError: errors.New("invalid table name: table name is empty"),
        Input: structs.Input_cr_tbl{
            TableName: "",
            TableRows: []structs.TableRow{
                {Name: "id", Data_type: "INT"},
            },
        },
    },
    {
        Name: "Table name starting with digit",
        ExpectedOutcome: false,
        ExpectedError: errors.New("invalid table name: table name cannot start with a digit"),
        Input: structs.Input_cr_tbl{
            TableName: "1invalid",
            TableRows: []structs.TableRow{
                {Name: "id", Data_type: "INT"},
            },
        },
    },
    {
        Name: "Table name with hyphen",
        ExpectedOutcome: false,
        ExpectedError: errors.New("invalid table name: table name contains invalid character '-'"),
        Input: structs.Input_cr_tbl{
            TableName: "table-name",
            TableRows: []structs.TableRow{
                {Name: "id", Data_type: "INT"},
            },
        },
    },
    {
        Name: "Table name with space",
        ExpectedOutcome: false,
        ExpectedError: errors.New("invalid table name: table name contains invalid character ' '"),
        Input: structs.Input_cr_tbl{
            TableName: "table name",
            TableRows: []structs.TableRow{
                {Name: "id", Data_type: "INT"},
            },
        },
    },
    {
        Name: "Table name that's a SQL keyword",
        ExpectedOutcome: false,
        ExpectedError: errors.New("invalid table name: table name 'select' is a SQL keyword"),
        Input: structs.Input_cr_tbl{
            TableName: "select",
            TableRows: []structs.TableRow{
                {Name: "id", Data_type: "INT"},
            },
        },
    },
    {
        Name: "Table name exceeding max length",
        ExpectedOutcome: false,
        ExpectedError: errors.New("invalid table name: table name exceeds maximum length of 64 characters"),
        Input: structs.Input_cr_tbl{
            TableName: "a_very_long_table_name_that_is_longer_than_64_characters_12345678901234567890",
            TableRows: []structs.TableRow{
                {Name: "id", Data_type: "INT"},
            },
        },
    },

    // Column validation cases
    {
        Name: "Empty columns",
        ExpectedOutcome: false,
        ExpectedError: fmt.Errorf("the table doesn't have any columns"),
        Input: structs.Input_cr_tbl{
            TableName: "empty_columns",
            TableRows: []structs.TableRow{},
        },
    },
}

func TestCreateTable(t *testing.T) {
	// opening a connection to the database
	dsn := valid_dsn
	testdb, _, _ := connectToDb(dsn)

	for _, testCase := range testCases_CreateTable {
		t.Run(testCase.Name, func(t *testing.T){
			assert := assert.New(t)

			outcome, err := createTable(testdb, testCase.Input.TableName, testCase.Input.TableRows, testCase.Input.Addition)

			assert.Equal(testCase.ExpectedOutcome, outcome)
			assert.Equal(testCase.ExpectedError, err)
		})
	}
	testdb.Close()
	if !t.Failed() {
		os.Setenv("cr_tab", "true") 
	} else {
		os.Setenv("cr_tab", "false")
	}
}

// ------------------- FILE: addEntry --------------------- //

var testCases_AddCountry = []structs.Testcase[structs.Input_add_coun]{
	{
		Name: "Successful country addition",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: structs.Input_add_coun{
			ISO2: "AA",
			Name: "aaa",
			TimeZone: "UTC+1",
		},
	},
	{
		Name: "Country already exists",
		ExpectedOutcome: false,
		ExpectedError: errors.New("the country already exists"),
		Input: structs.Input_add_coun{
			ISO2: "AA",
			Name: "aaa",
			TimeZone: "UTC+1",
		},
	},
	{
		Name: "Database error during country addition",
		ExpectedOutcome: false,
		ExpectedError:   errors.New("invalid ISO2 value"),
		Input: structs.Input_add_coun{
			ISO2: "FFF",
			Name: "FFF",
			TimeZone: "UTC+3",
		},
	},
}

func TestAddCountry(t *testing.T) {
	if os.Getenv("cr_tab") == "false" || os.Getenv("db_conn") == "false" {
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	t.Run("Setting ut the environment for adding", func(t *testing.T) {
		out := SetUpBeforeAdd(valid_dsn)
		assert.Equal(t, true, out)
	})

	for _, tc := range testCases_AddCountry {
		t.Run(tc.Name, func(t *testing.T) {
			db, _, _ := connectToDb(valid_dsn)
			result, err := addCountry(db, tc.Input.ISO2, tc.Input.Name, tc.Input.TimeZone)

			assert.Equal(t, tc.ExpectedOutcome, result, "Unexpected outcome")
			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error(), "Unexpected error")
			} else {
				assert.Equal(t, err, nil, "Unexpected error")	
			}

			db.Close()
		})
	}

	if !t.Failed() {
		os.Setenv("add_country", "true")
	} else {
		os.Setenv("add_country", "false")
	}
}

var testCases_AddHeadquarter = []structs.Testcase[structs.ReqBranch]{
	{
		Name: "Successful headquarter addition",
		ExpectedOutcome:true,
		ExpectedError: nil,
		Input: structs.ReqBranch{
			SwiftCode: "AAAAAAAAXXX",
			BankName: "aaaa",
			Address: "aa aaa aaa",
			CountryISO2: "AA",
			CountryName: "aaa",
		},
	},
	{
		Name: "Headquarter already exists",
		ExpectedOutcome: false,
		ExpectedError: fmt.Errorf("the headquarter with this swift code already exists"),
		Input: structs.ReqBranch{
			SwiftCode: "AAAAAAAAXXX",
			BankName: "aaaa",
			Address: "aa aaa aaa",
			CountryISO2: "AA",
			CountryName: "aaa",
		},
	},
	{
		Name: "Country needs to be added",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: structs.ReqBranch{
			SwiftCode: "BBBBBBBBXXX",
			BankName: "bbbbbb",
			Address: "bb bbbb bbb",
			CountryISO2: "BB",
			CountryName: "Bbbb",
		},
	},
}

func TestAddHeadquarter(t *testing.T) {
	// if the previous tests failed - the test fails automatically
	if os.Getenv("cr_tab") == "false" || os.Getenv("db_conn") == "false" {
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	for _, tc := range testCases_AddHeadquarter {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := AddHeadquarter(valid_dsn, tc.Input)

			assert.Equal(t, tc.ExpectedOutcome, result, "Unexpected outcome")
			assert.Equal(t, tc.ExpectedError, err, "Unexpected error")
		})
	}

	if !t.Failed() {
		os.Setenv("add_hq", "true") 
	} else {
		os.Setenv("add_hq", "false")
	}

	if !t.Failed() {
		os.Setenv("add_hq", "true") 
	} else {
		os.Setenv("add_hq", "false")
	}
}

var testCases_AddBranch = []structs.Testcase[structs.ReqBranch]{
	{
		Name: "Successful branch addition",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: structs.ReqBranch{
			SwiftCode: "ACCCCCCCCCC",
			BankName: "ccc ccc",
			Address: "cc cccc ccc",
			CountryISO2: "AA",
			CountryName: "aaa",
		},
	},
	{
		Name: "Branch already exists",
		ExpectedOutcome: false,
		ExpectedError: errors.New("the branch with this swift code already exists"),
		Input: structs.ReqBranch{
			SwiftCode: "ACCCCCCCCCC",
			BankName: "ccc ccc",
			Address: "cc cccc ccc",
			CountryISO2: "AA",
			CountryName: "aaa",
		},
	},
	{
		Name: "Country needs to be added",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: structs.ReqBranch{
			SwiftCode: "DDDDDDDDDDD",
			BankName: "dddd dd",
			Address: "dd dddd ddd",
			CountryISO2: "DD",
			CountryName: "dddd",
		},
	},
	{
		Name: "Headquarter exists - link branch to HQ",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input: structs.ReqBranch{
			SwiftCode: "AAAAAAAAAAA",
			BankName: "aa a a a",
			Address: "aaaaaaaa",
			CountryISO2: "aa",
			CountryName: "AaaaAa",
		},
	},
	{
		Name: "Headquarter doesn't exist - standalone branch",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: structs.ReqBranch{
			SwiftCode: "D2DDDDDDDDD",
			BankName: "dddd ddd",
			Address: "dd ddddd dd",
			CountryISO2: "DD",
			CountryName: "dddd",
		},
	},
}

func TestAddBranch(t *testing.T) {
	if os.Getenv("cr_tab") == "false" || os.Getenv("db_conn") == "false" {
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	for _, tc := range testCases_AddBranch {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := AddBranch(valid_dsn, tc.Input)

			assert.Equal(t, tc.ExpectedOutcome, result, "Unexpected outcome")
			assert.Equal(t, tc.ExpectedError, err, "Unexpected error")

		})
	}

	if !t.Failed() {
		os.Setenv("add_br", "true") 
	} else {
		os.Setenv("add_br", "false")
	}
}

// ------------------- FILE: getEntry --------------------- //

var testCases_GetBranch = []structs.Testcase[struct{ SwiftCode string }]{
	{
		Name: "Successfully retrieve a branch",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input:  struct{ SwiftCode string }{SwiftCode: "DDDDDDDDDDD"},
	},
	{
		Name: "Branch not found",
		ExpectedOutcome: false,
		ExpectedError: fmt.Errorf("no branch found with swift code: EEEEEEEEEEE"),
		Input: struct{ SwiftCode string }{SwiftCode: "EEEEEEEEEEE"},
	},
}

func TestGetBranch(t *testing.T) {
	if os.Getenv("add_hq") == "false" || os.Getenv("add_br") == "false" || os.Getenv("add_count") == "false"{
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	for _, tc := range testCases_GetBranch {
		t.Run(tc.Name, func(t *testing.T) {

			_, result, err := GetBranch(valid_dsn, tc.Input.SwiftCode)

			assert.Equal(t, tc.ExpectedOutcome, result, "Unexpected outcome")

			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error(), "Unexpected error")
			} else {
				assert.Equal(t, err, nil, "Unexpected error")
			}
		})
	}

	if !t.Failed() {
		os.Setenv("get_br", "true")
	} else {
		os.Setenv("get_br", "false")
	}
}

var testCases_GetHeadquarter = []structs.Testcase[struct{ SwiftCode string }]{
	{
		Name: "Successfully retrieve a headquarter",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: struct{ SwiftCode string }{SwiftCode: "AAAAAAAAXXX"},
	},
	{
		Name: "Headquarter not found",
		ExpectedOutcome: false,
		ExpectedError: fmt.Errorf("no headquarter found with swift code: DDDDDDDDXXX"),
		Input: struct{ SwiftCode string }{SwiftCode: "DDDDDDDDXXX"},
	},
}

func TestGetHeadquarter(t *testing.T) {
	if os.Getenv("add_hq") == "false" || os.Getenv("add_br") == "false" || os.Getenv("add_count") == "false"{
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	for _, tc := range testCases_GetHeadquarter {
		t.Run(tc.Name, func(t *testing.T) {
			_, result, err := GetHeadquarter(valid_dsn, tc.Input.SwiftCode)

			assert.Equal(t, tc.ExpectedOutcome, result, "Unexpected outcome")
			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error(), "Unexpected error")
			} else {
				assert.Equal(t, err, nil, "Unexpected error")
			}

		})
	}

	if !t.Failed() {
		os.Setenv("get_hq", "true")
	} else {
		os.Setenv("get_hq", "false")
	}
}

var testCases_GetCountry = []structs.Testcase[string]{
	{
		Name: "Successfully retrieve a country",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: "AA",
	},
	{
		Name: "Country not found",
		ExpectedOutcome: false,
		ExpectedError: fmt.Errorf("no country found with ISO2 : KB"),
		Input: "KB",
	},
}

func TestGetCountry(t *testing.T) {
	if os.Getenv("add_hq") == "false" || os.Getenv("add_br") == "false" || os.Getenv("add_count") == "false"{
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	for _, tc := range testCases_GetCountry {
		t.Run(tc.Name, func(t *testing.T) {
			_, result, err := GetCountry(valid_dsn, tc.Input)

			assert.Equal(t, tc.ExpectedOutcome, result, "Unexpected outcome")
			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error(), "Unexpected error")
			} else {
				assert.Equal(t, err, nil, "Unexpected error")
			}
		})
	}

	if !t.Failed() {
		os.Setenv("get_count", "true")
	} else {
		os.Setenv("get_count", "false")
	}
}

// ------------------- FILE: deleteEntry --------------------- //

var testCases_DeleteEntry = []structs.Testcase[string]{
	{
		Name: "Delete existing branch",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: "AAAAAAAAAAA",
	},
	{
		Name: "Delete existing headquarter",
		ExpectedOutcome: true,
		ExpectedError: nil,
		Input: "AAAAAAAAXXX",
	},
	{
		Name: "Delete non-existing branch",
		ExpectedOutcome: false,
		ExpectedError: errors.New("no entry found with swift code : GGGGGGGGGGG"),
		Input: "GGGGGGGGGGG",
	},
	{
		Name: "Delete non-existing headquarter",
		ExpectedOutcome: false,
		ExpectedError: errors.New("no entry found with swift code : GGGGGGGGXXX"),
		Input: "GGGGGGGGXXX",
	},
}

func TestDeleteEntry(t *testing.T) {
	if os.Getenv("add_hq") == "false" || os.Getenv("add_br") == "false" || os.Getenv("add_count") == "false"{
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	for _, tc := range testCases_DeleteEntry {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := DeleteEntry(valid_dsn, tc.Input)
			
			assert.Equal(t, tc.ExpectedOutcome, result, "unexpected outcome")
			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error(), "unexpected error")
			} else {
				assert.Equal(t, err, tc.ExpectedError, "unexpected error")
			}
		})
	}

	if !t.Failed() {
		os.Setenv("del_ent", "true")
	} else {
		os.Setenv("del_ent", "false")
	}
}

// EVALUATION - DID ALL THE TESTS PASS?
func TestAllPassed(t * testing.T) {
	t.Run("Package 'databaseControl' - successfull testing", func(t *testing.T) {
		if os.Getenv("db_conn") == "false" || os.Getenv("cr_tab") == "false" ||
		os.Getenv("add_hq") == "false" || os.Getenv("add_br") == "false" || os.Getenv("add_count") == "false" ||
		os.Getenv("get_hq") == "false" || os.Getenv("get_br") == "false" || os.Getenv("get_count") == "false" ||
		os.Getenv("del_ent") == "false" {
			os.Setenv("db_cont", "false")
			t.Fatalf("The package 'databaseControl' didn't pass all the tests")
		}
		os.Setenv("db_cont", "true")
	})
}

