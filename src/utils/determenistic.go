// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package utils

import (
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
		panic(err)
	}

	// Bestimme den Index durch Modulo-Operation
	index := hashValue % uint64(len(colors))

	// Gebe die entsprechende Farbe zurück
	return colors[index]
}
