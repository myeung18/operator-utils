package platform

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-semver/semver"
	openapi_v2 "github.com/googleapis/gnostic/OpenAPIv2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	log                   = logf.Log.WithName("utils")
	ClusterVersionApiPath = "apis/config.openshift.io/v1/clusterversions/version"
)

type PlatformVersioner interface {
	GetPlatformInfo(discoverer Discoverer, cfg *rest.Config) (PlatformInfo, error)
}

type Discoverer interface {
	ServerVersion() (*version.Info, error)
	ServerGroups() (*v1.APIGroupList, error)
	OpenAPISchema() (*openapi_v2.Document, error)
	RESTClient() rest.Interface
}

type K8SBasedPlatformVersioner struct{}

// deal with cfg coming from legacy method signature and allow injection for client testing
func (K8SBasedPlatformVersioner) DefaultArgs(client Discoverer, cfg *rest.Config) (Discoverer, *rest.Config, error) {
	if cfg == nil {
		var err error
		cfg, err = config.GetConfig()
		if err != nil {
			return nil, nil, err
		}
	}
	if client == nil {
		var err error
		client, err = discovery.NewDiscoveryClientForConfig(cfg)
		if err != nil {
			return nil, nil, err
		}
	}
	return client, cfg, nil
}

func (pv K8SBasedPlatformVersioner) GetPlatformInfo(client Discoverer, cfg *rest.Config) (PlatformInfo, error) {
	log.Info("detecting platform version...")
	info := PlatformInfo{Name: Kubernetes}

	var err error
	client, cfg, err = pv.DefaultArgs(client, cfg)
	if err != nil {
		log.Info("issue occurred while defaulting client/cfg args")
		return info, err
	}

	k8sVersion, err := client.ServerVersion()
	if err != nil {
		log.Info("issue occurred while fetching ServerVersion")
		return info, err
	}
	info.K8SVersion = k8sVersion.Major + "." + k8sVersion.Minor
	info.OS = k8sVersion.Platform

	apiList, err := client.ServerGroups()
	if err != nil {
		log.Info("issue occurred while fetching ServerGroups")
		return info, err
	}

	for _, v := range apiList.Groups {
		if v.Name == "route.openshift.io" {

			log.Info("route.openshift.io found in apis, platform is OpenShift")
			info.Name = OpenShift
			break
		}
	}
	log.Info(info.String())
	return info, nil
}

/*
OCP4.1+ requires elevated cluster configuration user security permissions for version fetch
REST call URL requiring permissions: /apis/config.openshift.io/v1/clusterversions
*/
func (pv K8SBasedPlatformVersioner) LookupOpenShiftVersion(client Discoverer, cfg *rest.Config) (OpenShiftVersion, error) {

	osv := OpenShiftVersion{}
	client, _, err := pv.DefaultArgs(nil, nil)
	if err != nil {
		log.Info("issue occurred while defaulting args for version lookup")
		return osv, err
	}
	doc, err := client.OpenAPISchema()
	if err != nil {
		log.Info("issue occurred while fetching OpenAPISchema")
		return osv, err
	}
	log.Info("doc info >>>+: ", "doc", doc.Info)
	fmt.Println("doc info >>>: ", doc.Info)

	switch doc.Info.Version[:4] {
	case "v3.1":
		osv.Version = doc.Info.Version

	// OCP4 returns K8S major/minor from old API endpoint [bugzilla-1658957]
	case "v1.1":
		rest := client.RESTClient().Get().AbsPath(ClusterVersionApiPath)

		result := rest.Do()
		if result.Error() != nil {
			log.Info("issue making API version rest call: " + result.Error().Error())
			return osv, result.Error()
		}

		// error handling before/after Raw() seems redundant, but error detail can be lost in convert
		body, err := result.Raw()
		if err != nil {
			log.Info("issue pulling raw result from API call")
			return osv, err
		}

		var cvi PlatformClusterInfo
		err = json.Unmarshal(body, &cvi)
		if err != nil {
			log.Info("issue occurred while unmarshalling PlatformClusterInfo")
			return osv, err
		}
		osv.Version = cvi.Status.Desired.Version
	}
	return osv, nil
}

