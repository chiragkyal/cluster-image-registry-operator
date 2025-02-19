package framework

import (
	"fmt"
	"testing"

	clientappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	clientbatchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	clientcoordinationv1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	clientstoragev1 "k8s.io/client-go/kubernetes/typed/storage/v1"
	restclient "k8s.io/client-go/rest"

	buildv1 "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	clientconfigv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	imagev1 "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	clientimageregistryv1 "github.com/openshift/client-go/imageregistry/clientset/versioned/typed/imageregistry/v1"
	clientroutev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"

	"github.com/openshift/cluster-image-registry-operator/pkg/client"
)

// Clientset is a set of Kubernetes clients.
type Clientset struct {
	clientcorev1.CoreV1Interface
	clientappsv1.AppsV1Interface
	clientconfigv1.ConfigV1Interface
	clientimageregistryv1.ImageregistryV1Interface
	clientroutev1.RouteV1Interface
	clientstoragev1.StorageV1Interface
	clientbatchv1.BatchV1Interface
	clientcoordinationv1.CoordinationV1Interface
	ImageInterface imagev1.ImageV1Interface
	BuildInterface buildv1.BuildV1Interface
}

// NewClientset creates a set of Kubernetes clients. The default kubeconfig is
// used if not provided.
func NewClientset(kubeconfig *restclient.Config) (clientset *Clientset, err error) {
	if kubeconfig == nil {
		kubeconfig, err = client.GetConfig()
		if err != nil {
			return nil, fmt.Errorf("unable to get kubeconfig: %s", err)
		}
	}

	clientset = &Clientset{}
	clientset.BatchV1Interface, err = clientbatchv1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.CoreV1Interface, err = clientcorev1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.AppsV1Interface, err = clientappsv1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.ConfigV1Interface, err = clientconfigv1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.ImageregistryV1Interface, err = clientimageregistryv1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.RouteV1Interface, err = clientroutev1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.StorageV1Interface, err = clientstoragev1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.CoordinationV1Interface, err = clientcoordinationv1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.BuildInterface, err = buildv1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	clientset.ImageInterface, err = imagev1.NewForConfig(kubeconfig)
	if err != nil {
		return
	}
	return
}

// MustNewClientset is like NewClienset but aborts the test if clienset cannot
// be constructed.
func MustNewClientset(t *testing.T, kubeconfig *restclient.Config) *Clientset {
	clientset, err := NewClientset(kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	return clientset
}
