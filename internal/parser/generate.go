package parser

import (
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

func Generate(reader io.Reader, dir string) error {
	//reader := io.Reader(os.Stdin)
	dec := yaml.NewDecoder(reader)

	for {
		var node yaml.Node
		err := dec.Decode(&node)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		content, err := yaml.Marshal(&node)
		if err != nil {
			return err
		}
		data := map[string]interface{}{}
		if err := yaml.Unmarshal(content, data); err != nil {
			return err
		}

		if err := processYAML(dir, content, data); err != nil {
			return err
		}
	}
	return nil
}
