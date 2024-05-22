package vmdb

type Manifest struct {
	Name         string      `json:"name"`
	Version      string      `json:"version"`
	Owner        string      `json:"owner"`
	RepoURL      string      `json:"repourl"`
	Mode         string      `json:"mode"`
	Whitelist    []Whitelist `json:"whitelist"`
	HostCAMember []CAMember  `json:"hostcamember"`
	Databases    []Database  `json:"databases"`
	Services     Services    `json:"services"`
}

type Services struct {
	Experimental ExperimentalServices `json:"experimental"`
	External     []ExternalService    `json:"external"`
}

type ExternalService struct {
	Name       string `json:"name"`
	MinVersion int    `json:"min_version"`
	Required   bool   `json:"required"`
}

type ExperimentalServices struct {
	Webservice []WebService `json:"webservice"`
}

type Whitelist struct {
	Endpoint struct {
		Domain struct {
			Wildcards []string `json:"wildcards"`
			Exact     []string `json:"exact"`
		} `json:"domain"`
		IPv4List []string `json:"ipv4list"`
		IPv6List []string `json:"ipv6list"`
	} `json:"endpoint"`
	Methods []string `json:"methods"`
}

type CAMember struct {
	Fingerprint string `json:"fingerprint"`
	Type        string `json:"type"`
	ID          string `json:"id"`
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

type WebService struct {
	ID       string `json:"id"`
	Port     int    `json:"port"`
	Domain   string `json:"domain"`
	SSLOwner struct {
		ByID string `json:"byid"`
	} `json:"sslowner"`
	PHP struct {
		Use     bool   `json:"use"`
		Version string `json:"version"`
	} `json:"php"`
}

type Config struct {
	Scripts map[string]string `json:"scripts"`
}
