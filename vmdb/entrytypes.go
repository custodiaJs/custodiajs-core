package vmdb

// Definiere die Strukturen, die dem JSON-Format entsprechen
type Manifest struct {
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Owner        string       `json:"owner"`
	RepoURL      string       `json:"repourl"`
	Mode         string       `json:"mode"`
	Whitelist    []Whitelist  `json:"whitelist"`
	HostCAMember []CAMember   `json:"hostcamember"`
	Databases    []Database   `json:"databases"`
	Scripts      ScriptConfig `json:"scripts"`
}

type Whitelist struct {
	Alias   string   `json:"alias"`
	URL     string   `json:"url"`
	Methods []string `json:"methods"`
}

type CAMember struct {
	Fingerprint string `json:"fingerprint"`
	Type        string `json:"type"`
}

type Database struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Alias    string `json:"alias"`
}

type ScriptConfig struct {
	Python ScriptDetail `json:"python"`
	NodeJS ScriptDetail `json:"nodejs"`
}

type ScriptDetail struct {
	Enable  bool           `json:"enable"`
	Modules []ScriptModule `json:"modules"`
}

type ScriptModule struct {
	Alias          string `json:"alias"`
	StartCommand   string `json:"startcommand"`
	InstallCommand string `json:"installcommand,omitempty"`
	Name           string `json:"name"`
}
