//go:build windows

package main

import (
	"SrtComparator/dialog"
	"SrtComparator/srt"
	"fmt"
	"os"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/ncruces/zenity"
	"github.com/xuri/excelize/v2"
)

const (
	APPTITLE = "SrtCompare"
	FONT     = "Calibri 14"
	ICON     = "res/srt.ico"
)

func main() {
	startWindow()
}

func startWindow() {
	var outTE1, outTE2 *walk.TextEdit
	var file1, file2 string

	declarative.MainWindow{
		Title:   APPTITLE,
		Icon:    ICON,
		MinSize: declarative.Size{Width: 600, Height: 400},
		Layout:  declarative.VBox{},
		Children: []declarative.Widget{
			declarative.HSplitter{
				Children: []declarative.Widget{
					declarative.VSplitter{
						Children: []declarative.Widget{
							declarative.TextEdit{AssignTo: &outTE1, ReadOnly: true, VScroll: true, Font: declarative.Font{Family: FONT}, Background: declarative.SolidColorBrush{Color: walk.RGB(255, 255, 255)}},
							declarative.PushButton{
								Text: "File1",
								Font: declarative.Font{Family: FONT},
								OnClicked: func() {
									file1 = dialog.PrintDialog("")
									text1, _ := os.ReadFile(file1)
									outTE1.SetText(string(text1))
								},
							},
						},
					},
					declarative.VSplitter{
						Children: []declarative.Widget{
							declarative.TextEdit{AssignTo: &outTE2, ReadOnly: true, VScroll: true, Font: declarative.Font{Family: FONT}, Background: declarative.SolidColorBrush{Color: walk.RGB(255, 255, 255)}},
							declarative.PushButton{
								Text: "File2",
								Font: declarative.Font{Family: FONT},
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
			declarative.PushButton{
				Text: "Compara",
				Font: declarative.Font{Family: FONT},
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

					// Crea il file XLSX
					filePath := dialog.GetSavePath("Salva file")
					outputFile, err := os.Create(filePath)
					if err != nil {
						fmt.Println("Errore nella creazione del file XLSX:", err)
						return
					}
					defer outputFile.Close()

					file := excelize.NewFile()
					defer func() {
						if err := file.Close(); err != nil {
							fmt.Println(err)
						}
					}()
					sheet := "Sheet1"
					index, err := file.NewSheet(sheet)
					if err != nil {
						fmt.Println(err)
						return
					}

					maxLength := len(subtitles1)
					if len(subtitles2) > maxLength {
						maxLength = len(subtitles2)
					}

					// Scrivi le righe della tabella
					for i := 0; i < maxLength; i++ {
						var subtitle1, subtitle2 srt.Subtitle
						if i < len(subtitles1) {
							subtitle1 = subtitles1[i]
						}
						if i < len(subtitles2) {
							subtitle2 = subtitles2[i]
						}

						a := "A" + fmt.Sprint(i+1)
						b := "B" + fmt.Sprint(i+1)
						c := "C" + fmt.Sprint(i+1)
						d := "D" + fmt.Sprint(i+1)
						file.SetCellValue(sheet, a, subtitle1.Index)
						file.SetCellValue(sheet, b, subtitle1.Start+" --->\n "+subtitle1.End)
						file.SetCellValue(sheet, c, subtitle1.Text)
						file.SetCellValue(sheet, d, subtitle2.Text)

					}
					file.SetActiveSheet(index)
					// Scrivi la riga nel file XLSX
					if err := file.SaveAs(filePath); err != nil {
						fmt.Println(err)
					}
					zenity.Info("File generato correttamente", zenity.Title("Comparatore SRT"), zenity.NoIcon)
					os.Exit(0)
				},
			},
		},
	}.Run()
}
