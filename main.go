package main

import (
	"SrtCompare/dialog"
	"SrtCompare/srt"
	"fmt"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/ncruces/zenity"
	"github.com/xuri/excelize/v2"
)

const (
	APPTITLE = "SrtCompare"
)

type UI struct {
	th               *material.Theme
	file1ButtonOp    *widget.Clickable
	file2ButtonOp    *widget.Clickable
	generateButtonOp *widget.Clickable
	editor1          *widget.Editor
	editor2          *widget.Editor
	file1Path        string
	file2Path        string
	window           *app.Window
}

func main() {
	ui := NewUI()

	go func() {
		// Build the window
		w := new(app.Window)
		w.Option(app.Title(APPTITLE))
		ui.window = w
		if err := ui.Run(w); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}

func NewUI() *UI {
	ui := &UI{
		th:               material.NewTheme(),
		file1ButtonOp:    new(widget.Clickable),
		file2ButtonOp:    new(widget.Clickable),
		generateButtonOp: new(widget.Clickable),
		editor1:          &widget.Editor{ReadOnly: true, Submit: true, SingleLine: false},
		editor2:          &widget.Editor{ReadOnly: true, Submit: true, SingleLine: false},
	}
	return ui
}

func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops
	var locked bool
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Save if buttons were clicked
			file1Clicked := ui.file1ButtonOp.Clicked(gtx)
			file2Clicked := ui.file2ButtonOp.Clicked(gtx)
			generateClicked := ui.generateButtonOp.Clicked(gtx)

			ui.Layout(gtx)

			if file1Clicked {
				go ui.openFile1Dialog(&locked)
			}

			if file2Clicked {
				go ui.openFile2Dialog(&locked)
			}

			if generateClicked {
				go ui.generateExcel(&locked)
			}

			e.Frame(gtx.Ops)
		}
	}
}

func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(
				gtx,
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(
						gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								border := widget.Border{
									Color: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
									Width: unit.Dp(1),
								}
								return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return material.Editor(ui.th, ui.editor1, "").Layout(gtx)
								})
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							button := material.Button(ui.th, ui.file1ButtonOp, "File1")
							button.TextSize = unit.Sp(14)
							return layout.UniformInset(unit.Dp(5)).Layout(gtx, button.Layout)
						}),
					)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(
						gtx,
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								border := widget.Border{
									Color: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
									Width: unit.Dp(1),
								}
								return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return material.Editor(ui.th, ui.editor2, "").Layout(gtx)
								})
							})
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							button := material.Button(ui.th, ui.file2ButtonOp, "File2")
							button.TextSize = unit.Sp(14)
							return layout.UniformInset(unit.Dp(5)).Layout(gtx, button.Layout)
						}),
					)
				}),
			)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			button := material.Button(ui.th, ui.generateButtonOp, "Generate")
			button.TextSize = unit.Sp(14)
			return layout.UniformInset(unit.Dp(10)).Layout(gtx, button.Layout)
		}),
	)
}

func (ui *UI) openFile1Dialog(locked *bool) {
	if !*locked {
		*locked = true
		defer func() {
			*locked = false
			ui.window.Invalidate()
		}()

		filePath := dialog.SelectFile("Open File1")
		ui.file1Path = filePath

		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		ui.editor1.SetText(string(content))
	}
}

func (ui *UI) openFile2Dialog(locked *bool) {
	if !*locked {
		*locked = true
		defer func() {
			*locked = false
			ui.window.Invalidate()
		}()

		filePath := dialog.SelectFile("Open File2")
		ui.file2Path = filePath

		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		ui.editor2.SetText(string(content))
	}
}

func (ui *UI) generateExcel(locked *bool) {
	if !*locked {
		*locked = true
		if ui.file1Path == "" || ui.file2Path == "" {
			zenity.Error("Please select both files first", zenity.Title(APPTITLE))
			return
		}

		subtitles1, err := srt.ReadSRTFile(ui.file1Path)
		if err != nil {
			zenity.Error(fmt.Sprintf("Error reading first SRT file: %v", err), zenity.Title(APPTITLE))
			return
		}

		subtitles2, err := srt.ReadSRTFile(ui.file2Path)
		if err != nil {
			zenity.Error(fmt.Sprintf("Error reading second SRT file: %v", err), zenity.Title(APPTITLE))
			return
		}

		filePath := dialog.GetSavePath("Save file")
		if filePath == "" {
			return
		}

		outputFile, err := os.Create(filePath)
		if err != nil {
			zenity.Error(fmt.Sprintf("Error creating XLSX file: %v", err), zenity.Title(APPTITLE))
			return
		}
		defer outputFile.Close()

		file := excelize.NewFile()
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Println(err)
			}
			*locked = false
		}()

		sheet := "Sheet1"
		index, err := file.NewSheet(sheet)
		if err != nil {
			zenity.Error(fmt.Sprintf("Error creating sheet: %v", err), zenity.Title(APPTITLE))
			return
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

		if err := file.SaveAs(filePath); err != nil {
			zenity.Error(fmt.Sprintf("Error saving file: %v", err), zenity.Title(APPTITLE))
			return
		}

		zenity.Info("File generated successfully", zenity.Title(APPTITLE), zenity.NoIcon)
	}
	os.Exit(0)
}
