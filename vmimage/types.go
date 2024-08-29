package vmimage

import "github.com/CustodiaJS/custodiajs-core/types"

type VMEntryBaseData struct {
	Type     string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Alias    string
}

type MainJsFile struct {
	content  string
	fileSize uint64
	fileHash string
}

type VmImage struct {
	mainFile  *MainJsFile
	signature *ImageSignature
	manifest  *types.Manifest
}

type ImageSignature struct {
}
