package main

import (
	"comparatore_sottotitoli/dialog"
	"comparatore_sottotitoli/srt"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/ncruces/zenity"
)

func main() {

	// Leggi i sottotitoli da due file SRT
	res := dialog.MainPage()
	if !res {
		return
	}
	file1, file2 := dialog.SecondPage()

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
	csvFile, err := os.Create("sottotitoli.csv")
	if err != nil {
		fmt.Println("Errore nella creazione del file CSV:", err)
		return
	}
	defer csvFile.Close()

	// Crea un writer CSV con separatore ';'
	writer := csv.NewWriter(csvFile)
	writer.Comma = ';'

	// Scrivi l'intestazione del CSV
	err = writer.Write([]string{"FILE1", "FILE2"})
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
			srt.FormatSubtitle(&subtitle1),
			srt.FormatSubtitle(&subtitle2),
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
}
