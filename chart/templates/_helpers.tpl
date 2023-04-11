{{ define "webhook.image" -}}
{{ printf "%s:%s" .Values.image.repository .Values.image.tag }}
{{- end }}