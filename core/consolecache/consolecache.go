package consolecache

import (
	"fmt"
	"sync"
	"time"
)

type ConsoleLogEntry struct {
	otype     string
	value     string
	timestamp int64
}

type ConsoleOutputCache struct {
	logs        []*ConsoleLogEntry
	liveConsole []chan ConsoleLogEntry
	syncLock    *sync.Mutex
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
	o.ConsoleWrite("error", value)
}

func (o *ConsoleOutputCache) InfoLog(value string) {
	o.ConsoleWrite("info", value)
}

func (o *ConsoleOutputCache) Log(value string) {
	o.ConsoleWrite("log", value)
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

func NewConsoleOutputCache() *ConsoleOutputCache {
	return &ConsoleOutputCache{
		logs:        make([]*ConsoleLogEntry, 0),
		liveConsole: make([]chan ConsoleLogEntry, 0),
		syncLock:    &sync.Mutex{},
	}
}
