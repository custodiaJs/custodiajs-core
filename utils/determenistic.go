package utils

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
)

func DetermineColorFromHex(hex string) func(a ...interface{}) string {
	// Liste der möglichen Farben
	colors := []func(a ...interface{}) string{
		color.New(color.FgRed).SprintFunc(),
		color.New(color.FgYellow).SprintFunc(),
		color.New(color.FgBlue).SprintFunc(),
		color.New(color.FgHiGreen).SprintFunc(),
		color.New(color.FgMagenta).SprintFunc(),
	}

	// Konvertiere den Hex-String in eine Ganzzahl
	hashValue, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		fmt.Println("Fehler beim Parsen des Hex-Strings:", err)
		return nil
	}

	// Bestimme den Index durch Modulo-Operation
	index := hashValue % uint64(len(colors))

	// Gebe die entsprechende Farbe zurück
	return colors[index]
}
