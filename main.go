package main

import (
	"fmt"

	"SWIFT/xlsxParser"
	"SWIFT/REST"
	"SWIFT/databaseControl"

)

func main() {
	// PARSE THE .XLSX FILE 

	SWIFTdata := xlsxParser.Parse("Interns_2025_SWIFT_CODES.xlsx")

	fmt.Println(SWIFTdata[22].Address)

	// STORE THE DATA IN THE MySQL DATABASE

	databaseControl.AddTheInitialData(SWIFTdata)

	// PROVIDE A REST API
	server := REST.RunTheServer()

	// Run the server on port 8080
	server.Run(":8080")
}