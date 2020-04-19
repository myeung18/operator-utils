package openshift

import (
	"k8s.io/client-go/rest"
)

/*
GetPlatformInfo examines the Kubernetes-based environment and determines the running platform, version, & OS.
Accepts <nil> or instantiated 'cfg' rest config parameter.

Result: PlatformInfo{ Name: OpenShift, K8SVersion: 1.13+, OS: linux/amd64 }
*/
func GetPlatformInfo(cfg *rest.Config) (PlatformInfo, error) {
	return K8SBasedPlatformVersioner{}.GetPlatformInfo(nil, cfg)
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
Deprecated:
LookupOpenShiftVersion fetches OpenShift version info from API endpoints
*** NOTE: OCP 4.1+ requires elevated user permissions, see PlatformVersioner for details
Accepts <nil> or instantiated 'cfg' rest config parameter.

Result: OpenShiftVersion{ Version: 4.1.2 }
*/
func LookupOpenShiftVersion(cfg *rest.Config) (OpenShiftVersion, error) {
	return K8SBasedPlatformVersioner{}.LookupOpenShiftVersion(nil, cfg)
}

/*
Compare the runtime OpenShift with the version passed in.
version: Semantic format
cfg : OpenShift platform config, use runtime config if nil is passed in.
return:
	-1 : if ver1 < ver2
	 0 : if ver1 == ver2
     1 : if ver1 > ver2
The int value returned should be discarded if err is not nil
*/
func CompareOpenShiftVersions(version string, cfg ...*rest.Config) (int, error) {
	return K8SBasedPlatformVersioner{}.compareOpenShiftVersions(nil, version, cfg)
}

/*
version: Semantic format
return:
	-1 : if ver1 < ver2
	 0 : if ver1 == ver2
     1 : if ver1 > ver2
The int value returned should be discarded if err is not nil
 */
func CompareVersions(version1 string, version2 string) int {
	return K8SBasedPlatformVersioner{}.CompareVersions(version1, version2)
}

/*
MapKnownVersion maps from K8S version of PlatformInfo to equivalent OpenShift version

Result: OpenShiftVersion{ Version: v4.1 }
*/
func MapKnownVersionx(info PlatformInfo) OpenShiftVersion {
	return MapKnownVersion(info)
}
