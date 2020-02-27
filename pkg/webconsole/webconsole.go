package webconsole

import (
	"fmt"
	creator2 "github.com/RHsyseng/operator-utils/pkg/webconsole/creator"
	"github.com/RHsyseng/operator-utils/pkg/webconsole/factory"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr/v2"
	"strings"
)

func LoadWebConsoleYamlSamples(path string, folder string) (map[string]string, error) {
	return loadFiles(path, folder)
}

func loadFiles(path string, folder string) (map[string]string, error) {
	filename := strings.Join([]string{path, folder}, "/")

	box := packr.New("folder name", filename)
	if box.List() == nil {
		fmt.Println("file not found")
		return nil, fmt.Errorf("%s not found ", filename)
	}

	resMap := make(map[string]string)
	for _, filename := range box.List() {
		yamlStr, err := box.FindString(filename)
		if err != nil {
			resMap[filename] = err.Error()
		}
		obj := &creator2.CustomResourceDefinition{}
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

func LoadWebConsoleEnrichment(path string, folder string) (map[string]string, error) {
	return loadFiles(path, folder)
}
