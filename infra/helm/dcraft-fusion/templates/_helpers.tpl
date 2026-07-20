{{/*
DCraft Fusion — Helm helpers (Airbyte-style)
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
dcraft.edition: {{ .Values.global.edition | default "community" | quote }}
{{- with .Values.global.labels }}
{{- toYaml . }}
{{- end }}
{{- end }}

{{- define "dcraft-fusion.selectorLabels" -}}
app.kubernetes.io/name: {{ include "dcraft-fusion.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Image helper: registry + repository + tag or digest.
Usage: include "dcraft-fusion.image" (dict "global" .Values.global "image" .Values.controlPlaneKernel.image)
*/}}
{{- define "dcraft-fusion.image" -}}
{{- $registry := .global.imageRegistry | default "" -}}
{{- $repo := .image.repository -}}
{{- if and .image.digest (ne .image.digest "") -}}
{{- if $registry -}}
{{- printf "%s/%s@%s" (trimSuffix "/" $registry) $repo .image.digest -}}
{{- else -}}
{{- printf "%s@%s" $repo .image.digest -}}
{{- end -}}
{{- else -}}
{{- $tag := .image.tag | default "latest" -}}
{{- if $registry -}}
{{- printf "%s/%s:%s" (trimSuffix "/" $registry) $repo $tag -}}
{{- else -}}
{{- printf "%s:%s" $repo $tag -}}
{{- end -}}
{{- end -}}
{{- end }}

{{- define "dcraft-fusion.secretName" -}}
{{- if .Values.global.secrets.existingSecret -}}
{{- .Values.global.secrets.existingSecret -}}
{{- else -}}
{{- printf "%s-secrets" (include "dcraft-fusion.fullname" .) -}}
{{- end -}}
{{- end }}

{{- define "dcraft-fusion.redisAddr" -}}
{{- .Values.externalRedis.addr | default "" -}}
{{- end }}

{{- define "dcraft-fusion.controlPlaneKernel.serviceAccountName" -}}
{{- if .Values.controlPlaneKernel.serviceAccount.create -}}
{{- default (printf "%s-control-plane-kernel" (include "dcraft-fusion.fullname" .)) .Values.controlPlaneKernel.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.controlPlaneKernel.serviceAccount.name -}}
{{- end -}}
{{- end }}

{{- define "dcraft-fusion.web.serviceAccountName" -}}
{{- if .Values.web.serviceAccount.create -}}
{{- default (printf "%s-web" (include "dcraft-fusion.fullname" .)) .Values.web.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.web.serviceAccount.name -}}
{{- end -}}
{{- end }}

{{/*
Render env from map (.env) and list (.extraEnv).
Usage: include "dcraft-fusion.renderEnv" (dict "env" .Values.web.env "extraEnv" .Values.web.extraEnv)
*/}}
{{- define "dcraft-fusion.renderEnv" -}}
{{- range $key, $value := .env }}
- name: {{ $key }}
  value: {{ $value | quote }}
{{- end }}
{{- with .extraEnv }}
{{- toYaml . }}
{{- end }}
{{- end }}
