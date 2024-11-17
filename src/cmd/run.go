package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/custodia-cenv/cenvx-core/src/core"
	"github.com/custodia-cenv/cenvx-core/src/log"
)

// Wird verwendet um den Core Service Offen zu halten
func RunCoreConsoleOrBackgroundService() {
	// Es wird eine neuer Waitgroup erzeugt,
	// diese Waitgroup wird verwendet um zu ermitteln
	// ob die Core Instanz ausgeführt wird
	var waitGroupForServing sync.WaitGroup
	waitGroupForServing.Add(1)

	// Der Core wird in einer eigenen Gouroutine ausgeführt
	go func() {
		// Der Core wird gestartet
		core.Serve()

		// Es wird Signalisiert dass der Core beendet wurde
		waitGroupForServing.Done()
	}()

	// Ein Channel, um Signale zu empfangen.
	sigChan := make(chan os.Signal, 1)

	// Wählen Sie die Signale basierend auf dem Betriebssystem aus
	switch runtime.GOOS {
	case "linux", "darwin": // macOS und Linux verwenden ähnliche Signale
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	case "windows":
		// Windows verwendet andere Signale und hat keine direkte Entsprechung für alle UNIX-Signale
		// Windows behandelt Signale wie Ctrl+C über os.Interrupt und kann das Herunterfahren durch andere APIs behandeln
		signal.Notify(sigChan, os.Interrupt)
	}

	// Diese Schleife wird ausgeführt um
	for {
		// Es wird auf das Signal gewartet
		sig := <-sigChan

		// Es wird anahnd des OS das Signal ausgewertet
		if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
			// Es wird versucht das Signal zu ermitteln
			if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGHUP || sig == syscall.SIGQUIT {
				// Bei Beenden-Signalen Core ordnungsgemäß herunterfahren
				core.SignalShutdown()

				// Es wird auf die Wait Group gewartet
				waitGroupForServing.Wait()

				// Die Schleife wird beenndet
				break
			} else {
				// Benutzerdefinierte Signalbehandlung
				log.InfoLogPrint("Received user-defined signal. Performing custom action...")
			}
		} else if runtime.GOOS == "windows" {
			// Es wird versucht das Signal zu ermitteln
			if sig == os.Interrupt {
				// Bei Beenden-Signalen Core ordnungsgemäß herunterfahren
				core.SignalShutdown()

				// Es wird auf die Wait Group gewartet
				waitGroupForServing.Wait()

				// Die Schleife wird beenndet
				break
			} else {
				log.LogError("Unhandled signal: %s\n", sig)
			}
		} else {
			panic("critical error 1")
		}
	}

	// Neue Zeile
	fmt.Printf("\n")
}
