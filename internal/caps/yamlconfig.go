package caps

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	HostName         string `yaml:"host_name"`
	Provider         string `yaml:"provider"`
	Region           string `yaml:"region"`
	HostIp           string `yaml:"host_ip"`
	ReflectorUdpPort int    `yaml:"reflector_udp_port"`
	TcpPort          int    `yaml:"tcp_port"`
	UrlApimonitor    string `yaml:"url_apimonitor"`
}

// config setup for on server

func SetupConfig(path string) (Config, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var config Config
	fmt.Printf("%s", yamlFile)
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil

}
