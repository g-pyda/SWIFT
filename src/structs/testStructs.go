package structs

type Testcase [T any] struct {
	Name string
	ExpectedOutcome bool
	ExpectedError error
	Input T
}

type Input_x_parse struct {
	FileName string
	SheetName string
}

type Input_cr_tbl struct {
	TableName string
	TableRows []TableRow
	Addition string
}

type Input_add_coun struct {
	ISO2 string
	Name string
	TimeZone string
}