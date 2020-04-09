package openshift

import (
	"errors"
	"github.com/myeung18/operator-utils/internal/platform"
	"k8s.io/client-go/rest"
	"strconv"
	"strings"
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

func LookupOpenShiftSemVer(cfg *rest.Config) (platform.OpenShiftVersion, error) {
	return platform.K8SBasedPlatformVersioner{}.LookupClusterVersionSemVer(nil, cfg)
}

func LookupOpenShiftVersion4(cfg *rest.Config) (platform.OpenShiftVersion, error) {
	return platform.K8SBasedPlatformVersioner{}.TestLookupVersion4(nil, cfg)
}
func LookupOpenShiftVersion3(cfg *rest.Config) (platform.OpenShiftVersion, error) {
	return platform.K8SBasedPlatformVersioner{}.TestLookupVersion3(nil, cfg)
}
/*
MapKnownVersion maps from K8S version of PlatformInfo to equivalent OpenShift version

Result: OpenShiftVersion{ Version: 4.1.2 }
*/
func MapKnownVersion(info platform.PlatformInfo) platform.OpenShiftVersion {
	k8sToOcpMap := map[string]string{
		"1.10+": "3.10",
		"1.11+": "3.11",
		"1.13+": "4.1",
		"1.14+": "4.2",
		"1.16+": "4.3",
	}
	return platform.OpenShiftVersion{Version: k8sToOcpMap[info.K8SVersion]}
}

func CompareOpenShiftVersions(cfg *rest.Config, version string) (int, error) {
	isOcp, err := IsOpenShift(cfg)
	if err != nil {
		return -1, err
	}
	if !isOcp {
		return -1, errors.New("There is no OpenShift platform detected.")
	}
	info, err := LookupOpenShiftVersion(cfg)
	if err != nil {
		return -1, err
	}
	return CompareVersions(info.Version, version)
}

/*
Supported version format : Major.Minor.Patch
e.g.: 2.3.4
return:
	-1 : if ver1 < ver2
	 0 : if ver1 == ver2
     1 : if ver1 > ver2
The int value returned should be discarded if err is not nil
 */
func CompareVersions(ver1 string, ver2 string) (int, error) {
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
