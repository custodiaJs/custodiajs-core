package consolecache

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type ConsoleOutputFormatter uint8

type ConsoleLogEntry struct {
	otype     string
	value     string
	timestamp int64
}

type ConsoleOutputCache struct {
	logs        []*ConsoleLogEntry
	liveConsole []chan ConsoleLogEntry
	logChan     chan *ConsoleLogEntry
	syncLock    *sync.Mutex
	logFile     *os.File
}

func (f *ConsoleOutputFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Erstelle die Basis-Nachricht
	msg := fmt.Sprintf("%s - [%s] %s", entry.Time.Format("2006-01-02 15:04:05"), entry.Level, entry.Message)

	// Füge nur die Datenfelder hinzu, wenn sie nicht leer sind
	if len(entry.Data) > 0 {
		msg += fmt.Sprintf(": %v", entry.Data)
	}

	// Füge einen Zeilenumbruch hinzu
	msg += "\n"

	return []byte(msg), nil
}

func (o *ConsoleOutputCache) enterToAllClient(clogentry *ConsoleLogEntry) {
	// Der Synclock wird verwendet
	o.syncLock.Lock()
	defer o.syncLock.Unlock()

	// Es werden alle Clients abgearbeitet
	for _, item := range o.liveConsole {
		item <- *clogentry
	}
}

func (o *ConsoleOutputCache) ConsoleWrite(typev string, value string) {
	// Der Synclock wird verwendet
	o.syncLock.Lock()
	defer o.syncLock.Unlock()

	// Es wird geprüft ob mehr als 1024 Einträge vorhanden sind, wenn ja wird der erste gelöscht
	if len(o.logs) >= 1024 {
		o.logs = o.logs[1:]
	}

	// Die Aktuelle Zeit wird erfasst
	time := time.Now()
	unixtime := time.Unix()

	// Das neue Objekt wird erzeugt
	newObj := &ConsoleLogEntry{
		otype:     typev,
		value:     value,
		timestamp: unixtime,
	}

	// Der Eintrag wird hinzugefügt
	o.logs = append(o.logs, newObj)

	// Der Eintrag wird an alle Clients gesendet
	go o.enterToAllClient(newObj)
}

func (o *ConsoleOutputCache) ErrorLog(value string) {
	o.logChan <- &ConsoleLogEntry{otype: "error", value: value, timestamp: time.Now().Unix()}
	go o.ConsoleWrite("error", value)
}

func (o *ConsoleOutputCache) InfoLog(value string) {
	o.logChan <- &ConsoleLogEntry{otype: "info", value: value, timestamp: time.Now().Unix()}
	go o.ConsoleWrite("info", value)
}

func (o *ConsoleOutputCache) Log(value string) {
	o.logChan <- &ConsoleLogEntry{otype: "log", value: value, timestamp: time.Now().Unix()}
	go o.ConsoleWrite("log", value)
}

func (o *ConsoleOutputCache) GetOutputStream() *Watcher {
	marstring := make(chan ConsoleLogEntry)

	w := &Watcher{
		transport: marstring,
	}

	o.liveConsole = append(o.liveConsole, marstring)
	fmt.Println("CONSOLE_CACHE_STREAM_CREATED")

	go func() {
		for _, item := range o.logs {
			marstring <- *item
		}
	}()

	return w
}

func (o *ConsoleOutputCache) Close() {
	o.InfoLog("Closing instance")
}

func (o *ConsoleOutputCache) _logRoutine(logger *logrus.Logger) {
	for {
		vale := <-o.logChan
		switch vale.otype {
		case "error":
			logger.Error(vale.value)
		case "info":
			logger.Info(vale.value)
		default:
			logger.Info(vale.value)
		}
		o.ConsoleWrite(vale.otype, vale.value)
	}
}

func (o *ConsoleOutputCache) _init() error {
	// Der Log Loggers wird erzeugt
	logger := logrus.New()
	logger.Out = o.logFile
	logger.SetFormatter(new(ConsoleOutputFormatter))

	// Init Message
	logger.Info("New instance started")

	// Log Routine
	go o._logRoutine(logger)

	// Rückgabe
	return nil
}

func NewConsoleOutputCache(loggingPath string) (*ConsoleOutputCache, error) {
	// Die Neuen Paths werden erezeugt
	logConsoleFilePath := path.Join(loggingPath, "log.console.txt")

	// Die Log Dateien werden erzeugt
	logFile, err := os.OpenFile(logConsoleFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("ConsoleOutputCache->NewConsoleOutputCache: " + err.Error())
	}

	// Das Objekt wird erstellt
	obj := &ConsoleOutputCache{
		logs:        make([]*ConsoleLogEntry, 0),
		liveConsole: make([]chan ConsoleLogEntry, 0),
		logChan:     make(chan *ConsoleLogEntry),
		syncLock:    &sync.Mutex{},
		logFile:     logFile,
	}

	// Das Objekt wird Initialisiert
	if err := obj._init(); err != nil {
		return nil, fmt.Errorf("ConsoleOutputCache->NewConsoleOutputCache: " + err.Error())
	}

	// Das Objelt wird zurückgegeben
	return obj, nil
}
