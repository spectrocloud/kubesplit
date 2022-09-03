package parser

import (
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

// Mutator defines a mutator which changes the content to be written after parsing.
type Mutator func(content []byte, data map[string]interface{}) []byte

type Config struct {
	Mutators []Mutator
}

type Option func(c *Config) error

func (c *Config) Apply(opts ...Option) error {
	for _, o := range opts {
		if err := o(c); err != nil {
			return err
		}
	}
	return nil
}

// WithMutators adds mutator to the parsing process.
func WithMutators(m ...Mutator) Option {
	return func(c *Config) error {
		c.Mutators = append(c.Mutators, m...)

		return nil
	}
}

// Generates a set of files ordered by type.
func Generate(reader io.Reader, dir string, opts ...Option) error {

	c := &Config{}
	if err := c.Apply(opts...); err != nil {
		return err
	}

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

		for _, m := range c.Mutators {
			content = m(content, data)
		}

		if err := processYAML(dir, content, data); err != nil {
			return err
		}
	}
	return nil
}
