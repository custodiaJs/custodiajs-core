# vnh1

( EN | [DE](../main/README_DE.md) )

## Beschreibung

**vnh1** bietet eine robuste Lösung für das sichere und isolierte Ausführen von Skripten durch die Kombination von Containerisierung und Sandboxing-Techniken. Mit vnh1 können V8go-Instanzen in separaten Prozessen und Containern ausgeführt werden, was eine starke Isolation und Sicherheit gewährleistet.

## Features

- **Containerisierte Sandbox-Umgebung**: Jede V8go-Instanz läuft in ihrem eigenen Container und sorgt so für vollständige Isolation von anderen Instanzen und dem Host-System.
- **Prozessisolierung**: V8go-Interpreter werden in separaten Prozessen ausgeführt, sodass jeder Prozess seine eigenen Ressourcen verwaltet.
- **Eingeschränkter Zugriff**: Die ausgeführten Skripte haben keinen Zugriff auf das Host-Dateisystem oder Netzwerk, wodurch das Risiko von Schäden oder Datenexfiltration minimiert wird.

### Containerisierung Details

- **Linux**:
  Unter Linux nutzt vnh1 Namespaces, um eine isolierte Containerumgebung für jede V8go-Instanz zu erstellen. Dies stellt sicher, dass jeder Container sein eigenes Set von Benutzern, Prozessen und Dateisystemen hat. Durch die Verwendung von Linux-Namespaces kann vnh1 ein hohes Maß an Sicherheit und Isolation bieten, wodurch jede Interferenz zwischen laufenden Skripten und dem Host-System effektiv verhindert wird.

- **macOS**:
  Für macOS verwendet vnh1 launchd, um die Containerisierung der VM-Prozesse zu verwalten. Launchd ist ein Service-Management-Framework für macOS, das Daemons, Anwendungen, Prozesse und Skripte starten, stoppen und verwalten kann. Durch die Nutzung von launchd stellt vnh1 sicher, dass jeder VM-Prozess ordnungsgemäß isoliert und verwaltet wird, was eine sichere Umgebung für die Skriptausführung bietet.

- **Windows**:
  Unter Windows verwendet vnh1 keine traditionellen Containerisierungstechniken. Stattdessen wird jede V8go-Instanz als separater Prozess ausgeführt. Obwohl diese Prozesse voneinander isoliert sind, profitieren sie nicht von der zusätzlichen Sicherheitsschicht, die durch Containerisierung bereitgestellt wird. Das bedeutet, dass Ressourcen zwar unabhängig für jede VM verwaltet werden, aber nicht so sicher isoliert sind wie unter Linux oder macOS.

- **BSD**:
  Für BSD-Systeme plant vnh1 die Implementierung einer Containerisierung, die der unter Linux ähnelt. BSD bietet mehrere Containerisierungstechnologien wie Jails, die eine Möglichkeit bieten, das Betriebssystem in mehrere unabhängige Mini-Systeme zu partitionieren. Jede Jail hat ihr eigenes Dateisystem, ihre eigenen Benutzer und Prozesse und gewährleistet so starke Isolation und Sicherheit. Durch die Nutzung von BSD Jails zielt vnh1 darauf ab, eine sichere und isolierte Umgebung für jede V8go-Instanz zu schaffen, ähnlich dem Ansatz, der unter Linux mit Namespaces verwendet wird.

## Funktionalität

- **RPC-Funktionsaufrufe** ✅
- **Netzwerkfunktionen, Sockets, etc.** ⚠️
- **HTTP-Client/Server-Funktionen** ⚠️
- **Electrum-Support** ⚠️
- **Lightning-Support** ⚠️
- **Nostr-Support** ⚠️
- **Datenbank-Support (MySQL/SQLite, MongoDB)** ⚠️
- **Krypto-Funktionen (SSL, ECC, RSA, PGP)** ⚠️
- **NodeJS (JS)-Konsolenfunktionen** ⚠️
- **Dateisystemzugriff** ⚠️
- **Wireguard-Verwaltung** ⚠️

## Roadmap

- **Linux**: Wird aktuell ausgebaut ⚠️
- **Windows**: Geplant mit eingeschränktem Support ⚠️
- **MacOS**: Geplanter Support ⚠️
- **BSD**: Geplanter Support ⚠️
