package openshift

import (
	"errors"
	"github.com/RHsyseng/operator-utils/internal/platform"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

/*
GetPlatformInfo examines the Kubernetes-based environment and determines the running platform, version, & OS.
Accepts <nil> or instantiated 'cfg' rest config parameter.

Result: PlatformInfo{ Name: OpenShift, K8SVersion: 1.13+, OS: linux/amd64 }
*/
func GetPlatformInfo(cfg *rest.Config) (platform.PlatformInfo, error) {
	return platform.K8SBasedPlatformVersioner{}.GetPlatformInfo(nil, cfg)
}

/*
IsOpenShift is a helper method to simplify boolean OCP checks against GetPlatformInfo results
Accepts <nil> or instantiated 'cfg' rest config parameter.
*/
func IsOpenShift(cfg *rest.Config) (bool, error) {
	info, err := GetPlatformInfo(cfg)
	if err != nil {
		return false, err
	}
	return info.IsOpenShift(), nil
}

/*
LookupOpenShiftVersion fetches OpenShift version info from API endpoints
*** NOTE: OCP 4.1+ requires elevated user permissions, see PlatformVersioner for details
Accepts <nil> or instantiated 'cfg' rest config parameter.

Result: OpenShiftVersion{ Version: 4.1.2 }
*/
func LookupOpenShiftVersion(cfg *rest.Config) (platform.OpenShiftVersion, error) {
	return platform.K8SBasedPlatformVersioner{}.LookupOpenShiftVersion(nil, cfg)
}

/*
Supported platform: OpenShift
cfg : OpenShift platform config, use runtime config if nil is passed in.
version: Supported version format : Major.Minor
	       e.g.: 4.3
*/
func CompareOpenShiftVersion(cfg *rest.Config, version string) (int, error) {
	return platform.K8SBasedPlatformVersioner{}.CompareOpenShiftVersion(nil, cfg, version)
}

/*
MapKnownVersion maps from K8S version of PlatformInfo to equivalent OpenShift version

Result: OpenShiftVersion{ Version: 4.1.2 }
*/
func MapKnownVersion(info platform.PlatformInfo) platform.OpenShiftVersion {
	return platform.MapKnownVersion(info)
}

func CustomResourceExistsDirect(groupVersion string, apiResource string, cfg ...*rest.Config) (bool, error) {
	var client *discovery.DiscoveryClient
	var err error
	if len(cfg) > 0 {
		client, err = getDiscoveryClient(client, cfg[0])
	} else {
		client, err = getDiscoveryClient(client, nil)
	}
	if err != nil {
		return false, errors.New("issue occurred while defaulting args for groupVersion lookup:" + err.Error())
	}
	return getCustomResource(groupVersion, apiResource, client)
}

func getDiscoveryClient(client *discovery.DiscoveryClient, cfg *rest.Config) (*discovery.DiscoveryClient, error) {
	if cfg == nil {
		var err error
		cfg, err = config.GetConfig()
		if err != nil {
			return nil, err
		}
	}
	if client == nil {
		var err error
		client, err = discovery.NewDiscoveryClientForConfig(cfg)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

/***********************************************************************************************/

func CustomResourceExists(groupVersion string, apiResource string, cfg ...*rest.Config) (bool, error) {
	return OCPPlatformService{}.customResourceExists(groupVersion, apiResource, nil, cfg)
}

func (ocp OCPPlatformService) customResourceExists(groupVersion string, apiResource string, client PlatformService, cfg []*rest.Config) (bool, error) {
	var err error
	if len(cfg) > 0 {
		client, err = ocp.getDiscoveryClient(client, cfg[0])
	} else {
		client, err = ocp.getDiscoveryClient(client, nil)
	}
	if err != nil {
		return false, errors.New("issue occurred while defaulting args for groupVersion lookup:" + err.Error())
	}
	return getCustomResource(groupVersion, apiResource, client)
}

func getCustomResource(groupVersion string, apiResource string, client PlatformService) (bool, error) {
	apis, err := client.ServerResourcesForGroupVersion(groupVersion)
	if err != nil {
		return false, errors.New("not getting group version: " + err.Error())
	}
	for _, api := range apis.APIResources {
		if api.Name == apiResource {
			return true, nil
		}
	}
	return false, errors.New(apiResource + " is not defined in this group:" + groupVersion)
}

func (ocp OCPPlatformService) getDiscoveryClient(client PlatformService, cfg *rest.Config) (PlatformService, error) {
	if cfg == nil {
		var err error
		cfg, err = config.GetConfig()
		if err != nil {
			return nil, err
		}
	}
	if client == nil {
		var err error
		client, err = discovery.NewDiscoveryClientForConfig(cfg)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

type PlatformService interface {
	ServerResourcesForGroupVersion(groupVersion string) (resources *metav1.APIResourceList, err error)
}

type OCPPlatformService struct{}
