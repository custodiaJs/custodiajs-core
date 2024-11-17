// Author: fluffelpuff
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package filesystem

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
)

// GetFileSize gibt die Größe einer Datei in Bytes zurück.
func GetFileSize(filePath string) (int64, error) {
	// Öffne die Datei
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		// Gibt einen Fehler zurück, falls das Öffnen fehlschlägt
		return 0, err
	}

	// Die Größe der Datei in Bytes
	return fileInfo.Size(), nil
}

// ExtractFileName nimmt einen Dateipfad als Eingabe und gibt den Dateinamen zurück.
func ExtractFileName(filePath string) string {
	return filepath.Base(filePath)
}

// IsUnixSOFile überprüft, ob eine Datei eine Unix Shared Object-Datei ist
func IsUnixSOFile(filePath string) bool {
	// Öffne die Datei
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Lese den ELF-Header
	var elfHeader elfheader
	if err := binary.Read(file, binary.LittleEndian, &elfHeader); err != nil {
		return false
	}

	// Überprüfe die ELF-Magic-Number
	if !bytes.Equal(elfHeader.Ident[:4], []byte{0x7f, 'E', 'L', 'F'}) {
		return false
	}

	// Überprüfe den Typ der ELF-Datei
	if elfHeader.Type != 3 { // 3 entspricht ET_DYN, dem Typ für Shared Objects
		return false
	}

	return true
}

// IsDotNetDLL überprüft, ob eine Datei eine .NET-DLL ist
func IsDotNetDLL(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Lese den DOS-Header
	var dosHeader DosHeader
	if err := binary.Read(file, binary.LittleEndian, &dosHeader); err != nil {
		return false
	}

	// Überprüfe, ob es sich um eine PE-Datei handelt
	if dosHeader.Magic != 0x5A4D { // "MZ"
		return false
	}

	// Springe zum PE-Header
	if _, err := file.Seek(int64(dosHeader.LfaNew), 0); err != nil {
		return false
	}

	// Lese den PE-Header
	var peHeader PeHeader
	if err := binary.Read(file, binary.LittleEndian, &peHeader); err != nil {
		return false
	}

	// Überprüfe die PE-Signatur
	if peHeader.PeSignature != 0x00004550 { // "PE\0\0"
		return false
	}

	// Springe zum CLI-Header
	cliHeaderOffset := int64(dosHeader.LfaNew) + int64(binary.Size(peHeader)) + 0x60 // 0x60 ist die Offset-Adresse des CLI-Headers im PE-Header
	if _, err := file.Seek(cliHeaderOffset, 0); err != nil {
		return false
	}

	// Lese den CLI-Header
	var cliHeader CliHeader
	if err := binary.Read(file, binary.LittleEndian, &cliHeader); err != nil {
		return false
	}

	// Überprüfe die Signatur des CLI-Headers
	if cliHeader.Signature != 0x424A5342 { // "BSJB"
		return false
	}

	return true
}

// IsWindowsDLL überprüft, ob eine Datei eine normale Windows-DLL ist
func IsWindowsDLL(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Lese den DOS-Header
	var dosHeader DosHeader
	if err := binary.Read(file, binary.LittleEndian, &dosHeader); err != nil {
		return false
	}

	// Überprüfe, ob es sich um eine PE-Datei handelt
	if dosHeader.Magic != 0x5A4D { // "MZ"
		return false
	}

	// Springe zum PE-Header
	if _, err := file.Seek(int64(dosHeader.LfaNew), 0); err != nil {
		return false
	}

	// Lese den PE-Header
	var peHeader PeHeader
	if err := binary.Read(file, binary.LittleEndian, &peHeader); err != nil {
		return false
	}

	// Überprüfe die PE-Signatur
	if peHeader.PeSignature != 0x00004550 { // "PE\0\0"
		return false
	}

	// Überprüfe den Dateityp im PE-Header
	// Der Dateityp 0x02 entspricht IMAGE_FILE_DLL (DLL-Datei)
	if binary.LittleEndian.Uint16(peHeader.FileHeader[0:2]) == 0x02 {
		return true
	}

	return false
}

// IsDylib überprüft, ob eine Datei eine dylib unter macOS ist
func IsDylib(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	// Lese die ersten 4 Bytes, um den Dateityp zu erkennen
	magic := make([]byte, 4)
	_, err = file.Read(magic)
	if err != nil {
		return false
	}

	// Überprüfe die Mach-O-Magic-Number
	if binary.LittleEndian.Uint32(magic) == 0xfeedface {
		// Wenn die Magic Number übereinstimmt, ist es eine Mach-O-Datei
		// Überprüfe nun den Dateityp im Mach-O-Header
		var machoHeader MachOHeader
		if err := binary.Read(file, binary.LittleEndian, &machoHeader); err != nil {
			return false
		}

		// Der Dateityp 6 entspricht MH_DYLIB (dylib-Datei)
		if machoHeader.Filetype == 6 {
			return true
		}
	}

	return false
}
