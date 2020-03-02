package webconsole

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr/v2"
	creator "github.com/myeung18/operator-utils/pkg/webconsole/creator"
	"github.com/myeung18/operator-utils/pkg/webconsole/factory"
	"strings"
)

func LoadWebConsoleYamlSamples(name string, path string, folder string) (map[string]string, error) {
	return loadFiles(name, path, folder)
}

func loadFiles(name string, path string, folder string) (map[string]string, error) {
	filename := strings.Join([]string{path, folder}, "/")

	box := packr.New(name, filename)
	if box.List() == nil {
		return nil, fmt.Errorf("%s not found ", filename)
	}

	resMap := make(map[string]string)
	for _, filename := range box.List() {
		yamlStr, err := box.FindString(filename)
		if err != nil {
			resMap[filename] = err.Error()
		}
		obj := &creator.CustomResourceDefinition{}
		err = yaml.Unmarshal([]byte(yamlStr), obj)
		if err != nil {
			resMap[filename] = err.Error()
		}

		//check for any non ConsoleYAMLsamples
		creator := factory.GetCreator(obj.Kind)
		if creator == factory.NullCreatorImpl {
			kind := obj.Annotations["consolekind"]
			creator = factory.GetCreator(kind)
		}
		_, err = creator.Create(yamlStr)
		if err != nil {
			resMap[filename] = err.Error()
		} else if creator == factory.NullCreatorImpl {
			resMap[filename] =  "Unknown web console yaml"
		} else {
			resMap[filename] =  "processed"
		}
	}
	return resMap, nil
}

func LoadWebConsoleEnrichment(name string, path string, folder string) (map[string]string, error) {
	return loadFiles(name, path, folder)
}
