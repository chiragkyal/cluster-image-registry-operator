package v1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorsv1api "github.com/openshift/api/operator/v1"
)

const (

	// DefaultRouteName is the name of the default route created for the registry
	// when a default route is requested from the operator
	DefaultRouteName = "default-route"
	// ImageRegistryName is the name of the image-registry workload resource (deployment)
	ImageRegistryName = "image-registry"

	// ImageRegistryResourceName is the name of the image registry config instance
	ImageRegistryResourceName = "instance"

	// ImageRegistryCertificatesName is the name of the configmap that is managed by the
	// registry operator and mounted into the registry pod, to provide additional
	// CAs to be trusted during image pullthrough
	ImageRegistryCertificatesName = ImageRegistryName + "-certificates"

	// ImageRegistryPrivateConfiguration is the name of a secret that is managed by the
	// registry operator and which provides credentials to the registry for things like
	// accessing S3 storage
	ImageRegistryPrivateConfiguration = ImageRegistryName + "-private-configuration"

	// ImageRegistryPrivateConfigurationUser is the name of a secret that is managed by
	// the administrator and which provides credentials to the registry for things like
	// accessing S3 storage.  This content takes precedence over content the operator
	// automatically pulls from other locations, and it is merged into ImageRegistryPrivateConfiguration
	ImageRegistryPrivateConfigurationUser = ImageRegistryPrivateConfiguration + "-user"

	// ImageRegistryOperatorNamespace is the namespace containing the registry operator
	// and the registry itself
	ImageRegistryOperatorNamespace = "openshift-image-registry"

	// Status Conditions

	// OperatorStatusTypeRemoved denotes that the image-registry instance has been
	// removed
	// TODO: does this serve any purpose?  As soon as it's removed we'll bootstrap
	// a new one.
	OperatorStatusTypeRemoved = "Removed"

	// StorageExists denotes whether or not the registry storage medium exists
	StorageExists = "StorageExists"

	// StorageTagged denotes whether or not the registry storage medium
	// that we created was tagged correctly
	StorageTagged = "StorageTagged"

	// StorageEncrypted denotes whether or not the registry storage medium
	// that we created has encryption enabled
	StorageEncrypted = "StorageEncrypted"

	// StorageIncompleteUploadCleanupEnabled denotes whethere or not the registry storage
	// medium is configured to automatically cleanup incomplete uploads
	StorageIncompleteUploadCleanupEnabled = "StorageIncompleteUploadCleanupEnabled"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Config `json:"items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Config is the configuration object for a registry instance managed by
// the registry operator
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ImageRegistrySpec   `json:"spec"`
	Status            ImageRegistryStatus `json:"status,omitempty"`
}

type ImageRegistrySpec struct {
	// ManagementState indicates whether the registry instance represented
	// by this config instance is under operator management or not.  Valid
	// values are Managed, Unmanaged, and Removed
	ManagementState operatorsv1api.ManagementState `json:"managementState"`

	// HTTPSecret is the value needed by the registry to secure uploads, generated by default
	HTTPSecret string `json:"httpSecret,omitempty"`

	// Proxy defines the proxy to be used when calling master api, upstream registries, etc
	Proxy ImageRegistryConfigProxy `json:"proxy,omitempty"`

	// Storage details for configuring registry storage, e.g. S3 bucket coordinates.
	Storage ImageRegistryConfigStorage `json:"storage,omitempty"`

	// Requests controls how many parallel requests a given registry instance will handle before queuing additional requests
	Requests ImageRegistryConfigRequests `json:"requests,omitempty"`

	// CAConfigName identifies the configmap in the openshift-config namespace which contains
	// additional CAs to be trusted by the registry when doing pullthrough
	CAConfigName string `json:"caConfigName,omitempty"`

	// DefaultRoute indicates whether an external facing route for the registry
	// should be created using the default generated hostname
	DefaultRoute bool `json:"defaultRoute,omitempty"`

	// Routes defines additional external facing routes which should be created for the registry
	Routes []ImageRegistryConfigRoute `json:"routes,omitempty"`

	// Replicas determines the number of registry instances to run
	Replicas int32 `json:"replicas,omitempty"`

	// LogLevel determines the level of logging enabled in the registry
	LogLevel int64 `json:"logging,omitempty"`

	// Resources defines the resource requests+limits for the registry pod
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// NodeSelector defines the node selection constraints for the registry pod
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

type ImageRegistryStatus struct {
	operatorsv1api.OperatorStatus `json:",inline"`

	// StorageManaged is a boolean which denotes whether or not
	// we created the registry storage medium (such as an
	// S3 bucket)
	StorageManaged bool `json:"storageManaged"`

	// Storage indicates the current applied storage configuration of the registry
	Storage ImageRegistryConfigStorage `json:"storage"`
}

type ImageRegistryConfigProxy struct {
	HTTP    string `json:"http,omitempty"`
	HTTPS   string `json:"https,omitempty"`
	NoProxy string `json:"noProxy,omitempty"`
}

// ImageRegistryConfigStorageS3 holds the information to configure
// the registry to use the AWS S3 service for backend storage
// https://docs.docker.com/registry/storage-drivers/s3/
type ImageRegistryConfigStorageS3 struct {
	// Bucket is the bucket name in which you want to store the registry's data
	// Optional, will be generated if not provided
	Bucket string `json:"bucket,omitempty"`
	// Region is the AWS region in which your bucket exists
	// Optional, will be set based on the installed AWS Region
	Region string `json:"region,omitempty"`
	// RegionEndpoint is the endpoint for S3 compatible storage services
	// Optional, defaults based on the Region that is provided
	RegionEndpoint string `json:"regionEndpoint,omitempty"`
	// Encrypt specifies whether the registry stores the image in encrypted format or not
	// Optional, defaults to false
	Encrypt bool `json:"encrypt,omitempty"`
	// KeyID is the KMS key ID to use for encryption
	// Optional, Encrypt must be true, or this parameter is ignored
	KeyID string `json:"keyID,omitempty"`
}

type ImageRegistryConfigStorageAzure struct {
	Container string `json:"container,omitempty"`
}

type ImageRegistryConfigStorageGCS struct {
	Bucket string `json:"bucket,omitempty"`
}

type ImageRegistryConfigStorageSwift struct {
	AuthURL   string `json:"authURL,omitempty"`
	Container string `json:"container,omitempty"`
}

type ImageRegistryConfigStorageFilesystem struct {
	VolumeSource corev1.VolumeSource `json:"volumeSource,omitempty"`
}

type ImageRegistryConfigStorage struct {
	Azure      *ImageRegistryConfigStorageAzure      `json:"azure,omitempty"`
	Filesystem *ImageRegistryConfigStorageFilesystem `json:"filesystem,omitempty"`
	GCS        *ImageRegistryConfigStorageGCS        `json:"gcs,omitempty"`
	S3         *ImageRegistryConfigStorageS3         `json:"s3,omitempty"`
	Swift      *ImageRegistryConfigStorageSwift      `json:"swift,omitempty"`
}

type ImageRegistryConfigRequests struct {
	Read  ImageRegistryConfigRequestsLimits `json:"read,omitempty"`
	Write ImageRegistryConfigRequestsLimits `json:"write,omitempty"`
}

type ImageRegistryConfigRequestsLimits struct {

	// MaxRunning sets the maximum in flight api requests to the registry
	MaxRunning int `json:"maxRunning,omitempty"`

	// MaxInQueue sets the maximum queued api requests to the registry
	MaxInQueue int `json:"maxInQueue,omitempty"`

	// MaxWaitInQueue sets the maximum time a request can wait in the queue
	// before being rejected
	MaxWaitInQueue time.Duration `json:"maxWaitInQueue,omitempty"`
}

type ImageRegistryConfigRoute struct {

	// Name of the route to be created
	Name string `json:"name"`

	// Hostname for the route
	Hostname string `json:"hostname"`

	// SecretName points to secret containing the certificates to be used
	// by the route
	SecretName string `json:"secretName"`
}
