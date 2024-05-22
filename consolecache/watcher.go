package consolecache

type Watcher struct {
	transport chan ConsoleLogEntry
}

func (o *Watcher) Read() string {
	ch := <-o.transport
	return ch.value
}
