package openshift

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/myeung18/operator-utils/pkg/utils"
	oappsv1 "github.com/openshift/api/apps/v1"
	v1 "github.com/openshift/api/console/v1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	"strconv"
	"strings"
	"testing"
)

func TestOpenshift(t *testing.T) {
	kubeconfig, err := utils.GetConfig()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(kubeconfig)
	info, err := GetPlatformInfo(kubeconfig)
	versionStr := MapKnownVersion(info)
	fmt.Println("info:", info, "ocp:", versionStr)

	var comp int = -1
	if comp, err = compareVersion_xx(versionStr.Version, "4.3"); err == nil {
		fmt.Println("comp: ", comp, ", res: ", comp >= 0);
	}

	fmt.Println(info.K8SVersion)
	fmt.Println(IsOpenShift(kubeconfig))

	fmt.Println("------")
	info2, err := LookupOpenShiftVersion(kubeconfig)
	fmt.Println("info2:", info2, MapKnownVersion(info))

}

func compareVersion_xx(ver1 string, ver2 string) (int, error) {
	ver1Nums := strings.Split(ver1, ".")
	ver2Nums := strings.Split(ver2, ".")
	length := len(ver1Nums)
	if length < len(ver2Nums) {
		length = len(ver2Nums)
	}
	for i := 0; i < length; i++ {
		v1 := 0
		if i < len(ver1Nums) {
			if v1, err := strconv.Atoi(ver1Nums[i]); err != nil {
				_ = v1
				return -1, err
			}
		}
		v2 := 0
		if i < len(ver2Nums) {
			if v2, err := strconv.Atoi(ver2Nums[i]); err != nil {
				_ = v2
				return -1, err
			}
		}
		if v1 > v2 {
			return 1, nil
		} else if v2 > v1 {
			return -1, nil
		}
	}
	return 0, nil
}

func TestGetConsoleYAMLSample(t *testing.T) {
	var inputYaml = `
apiVersion: v1
kind: DeploymentConfig
metadata:
 name: sample-dc
 annotations:
   consoleName: sample-deploymentconfig
   consoleDesc: Sample Deployment Config
   consoleTitle: Sample Deployment Config
spec:
   replicas: 2
`
	original := &oappsv1.DeploymentConfig{}
	assert.NoError(t, yaml.Unmarshal([]byte(inputYaml), original))

	yamlSample, err := GetConsoleYAMLSample(original)
	assert.NoError(t, err)

	assert.Equal(t, "sample-deploymentconfig", yamlSample.ObjectMeta.Name)
	assert.Equal(t, "openshift-console", yamlSample.ObjectMeta.Namespace)
	assert.Equal(t, "v1", yamlSample.Spec.TargetResource.APIVersion)
	assert.Equal(t, "DeploymentConfig", yamlSample.Spec.TargetResource.Kind)
	assert.Equal(t, v1.ConsoleYAMLSampleTitle("Sample Deployment Config"), yamlSample.Spec.Title)
	assert.Equal(t, v1.ConsoleYAMLSampleDescription("Sample Deployment Config"), yamlSample.Spec.Description)

	yamlContent := yamlSample.Spec.YAML
	actual := &oappsv1.DeploymentConfig{}
	assert.NoError(t, yaml.Unmarshal([]byte(string(yamlContent)), actual))

	original.SetAnnotations(nil)
	assert.EqualValues(t, original, actual, "original yaml should be the same as the actual yaml")
}

func TestGetConsoleYAMLSampleWithNoAnnotations(t *testing.T) {
	var inputYaml = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
`
	original := &appsv1.Deployment{}
	assert.NoError(t, yaml.Unmarshal([]byte(inputYaml), original))

	yamlSample, err := GetConsoleYAMLSample(original)
	assert.NoError(t, err)

	assert.Equal(t, "nginx-deployment-yamlsample", yamlSample.ObjectMeta.Name)
	assert.Equal(t, "openshift-console", yamlSample.ObjectMeta.Namespace)
	assert.Equal(t, "apps/v1", yamlSample.Spec.TargetResource.APIVersion)
	assert.Equal(t, "Deployment", yamlSample.Spec.TargetResource.Kind)
	assert.Equal(t, v1.ConsoleYAMLSampleTitle("nginx-deployment-yamlsample"), yamlSample.Spec.Title)
	assert.Equal(t, v1.ConsoleYAMLSampleDescription("nginx-deployment-yamlsample"), yamlSample.Spec.Description)

	yamlContent := yamlSample.Spec.YAML
	actual := &appsv1.Deployment{}
	assert.NoError(t, yaml.Unmarshal([]byte(string(yamlContent)), actual))

	assert.EqualValues(t, original, actual, "original yaml should be the same as the actual yaml")
}
