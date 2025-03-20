package xlsxParser

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

type SWIFT struct {
	ISO2 string
	SWIFTcode string
	codeType string
	name string
	address string
	townName string
	countryName string
	timeZone string
}

func Parse(fileName string) []SWIFT {
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
		return []SWIFT{}
	}

	// getting all the data from the specific rows
	SWIFTdata := []SWIFT{}

	for _, row := range rows {
		SWIFTdata = append(SWIFTdata, SWIFT{
			ISO2: row[0],
			SWIFTcode: row[1],
			codeType: row[2],
			name: row[3],
			address: row[4],
			townName: row[5],
			countryName: row[6],
			timeZone: row[7],
		})
	}

	return SWIFTdata
}