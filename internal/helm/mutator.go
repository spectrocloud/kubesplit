package helm

import (
	"fmt"
	"strings"

	yamled "github.com/vmware-labs/go-yaml-edit"
	yptr "github.com/vmware-labs/yaml-jsonpointer"
	"golang.org/x/text/transform"
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
		Condition:   hasKey("/metadata/namespace"),
		Replacement: replaceKeyValue(`/metadata/namespace`, "{{.Release.Namespace}}"),
	},
	{
		Condition:   hasKeyWithValue("/kind", "Deployment"),
		Replacement: replaceKeyValue(`/metadata/name`, `{{ include "helm-chart.fullname" . }}`),
	},
	{
		Condition:   hasKeyWithValue("/kind", "Deployment"),
		Replacement: replaceKeyValue(`/spec/template/spec/containers/1/image`, `{{ .Values.image.repository | default "" }}:{{ .Values.image.tag | default .Chart.AppVersion }}`),
	},
	{
		Condition:   hasKeyWithValue("/kind", "Deployment"),
		Replacement: replaceKeyValue(`/spec/template/spec/serviceAccountName`, `{{ include "helm-chart.serviceAccountName" . }}`),
	},
	{
		Condition:   hasKeyWithValue("/kind", "ServiceAccount"),
		Replacement: replaceKeyValue(`/metadata/name`, `{{ include "helm-chart.serviceAccountName" . }}`),
	},
	{
		Condition:   keyContainsValue("/metadata/name", "metrics-service"),
		Replacement: replaceKeyValue(`/metadata/name`, `{{ include "helm-chart.fullname" . }}-metrics-service`),
	},
	{
		Condition:   keyContainsValue("/metadata/name", "webhook-service"),
		Replacement: replaceKeyValue(`/metadata/name`, `{{ include "helm-chart.fullname" . }}-webhook-service`),
	},
	{
		Condition:   keyContainsValue("/metadata/name", "serving-cert"),
		Replacement: replaceKeyValue(`/metadata/name`, `{{ include "helm-chart.fullname" . }}-serving-cert`),
	},
	{
		Condition:   hasKeyWithValue("/kind", "ClusterRoleBinding"),
		Replacement: replaceKeyValue(`/subjects/0/namespace`, `{{.Release.Namespace}}`),
	},
	{
		Condition:   hasKeyWithValue("/kind", "RoleBinding"),
		Replacement: replaceKeyValue(`/subjects/0/namespace`, `{{.Release.Namespace}}`),
	},
	{
		Condition:   hasKeyWithValue("/kind", "ClusterRoleBinding"),
		Replacement: replaceKeyValue(`/subjects/0/name`, `{{ include "helm-chart.serviceAccountName" . }}`),
	},
	{
		Condition:   hasKeyWithValue("/kind", "RoleBinding"),
		Replacement: replaceKeyValue(`/subjects/0/name`, `{{ include "helm-chart.serviceAccountName" . }}`),
	},
	// {
	// 	Condition:   hasKeyWithValue("/kind", "Deployment"),
	// 	Replacement: replaceKeyValue(`/spec/selector/matchLabels`, `{{- include "helm-chart.selectorLabels" . | nindent 6 }}`),
	// },
	// {
	// 	Condition:   hasKeyWithValue("/kind", "Deployment"),
	// 	Replacement: replaceKeyValue(`/spec/template/metadata/labels`, `{{- include "helm-chart.selectorLabels" . | nindent 10 }}`),
	// },
}

func hasKey(key string) func(content []byte, data map[string]interface{}) bool {
	return func(content []byte, data map[string]interface{}) bool {
		var root yaml.Node
		if err := yaml.Unmarshal(content, &root); err != nil {
			return false
		}
		_, err := yptr.Find(&root, key)
		if err != nil {
			return false
		}

		return true
	}
}

func hasKeyWithValue(key, value string) func(content []byte, data map[string]interface{}) bool {
	return func(content []byte, data map[string]interface{}) bool {
		var root yaml.Node
		if err := yaml.Unmarshal(content, &root); err != nil {
			return false
		}
		v, err := yptr.Find(&root, key)
		if err != nil {
			return false
		}

		return v.Value == value
	}
}

func keyContainsValue(key, value string) func(content []byte, data map[string]interface{}) bool {
	return func(content []byte, data map[string]interface{}) bool {
		var root yaml.Node
		if err := yaml.Unmarshal(content, &root); err != nil {
			return false
		}
		v, err := yptr.Find(&root, key)
		if err != nil {
			return false
		}

		return strings.Contains(v.Value, value)
	}
}

func replaceKeyValue(key, value string) func(content []byte) []byte {
	return func(content []byte) []byte {

		var root yaml.Node
		if err := yaml.Unmarshal(content, &root); err != nil {
			fmt.Println(err)
			return content
		}

		nameNode, err := yptr.Find(&root, key)
		if err != nil {
			fmt.Println(err)
			return content
		}

		out, _, err := transform.Bytes(yamled.T(
			yamled.Node(nameNode).With(value),
		), content)

		if err == nil {
			return out
		}

		fmt.Println(err)
		return content
	}
}
