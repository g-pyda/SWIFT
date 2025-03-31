package main

import (

	"SWIFT/src-no-docker/xlsxParser"
	"SWIFT/src-no-docker/REST"
	"SWIFT/src-no-docker/databaseControl"

)

func main() {
	// PARSE THE .XLSX FILE 

	SWIFTdata := xlsxParser.Parse("./data/Interns_2025_SWIFT_CODES.xlsx")

	// STORE THE DATA IN THE MySQL DATABASE

	databaseControl.AddTheInitialData(SWIFTdata)

	// PROVIDE A REST API
	server := REST.RunTheServer()

	// run the server on port 8080
	server.Run(":8080")
}
