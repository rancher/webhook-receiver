package tmpl

import (
	"bytes"
	"text/template"
)

var tpl *template.Template

func init() {
	var err error
	tpl, err = template.New("").Option("missingkey=zero").Parse(NotificationTmpl)
	if err != nil {
		panic(err)
	}
}

const (
	NotificationTmpl = `
{{- if eq .Status "resolved" -}}
[Resolved]
{{- end -}}
{{- if eq .CommonLabels.alert_type "event" }}
{{ .CommonLabels.event_type}} event of {{.GroupLabels.resource_kind}} occurred
{{- else if eq .CommonLabels.alert_type "systemService" }}
The system component {{ .GroupLabels.component_name}} is not running
{{- else if eq .CommonLabels.alert_type "nodeHealthy" }}
The kubelet on the node {{ .GroupLabels.node_name}} is not healthy
{{- else if eq .CommonLabels.alert_type "nodeCPU" }}
The CPU usage on the node {{ .GroupLabels.node_name}} is over {{ .CommonLabels.cpu_threshold}}%
{{- else if eq .CommonLabels.alert_type "nodeMemory" }}
The memory usage on the node {{ .GroupLabels.node_name}} is over {{ .CommonLabels.mem_threshold}}%
{{- else if eq .CommonLabels.alert_type "podNotScheduled" }}
The Pod {{ if .GroupLabels.namespace}}{{.GroupLabels.namespace}}:{{end}}{{.GroupLabels.pod_name}} is not scheduled
{{- else if eq .CommonLabels.alert_type "podNotRunning" }}
The Pod {{ if .GroupLabels.namespace}}{{.GroupLabels.namespace}}:{{end}}{{.GroupLabels.pod_name}} is not running
{{- else if eq .CommonLabels.alert_type "podRestarts" }}
The Pod {{ if .GroupLabels.namespace}}{{.GroupLabels.namespace}}:{{end}}{{.GroupLabels.pod_name}} restarts {{ .CommonLabels.restart_times}} times in {{ .CommonLabels.restart_interval}} sec
{{- else if eq .CommonLabels.alert_type "workload" }}
The workload {{ if .GroupLabels.workload_namespace}}{{.GroupLabels.workload_namespace}}:{{end}}{{.GroupLabels.workload_name}} has available replicas less than {{ .CommonLabels.available_percentage}}%
{{- else if eq .CommonLabels.alert_type "metric" }}
The metric {{ .CommonLabels.alert_name}} crossed the threshold
{{ end -}}

{{- if eq .Status "resolved" -}}
{{ range .Alerts.Resolved }}
{{ template "__text_single" . }}
{{ end -}}
{{- else}}
{{ range .Alerts.Firing }}
{{ template "__text_single" . }}
{{ end -}}
{{ end -}}

{{- define "__text_single" -}}
Alert Name: {{ .Labels.alert_name}}
Severity: {{ .Labels.severity}}
Cluster Name: {{.Labels.cluster_name}}
{{- if eq .Labels.alert_type "event" }}
{{- if .Labels.workload_name }}
Workload Name: {{.Labels.workload_name}}{{ end }}
Target: {{ if .Labels.target_namespace -}}{{.Labels.target_namespace}}:{{ end -}}{{.Labels.target_name}}
Count: {{ .Labels.event_count}}
Event Message: {{ .Labels.event_message}}
First Seen: {{ .Labels.event_firstseen}}
Last Seen: {{ .Labels.event_lastseen}}
{{- else if eq .Labels.alert_type "nodeCPU" }}
Used CPU: {{ .Labels.used_cpu}} m
Total CPU: {{ .Labels.total_cpu}} m
{{- else if eq .Labels.alert_type "nodeMemory" }}
Used Memory: {{ .Labels.used_mem}}
Total Memory: {{ .Labels.total_mem}}
{{- else if eq .Labels.alert_type "podRestarts" }}
Project Name: {{ .Labels.project_name}}
Namespace: {{ .Labels.namespace}}
{{- if .Labels.workload_name }}
Workload Name: {{.Labels.workload_name}}
{{ end -}}
Container Name: {{ .Labels.container_name}}
{{- else if eq .Labels.alert_type "podNotRunning" }}
Project Name: {{ .Labels.project_name}}
Namespace: {{ .Labels.namespace}}
{{- if .Labels.workload_name }}
Workload Name: {{.Labels.workload_name}}
{{ end -}}
Container Name: {{ .Labels.container_name}}
{{- else if eq .Labels.alert_type "podNotScheduled" }}
Project Name: {{ .Labels.project_name}}
Namespace: {{ .Labels.namespace}}
Pod Name: {{ .Labels.pod_name}}
{{- if .Labels.workload_name }}
Workload Name: {{.Labels.workload_name}}
{{ end -}}
{{- else if eq .Labels.alert_type "workload" }}
Project Name: {{ .Labels.project_name}}
Available Replicas: {{ .Labels.available_replicas}}
Desired Replicas: {{ .Labels.desired_replicas}}
{{- else if eq .Labels.alert_type "metric" }}
{{- if .Labels.namespace }}
Namespace: {{ .Labels.namespace}}{{ end }}
{{- if .Labels.project_name }}
Project Name: {{ .Labels.project_name}}{{ end }}
{{- if .Labels.pod_name }}
Pod Name: {{ .Labels.pod_name}}{{ else if .Labels.pod -}}Pod Name: {{ .Labels.pod}}{{ end }}
Expression: {{ .Labels.expression}}
Description: Threshold Crossed: datapoint value {{ .Annotations.current_value}} was {{ .Labels.comparison}} to the threshold ({{ .Labels.threshold_value}}) for ({{ .Labels.duration}})
{{ end -}}
{{- if .Labels.logs }}
Logs: {{ .Labels.logs}}
{{ end -}}
{{ end -}}
`
)

func ExecuteTextString(data interface{}) (string, error) {
	buf := bytes.Buffer{}
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
