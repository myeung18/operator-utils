package platform

import (
	"fmt"
	"golang.org/x/mod/semver"
	"strings"
)

type PlatformType string

const (
	OpenShift  PlatformType = "OpenShift"
	Kubernetes PlatformType = "Kubernetes"
)

type PlatformInfo struct {
	Name       PlatformType `json:"name"`
	K8SVersion string       `json:"k8sVersion"`
	OS         string       `json:"os"`
}

func (info PlatformInfo) K8SMajorVersion() string {
	return strings.Split(info.K8SVersion, ".")[0]
}

func (info PlatformInfo) K8SMinorVersion() string {
	return strings.Split(info.K8SVersion, ".")[1]
}

func (info PlatformInfo) IsOpenShift() bool {
	return info.Name == OpenShift
}

func (info PlatformInfo) IsKubernetes() bool {
	return info.Name == Kubernetes
}

func (info PlatformInfo) String() string {
	return "PlatformInfo [" +
		"Name: " + fmt.Sprintf("%v", info.Name) +
		", K8SVersion: " + info.K8SVersion +
		", OS: " + info.OS + "]"
}

type OpenShiftVersion struct {
	Version string `json:"ocpVersion"`
}

func (info OpenShiftVersion) MajorVersion() string {
	return semver.Major(info.Version)
}

func (info OpenShiftVersion) MinorVersion() string {
	ver := semver.MajorMinor(info.Version)
	if ver != "" {
		return strings.Split(info.Version, ".")[1]
	}
	return ""
}

func (info OpenShiftVersion) PrereleaseVersion() string {
	return semver.Prerelease(info.Version)
}

func (info OpenShiftVersion) BuildVersion() string {
	return semver.Build(info.Version)
}

func (info OpenShiftVersion) String() string {
	return "OpenShiftVersion [" +
		"Version: " + info.Version + "]"
}

func (v OpenShiftVersion) Compare(o OpenShiftVersion) int {
	return semver.Compare(v.Version, o.Version)
}

// full generated 'version' API fetch result struct @
// gist.github.com/jeremyary/5a66530611572a057df7a98f3d2902d5
type PlatformClusterInfo struct {
	Status struct {
		Desired struct {
			Version string `json:"version"`
		} `json:"desired"`
	} `json:"status"`
}
