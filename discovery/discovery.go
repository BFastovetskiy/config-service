package discovery

import (
	"config-service/about"
	"config-service/core"
	"fmt"
	"net/http"
	"sync"
	"time"

	_ "net/http/pprof"
)

func InitDiscoveryService(app core.IApplication, duration time.Duration) *DiscoveryService {
	d := DiscoveryService{
		app:       app,
		log:       app.GetLogger(),
		db:        app.GetDatabase(),
		duration:  duration,
		discovery: make(map[string][]core.DiscoveryItem),
	}
	d.log.Info("Discovery function init")
	if !d.app.UseSSL() {
		go checkServiceAvailability(&d, duration)
	}
	return &d
}

// GetService implements IDiscoveryService.
func (d DiscoveryService) GetService(service string) []core.DiscoveryItem {
	si := d.discovery[service]
	return si
}

func (d DiscoveryService) addService(service core.DiscoveryItem) {
	service.LastAccessTime = time.Now()
	si := d.GetService(service.Name)
	si = append(si, service)
	var mu sync.Mutex
	mu.Lock()
	d.discovery[service.Name] = si
	mu.Unlock()
	d.log.Infof("Discovery register new item service [%v]", service.Name)
}

// SetService implements IDiscoveryService.
func (d DiscoveryService) SetService(service core.DiscoveryItem) {
	var notFoundService bool = false
	si := d.discovery[service.Name]
	if len(si) == 0 {
		d.addService(service)
		return
	}
	for index, item := range si {
		if (item.Address == service.Address) && (item.Port == service.Port) {
			item.LastAccessTime = time.Now()
			si[index] = item
			notFoundService = true
			break
		}
	}
	var mu sync.Mutex
	mu.Lock()
	d.discovery[service.Name] = si
	mu.Unlock()
	d.db.SetApplicationData(service.Name, service)
	if !notFoundService {
		d.addService(service)
	}
}

// DelService implements IDiscoveryService.
func (d DiscoveryService) DelService(service core.DiscoveryItem) {
	var mu sync.Mutex
	si := d.discovery[service.Name]
	if len(si) == 1 {
		mu.Lock()
		delete(d.discovery, service.Name)
		mu.Unlock()
		return
	}
	var index int
	for idx, item := range si {
		if (item.Address == service.Address) && (item.Port == service.Port) {
			index = idx
			break
		}
	}
	mu.Lock()
	si[index] = si[len(si)-1]
	d.discovery[service.Name] = si[:len(si)-1]
	mu.Unlock()
}

// GetSecret implements IDiscoveryService.
func (d DiscoveryService) GetSecret() string {
	return d.secret
}

// GetAllService implements IDiscoveryService.
func (d DiscoveryService) GetAllService() map[string][]core.DiscoveryItem {
	return d.discovery
}

func checkServiceAvailability(d *DiscoveryService, duration time.Duration) {
	d.log.Debugf("Duration for check=[%v]", duration)
	for {
		time.Sleep(duration)
		d.log.Infof("Running a service availability check")
		for k, v := range d.discovery {
			var si []core.DiscoveryItem
			for _, item := range v {
				result := ping(item, d.log)
				if result {
					item.LastAccessTime = time.Now()
					si = append(si, item)
				}
			}
			var mu sync.Mutex
			mu.Lock()
			d.discovery[k] = si
			mu.Unlock()
		}
	}
}

func ping(si core.DiscoveryItem, log core.ILogger) bool {
	client := &http.Client{}
	result := false
	uri := fmt.Sprintf("http://%s:%d%s", si.Address, si.Port, si.Health)
	req, _ := http.NewRequest(
		"GET", uri, nil,
	)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", about.Application_Title)
	resp, err := client.Do(req)

	if err != nil {
		log.Errorf("Error: [%v]", err)
		return result
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		log.Debug(si.Name + " service is available")
		result = true
	}
	return result
}
