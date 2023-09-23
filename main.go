package main

import (
	"config-service/about"
	"config-service/appconfig"
	"config-service/logger"
	"config-service/utils"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"net/http"
	_ "net/http/pprof"

	"gopkg.in/yaml.v2"
)

var (
	cfg               *appconfig.AppConfig
	interactiveLaunch bool = false

	srvPort     int
	srvSslPort  int
	clusterPort int
	onlySsl     bool
	db          string
	pem         string
	key         string
	secret      string
	frequency   int
	join        string
	dbg         bool
)

func getArgs() {
	flag.BoolFunc("i", "Interactive mode the service setup", interactiveConfiguration)
	srvPort = *flag.Int("listener-port", 3000, "Listener HTTP port")
	srvSslPort = *flag.Int("listener-ssl-port", 3001, "Listener SSL port")
	clusterPort = *flag.Int("listener-cluster-port", 4000, "Listener cluster port")
	db = *flag.String("database-name", "database.db", "Database name")
	onlySsl = *flag.Bool("only-ssl", false, "Listener only SSL")
	pem = *flag.String("certificate-pem", "public.crt", "PEM file certificate")
	key = *flag.String("certificate-key", "private.key", "Key file certificate")
	secret = *flag.String("secret", "EIaUPvdI1ONo6IQowmo6HsSRZBxUv4Hb", "Secret for connection to service")
	frequency = *flag.Int("frequency", 60, "Frequency timeout of checking the availability of services. Seconds")
	join = *flag.String("join-to", "", "Address to join on cluster")
	dbg = *flag.Bool("debug", false, "Use http listener for profilers and debug")

	flag.Parse()
}

func main() {
	getArgs()
	if interactiveLaunch {
		return
	}

	log := logger.InitLogger(about.Application_Title)
	confFile := filepath.Join(filepath.Dir(os.Args[0]), about.Application_Configuration_File)
	if utils.ExistFileOrDir(confFile) {
		var err error
		cfg, err = appconfig.LoadConfig(confFile)
		if err != nil {
			log.Error(err.Error())
			return
		}
		if cfg.Discovery.Frequency == 0 {
			cfg.Discovery.Frequency = 60
		}
	} else {
		cfg = &appconfig.AppConfig{
			Discovery: appconfig.Discovery{
				//Health:    health,
				Secret:    secret,
				Frequency: frequency,
			},
			Cluster: appconfig.Cluster{
				Join: join,
			},
			ServerPort:    srvPort,
			ServerSSLPort: srvSslPort,
			ClusterPort:   clusterPort,
			DatabaseName:  db,
			PemFile:       pem,
			KeyFile:       key,
			OnlySSL:       onlySsl,
		}
	}

	if dbg {
		go http.ListenAndServe("0.0.0.0:8888", nil)
	}

	app := Init(log, cfg)
	if app == nil {
		return
	}

	app.Run()

}

func interactiveConfiguration(interactive string) error {
	interactiveLaunch = true

	srvPort := 3000
	fmt.Printf("Listener HTTP port. Default value [%d] ", srvPort)
	fmt.Scanln(&srvPort)
	fmt.Printf("Input [%d]\n", srvPort)

	srvSslPort := 3001
	fmt.Printf("Listener SSL port. Default value [%d] ", srvSslPort)
	fmt.Scanln(&srvSslPort)
	fmt.Printf("Input [%d]\n", srvSslPort)

	clusterPort := 4000
	fmt.Printf("Listener cluster port [%d] ", clusterPort)
	fmt.Scanln(&clusterPort)
	fmt.Printf("Input [%d]\n", clusterPort)

	db := "database.db"
	fmt.Printf("Database name [%s] ", db)
	fmt.Scanln(&db)
	fmt.Printf("Input [%s]\n", db)

	pem := "public.crt"
	fmt.Printf("PEM file certificate [%s] ", pem)
	fmt.Scanln(&pem)
	fmt.Printf("Input [%s]\n", pem)

	key := "private.key"
	fmt.Printf("Key file certificate [%s] ", key)
	fmt.Scanln(&key)
	fmt.Printf("Input [%s]\n", key)

	useSSL := utils.CheckExistSSL(filepath.Dir(os.Args[0]), pem, key)
	onlySsl := false
	if useSSL {
		fmt.Println("SSL certificates exists success")
		fmt.Printf("Listener only SSL [%v] ", onlySsl)
		fmt.Scanln(&onlySsl)
		fmt.Printf("Input [%v]\n", onlySsl)
	}

	secret := "EIaUPvdI1ONo6IQowmo6HsSRZBxUv4Hb"
	fmt.Printf("Secret for connection to Discovery service [%s]", secret)
	fmt.Scanln(&secret)
	fmt.Printf("Input [%s]\n", secret)

	frequency := 60
	fmt.Printf("Frequency of checking the availability of services. Default value [1 in %d] seconds ", frequency)
	fmt.Scanln(&frequency)
	fmt.Printf("Input [%d]\n", frequency)

	join := ":4000"
	fmt.Printf("Address to join on cluster [%s] ", join)
	fmt.Scanln(&join)
	fmt.Printf("Input [%s]\n", join)

	cfg = &appconfig.AppConfig{
		Discovery: appconfig.Discovery{
			Secret:    secret,
			Frequency: frequency,
		},
		Cluster: appconfig.Cluster{
			Join: join,
		},
		ServerPort:    srvPort,
		ServerSSLPort: srvSslPort,
		ClusterPort:   clusterPort,
		DatabaseName:  db,
		PemFile:       pem,
		KeyFile:       key,
		OnlySSL:       onlySsl,
	}

	cfgApp := appconfig.ConfFile{
		AppConfig: *cfg,
	}
	log := logger.InitLogger(about.Application_Title)
	confFile := filepath.Join(filepath.Dir(os.Args[0]), about.Application_Configuration_File)
	buf, err := yaml.Marshal(cfgApp)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	err = os.WriteFile(confFile, buf, 0)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	fmt.Println("The configuration is saved to a file. Restart the application.")
	return nil
}
