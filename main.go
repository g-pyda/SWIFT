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

	fmt.Println(SWIFTdata[22].Address)

	// STORE THE DATA IN THE MySQL DATABASE

	databasecontrol.Test(SWIFTdata)

	// PROVIDE A REST API
}