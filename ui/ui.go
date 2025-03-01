package ui

import (
	"SrtCompare/constants"
	"SrtCompare/dialog"
	"SrtCompare/srt"
	"SrtCompare/xlsx"
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

type UI struct {
	th               *material.Theme
	file1ButtonOp    *widget.Clickable
	file2ButtonOp    *widget.Clickable
	generateButtonOp *widget.Clickable
	editor1          *widget.Editor
	editor2          *widget.Editor
	file1Path        string
	file2Path        string
	list1            *widget.List
	list2            *widget.List
	MainWindow       *app.Window
}

func NewUI() *UI {
	ui := &UI{
		th:               material.NewTheme(),
		file1ButtonOp:    new(widget.Clickable),
		file2ButtonOp:    new(widget.Clickable),
		list1:            new(widget.List),
		list2:            new(widget.List),
		generateButtonOp: new(widget.Clickable),
		editor1:          &widget.Editor{ReadOnly: true, Submit: true, SingleLine: false},
		editor2:          &widget.Editor{ReadOnly: true, Submit: true, SingleLine: false},
	}
	ui.list1.Axis = layout.Vertical
	ui.list2.Axis = layout.Vertical
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

// Creates the main window layout
func (ui *UI) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Flexed(1, ui.layoutTextAreas),
		layout.Rigid(ui.layoutGenerateButton),
	)
}

// Displays two text areas side by side
func (ui *UI) layoutTextAreas(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(
		gtx,
		layout.Flexed(1, ui.layoutLeftPanel),
		layout.Flexed(1, ui.layoutRightPanel),
	)
}

// Layout fot left Panel
func (ui *UI) layoutLeftPanel(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Flexed(1, ui.layoutEditor(ui.editor1, ui.list1)),
		layout.Rigid(ui.layoutButton(ui.file1ButtonOp, "File1")),
	)
}

// Layut for right panel
func (ui *UI) layoutRightPanel(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(
		gtx,
		layout.Flexed(1, ui.layoutEditor(ui.editor2, ui.list2)),
		layout.Rigid(ui.layoutButton(ui.file2ButtonOp, "File2")),
	)
}

// Function that creates a reusable editor with borders
func (ui *UI) layoutEditor(editor *widget.Editor, list *widget.List) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(5)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			border := widget.Border{
				Color: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
				Width: unit.Dp(1),
			}
			return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.List(ui.th, list).Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions {
					return material.Editor(ui.th, editor, "").Layout(gtx)
				})
			})
		})
	}
}

// Function that creates a reusable button
func (ui *UI) layoutButton(button *widget.Clickable, text string) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		btn := material.Button(ui.th, button, text)
		btn.TextSize = unit.Sp(14)
		return layout.UniformInset(unit.Dp(5)).Layout(gtx, btn.Layout)
	}
}

// Layout for generation button
func (ui *UI) layoutGenerateButton(gtx layout.Context) layout.Dimensions {
	button := material.Button(ui.th, ui.generateButtonOp, "Generate")
	button.TextSize = unit.Sp(14)
	return layout.UniformInset(unit.Dp(5)).Layout(gtx, button.Layout)
}

func (ui *UI) openFile1Dialog(locked *bool) {
	if !*locked {
		*locked = true
		defer func() {
			*locked = false
			ui.MainWindow.Invalidate()
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
			ui.MainWindow.Invalidate()
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
			zenity.Error("Please select both files first", zenity.Title(constants.APPTITLE))
			return
		}

		subtitles1, err := srt.ReadSRTFile(ui.file1Path)
		if err != nil {
			zenity.Error(fmt.Sprintf("Error reading first SRT file: %v", err), zenity.Title(constants.APPTITLE))
			return
		}

		subtitles2, err := srt.ReadSRTFile(ui.file2Path)
		if err != nil {
			zenity.Error(fmt.Sprintf("Error reading second SRT file: %v", err), zenity.Title(constants.APPTITLE))
			return
		}

		filePath := dialog.GetSavePath("Save file")
		if filePath == "" {
			return
		}

		outputFile, err := os.Create(filePath)
		if err != nil {
			zenity.Error(fmt.Sprintf("Error creating XLSX file: %v", err), zenity.Title(constants.APPTITLE))
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
		xlsx.FormatSheet(file, subtitles1, subtitles2)
		if err := file.SaveAs(filePath); err != nil {
			zenity.Error(fmt.Sprintf("Error saving file: %v", err), zenity.Title(constants.APPTITLE))
			return
		}

		zenity.Info("File generated successfully", zenity.Title(constants.APPTITLE), zenity.NoIcon)
	}
	os.Exit(0)
}
