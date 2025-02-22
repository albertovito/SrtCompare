package srt

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Struttura per memorizzare un sottotitolo
type Subtitle struct {
	Index int
	Start string
	End   string
	Text  string
}

// Funzione per leggere i file SRT e restituire un array di sottotitoli
func ReadSRTFile(filename string) ([]Subtitle, error) {
	// Apri il file SRT
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Scanner per leggere il file riga per riga
	scanner := bufio.NewScanner(file)

	// Variabili per gestire la lettura
	var subtitles []Subtitle
	var currentSubtitle Subtitle
	var currentText []string

	// Regex per identificare i timestamp
	timestampRegex := regexp.MustCompile(`(\d{2}:\d{2}:\d{2},\d{3}) --> (\d{2}:\d{2}:\d{2},\d{3})`)

	// Leggi il file riga per riga
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, "\ufeff", "", 1)
		line = strings.Replace(line, ";", ".", 1)
		// Se la riga è vuota, significa che il sottotitolo è completo
		if line == "" {
			if currentSubtitle.Index > 0 {
				// Aggiungi il sottotitolo completo alla lista
				currentSubtitle.Text = strings.Join(currentText, "\n")
				subtitles = append(subtitles, currentSubtitle)
			}

			// Resetta per il prossimo sottotitolo
			currentText = nil
			currentSubtitle = Subtitle{}
		} else if currentSubtitle.Index == 0 {
			// Prima riga: indice del sottotitolo
			fmt.Sscanf(line, "%d", &currentSubtitle.Index)
		} else if timestampRegex.MatchString(line) {
			// Seconda riga: timestamp (Start --> End)
			matches := timestampRegex.FindStringSubmatch(line)
			currentSubtitle.Start = matches[1]
			currentSubtitle.End = matches[2]
		} else {
			// Righe successive: testo del sottotitolo
			currentText = append(currentText, line)
		}
	}

	// Gestione errori durante la lettura
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Gestisci l'ultimo sottotitolo se manca l'ultima riga vuota
	if currentSubtitle.Index > 0 {
		currentSubtitle.Text = strings.Join(currentText, "\n")
		subtitles = append(subtitles, currentSubtitle)
	}

	return subtitles, nil
}

// Funzione per formattare il sottotitolo come testo completo
func FormatSubtitle(subtitle *Subtitle) string {
	// Unifica il sottotitolo come una singola stringa
	if subtitle.Index == 0 {
		return "" // Non scrivere nulla se non ci sono sottotitoli
	}
	return fmt.Sprintf("%s\n%s --> %s\n%s", fmt.Sprint(subtitle.Index), subtitle.Start, subtitle.End, subtitle.Text)
}
