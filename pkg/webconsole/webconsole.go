package webconsole

import (
	"fmt"
	creator2 "github.com/myeung18/operator-utils/pkg/webconsole/creator"
	"github.com/myeung18/operator-utils/pkg/webconsole/factory"
	"github.com/ghodss/yaml"
	"github.com/gobuffalo/packr/v2"
	"io/ioutil"
	"strings"
)

func LoadWebConsoleYamlSamples(path string, folder string) (map[string]string, error) {
	return loadFilesWithIO(path, folder)
}

func loadFilesWithIO(path string, folder string) (map[string]string, error) {
	fullpath := strings.Join([]string{path, folder}, "/")
	fileList, err := ioutil.ReadDir(fullpath)
	if err != nil {
		return nil, fmt.Errorf("%s not found with io ", fullpath)
	}

	resMap := make(map[string]string)
	for _, filename := range fileList {
		process(fullpath, filename.Name(), resMap)
	}
	return resMap, nil
}

func process(fullpath string, filename string, resMap map[string]string) {
	yamlStr, err := ioutil.ReadFile(fullpath + "/" + filename)
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
	_, err = creator.Create(string(yamlStr))
	if err != nil {
		resMap[filename] = err.Error()
	} else if creator == factory.NullCreatorImpl {
		resMap[filename] =  "Unknown web console yaml"
	} else {
		resMap[filename] =  "processed"
	}

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
