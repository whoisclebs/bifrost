package discovery

import (
	"bifrost/proxy/internal/config"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Applications struct {
	XMLName xml.Name      `xml:"applications"`
	Apps    []Application `xml:"application"`
}

type Application struct {
	Name      string     `xml:"name"`
	Instances []Instance `xml:"instance"`
}

type Instance struct {
	HostName string `xml:"hostName"`
	IPAddr   string `xml:"ipAddr"`
	Port     int    `xml:"port"`
	Status   string `xml:"status"`
}

func GetInstances(cfg config.ProxyConfig, serviceId string) ([]Instance, error) {
	if len(cfg.Eureka.Addresses) == 0 {
		return nil, fmt.Errorf("no eureka addresses configured")
	}
	url := fmt.Sprintf("%s/apps/%s", cfg.Eureka.Addresses[0], strings.ToUpper(serviceId))
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get instances, status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var app Application
	err = xml.Unmarshal(body, &app)
	if err != nil {
		return nil, err
	}
	return app.Instances, nil
}
