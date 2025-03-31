package xlsxParser

import (
	"fmt"
	"strings"

	"SWIFT/src/structs"

	"github.com/xuri/excelize/v2"
)

func Parse(fileName string, sheetName string) ([]structs.Xlsx_data, bool, error) {
	// open the file and reading all the rows
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return []structs.Xlsx_data{}, false, fmt.Errorf("failed to open the '%s' file", fileName)
	}
	defer file.Close()

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return []structs.Xlsx_data{}, false, fmt.Errorf("failed to open the '%s' sheet in a file", sheetName)
	}

	// removing the first row (it only contains the headers to the data)
	if len(rows) > 1 {
		rows = rows[1:]
	} else {
		return []structs.Xlsx_data{}, false, fmt.Errorf("the file contains one or less lines, so is invalid")
	}

	// getting all the data from the specific rows
	SWIFTdata := []structs.Xlsx_data{}

	for _, row := range rows {
		SWIFTdata = append(SWIFTdata, structs.Xlsx_data{
			ISO2: strings.ToUpper(row[0]), 
			SWIFTcode: strings.ToUpper(row[1]), 
			Name: strings.ToUpper(row[3]), 
			Address: strings.ToUpper(row[4]), 
			TownName: strings.ToUpper(row[5]), 
			CountryName: strings.ToUpper(row[6]), 
			TimeZone: row[7],
		})
	}

	return SWIFTdata, true, nil
}