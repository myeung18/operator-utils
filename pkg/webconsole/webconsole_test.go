package webconsole

import (
	"fmt"
	"testing"
)

func TestLoadYaml(t *testing.T) {
	resMap, err := LoadWebConsoleYamlSamples("../../examples", "resources")
	if err != nil {
		fmt.Println(err)
	}

	for k, v := range resMap {
		fmt.Println(k, " - ", v)
	}
}

func TestWebConsoleEnrichment(t *testing.T) {
	resMap, err := LoadWebConsoleEnrichment("../../examples", "console")
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range resMap {
		fmt.Println(k, " - ", v)
	}
}
