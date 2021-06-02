package caps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type Metric struct {
	Ctx      context.Context
	Metricch chan Output
}

var ssUrl = os.Getenv("SS_APIMONITOR_URL")

func (m Metric) MetricSender() {

	for {
		select {
		case metric := <-m.Metricch:
			bs, err := json.Marshal(metric)
			if err != nil {
				logrus.Errorf("Marshaling metric:%s", err)
			}
			SendMetrics(bs)
			logrus.Infof("SENT METRIC TO APIMONITOR:%v SRCHOST:%s DESTHOST:%s DESTIP:%s ", metric, metric.SrcHost, metric.DestName, metric.DestIp)
		}
	}

}

func SendMetrics(m []byte) {
	url := ssUrl
	method := "POST"
	client := &http.Client{}
	metric := bytes.NewReader(m)
	req, err := http.NewRequest(method, url, metric)
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
	logrus.Infof("SendMetrics Body Response:%s", string(body))
}
