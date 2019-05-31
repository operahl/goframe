package excel

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

type Xlsx struct {
	file *excelize.File
	rows [][]string
}

func LoadXlsx(path string) Xlsx {
	var xlsx Xlsx
	f, err := excelize.OpenFile(path)
	if err != nil {
		return xlsx
	}
	xlsx.file = f
	return xlsx
}

func (x *Xlsx) GetRow(sheet int, row int) []string {
	if x.rows == nil || len(x.rows) < 1 {
		x.rows, _ = x.file.GetRows(x.file.GetSheetName(sheet))
	}
	return x.rows[row]
}
