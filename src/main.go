package main

import (

	"SWIFT/src/xlsxParser"
	"SWIFT/src/REST"
	"SWIFT/src/databaseControl"

)

func main() {
	// PARSE THE .XLSX FILE 

	SWIFTdata := xlsxParser.Parse("/app/data/Interns_2025_SWIFT_CODES.xlsx")

	// STORE THE DATA IN THE MySQL DATABASE

	databaseControl.AddTheInitialData(SWIFTdata)

	// PROVIDE A REST API
	server := REST.RunTheServer()

	// run the server on port 8080
	server.Run(":8080")
}