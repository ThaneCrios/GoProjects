{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "delivery.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "delivery.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "delivery.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "delivery.labels" -}}
helm.sh/chart: {{ include "delivery.chart" . }}
app.kubernetes.io/name: {{ include "delivery.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Common selector
*/}}
{{- define "delivery.selector" -}}
app.kubernetes.io/name: {{ include "delivery.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/* deliverycore */}}

{{- define "deliverycore.fullname" -}}
{{ include "delivery.fullname" . }}
{{- end -}}

{{- define "deliverycore.labels" -}}
{{ include "delivery.labels" . }}
app.kubernetes.io/component: deliverycore
{{- end -}}

{{- define "deliverycore.selector" -}}
{{ include "delivery.selector" . }}
app.kubernetes.io/component: deliverycore
{{- end -}}