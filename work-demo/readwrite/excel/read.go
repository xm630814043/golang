package excel

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func ReadExcel() error {
	f, err := excelize.OpenFile("tms_driver.xlsx")
	if err != nil {
		return err
	}
	fmt.Println("f：", f)
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		return err
	}
	fmt.Println("cell：", cell)
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}
	for _, row := range rows {
		fmt.Print("row：", row, "\t")
		for _, colCell := range row {
			fmt.Print("colCell：", colCell, "\t")
		}
	}
	return nil
}
