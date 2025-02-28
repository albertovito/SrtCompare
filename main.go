package main

import (
	"SrtCompare/constants"
	"SrtCompare/ui"
	"fmt"
	"os"

	"gioui.org/app"
)

func main() {
	window := ui.NewUI()

	go func() {
		// Build the window
		w := new(app.Window)
		w.Option(app.Title(constants.APPTITLE))
		window.MainWindow = w
		if err := window.Run(w); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}
