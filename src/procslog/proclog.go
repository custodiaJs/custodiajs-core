package procslog

import (
	"fmt"
	"log"
	"strings"

	cenvxcore "github.com/custodia-cenv/cenvx-core/src"
	"github.com/fatih/color"
)

var (
	foldEnd = color.New(color.Bold).SprintfFunc()(":-$")
)

func (o *ProcLogSession) LogPrint(header, format string, value ...interface{}) {
	// Die Eingabe wird formatiert
	userinput := fmt.Sprintf(format, value...)

	// Der Header wird ermittelt
	var foldedHeader string
	if header != "" {
		foldedHeader = color.New(color.Bold).SprintFunc()(header)
	} else {
		foldedHeader = color.New(color.Bold).SprintFunc()(o.header)
	}

	// Es wird geprüft ob Merged Procs vorhanden sind
	var finalValue string
	if len(o.merged) > 0 {
		// Speichert die Extrahierten Elemente ab
		elements := []string{}

		// Die Einzelnen Merged Elememente werden extrahiert
		for _, item := range o.merged {
			elements = append(elements, fmt.Sprintf("[%s] %s", item.sessionColorFunc(strings.ToUpper(item.id)), item.header))
		}

		// Der Neu Formatierte Text wird erstellt
		newFormated := strings.Join(elements, " > ")

		// Der Neue Finale Wert wird erzeugt
		finalValue = fmt.Sprintf("%s%s%s %s", newFormated, foldedHeader, foldEnd, userinput)
	} else {
		// Der Neue Finale Wert wird erzeugt
		finalValue = fmt.Sprintf("[%s] %s%s %s", o.sessionColorFunc(strings.ToUpper(o.id)), foldedHeader, foldEnd, userinput)
	}

	// Der Text wird angezeigt
	if o.printFunction != nil {
		o.printFunction(finalValue)
	} else {
		logPrint(finalValue)
	}
}

func (o *ProcLogSession) Log(format string, value ...interface{}) {
	o.LogPrint("", format, value...)
}

func (o *ProcLogSession) Debug(format string, value ...interface{}) {
	o.LogPrint("", format, value...)
}

func (o *ProcLogSession) LogPrintSuccs(format string, value ...interface{}) {
	o.LogPrint("", format, value...)
}

func (o *ProcLogSession) LogPrintError(format string, value ...interface{}) {
	o.LogPrint("", format, value...)
}

func (o *ProcLogSession) GetChildLog(header string) cenvxcore.ProcessLogSessionInterface {
	newProcLog := NewProcLogSessionWithHeader(header)
	merged := NewChainMergedProcLog(o, newProcLog)
	return merged
}

func (o *ProcLogSession) GetID() string {
	return o.id
}

func logPrint(text string) {
	log.Print(text)
}
