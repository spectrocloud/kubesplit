package helm

import (
	"fmt"

	yamled "github.com/vmware-labs/go-yaml-edit"
	yptr "github.com/vmware-labs/yaml-jsonpointer"
	"golang.org/x/text/transform"
	"gopkg.in/yaml.v3"
)

type rules struct {
	Condition   func(content []byte, data map[string]interface{}) bool
	Replacement func(content []byte) []byte
}

// Mutator mutates content with helm rules
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

// unused
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

		return content
	}
}
