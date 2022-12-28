## Discussion

After going through the code flow of [cluster-image-registry-operator](https://github.com/openshift/cluster-image-registry-operator) to create azure storage account, here are a few findings




* [CreateStorage()](https://github.com/openshift/cluster-image-registry-operator/blob/7e5dd5d9132d0908d750c2c65ac820fbd791f9d1/pkg/storage/azure/azure.go#L630) is the function which creates the azure storage account and storage container when being called by [generator.go](https://github.com/openshift/cluster-image-registry-operator/blob/7e5dd5d9132d0908d750c2c65ac820fbd791f9d1/pkg/resource/generator.go#L153)

* [CreateStorage()](https://github.com/openshift/cluster-image-registry-operator/blob/7e5dd5d9132d0908d750c2c65ac820fbd791f9d1/pkg/storage/azure/azure.go#L630) internally calls [assureStorageAccount()](https://github.com/openshift/cluster-image-registry-operator/blob/7e5dd5d9132d0908d750c2c65ac820fbd791f9d1/pkg/storage/azure/azure.go#L491) function, which makes sure to create a storage account. 
**NOTE**: Regardless of the storage account name was provided by the user or was generated in runtime by us, this function always attempt to create one. 

* [assureStorageAccount()](https://github.com/openshift/cluster-image-registry-operator/blob/7e5dd5d9132d0908d750c2c65ac820fbd791f9d1/pkg/storage/azure/azure.go#L491) finally calls the actual [createStorageAccount()](https://github.com/openshift/cluster-image-registry-operator/blob/7e5dd5d9132d0908d750c2c65ac820fbd791f9d1/pkg/storage/azure/azure.go#L157) function, which then calls the Azure's `storageAccountsClient.Create()` API passing all the required configurations to create a storage account. 

* The above API takes [AccountCreateParameters] object, where we can pass all the essential parameters for the creation. This object has the field **Tags** , which we potentially want to update.  

Here is the `AccountCreateParameters` stucture details: 
```go
type AccountCreateParameters struct {
	Sku *Sku `json:"sku,omitempty"`
	Kind Kind `json:"kind,omitempty"`
	Location *string `json:"location,omitempty"`
	Tags map[string]*string `json:"tags"`
	Identity *Identity `json:"identity,omitempty"`
	*AccountPropertiesCreateParameters `json:"properties,omitempty"`
}
```
**NOTE:** A maximum of 15 tags can be provided for a resource. Each tag must have a key with a length no greater than 128 characters and a value with a length no greater than 256 characters.


### TODO add diagram

## Impementation Thoughts

The function definition of the above three functions look like this:

```go
func (d *driver) CreateStorage(cr *imageregistryv1.Config) error

func (d *driver) assureStorageAccount(cfg *Azure, infra *configv1.Infrastructure) (string, bool, error) 

func (d *driver) createStorageAccount(storageAccountsClient storage.AccountsClient, resourceGroupName, accountName, location, cloudName string) error
```

The azure userTags is available in the status sub resource of infrastructure CR. The access of the userTags from the infrastructure CR would be through `infra.Status.PlatformStatus.Azure.ResourceTags`. 

Now, the `infra` object initially gets created in CreateStorage() function and passed as argument to assureStorageAccount() function. So, `infra.Status.PlatformStatus.Azure.ResourceTags` filed can be accessed inside assureStorageAccount().

A new map `tagset` of type `map[string]*string` can be created inside assureStorageAccount() function, which will hold all the userTags (if provided) and default `"kubernetes.io_cluster.${cluster_id}" = "owned"` tag.

```go
// Code Snippet
```

Then, the final `tagset` will be passed as argument to createStorageAccount() function. The updated function definition would be

```go
func (d *driver) createStorageAccount(storageAccountsClient storage.AccountsClient, resourceGroupName, accountName, location, cloudName string, tagset map[string]*string) error
```


`createStorageAccount()` will then pass the `tagset` variable as value to the `Tags` field of `AccountCreateParameters` object while calling the `storageAccountsClient.Create()` API. This can be added like follwing:
```go
storageAccountsClient.Create(
    d.Context,
    resourceGroupName,
    accountName,
    storage.AccountCreateParameters{
        Kind:     kind,
        Location: to.StringPtr(location),
        Sku: &storage.Sku{
            Name: storage.StandardLRS,
        },
        AccountPropertiesCreateParameters: params,
        Tags:                              tagset,
    },
)
```
The Azure API should takes care the rest to add tags to the storage resource.

