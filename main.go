package main

import (
	"SWIFT/xlsxParser"
	"fmt"
	//"SWIFT/REST"
	"SWIFT/databaseControl"
)

func main() {
	// PARSE THE .XLSX FILE 

	SWIFTdata := xlsxParser.Parse("../Interns_2025_SWIFT_CODES.xlsx")

	fmt.Println(SWIFTdata[0].ISO2)

	// STORE THE DATA IN THE MySQL DATABASE

	databasecontrol.Test()

	// PROVIDE A REST API
}