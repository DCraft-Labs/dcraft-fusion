{{/*
Fusion CDC — Helm helpers (image-only chart)
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
dcraft.edition: {{ .Values.global.edition | default "community" | quote }}
{{- with .Values.global.labels }}
{{- toYaml . }}
{{- end }}
{{- end }}

{{- define "fusion-cdc.selectorLabels" -}}
app.kubernetes.io/name: {{ include "fusion-cdc.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Image helper: registry + repository + tag or digest.
v1.3.1 Fix 3: the tag is rendered with ``tpl`` against the chart context
so values.yaml can default ``image.tag`` to ``{{ .Chart.AppVersion }}``
and have the chart default track the app version without a stale
hardcoded number. Plain tags (e.g. ``latest``) pass through unchanged.
*/}}
{{- define "fusion-cdc.image" -}}
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
{{- if .ctx -}}
{{- $tag = tpl $tag .ctx -}}
{{- end -}}
{{- if $registry -}}
{{- printf "%s/%s:%s" (trimSuffix "/" $registry) $repo $tag -}}
{{- else -}}
{{- printf "%s:%s" $repo $tag -}}
{{- end -}}
{{- end -}}
{{- end }}

{{- define "fusion-cdc.secretName" -}}
{{- if .Values.global.secrets.existingSecret -}}
{{- .Values.global.secrets.existingSecret -}}
{{- else -}}
{{- printf "%s-secrets" (include "fusion-cdc.fullname" .) -}}
{{- end -}}
{{- end }}

{{- define "fusion-cdc.redisUrl" -}}
{{- if .Values.global.secrets.redisUrl -}}
{{- .Values.global.secrets.redisUrl -}}
{{- else if .Values.externalRedis.url -}}
{{- .Values.externalRedis.url -}}
{{- else -}}
{{- "" -}}
{{- end -}}
{{- end }}

{{- define "fusion-cdc.controlPlaneURL" -}}
{{- printf "http://%s-control-plane:%v" (include "fusion-cdc.fullname" .) .Values.controlPlane.service.port -}}
{{- end }}
{{- define "fusion-cdc.controlPlane.serviceAccountName" -}}
{{- if .Values.controlPlane.serviceAccount.create -}}
{{- default (printf "%s-control-plane" (include "fusion-cdc.fullname" .)) .Values.controlPlane.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.controlPlane.serviceAccount.name -}}
{{- end -}}
{{- end }}

{{- define "fusion-cdc.cdcWorkers.serviceAccountName" -}}
{{- if .Values.cdcWorkers.serviceAccount.create -}}
{{- default (printf "%s-cdc-worker" (include "fusion-cdc.fullname" .)) .Values.cdcWorkers.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.cdcWorkers.serviceAccount.name -}}
{{- end -}}
{{- end }}

{{- define "fusion-cdc.frontend.serviceAccountName" -}}
{{- if .Values.frontend.serviceAccount.create -}}
{{- default (printf "%s-frontend" (include "fusion-cdc.fullname" .)) .Values.frontend.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.frontend.serviceAccount.name -}}
{{- end -}}
{{- end }}

{{- define "fusion-cdc.workerServiceAccountName" -}}
{{- include "fusion-cdc.cdcWorkers.serviceAccountName" . -}}
{{- end }}

{{- define "fusion-cdc.renderEnv" -}}
{{- range $key, $value := .env }}
- name: {{ $key }}
  value: {{ $value | quote }}
{{- end }}
{{- with .extraEnv }}
{{- toYaml . }}
{{- end }}
{{- end }}
