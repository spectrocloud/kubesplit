package helm

import (
	"fmt"
	"strings"

	unstructured "github.com/kairos-io/kairos/sdk/unstructured"
	"gopkg.in/yaml.v3"
)

type rules struct {
	Condition   func(content []byte, data map[string]interface{}) bool
	Replacement func(content []byte) []byte
}

// Mutator mutates content with helm rules.
func Mutator(content []byte, data map[string]interface{}) []byte {
	for _, rule := range helmRules {
		if rule.Condition(content, data) {
			content = rule.Replacement(content)
		}
	}

	return content
}

var helmRules = []rules{
	{
		Condition:   hasKey("metadata.namespace"),
		Replacement: replaceKeyValue(`.metadata.namespace`, `"{{.Release.Namespace}}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "Service"),
		Replacement: replaceKeyValue(`.metadata.name`, `"{{ include \"helm-chart.fullname\" . }}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "Service"),
		Replacement: replaceKeyValue(`.spec.selector`, `"{{- include \"helm-chart.selectorLabels\" . | nindent 6 }}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "Deployment"),
		Replacement: replaceKeyValue(`.spec.template.spec.containers[1].image`, `"{{ .Values.image.repository | default \"\" }}:{{ .Values.image.tag | default .Chart.AppVersion }}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "Deployment"),
		Replacement: replaceKeyValue(`.spec.template.spec.serviceAccountName`, `"{{ include \"helm-chart.serviceAccountName\" . }}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "ServiceAccount"),
		Replacement: replaceKeyValue(`.metadata.name`, `"{{ include \"helm-chart.serviceAccountName\" . }}"`),
	},
	{
		Condition:   keyContainsValue(".metadata.name", "metrics-service"),
		Replacement: replaceKeyValue(`.metadata.name`, `"{{ include \"helm-chart.fullname\" . }}-metrics-service"`),
	},
	{
		Condition:   keyContainsValue(".metadata.name", "webhook-service"),
		Replacement: replaceKeyValue(`.metadata.name`, `"{{ include \"helm-chart.fullname\" . }}-webhook-service"`),
	},
	{
		Condition:   keyContainsValue(".metadata.name", "serving-cert"),
		Replacement: replaceKeyValue(`.metadata.name`, `"{{ include \"helm-chart.fullname\" . }}-serving-cert"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "ClusterRoleBinding"),
		Replacement: replaceKeyValue(`.subjects[0].namespace`, `"{{.Release.Namespace}}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "RoleBinding"),
		Replacement: replaceKeyValue(`.subjects[0].namespace`, `"{{.Release.Namespace}}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "ClusterRoleBinding"),
		Replacement: replaceKeyValue(`.subjects[0].name`, `"{{ include \"helm-chart.serviceAccountName\" . }}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "RoleBinding"),
		Replacement: replaceKeyValue(`.subjects[0].name`, `"{{ include \"helm-chart.serviceAccountName\" . }}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "Deployment"),
		Replacement: replaceKeyValue(`.spec.labels`, `"{{- include \"helm-chart.labels\" . | nindent 8 }}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "Deployment"),
		Replacement: replaceKeyValue(`.spec.selector.matchLabels`, `"{{- include \"helm-chart.selectorLabels\" . | nindent 10 }}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "Deployment"),
		Replacement: replaceKeyValue(`.spec.template.metadata.labels`, `"{{- include \"helm-chart.selectorLabels\" . | nindent 14 }}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "Deployment"),
		Replacement: replaceKeyValue(`.spec.template.metadata.annotations`, `"{{- range keys .Values.podAnnotations }}{{ . | quote }}: {{ get $.Values.podAnnotations . | quote}}{{- end }}"`),
	},
	{
		Condition:   hasKeyWithValue(".kind", "Deployment"),
		Replacement: replaceKeyValue(`.imagePullSecrets`, `"{{ toYaml .Values.imagePullSecrets | nindent 14 }}"`),
	},
}

func hasKey(key string) func(content []byte, data map[string]interface{}) bool {
	return func(content []byte, data map[string]interface{}) bool {
		result, _ := unstructured.YAMLHasKey(key, content)

		return result
	}
}

func hasKeyWithValue(key, value string) func(content []byte, data map[string]interface{}) bool {
	return func(content []byte, data map[string]interface{}) bool {
		v, err := unstructured.LookupString(key, data)
		return err == nil && v == value
	}
}

func keyContainsValue(key, value string) func(content []byte, data map[string]interface{}) bool {
	return func(content []byte, data map[string]interface{}) bool {
		v, err := unstructured.LookupString(key, data)
		return err == nil && strings.Contains(v, value)
	}
}

func replaceKeyValue(key, value string) func(content []byte) []byte {
	return func(content []byte) []byte {
		data := map[string]interface{}{}
		err := yaml.Unmarshal(content, &data)
		if err != nil {
			panic(err)
		}
		v, err := unstructured.ReplaceValue(fmt.Sprintf("%s=%s", key, value), data)
		if err != nil {
			panic(err.Error())
		}
		return []byte(v)
	}
}