func (pv K8SBasedPlatformVersioner) TestLookupVersion3(client Discoverer, cfg *rest.Config) (OpenShiftVersion, error) {
	osv := OpenShiftVersion{}
	client, _, err := pv.DefaultArgs(nil, nil)
	if err != nil {
		log.Info("issue occurred while defaulting args for version lookup")
		return osv, err
	}
	doc, err := client.OpenAPISchema()
	if err != nil {
		log.Info("issue occurred while fetching OpenAPISchema")
		return osv, err
	}
	log.Info("doc info >>>+: ", "doc", doc.Info)

	switch doc.Info.Version[:4] {
	case "v3.1":
		osv.Version = doc.Info.Version
	default:
		log.Info("default ver3 Info: ", "doc", doc.Info)
	}
	return osv, nil
}

func (pv K8SBasedPlatformVersioner) TestLookupVersion4(client Discoverer, cfg *rest.Config) (OpenShiftVersion, error) {

	osv := OpenShiftVersion{}
	client, _, err := pv.DefaultArgs(nil, nil)
	if err != nil {
		log.Info("issue occurred while defaulting args for version lookup")
		return osv, err
	}
	rest := client.RESTClient().Get().AbsPath(ClusterVersionApiPath)

	result := rest.Do()
	if result.Error() != nil {
		log.Info("issue making API version rest call: " + result.Error().Error())
		return osv, result.Error()
	}

	// error handling before/after Raw() seems redundant, but error detail can be lost in convert
	body, err := result.Raw()
	if err != nil {
		log.Info("issue pulling raw result from API call")
		return osv, err
	}

	var cvi PlatformClusterInfo
	err = json.Unmarshal(body, &cvi)
	if err != nil {
		log.Info("issue occurred while unmarshalling PlatformClusterInfo")
		return osv, err
	}
	osv.Version = cvi.Status.Desired.Version
	log.Info("Infov4 : ", osv)

	return osv, nil
}

func (pv K8SBasedPlatformVersioner)  LookupClusterVersionSemVer(client Discoverer, config *rest.Config) (OpenShiftVersion, error) {
	osv := OpenShiftVersion{}
	//client, _, err := pv.DefaultArgs(nil, nil)
	//if err != nil {
	//	log.Info("issue occurred while defaulting args for version lookup")
	//	return osv, err
	//}

	configClient, err := configv1client.NewForConfig(config)
	if err != nil {
		log.Error(err, "Failed to create config client")
		return osv, err
	}

	var openShiftSemVer *semver.Version
	clusterVersion, err := configClient.
		ConfigV1().
		ClusterVersions().
		Get("version", metav1.GetOptions{})

	log.Info("print1: ", "ver: " , configClient.ConfigV1())
	log.Info("print2: ", "ver: " , configClient.ConfigV1().ClusterVersions())
	if err != nil {
		log.Info("err != nil : ", "ver: " , semver.Version{})
		if errors.IsNotFound(err) {
			// default to OpenShift 3 as ClusterVersion API was introduced in OpenShift 4
			openShiftSemVer, _ = semver.NewVersion("3")
		} else {
			log.Error(err, "Failed to get OpenShift cluster version")
			return osv, err
		}
	} else {
		//latest version from the history
		log.Info("err == nil : ", "ver: " , semver.Version{})
		v := clusterVersion.Status.History[0].Version
		openShiftSemVer, err = semver.NewVersion(v)
		if err != nil {
			log.Error(err, "Failed to get OpenShift cluster version")
			return osv, err
		}
	}

	fmt.Println("openShiftSemVer-config semVer ", openShiftSemVer)
	log.Info("openShiftSemVer-config semVer:", "ver: " , openShiftSemVer)
	return osv, nil
}