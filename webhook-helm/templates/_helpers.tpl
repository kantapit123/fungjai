{{ define "webhook.image" -}}
{{ printf "%s:%s" .Values.image.repository .Values.image.tag }}
{{- end }}

{{ define "endpoint.env" -}}
{{ printf "%s:%s" .Values.image.env.endpoint .Values.image.env.endpointPort }}
{{- end }}