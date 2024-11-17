package utils

import (
	"fmt"
	"strings"
)

// Funktion zum Formatieren der Zahl mit Tausender-Trennzeichen
func FormatNumberWithDots(number int) string {
	// Wandelt die Zahl in einen String um
	numStr := fmt.Sprintf("%d", number)

	// Länge des Strings
	length := len(numStr)

	// Prüfen, ob alle Ziffern außer der ersten `0` sind
	if strings.TrimRight(numStr[1:], "0") == "" {
		return "1.0"
	}

	// Wenn die Länge kleiner oder gleich 3 ist, gibt die Zahl direkt zurück
	if length <= 3 {
		return numStr
	}

	// Initialisiert einen StringBuilder
	var result strings.Builder

	// Variable für den Zählindex
	count := 0

	// Durchläuft den String von hinten nach vorne
	for i := length - 1; i >= 0; i-- {
		if count > 0 && count%3 == 0 {
			result.WriteString(".")
		}
		result.WriteByte(numStr[i])
		count++
	}

	// Der resultierende String wird umgekehrt, da er von hinten aufgebaut wurde
	resultStr := result.String()
	formattedNumber := reverse(resultStr)

	return formattedNumber
}

// Hilfsfunktion zum Umkehren eines Strings
func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
