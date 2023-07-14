{{ define "webhook.image" -}}
{{ printf "%s:%s" .Values.webhook.image.repository .Values.webhook.image.tag }}
{{- end }}
{{ define "metabase-clickhouse.image" -}}
{{ printf "%s:%s" .Values.metabase.image.repository .Values.metabase.image.tag }}
{{- end }}
{{ define "consumer.image" -}}
{{ printf "%s:%s" .Values.consumer.image.repository .Values.consumer.image.tag }}
{{- end }}