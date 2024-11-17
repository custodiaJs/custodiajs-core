package cenvxcore

type VmWorkingDir string
type VmProcessId uint64
type QualifiedVmID string

// VM und Core Status Typen sowie Repo Datentypen
type ALTERNATIVE_SERVICE_PATH string        // Alternativer Socket Path
type VmState uint8                          // VM Status
type CoreState uint8                        // Core Status
type IPCRight uint8                         // CLI Benutzerrecht
type VERSION uint32                         // Version des Hauptpgrogrammes
type REPO string                            // URL der Sourccode Qeulle
type SOCKET_PATH string                     // Gibt einen Socket Path an
type LOG_DIR string                         // Gibt den Path des Log Dir's unter
type HOST_CRYPTOSTORE_WATCH_DIR_PATH string // Gibt den Ordner an, in dem sich alle Zertifikate und Schlüssel des Hosts befinden
type HOST_CONFIG_FILE_PATH string           // Gibt den Pfad der Config Datei an
type HOST_CONFIG_PATH string
type CHN_CORE_SOCKET_PATH string

// Gibt die QUID einer VM an
type QVMID string

// Gibt den Hash eines Scriptes an
type VmScriptHash string

// Gibt die ProcessId an
type ProcessId string

// Gibt die VmID an
type VmId string

// Gibt an ob es sich bei einer IP um eine TOR IP-handelt
type TorIpState bool
