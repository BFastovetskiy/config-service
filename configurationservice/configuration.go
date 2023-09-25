package configurationservice

import (
	"config-service/about"
	"config-service/core"
	"config-service/utils"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "net/http/pprof"
)

func InitConfigurationService(app core.IApplication) *ConfigurationService {
	c := &ConfigurationService{
		app:     app,
		log:     app.GetLogger(),
		confDir: filepath.Join(app.GetWorkDir(), about.Configurations_Directory),
	}
	files, err := c.loadConfigurations()
	if err != nil {
		c.files = make(map[string]core.ConfigFile)
	} else {
		c.files = files
	}

	c.log.Info("Configuration service is initialized")
	go watcher(c, 30*time.Second)

	return c
}

func (c ConfigurationService) loadConfigurations() (map[string]core.ConfigFile, error) {
	files := make(map[string]core.ConfigFile)

	dirEntries, err := os.ReadDir(c.confDir)
	if err != nil {
		c.log.Error(err.Error())
		return nil, err
	}

	for _, file := range dirEntries {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) == about.Configurations_Extension {

			fileInfo, err := os.Stat(filepath.Join(c.confDir, file.Name()))
			if err != nil {
				c.log.Error(err.Error())
			}
			cf := core.ConfigFile{}
			cf.Modify = fileInfo.ModTime()

			f, err := os.Open(filepath.Join(c.confDir, file.Name()))
			if err != nil {
				continue
			}
			cf.Body, err = io.ReadAll(f)
			if err != nil {
				continue
			}
			files[file.Name()] = cf
		}
	}

	return files, nil
}

func (c ConfigurationService) GetConfiguration(system, profile string) core.ConfigFile {
	var filename string
	if len(profile) == 0 {
		filename = system + about.Configurations_Extension
	} else {
		filename = system + "-" + profile + about.Configurations_Extension
	}

	return c.files[filename]
}

// TODO determine whether a method is needed
func (c ConfigurationService) SetConfiguration(system, profile string, data []byte) {
	var (
		filename string
		mu       sync.Mutex
	)
	if len(profile) == 0 {
		filename = system + about.Configurations_Extension
	} else {
		filename = system + "-" + profile + about.Configurations_Extension
	}

	mu.Lock()
	defer mu.Unlock()
	cf := c.files[filename]
	cf.Body = data
	c.files[filename] = cf
}

func watcher(c *ConfigurationService, duration time.Duration) {
	c.log.Info("Monitoring of service configuration files is running.")
	for {
		time.Sleep(duration)
		reloadConfigFiles(c)
	}
}

func reloadConfigFiles(c *ConfigurationService) {
	var mu sync.Mutex

	// if delete configuration
	files := make(map[string]core.ConfigFile)
	for filename, cf := range c.files {
		if utils.ExistFileOrDir(filepath.Join(c.confDir, filename)) {
			files[filename] = cf
		}
	}
	mu.Lock()
	c.files = files
	mu.Unlock()

	// if add new or modify configuration
	dirEntries, err := os.ReadDir(c.confDir)
	if err != nil {
		c.log.Error(err.Error())
		return
	}

	for _, file := range dirEntries {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if filepath.Ext(filename) == about.Configurations_Extension {
			fileInfo, err := os.Stat(filepath.Join(c.confDir, filename))
			if err != nil {
				c.log.Error(err.Error())
			}

			cf, ok := c.files[filename]
			if (ok) && (cf.Modify == fileInfo.ModTime()) {
				continue
			}

			cf = core.ConfigFile{}
			cf.Modify = fileInfo.ModTime()
			f, err := os.Open(filepath.Join(c.confDir, filename))
			if err != nil {
				c.log.Error(err.Error())
				continue
			}
			cf.Body, err = io.ReadAll(f)
			if err != nil {
				c.log.Error(err.Error())
				continue
			}

			mu.Lock()
			c.files[filename] = cf
			mu.Unlock()
			c.log.Infof("Reload new or modify configuration files [%s]", filename)
		}
	}

}
