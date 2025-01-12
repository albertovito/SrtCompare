//go:build windows

package main

import (
	"SrtComparator/dialog"
	"SrtComparator/srt"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/ncruces/zenity"
)

func main() {
	// Leggi i sottotitoli da due file SRT
	res := dialog.MainPage()
	if !res {
		return
	}
	// file1, file2 := dialog.SecondPage()
	startWindow()

}

func startWindow() {
	var outTE1, outTE2 *walk.TextEdit
	var file1, file2 string

	MainWindow{
		Title:   "SrtComparator",
		Icon:    "res/srt.ico",
		MinSize: Size{Width: 600, Height: 400},
		Layout:  VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					VSplitter{
						Children: []Widget{
							TextEdit{AssignTo: &outTE1, ReadOnly: true, VScroll: true, Font: Font{Family: "Calibri 14"}, Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)}},
							PushButton{
								Text: "File1",
								Font: Font{Family: "Calibri 14"},
								OnClicked: func() {
									file1 = dialog.PrintDialog("")
									text1, _ := os.ReadFile(file1)
									outTE1.SetText(string(text1))
								},
							},
						},
					},
					VSplitter{
						Children: []Widget{
							TextEdit{AssignTo: &outTE2, ReadOnly: true, VScroll: true, Font: Font{Family: "Calibri 14"}, Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)}},
							PushButton{
								Text: "File2",
								Font: Font{Family: "Calibri 14"},
								OnClicked: func() {
									file2 = dialog.PrintDialog("")
									text2, _ := os.ReadFile(file2)
									outTE2.SetText(string(text2))
								},
							},
						},
					},
				},
			},
			PushButton{
				Text: "Compara",
				Font: Font{Family: "Calibri 14"},
				OnClicked: func() {
					subtitles1, err := srt.ReadSRTFile(file1)
					if err != nil {
						fmt.Println("Errore nel leggere il primo file SRT:", err)
						return
					}
					subtitles2, err := srt.ReadSRTFile(file2)
					if err != nil {
						fmt.Println("Errore nel leggere il secondo file SRT:", err)
						return
					}

					// Crea il file CSV
					csvPath := dialog.GetSavePath("Salva file")
					csvFile, err := os.Create(csvPath)
					if err != nil {
						fmt.Println("Errore nella creazione del file CSV:", err)
						return
					}
					defer csvFile.Close()

					// Crea un writer CSV con separatore ';'
					writer := csv.NewWriter(csvFile)
					writer.Comma = ';'

					// Scrivi l'intestazione del CSV
					err = writer.Write([]string{"#", "TIMECODE", "FILE1", "FILE2"})
					if err != nil {
						fmt.Println("Errore nel scrivere l'intestazione CSV:", err)
						return
					}

					// Scrivi i sottotitoli combinati nel CSV
					maxLength := len(subtitles1)
					if len(subtitles2) > maxLength {
						maxLength = len(subtitles2)
					}

					// Scrivi le righe del CSV
					for i := 0; i < maxLength; i++ {
						var subtitle1, subtitle2 srt.Subtitle
						if i < len(subtitles1) {
							subtitle1 = subtitles1[i]
						}
						if i < len(subtitles2) {
							subtitle2 = subtitles2[i]
						}

						// Crea una riga per il CSV con il testo del sottotitolo per ogni file
						record := []string{
							fmt.Sprint(subtitle1.Index),
							subtitle1.Start + " --->\n " + subtitle1.End,
							subtitle1.Text,
							subtitle2.Text,
						}

						// Scrivi la riga nel file CSV
						err = writer.Write(record)
						if err != nil {
							fmt.Println("Errore nel scrivere una riga nel CSV:", err)
							return
						}
					}

					// Salva tutte le righe nel CSV
					writer.Flush()

					// Verifica se ci sono errori nel flush
					if err := writer.Error(); err != nil {
						fmt.Println("Errore nel flush del file CSV:", err)
						return
					}

					zenity.Info("File generato correttamente :)", zenity.Title("Comparatore SRT"), zenity.NoIcon)
					os.Exit(0)
				},
			},
		},
	}.Run()
}
