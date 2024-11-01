{{/*
Expand the name of the chart.
*/}}
{{- define "llmos-gpu-stack.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "llmos-gpu-stack.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "llmos-gpu-stack.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "llmos-gpu-stack.labels" -}}
helm.sh/chart: {{ include "llmos-gpu-stack.chart" . }}
{{ include "llmos-gpu-stack.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}


{{/*
Selector labels
*/}}
{{- define "llmos-gpu-stack.selectorLabels" -}}
app.kubernetes.io/name: {{ include "llmos-gpu-stack.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "llmos-gpu-stack.serviceAccountName" -}}
{{- if .Values.gpuStack.serviceAccount.create }}
{{- default (include "llmos-gpu-stack.fullname" .) .Values.gpuStack.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.gpuStack.serviceAccount.name }}
{{- end }}
{{- end }}


{{/*
Device manager templates
*/}}
{{- define "device-manager.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 48 | trimSuffix "-" }}-device-manager
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- printf "%s" .Release.Name | trunc 48 | trimSuffix "-" }}-device-manager
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 48 | trimSuffix "-" }}-device-manager
{{- end }}
{{- end }}
{{- end }}

{{/*
Device manager chart
*/}}
{{- define "device-manager.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Device manager labels
*/}}
{{- define "device-manager.labels" -}}
helm.sh/chart: {{ include "llmos-gpu-stack.chart" . }}
{{ include "llmos-gpu-stack.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Device manager selectors
*/}}
{{- define "device-manager.selectorLabels" -}}
app.kubernetes.io/name: {{ include "llmos-gpu-stack.name" . }}
app.kubernetes.io/instance: "device-manager"
{{- end }}

