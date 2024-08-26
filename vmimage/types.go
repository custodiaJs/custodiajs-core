package vmimage

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
	manifest  *Manifest
}

type ImageSignature struct {
}
