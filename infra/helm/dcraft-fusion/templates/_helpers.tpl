{{/*
DCraft Fusion — Helm helpers
*/}}

{{- define "dcraft-fusion.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "dcraft-fusion.fullname" -}}
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

{{- define "dcraft-fusion.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "dcraft-fusion.labels" -}}
helm.sh/chart: {{ include "dcraft-fusion.chart" . }}
{{ include "dcraft-fusion.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: dcraft-fusion
{{- end }}

{{- define "dcraft-fusion.selectorLabels" -}}
app.kubernetes.io/name: {{ include "dcraft-fusion.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "dcraft-fusion.image" -}}
{{- $registry := .global.imageRegistry | default "" -}}
{{- $repo := .image.repository -}}
{{- $tag := .image.tag | default "latest" -}}
{{- if $registry -}}
{{- printf "%s/%s:%s" $registry $repo $tag -}}
{{- else -}}
{{- printf "%s:%s" $repo $tag -}}
{{- end -}}
{{- end }}

{{- define "dcraft-fusion.secretName" -}}
{{- if .Values.secrets.existingSecret -}}
{{- .Values.secrets.existingSecret -}}
{{- else -}}
{{- printf "%s-secrets" (include "dcraft-fusion.fullname" .) -}}
{{- end -}}
{{- end }}

{{- define "dcraft-fusion.redisAddr" -}}
{{- if .Values.redis.enabled -}}
{{- printf "%s-redis-master:6379" .Release.Name -}}
{{- else if .Values.redis.externalAddr -}}
{{- .Values.redis.externalAddr -}}
{{- else -}}
redis:6379
{{- end -}}
{{- end }}

{{- define "dcraft-fusion.postgresDsn" -}}
{{- if .Values.secrets.postgresDsn -}}
{{- .Values.secrets.postgresDsn -}}
{{- else if .Values.postgresql.enabled -}}
{{- $user := .Values.postgresql.auth.username | default "fusion" -}}
{{- $pass := .Values.postgresql.auth.password | default "fusion" -}}
{{- $db := .Values.postgresql.auth.database | default "fusion" -}}
{{- printf "postgres://%s:%s@%s-postgresql:5432/%s?sslmode=disable" $user $pass .Release.Name $db -}}
{{- else -}}
{{- "" -}}
{{- end -}}
{{- end }}
