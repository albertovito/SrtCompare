package srt

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Structure for storing a subtitle
type Subtitle struct {
	Index int
	Start string
	End   string
	Text  string
}

// Function to read SRT files and return an array of subtitles
func ReadSRTFile(filename string) ([]Subtitle, error) {
	// Open the SRT file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Variables to manage reading
	var subtitles []Subtitle
	var currentSubtitle Subtitle
	var currentText []string

	// Regex to identify timestamps
	timestampRegex := regexp.MustCompile(`(\d{2}:\d{2}:\d{2},\d{3}) --> (\d{2}:\d{2}:\d{2},\d{3})`)

	// Read the file line by line
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, "\ufeff", "", 1)
		line = strings.Replace(line, ";", ".", 1)
		// If the line is empty, it means the subtitle is complete
		if line == "" {
			if currentSubtitle.Index > 0 {
				// Add the full subtitle to the list
				currentSubtitle.Text = strings.Join(currentText, "\n")
				subtitles = append(subtitles, currentSubtitle)
			}

			// Reset for next subtitle
			currentText = nil
			currentSubtitle = Subtitle{}
		} else if currentSubtitle.Index == 0 {
			// First line: subtitle index
			fmt.Sscanf(line, "%d", &currentSubtitle.Index)
		} else if timestampRegex.MatchString(line) {
			// Second line: timestamp (Start --> End)
			matches := timestampRegex.FindStringSubmatch(line)
			currentSubtitle.Start = matches[1]
			currentSubtitle.End = matches[2]
		} else {
			// Next lines: subtitle text
			currentText = append(currentText, line)
		}
	}

	// Error handling during reading
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Handle last subtitle if last blank line is missing
	if currentSubtitle.Index > 0 {
		currentSubtitle.Text = strings.Join(currentText, "\n")
		subtitles = append(subtitles, currentSubtitle)
	}

	return subtitles, nil
}

// Function to format the subtitle as full text
func FormatSubtitle(subtitle *Subtitle) string {
	// Unifies the subtitle as a single string
	if subtitle.Index == 0 {
		return "" // Don't write anything if there are no subtitles
	}
	return fmt.Sprintf("%s\n%s --> %s\n%s", fmt.Sprint(subtitle.Index), subtitle.Start, subtitle.End, subtitle.Text)
}
