{{/*
Fusion CDC Engine — Helm Helper Templates
*/}}

{{- define "fusion-cdc.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "fusion-cdc.fullname" -}}
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

{{- define "fusion-cdc.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "fusion-cdc.labels" -}}
helm.sh/chart: {{ include "fusion-cdc.chart" . }}
{{ include "fusion-cdc.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: fusion-cdc
{{- end }}

{{- define "fusion-cdc.selectorLabels" -}}
app.kubernetes.io/name: {{ include "fusion-cdc.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "fusion-cdc.image" -}}
{{- $registry := .global.imageRegistry | default "" -}}
{{- $repo := .image.repository -}}
{{- $tag := .image.tag | default "latest" -}}
{{- if $registry -}}
{{- printf "%s/%s:%s" $registry $repo $tag -}}
{{- else -}}
{{- printf "%s:%s" $repo $tag -}}
{{- end -}}
{{- end }}

{{- define "fusion-cdc.controlPlaneImage" -}}
{{- include "fusion-cdc.image" (dict "global" .Values.global "image" .Values.controlPlane.image) -}}
{{- end }}

{{- define "fusion-cdc.workerImage" -}}
{{- include "fusion-cdc.image" (dict "global" .Values.global "image" .Values.cdcWorkers.image) -}}
{{- end }}

{{- define "fusion-cdc.sparkConsumerImage" -}}
{{- include "fusion-cdc.image" (dict "global" .Values.global "image" .Values.sparkConsumer.image) -}}
{{- end }}

{{- define "fusion-cdc.frontendImage" -}}
{{- include "fusion-cdc.image" (dict "global" .Values.global "image" .Values.frontend.image) -}}
{{- end }}

{{- define "fusion-cdc.transformWorkerImage" -}}
{{- include "fusion-cdc.image" (dict "global" .Values.global "image" .Values.transformWorker.image) -}}
{{- end }}

{{- define "fusion-cdc.redisHost" -}}
{{- if .Values.redis.enabled -}}
{{- printf "%s-redis-master.%s.svc.cluster.local" .Release.Name .Release.Namespace -}}
{{- else -}}
{{- .Values.redis.externalHost | default "redis" -}}
{{- end -}}
{{- end }}

{{- define "fusion-cdc.controlPlaneURL" -}}
{{- printf "http://%s-control-plane-svc.%s.svc.cluster.local:%v" .Release.Name .Release.Namespace .Values.controlPlane.service.port -}}
{{- end }}
