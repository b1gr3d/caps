package caps

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Apimonitor struct {
	Ctx         context.Context
	LatencyChan chan Output
}

type Target struct {
	Destinationname string `json:"destination_name"`
	Protocol        string `json:"protocol"`
	Networkoneip        string `json:"networkone_ip"`
	Networkoneport      int    `json:"networkone_port"`
	Networktwoip      string `json:"networktwo_ip"`
	Networktwoport    int    `json:"networktwo_port"`
}

type Targets struct {
	targets []Target
	lock *sync.Mutex
}

func (a *Apimonitor) RunApimonitor() {

	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ticker.C:
			a.GetTargets()
		case <-a.Ctx.Done():
			return
		}

	}

}

func Register() {
	url := config.UrlApimonitor + "/registry"
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf(`{
    "host_name": "%s", 
    "provider": "%s", 
    "region": "%s",
    "host_ip": "%s",
    "udp_port": "%d",
	"tcp_port": "%d"
}`, config.HostName, config.Provider, config.Region, config.HostIp, config.ReflectorUdpPort, config.TcpPort))
	client := &http.Client{}
	logrus.Infof("%s", payload)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer network metrics")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	logrus.Infof(string(body))
}

func (a Apimonitor) GetTargets() {
	url := ssUrl
	method := "GET"
	payload := strings.NewReader(fmt.Sprintf(`{"host_name": "%s"}`, config.HostName))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		//return
	}
	req.Header.Add("Authorization", "Bearer network metrics")
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var target []Target
	err = json.Unmarshal(body, &target)
	if err != nil {
		logrus.Error("Error unmarshaling targets", err)
	}

	targets.lock.Lock()
	targets.targets = target
	targets.lock.Unlock()
}
