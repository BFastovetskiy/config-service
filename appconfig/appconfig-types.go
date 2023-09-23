package appconfig

type ConfFile struct {
	AppConfig AppConfig `yaml:"application"`
}

type AppConfig struct {
	ServerPort    int       `yaml:"srvPort"`
	ServerSSLPort int       `yaml:"srvSslPort"`
	ClusterPort   int       `yaml:"clusterProt"`
	DatabaseName  string    `yaml:"db"`
	PemFile       string    `yaml:"pem"`
	KeyFile       string    `yaml:"key"`
	OnlySSL       bool      `yaml:"onlySsl"`
	Discovery     Discovery `yaml:"discovery"`
	Cluster       Cluster   `yaml:"cluster"`
}

type Discovery struct {
	//Health    string `yaml:"health"`
	Secret    string `yaml:"secret"`
	Frequency int    `yaml:"frequency"`
}

type Cluster struct {
	Join string `yaml:"join"`
}
