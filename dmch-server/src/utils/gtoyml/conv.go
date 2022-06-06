package gtoyml

import (
	"errors"
	"reflect"

	"gopkg.in/yaml.v2"
)

type DocGenerator struct {
	modelsDocs map[string]any
}

func NewDocGenerator() *DocGenerator {
	return &DocGenerator{
		modelsDocs: make(map[string]any),
	}
}

func (c *DocGenerator) AddModel(models ...any) error {
	for i := range models {
		model := models[i]
		rt := reflect.TypeOf(model)
		if rt.Kind() != reflect.Struct {
			return errors.New("bad type")
		}

		desc, err := c.genDescForType(rt)
		if err != nil {
			return err
		}

		c.modelsDocs[rt.Name()] = desc
	}

	return nil
}

func (c *DocGenerator) EncodeYaml(models ...interface{}) (string, error) {
	d, err := yaml.Marshal(c.modelsDocs)
	if err != nil {
		return "", err
	}

	return string(d), nil
}
