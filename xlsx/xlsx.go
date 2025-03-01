package xlsx

import (
	"SrtCompare/constants"
	"SrtCompare/srt"
	"fmt"

	"github.com/ncruces/zenity"
	"github.com/xuri/excelize/v2"
)

// Formats the xlsx file with subtitles1 and subtitiles2
func FormatSheet(file *excelize.File, subtitles1 []srt.Subtitle, subtitles2 []srt.Subtitle) error {
	sheet := "Sheet1"
	index, err := file.NewSheet(sheet)
	if err != nil {
		zenity.Error(fmt.Sprintf("Error creating sheet: %v", err), zenity.Title(constants.APPTITLE))
		return err
	}

	maxLength := len(subtitles1)
	if len(subtitles2) > maxLength {
		maxLength = len(subtitles2)
	}

	file.SetCellValue(sheet, "A1", "Index")
	file.SetCellValue(sheet, "B1", "Timing")
	file.SetCellValue(sheet, "C1", "File1")
	file.SetCellValue(sheet, "D1", "File2")

	for i := 0; i < maxLength; i++ {
		var subtitle1, subtitle2 srt.Subtitle
		if i < len(subtitles1) {
			subtitle1 = subtitles1[i]
		}
		if i < len(subtitles2) {
			subtitle2 = subtitles2[i]
		}
		a := "A" + fmt.Sprint(i+2)
		b := "B" + fmt.Sprint(i+2)
		c := "C" + fmt.Sprint(i+2)
		d := "D" + fmt.Sprint(i+2)
		file.SetCellValue(sheet, a, subtitle1.Index)
		file.SetCellValue(sheet, b, subtitle1.Start+" --->\n "+subtitle1.End)
		file.SetCellValue(sheet, c, subtitle1.Text)
		file.SetCellValue(sheet, d, subtitle2.Text)
	}

	file.SetActiveSheet(index)
	return nil
}
