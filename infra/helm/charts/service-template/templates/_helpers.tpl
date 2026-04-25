{{/*
Expand the name of the chart.
*/}}
{{- define "complai-service.name" -}}
{{- default .Chart.Name .Values.service.name | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "complai-service.fullname" -}}
{{- if .Values.service.name }}
{{- .Values.service.name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{/*
Chart label value.
*/}}
{{- define "complai-service.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels.
*/}}
{{- define "complai-service.labels" -}}
helm.sh/chart: {{ include "complai-service.chart" . }}
{{ include "complai-service.selectorLabels" . }}
app.kubernetes.io/version: {{ .Values.image.tag | default .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: complai
{{- with .Values.labels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Selector labels.
*/}}
{{- define "complai-service.selectorLabels" -}}
app.kubernetes.io/name: {{ include "complai-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Service account name.
*/}}
{{- define "complai-service.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "complai-service.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
ConfigMap name.
*/}}
{{- define "complai-service.configMapName" -}}
{{- printf "%s-config" (include "complai-service.fullname" .) }}
{{- end }}
