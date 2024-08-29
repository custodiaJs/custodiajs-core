package image

func (o *MainJsFile) Content() string {
	return o.content
}

func (o *MainJsFile) GetHash() string {
	return o.fileHash
}
