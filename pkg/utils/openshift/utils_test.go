package openshift

import (
	"context"
	"fmt"
	"github.com/RHsyseng/operator-utils/internal/platform"
	"github.com/ash2k/stager"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"testing"
	"time"

	apiext_v1b1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	//apisvr "k8s.io/apiextensions-apiserver"
	apiExtClientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiext_v1b1inf "k8s.io/apiextensions-apiserver/pkg/client/informers/externalversions/apiextensions/v1beta1"
	apiext_v1b1list "k8s.io/apiextensions-apiserver/pkg/client/listers/apiextensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/tools/cache"
)

func TestGetSampleCrd(t *testing.T) {
	crd := &apiext_v1b1.CustomResourceDefinition{
		ObjectMeta: meta_v1.ObjectMeta{
			Name: "",
		},
	}
	_ = crd

	gvk := schema.GroupVersionKind{Group: "console.openshift.io", Version: "v1", Kind: "consoleyamssamples"}
	fmt.Println(gvk.GroupVersion(), "/" , gvk.Kind)
}


func TestClientGoCall(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
	}

	stgr := stager.New()
	defer stgr.Shutdown()

	ctxTest, cancel := context.WithTimeout(context.Background(), 200 * time.Second)
	defer cancel()

	apiClient, err := apiExtClientset.NewForConfig(cfg)
	if err != nil {
		fmt.Println(err)
	}
	crdInf := apiext_v1b1inf.NewCustomResourceDefinitionInformer(apiClient, 0, cache.Indexers{})
	stage := stgr.NextStage()
	stage.StartWithChannel(crdInf.Run)

	if !cache.WaitForCacheSync(ctxTest.Done(), crdInf.HasSynced) {
		t.Fatal("wait for CRD Informer was cancelled")
	}
	crdLister := apiext_v1b1list.NewCustomResourceDefinitionLister(crdInf.GetIndexer())
	obj, err := crdLister.Get("consoleyamlsamples.console.openshift.io")
	if err != nil  {
		fmt.Println("err", err)
	} else {
		fmt.Println(obj.Name)
	}
}






func TestOpenShiftVersion_MapKnownVersion(t *testing.T) {

	cases := []struct {
		label              string
		info               platform.PlatformInfo
		expectedOCPVersion string
	}{
		{
			label:              "case 1",
			info:               platform.PlatformInfo{K8SVersion: ""},
			expectedOCPVersion: "",
		},
		{
			label:              "case 2",
			info:               platform.PlatformInfo{K8SVersion: "1.10+"},
			expectedOCPVersion: "3.10",
		},
		{
			label:              "case 3",
			info:               platform.PlatformInfo{K8SVersion: "1.11+"},
			expectedOCPVersion: "3.11",
		},
		{
			label:              "case 4",
			info:               platform.PlatformInfo{K8SVersion: "1.13+"},
			expectedOCPVersion: "4.1",
		},
	}

	for _, v := range cases {
		assert.Equal(t, v.expectedOCPVersion, MapKnownVersion(v.info).Version, v.label+": expected OCP version to match")
	}
}

type MockDiscoverer struct {
	apiResourceList *metav1.APIResourceList
}

func (d MockDiscoverer) ServerResourcesForGroupVersion(groupVersion string) (resources *metav1.APIResourceList, err error) {
	return d.apiResourceList, nil
}

func TestCustomResourceExists(t *testing.T) {
	ts := []struct {
		label           string
		gv              string
		api             string
		discoveryClient *MockDiscoverer
		config          []*rest.Config
		expectetResult  bool
		expectErr       bool
	}{
		{
			label: "test 1",
			gv:    "console.openshift.io/v1",
			api:   "consoleyamlsamples",
			discoveryClient: &MockDiscoverer{
				apiResourceList: &metav1.APIResourceList{
					GroupVersion: "console.openshift.io/v1",
					APIResources: []metav1.APIResource{{Name: "consolelinks"}}},
			},
			config:         []*rest.Config{&rest.Config{}},
			expectetResult: false,
			expectErr:      true,
		},
		{
			label: "test 1",
			gv:    "console.openshift.io/v1",
			api:   "consoleyamlsamples",
			discoveryClient: &MockDiscoverer{
				apiResourceList: &metav1.APIResourceList{
					GroupVersion: "console.openshift.io/v1",
					APIResources: []metav1.APIResource{{Name: "consoleyamlsamples"}}},
			},
			config:         []*rest.Config{&rest.Config{}},
			expectetResult: true,
			expectErr:      false,
		},
	}

	ps := OCPPlatformService{}
	for _, test := range ts {
		res, err := ps.customResourceExists(test.gv, test.api, test.discoveryClient, test.config)
		if test.expectErr {
			assert.Error(t, err, "expeting error"+test.label)
		} else {
			assert.NoError(t, err, "unexpeted error"+test.label)
		}
		assert.Equal(t, test.expectetResult, res, "the expected and actual results are not the same "+test.label)
	}

	res, err := CustomResourceExistsDirect("console.openshift.io/v1", "consoleyamlsamples", nil)
	assert.Nil(t, err)
	assert.Equal(t, true, res, "")
	//
	//res, err = CustomResourceExists("console.openshift.io/v1", "consoleyamlsamxxx", nil)
	//assert.Error(t, err);

	//res, err := CustomResourceExists("console.openshift.io/v1", "consoleyamlsamples", nil)
	//if err != nil {
	//
	//}
	//fmt.Println(res)

}
