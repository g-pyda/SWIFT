package databaseControl

import (
	"fmt"
	"testing"
	"errors"
	"os"

	"SWIFT/src-no-docker/structs"

	"github.com/stretchr/testify/assert"
)

const dsn_test = "SWIFTuser:SWIFTpass@tcp(localhost:3306)/testswiftdb"

// ------------------- FILE: fundamentals -------------------- //

var testCases_connectToDb = []structs.Testcase[string]{
	{
		Name: "Valid dsn", 
		ExpectedOutcome: true, 
		ExpectedError: nil,
		Input: dsn_test,
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

func TestConnectToDb(t *testing.T) {
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
		t.Setenv("db_conn", "true") 
	} else {
		t.Setenv("db_conn", "false")
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
	dsn := dsn_test
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
		t.Setenv("cr_tab", "true") 
	} else {
		t.Setenv("cr_tab", "false")
	}
}

// ------------------- FILE: addEntry --------------------- //

// setup before the testing
func setUpBeforeAdd() bool{
	db, _, _ := connectToDb(dsn_test)
	defer db.Close()

	for _, table := range structs.Tables {
		out, err := createTable(db, table.Name, table.Rows, table.Addition)
		if !out || err != nil {
			return false
		}
	}

	queries := []string{"DELETE FROM branches", "DELETE FROM headquarters", "DELETE FROM countries"}
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return false
		}
	}

	return true
}

var testCases_AddHeadquarter = []structs.Testcase[structs.ReqBranch]{
	{
		Name:            "Successful headquarter addition",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input: structs.ReqBranch{
			SwiftCode:   "HQBANK11XXX",
			BankName:    "HQ Bank",
			Address:     "123 Main St",
			CountryISO2: "US",
			CountryName: "United States",
		},
	},
	{
		Name:            "Headquarter already exists",
		ExpectedOutcome: false,
		ExpectedError:   fmt.Errorf("the headquarter with this swift code already exists"),
		Input: structs.ReqBranch{
			SwiftCode:   "HQBANK11XXX",
			BankName:    "Existing HQ",
			Address:     "456 Elm St",
			CountryISO2: "US",
			CountryName: "United States",
		},
	},
	{
		Name:            "Country needs to be added",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input: structs.ReqBranch{
			SwiftCode:   "NEWCOUNTXXX",
			BankName:    "New Country HQ",
			Address:     "789 Oak St",
			CountryISO2: "ZZ",
			CountryName: "Zetaland",
		},
	},
}

func TestAddHeadquarter(t *testing.T) {
	// if the previous tests failed - the test fails automatically
	if os.Getenv("cr_tab") == "false" || os.Getenv("db_conn") == "false" {
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	t.Run("Setting ut the environment for adding", func(t *testing.T) {
		out := setUpBeforeAdd()
		assert.Equal(t, true, out)
	})

	for _, tc := range testCases_AddHeadquarter {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := AddHeadquarter(dsn_test, tc.Input)

			assert.Equal(t, tc.ExpectedOutcome, result, "Unexpected outcome")
			assert.Equal(t, tc.ExpectedError, err, "Unexpected error")
		})
	}

	if !t.Failed() {
		t.Setenv("add_hq", "true") 
	} else {
		t.Setenv("add_hq", "false")
	}

	if !t.Failed() {
		t.Setenv("add_hq", "true") 
	} else {
		t.Setenv("add_hq", "false")
	}
}

var testCases_AddBranch = []structs.Testcase[structs.ReqBranch]{
	{
		Name:            "Successful branch addition",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input: structs.ReqBranch{
			SwiftCode:   "BRANCHXX123",
			BankName:    "Branch Bank",
			Address:     "123 Main St",
			CountryISO2: "US",
			CountryName: "United States",
		},
	},
	{
		Name:            "Branch already exists",
		ExpectedOutcome: false,
		ExpectedError:   errors.New("the branch with this swift code already exists"),
		Input: structs.ReqBranch{
			SwiftCode:   "BRANCHXX123",
			BankName:    "Existing Branch",
			Address:     "456 Elm St",
			CountryISO2: "US",
			CountryName: "United States",
		},
	},
	{
		Name:            "Country needs to be added",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input: structs.ReqBranch{
			SwiftCode:   "NEWCOUNTRYB",
			BankName:    "New Country Branch",
			Address:     "789 Oak St",
			CountryISO2: "ZZ",
			CountryName: "Zetaland",
		},
	},
	{
		Name:            "Headquarter exists - link branch to HQ",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input: structs.ReqBranch{
			SwiftCode:   "HQBANK11EXI",
			BankName:    "Linked Branch",
			Address:     "HQ Linked St",
			CountryISO2: "US",
			CountryName: "United States",
		},
	},
	{
		Name:            "Headquarter doesn't exist - standalone branch",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input: structs.ReqBranch{
			SwiftCode:   "NOHQBRANCHS",
			BankName:    "Standalone Branch",
			Address:     "No HQ St",
			CountryISO2: "US",
			CountryName: "United States",
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
			result, err := AddBranch(dsn_test, tc.Input)

			assert.Equal(t, tc.ExpectedOutcome, result, "Unexpected outcome")
			assert.Equal(t, tc.ExpectedError, err, "Unexpected error")

		})
	}

	if !t.Failed() {
		t.Setenv("add_br", "true") 
	} else {
		t.Setenv("add_br", "false")
	}
}

