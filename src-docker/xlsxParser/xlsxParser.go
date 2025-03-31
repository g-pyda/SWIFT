package xlsxParser

import (
	"fmt"
	"log"
	"strings"

	"SWIFT/src/structs"

	"github.com/xuri/excelize/v2"
)

func Parse(fileName string) []structs.Xlsx_data {
	// open the file and reading all the rows
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	// removing the first row (it only contains the headers to the data)
	if len(rows) > 1 {
		rows = rows[1:]
	} else {
		fmt.Println("The file contains one or less lines, so is invalid!!!")
		return []structs.Xlsx_data{}
	}

	// getting all the data from the specific rows
	SWIFTdata := []structs.Xlsx_data{}

	for _, row := range rows {
		SWIFTdata = append(SWIFTdata, structs.Xlsx_data{
			strings.ToUpper(row[0]), strings.ToUpper(row[1]), 
			strings.ToUpper(row[3]), strings.ToUpper(row[4]), 
			strings.ToUpper(row[5]), strings.ToUpper(row[6]), 
			row[7],})
	}

	return SWIFTdata
}