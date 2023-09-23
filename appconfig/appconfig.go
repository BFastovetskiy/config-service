package appconfig

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig(path string) (*AppConfig, error) {
	fcfg := ConfFile{}
	osFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer osFile.Close()

	file, err := io.ReadAll(osFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &fcfg)
	if err != nil {
		return nil, err
	}

	cfg := cloneCfg(fcfg.AppConfig)

	return cfg, nil
}

func cloneCfg(oldCfg AppConfig) *AppConfig {
	d := Discovery{}
	//d.Health = oldCfg.Discovery.Health
	d.Secret = oldCfg.Discovery.Secret
	c := Cluster{}
	c.Join = oldCfg.Cluster.Join
	cfg := &AppConfig{
		Discovery:     d,
		Cluster:       c,
		ServerPort:    oldCfg.ServerPort,
		ServerSSLPort: oldCfg.ServerSSLPort,
		ClusterPort:   oldCfg.ClusterPort,
		OnlySSL:       oldCfg.OnlySSL,
		DatabaseName:  oldCfg.DatabaseName,
		PemFile:       oldCfg.PemFile,
		KeyFile:       oldCfg.KeyFile,
	}
	return cfg
}
