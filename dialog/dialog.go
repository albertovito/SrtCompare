package dialog

import (
	"github.com/ncruces/zenity"
)

const TITOLO string = "Comparatore SRT"

func PrintDialog(titolo string) (file string) {
	strings, _ := zenity.SelectFile(
		zenity.Title(titolo),
		zenity.Filename(""),
		zenity.FileFilters{
			{Name: "File SRT", Patterns: []string{"*.srt"}, CaseFold: false},
		})
	return strings
}

func PrintSingleDialog(titolo, button string) {
	zenity.Info("File da selezionare:", zenity.OKLabel(button), zenity.NoCancel(), zenity.Title(titolo), zenity.NoIcon)
}

func MainPage() bool {
	err := zenity.Question(
		"Benvenuto nel comparatore di SRT\nVuoi cominciare?",
		zenity.OKLabel("Si"),
		zenity.CancelLabel("No"),
		zenity.NoIcon,
		zenity.Title("Comparatore SRT"),
	)
	return err == nil
}

func SecondPage() (file1, file2 string) {
	firstFile := ""
	secondFile := ""
	err := zenity.Question("File da selezionare:", zenity.OKLabel("File1"), zenity.CancelLabel("File2"), zenity.Title("Comparatore SRT"), zenity.NoIcon)

	if err != nil {
		secondFile = PrintDialog("Seleziona il secondo file")
		PrintSingleDialog(TITOLO, "File1")
		firstFile = PrintDialog("Seleziona il primo file")
	} else {
		firstFile = PrintDialog("Seleziona il primo file")
		PrintSingleDialog(TITOLO, "File2")
		secondFile = PrintDialog("Seleziona il secondo file")
	}
	return firstFile, secondFile
}
