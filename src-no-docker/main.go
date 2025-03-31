// @title           SWIFTcode API
// @version         1.0
// @description     This is a sample server for maintenance of SWIFTcode API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /v1

// @securityDefinitions.basic  BasicAuth

package main

import (
	"SWIFT/src-no-docker/REST"
	"SWIFT/src-no-docker/databaseControl"
	"SWIFT/src-no-docker/xlsxParser"

	"fmt"

	"github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _ "SWIFT/src-no-docker/docs"
)

func main() {
	// PARSE THE .XLSX FILE 

	SWIFTdata, out, err := xlsxParser.Parse("./data/Interns_2025_SWIFT_CODES.xlsx", "Sheet1")
	if !out || err != nil {
		fmt.Println("Parsing of xlsx data failed")
		return
	}

	// STORE THE DATA IN THE MySQL DATABASE

	databaseControl.AddTheInitialData(SWIFTdata)

	// PROVIDE A REST API
	server := REST.RunTheServer()

	// run the server on port 8080
    url := ginSwagger.URL("/swagger/doc.json")
    server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
    
	server.Run(":8080")
}
