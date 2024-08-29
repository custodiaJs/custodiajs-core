package static

import (
	"path/filepath"

	"github.com/CustodiaJS/custodiajs-core/global/types"
)

var (
	// Gibt das Standard Config Verzeichniss an
	LINUX_HOST_CONFIG_DIR_PATH  types.HOST_CONFIG_PATH = types.HOST_CONFIG_PATH(filepath.Join("/", "etc", "CustodiaJS"))
	DARWIN_HOST_CONFIG_DIR_PATH types.HOST_CONFIG_PATH = types.HOST_CONFIG_PATH(filepath.Join("/", "Library", "Application Support", "CustodiaJS"))

	// Legt die Dateipfade für z.b Unix Sockets fest
	NONE_ROOT_UNIX_SOCKET types.SOCKET_PATH = types.SOCKET_PATH(filepath.Join("/", "tmp", "cusjs_none_root_sock"))
	ROOT_UNIX_SOCKET      types.SOCKET_PATH = types.SOCKET_PATH(filepath.Join("/", "tmp", "cusjs_root_sock"))

	// Die VM Datenbanken
	LINUX_DEFAULT_HOST_VM_DB_DIR_PATH  types.VM_DB_DIR_PATH = types.VM_DB_DIR_PATH(filepath.Join("var", "lib", "CustodiaJS", "vms"))
	DARWIN_DEFAULT_HOST_VM_DB_DIR_PATH types.VM_DB_DIR_PATH = types.VM_DB_DIR_PATH(filepath.Join("/", "Library", "Application Support", "CustodiaJS", "vms"))

	// Log Pfade
	LINUX_DEFAULT_LOGGING_DIR_PATH  types.LOG_DIR = types.LOG_DIR(filepath.Join("var", "log", "CustodiaJS"))
	DARWIN_DEFAULT_LOGGING_DIR_PATH types.LOG_DIR = types.LOG_DIR(filepath.Join("/", "Library", "Logs", "CustodiaJS"))

	// Legt den Pfad der Firmware für den Network Hypervisor fest
	LINUX_X86_DEFAULT_CNH_FIRMWARE    types.CNH_FIRMWARE_PATH = types.CNH_FIRMWARE_PATH(filepath.Join("var", "lib", "CustodiaJS", "hypervisorfrmw_x86"))
	LINUX_ARM_DEFAULT_CNH_FIRMWARE    types.CNH_FIRMWARE_PATH = types.CNH_FIRMWARE_PATH(filepath.Join("var", "lib", "CustodiaJS", "hypervisorfrmw_arm"))
	LINUX_RISCV_DEFAULT_CNH_FIRMWARE  types.CNH_FIRMWARE_PATH = types.CNH_FIRMWARE_PATH(filepath.Join("var", "lib", "CustodiaJS", "hypervisorfrmw_riscv"))
	LINUX_ARM64_DEFAULT_CNH_FIRMWARE  types.CNH_FIRMWARE_PATH = types.CNH_FIRMWARE_PATH(filepath.Join("var", "lib", "CustodiaJS", "hypervisorfrmw_arm64"))
	LINUX_AMD64_DEFAULT_CNH_FIRMWARE  types.CNH_FIRMWARE_PATH = types.CNH_FIRMWARE_PATH(filepath.Join("var", "lib", "CustodiaJS", "hypervisorfrmw_amd64"))
	DARWIN_ARM64_DEFAULT_CNH_FIRMWARE types.CNH_FIRMWARE_PATH = types.CNH_FIRMWARE_PATH(filepath.Join("/", "Library", "Application Support", "CustodiaJS", "hypervisorfrmw_arm64"))
	DARWIN_AMD64_DEFAULT_CNH_FIRMWARE types.CNH_FIRMWARE_PATH = types.CNH_FIRMWARE_PATH(filepath.Join("/", "Library", "Application Support", "CustodiaJS", "hypervisorfrmw_amd64"))

	// Gibt die Sockets für den Hypervisor an, wird verwendet damit der Hypervisor mit dem Host Kommunizieren kann
	LINUX_CNH_SOCKET_PATH  types.CHN_CORE_SOCKET_PATH = types.CHN_CORE_SOCKET_PATH("")
	DARWIN_CNH_SOCKET_PATH types.CHN_CORE_SOCKET_PATH = types.CHN_CORE_SOCKET_PATH("")
)
