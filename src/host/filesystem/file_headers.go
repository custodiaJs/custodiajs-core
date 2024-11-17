package filesystem

// ELFHeader enthält den Header einer ELF-Datei
type elfheader struct {
	Ident     [16]byte // ELF Identifikation
	Type      uint16   // Typ des Objekts
	Machine   uint16   // Zielrechner
	Version   uint32   // ELF Version
	Entry     uint64   // Eintragspunkt des Programms
	Phoff     uint64   // Offset des Program Header Table
	Shoff     uint64   // Offset des Section Header Table
	Flags     uint32   // Prozessor-spezifische Flags
	Ehsize    uint16   // Größe des ELF Header
	Phentsize uint16   // Größe eines Eintrags in der Program Header Table
	Phnum     uint16   // Anzahl der Einträge in der Program Header Table
	Shentsize uint16   // Größe eines Eintrags in der Section Header Table
	Shnum     uint16   // Anzahl der Einträge in der Section Header Table
	Shstrndx  uint16   // Index der Section Header Table, die die Sektionsnamen enthält
}

// DosHeader ist der DOS-Header einer PE-Datei
type DosHeader struct {
	Magic    uint16
	Used     [58]uint8
	LfaNew   uint32
	Reserved uint16
	ExeType  uint16
}

// PeHeader ist der PE-Header einer PE-Datei
type PeHeader struct {
	PeSignature uint32
	FileHeader  [20]byte
}

// CliHeader ist der CLI-Header einer .NET-DLL
type CliHeader struct {
	Signature           uint32
	HeaderSize          uint32
	MinorRuntimeVer     uint16
	MajorRuntimeVer     uint16
	MetaData            uint32
	Flags               uint32
	EntryPointToken     uint32
	Resources           uint32
	StrongNameSig       uint32
	CodeManagerTable    uint32
	VTableFixups        uint32
	ExportAddressTable  uint32
	ManagedNativeHeader uint32
}

// MachOHeader enthält den Header einer Mach-O-Datei
type MachOHeader struct {
	Magic      uint32 // Mach-O Magic Number
	Cputype    uint32 // CPU-Typ
	Cpusubtype uint32 // CPU-Subtyp
	Filetype   uint32 // Dateityp
	Ncmds      uint32 // Anzahl der Befehle
	Cmdsize    uint32 // Größe der Befehle
	Flags      uint32 // Flags
}
