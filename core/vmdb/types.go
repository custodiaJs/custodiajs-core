package vmdb

import (
	"os"
	"vnh1/utils"

	"github.com/gofrs/flock"
)

type VmDBEntry struct {
	Path                  string
	mainJSFile            *MainJsFile
	manifestFile          *ManifestFile
	signatureFile         *SignatureFile
	nodeJsModules         []*NodeJsModule
	vmContainerMerkleHash string
	containerBaseSize     uint64
}

type VmDatabase struct {
	vmMap     map[string]*VmDBEntry
	vmRootDir string
}

type SignatureFile struct {
	osFile         *os.File
	fileLock       *flock.Flock
	ownerCert      string
	ownerSignature string
}

type NodeJsModule struct {
	merkleRoot string
	files      []utils.FileInfo
	baseSize   uint64
	name       string
	alias      string
}

type ManifestFile struct {
	osFile   *os.File
	fileLock *flock.Flock
	fileHash string
	manifest *Manifest
	fSize    uint64
}

type MainJsFile struct {
	fileLock *flock.Flock
	fileSize uint64
	fileHash string
	filePath string
}

type CAMemberData struct {
	Fingerprint string
	Type        string
	ID          string
}

type VMDatabaseData struct {
	Type     string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Alias    string
}

type DatabaseFingerprint string
