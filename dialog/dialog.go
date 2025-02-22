package dialog

import (
	"strings"

	"github.com/ncruces/zenity"
)

// Function to select the path of a file
func SelectFile(titolo string) string {
	strings, _ := zenity.SelectFile(
		zenity.Title(titolo),
		zenity.Filename(""),
		zenity.FileFilters{
			{Name: "File SRT", Patterns: []string{"*.srt"}, CaseFold: false},
		})
	return strings
}

// Function to select the path where to save the processed file
func GetSavePath(titolo string) string {
	filePath, _ := zenity.SelectFileSave(
		zenity.Title(titolo),
		zenity.Filename(""),
		zenity.FileFilters{
			{Name: "File XLSX", Patterns: []string{"*.xlsx"}, CaseFold: false},
		})
	if !strings.Contains(filePath, ".") {
		filePath = filePath + ".xlsx"
	}
	return filePath
}