var testCases_AddCountry = []structs.Testcase[structs.Input_add_coun]{
	{
		Name:            "Successful country addition",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input: structs.Input_add_coun{
			ISO2:     "FR",
			Name:     "France",
			TimeZone: "UTC+1",
		},
	},
	{
		Name:            "Country already exists",
		ExpectedOutcome: false,
		ExpectedError:   errors.New("Error 1062 (23000): Duplicate entry 'US' for key 'countries.PRIMARY'"),
		Input: structs.Input_add_coun{
			ISO2:     "US",
			Name:     "United States",
			TimeZone: "UTC-5",
		},
	},
	{
		Name:            "Database error during country addition",
		ExpectedOutcome: false,
		ExpectedError:   errors.New("Error 1406 (22001): Data too long for column 'iso2' at row 1"),
		Input: structs.Input_add_coun{
			ISO2:     "XXF",
			Name:     "Unknownland",
			TimeZone: "UTC+3",
		},
	},
}

func TestAddCountry(t *testing.T) {
	if os.Getenv("cr_tab") == "false" || os.Getenv("db_conn") == "false" {
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	for _, tc := range testCases_AddCountry {
		t.Run(tc.Name, func(t *testing.T) {
			db, _, _ := connectToDb(dsn_test)
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
		t.Setenv("add_country", "true")
	} else {
		t.Setenv("add_country", "false")
	}
}

// ------------------- FILE: getEntry --------------------- //

var testCases_GetBranch = []structs.Testcase[struct{ SwiftCode string }]{
	{
		Name:            "Successfully retrieve a branch",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input:           struct{ SwiftCode string }{SwiftCode: "BRANCHXX123"},
	},
	{
		Name:            "Branch not found",
		ExpectedOutcome: false,
		ExpectedError:   fmt.Errorf("no branch found with swift code: INVALIDSWIF"),
		Input:           struct{ SwiftCode string }{SwiftCode: "INVALIDSWIF"},
	},
}

func TestGetBranch(t *testing.T) {
	if os.Getenv("add_hq") == "false" || os.Getenv("add_br") == "false" || os.Getenv("add_count") == "false"{
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	for _, tc := range testCases_GetBranch {
		t.Run(tc.Name, func(t *testing.T) {

			_, result, err := GetBranch(dsn_test, tc.Input.SwiftCode)

			assert.Equal(t, tc.ExpectedOutcome, result, "Unexpected outcome")

			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error(), "Unexpected error")
			} else {
				assert.Equal(t, err, nil, "Unexpected error")
			}
		})
	}

	if !t.Failed() {
		t.Setenv("get_branch", "true")
	} else {
		t.Setenv("get_branch", "false")
	}
}

var testCases_GetHeadquarter = []structs.Testcase[struct{ SwiftCode string }]{
	{
		Name:            "Successfully retrieve a headquarter",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input:           struct{ SwiftCode string }{SwiftCode: "HQBANK11XXX"},
	},
	{
		Name:            "Headquarter not found",
		ExpectedOutcome: false,
		ExpectedError:   fmt.Errorf("no headquarter found with swift code: HQBANK11XXX"),
		Input:           struct{ SwiftCode string }{SwiftCode: "INVALIDHQSWIFT"},
	},
}

func TestGetHeadquarter(t *testing.T) {
	if os.Getenv("add_hq") == "false" || os.Getenv("add_br") == "false" || os.Getenv("add_count") == "false"{
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	for _, tc := range testCases_GetHeadquarter {
		t.Run(tc.Name, func(t *testing.T) {
			db, _, _ := connectToDb(dsn_test)
			_, result, err := GetHeadquarter(dsn_test, tc.Input.SwiftCode)

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
		t.Setenv("get_headquarter", "true")
	} else {
		t.Setenv("get_headquarter", "false")
	}
}

var testCases_GetCountry = []structs.Testcase[struct{ ISO2 string }]{
	{
		Name:            "Successfully retrieve a country",
		ExpectedOutcome: true,
		ExpectedError:   nil,
		Input:           struct{ ISO2 string }{ISO2: "US"},
	},
	{
		Name:            "Country not found",
		ExpectedOutcome: false,
		ExpectedError:   fmt.Errorf("no country found with ISO2 : US"),
		Input:           struct{ ISO2 string }{ISO2: "ZZ"},
	},
}

func TestGetCountry(t *testing.T) {
	if os.Getenv("add_hq") == "false" || os.Getenv("add_br") == "false" || os.Getenv("add_count") == "false"{
		fmt.Println("The previous tests failed, so this one can't be executed safely")
		t.FailNow()
	}

	for _, tc := range testCases_GetCountry {
		t.Run(tc.Name, func(t *testing.T) {
			_, result, err := GetCountry(dsn_test, tc.Input.ISO2)

			assert.Equal(t, tc.ExpectedOutcome, result, "Unexpected outcome")
			if tc.ExpectedError != nil {
				assert.EqualError(t, err, tc.ExpectedError.Error(), "Unexpected error")
			} else {
				assert.Equal(t, err, nil, "Unexpected error")
			}
		})
	}

	if !t.Failed() {
		t.Setenv("get_country", "true")
	} else {
		t.Setenv("get_country", "false")
	}
}

// ------------------- FILE: deleteEntry --------------------- //



