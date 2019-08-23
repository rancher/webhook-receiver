package main

import (
	"encoding/json"
	"github.com/prometheus/alertmanager/template"
	"log"
)

const data = `{"receiver":"c-bktpt:kube-components-alert","status":"firing","alerts":[{"status":"firing","labels":{"alert_name":"Controller Manager is unavailable","alert_type":"systemService","cluster_name":"test (ID: c-bktpt)","component_name":"controller-manager","group_id":"c-bktpt:kube-components-alert","logs":"Get http://127.0.0.1:10252/healthz: dial tcp 127.0.0.1:10252: connect: connection refused","rule_id":"c-bktpt:kube-components-alert_controllermanager-system-service","severity":"critical"},"annotations":{},"startsAt":"2019-08-22T05:21:30.135535529Z","endsAt":"0001-01-01T00:00:00Z","generatorURL":""}],"groupLabels":{"component_name":"controller-manager","rule_id":"c-bktpt:kube-components-alert_controllermanager-system-service"},"commonLabels":{"alert_name":"Controller Manager is unavailable","alert_type":"systemService","cluster_name":"test (ID: c-bktpt)","component_name":"controller-manager","group_id":"c-bktpt:kube-components-alert","logs":"Get http://127.0.0.1:10252/healthz: dial tcp 127.0.0.1:10252: connect: connection refused","rule_id":"c-bktpt:kube-components-alert_controllermanager-system-service","severity":"critical"},"commonAnnotations":{},"externalURL":"http://alertmanager-cluster-alerting-0:9093","version":"4","groupKey":"{}/{group_id=\"c-bktpt:kube-components-alert\"}/{rule_id=\"c-bktpt:kube-components-alert_controllermanager-system-service\"}:{component_name=\"controller-manager\", rule_id=\"c-bktpt:kube-components-alert_controllermanager-system-service\"}"}`

func main() {
	alerts := template.Data{}
	if err := json.Unmarshal([]byte(data), &alerts); err != nil {
		log.Fatal(err)
	}

	log.Printf("%v", alerts)
}
