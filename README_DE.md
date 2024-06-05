# vnh1

( DE | [EN](../README) )

## Beschreibung

**vnh1** bietet eine robuste Lösung zur sicheren und isolierten Ausführung von Skripten durch die Kombination von Containerisierung und Sandboxing-Techniken. Mit vnh1 können V8go-Instanzen in separaten Prozessen und Containern ausgeführt werden, was eine starke Isolierung und Sicherheit gewährleistet.

## Features

- **Containerisierte Sandboxing-Umgebung**: Jede V8go-Instanz läuft in einem eigenen Container, was eine vollständige Isolierung von anderen Instanzen und vom Hostsystem gewährleistet.
- **Prozessisolation**: V8go-Interpreter werden in separaten Prozessen ausgeführt, um sicherzustellen, dass jeder Prozess seine eigenen Ressourcen verwaltet.
- **Eingeschränkter Zugriff**: Die ausgeführten Skripte haben keinen Zugriff auf das Host-Dateisystem oder das Netzwerk, was das Risiko von Schäden oder Datenexfiltration minimiert.

## Vorteile

- **Sicherheit**: Durch die Kombination von Containerisierung und Sandboxing wird eine sichere Ausführungsumgebung geschaffen.
- **Flexibilität**: Mehrere Skripte können parallel und unabhängig voneinander ausgeführt werden, ohne dass sie sich gegenseitig beeinflussen.
- **Isolation**: Strikte Trennung zwischen den Skripten und dem Hostsystem verhindert unerwünschte Interaktionen und erhöht die Sicherheit.

## Anwendungsfälle

- Sichere Ausführung von benutzerdefiniertem JavaScript-Code.
- Bereitstellung einer isolierten Umgebung für Skript-basierte Automatisierungen.
- Entwickeln und Testen von JavaScript-Code in einer kontrollierten Umgebung.
