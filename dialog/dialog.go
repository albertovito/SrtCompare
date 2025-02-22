package dialog

import (
	"strings"

	"github.com/ncruces/zenity"
)

// Funzione per selezionare il percorso di un file
func SelectFile(titolo string) string {
	strings, _ := zenity.SelectFile(
		zenity.Title(titolo),
		zenity.Filename(""),
		zenity.FileFilters{
			{Name: "File SRT", Patterns: []string{"*.srt"}, CaseFold: false},
		})
	return strings
}

// Funzione per selezionare il percorso dove salvare il file elaborato
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
