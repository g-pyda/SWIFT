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
	"SWIFT/src/REST"
	"SWIFT/src/databaseControl"
	"SWIFT/src/xlsxParser"

	"fmt"
	"os"

	"github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _ "SWIFT/docs"
)

func main() {
	fileAddress := "./data/Interns_2025_SWIFT_CODES.xlsx"
	// checking if the app is runing in Docker
	value, exists := os.LookupEnv("DOCKERIZED")
	if exists && value == "yes"{
		databaseControl.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		    os.Getenv("DB_USER"),
		    os.Getenv("DB_PASSWORD"),
		    os.Getenv("DB_HOST"),
		    os.Getenv("DB_PORT"),
		    os.Getenv("DB_NAME"),
	    )
		databaseControl.Dsn_test = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		    os.Getenv("DB_USER"),
		    os.Getenv("DB_PASSWORD"),
		    os.Getenv("DB_HOST"),
		    os.Getenv("DB_PORT"),
		    os.Getenv("DB_TESTNAME"),
	    )
		fileAddress = "/app/data/Interns_2025_SWIFT_CODES.xlsx"
		xlsxParser.TestFileAddress = "/app/data/test.xlsx"
	}


	// PARSE THE .XLSX FILE 

	SWIFTdata, out, err := xlsxParser.Parse(fileAddress, "Sheet1")
	if !out || err != nil {
		fmt.Println("Parsing of xlsx data failed")
		fmt.Println(err)
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
