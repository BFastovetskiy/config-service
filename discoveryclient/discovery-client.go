package discoveryclient

import (
	"config-service/about"
	"config-service/logger"
	"fmt"
	"io"
	"net/http"
)

type DiscoveryClient struct {
	title     string
	host      string
	port      int
	context   string
	secretKey string
	log       *logger.Logger
}

func InitDiscoveryClient(host string, port int, secret string, title string, log *logger.Logger) *DiscoveryClient {
	return &DiscoveryClient{
		title:     title,
		port:      port,
		host:      host,
		context:   about.API_Context,
		secretKey: secret,
		log:       log,
	}
}

func (cd DiscoveryClient) GetConfiguration(serviceName string, profile string) []byte {
	client := &http.Client{}

	uri := fmt.Sprintf("http://%s:%d/%s/%s/%s", cd.host, cd.port, cd.context, serviceName, profile)

	req, _ := http.NewRequest(
		"GET", uri, nil,
	)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", cd.title)
	req.Header.Add("Authorization", cd.secretKey)

	resp, err := client.Do(req)

	if err != nil {
		cd.log.Errorf("Error: %v", err)
		return nil
	}

	defer resp.Body.Close()

	conf, err := io.ReadAll(resp.Body)
	if err != nil {
		cd.log.Errorf("Error: %v", err)
		return nil
	}

	return conf
}
