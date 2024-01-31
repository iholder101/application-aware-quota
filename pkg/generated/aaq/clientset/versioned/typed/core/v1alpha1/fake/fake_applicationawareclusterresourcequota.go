/*
Copyright 2023 The AAQ Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "kubevirt.io/applications-aware-quota/staging/src/kubevirt.io/applications-aware-quota-api/pkg/apis/core/v1alpha1"
)

// FakeApplicationAwareClusterResourceQuotas implements ApplicationAwareClusterResourceQuotaInterface
type FakeApplicationAwareClusterResourceQuotas struct {
	Fake *FakeAaqV1alpha1
}

var applicationawareclusterresourcequotasResource = v1alpha1.SchemeGroupVersion.WithResource("applicationawareclusterresourcequotas")

var applicationawareclusterresourcequotasKind = v1alpha1.SchemeGroupVersion.WithKind("ApplicationAwareClusterResourceQuota")

// Get takes name of the applicationAwareClusterResourceQuota, and returns the corresponding applicationAwareClusterResourceQuota object, and an error if there is any.
func (c *FakeApplicationAwareClusterResourceQuotas) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ApplicationAwareClusterResourceQuota, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(applicationawareclusterresourcequotasResource, name), &v1alpha1.ApplicationAwareClusterResourceQuota{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationAwareClusterResourceQuota), err
}

// List takes label and field selectors, and returns the list of ApplicationAwareClusterResourceQuotas that match those selectors.
func (c *FakeApplicationAwareClusterResourceQuotas) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ApplicationAwareClusterResourceQuotaList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(applicationawareclusterresourcequotasResource, applicationawareclusterresourcequotasKind, opts), &v1alpha1.ApplicationAwareClusterResourceQuotaList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ApplicationAwareClusterResourceQuotaList{ListMeta: obj.(*v1alpha1.ApplicationAwareClusterResourceQuotaList).ListMeta}
	for _, item := range obj.(*v1alpha1.ApplicationAwareClusterResourceQuotaList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested applicationAwareClusterResourceQuotas.
func (c *FakeApplicationAwareClusterResourceQuotas) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(applicationawareclusterresourcequotasResource, opts))
}

// Create takes the representation of a applicationAwareClusterResourceQuota and creates it.  Returns the server's representation of the applicationAwareClusterResourceQuota, and an error, if there is any.
func (c *FakeApplicationAwareClusterResourceQuotas) Create(ctx context.Context, applicationAwareClusterResourceQuota *v1alpha1.ApplicationAwareClusterResourceQuota, opts v1.CreateOptions) (result *v1alpha1.ApplicationAwareClusterResourceQuota, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(applicationawareclusterresourcequotasResource, applicationAwareClusterResourceQuota), &v1alpha1.ApplicationAwareClusterResourceQuota{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationAwareClusterResourceQuota), err
}

// Update takes the representation of a applicationAwareClusterResourceQuota and updates it. Returns the server's representation of the applicationAwareClusterResourceQuota, and an error, if there is any.
func (c *FakeApplicationAwareClusterResourceQuotas) Update(ctx context.Context, applicationAwareClusterResourceQuota *v1alpha1.ApplicationAwareClusterResourceQuota, opts v1.UpdateOptions) (result *v1alpha1.ApplicationAwareClusterResourceQuota, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(applicationawareclusterresourcequotasResource, applicationAwareClusterResourceQuota), &v1alpha1.ApplicationAwareClusterResourceQuota{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationAwareClusterResourceQuota), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeApplicationAwareClusterResourceQuotas) UpdateStatus(ctx context.Context, applicationAwareClusterResourceQuota *v1alpha1.ApplicationAwareClusterResourceQuota, opts v1.UpdateOptions) (*v1alpha1.ApplicationAwareClusterResourceQuota, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(applicationawareclusterresourcequotasResource, "status", applicationAwareClusterResourceQuota), &v1alpha1.ApplicationAwareClusterResourceQuota{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationAwareClusterResourceQuota), err
}

// Delete takes name of the applicationAwareClusterResourceQuota and deletes it. Returns an error if one occurs.
func (c *FakeApplicationAwareClusterResourceQuotas) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(applicationawareclusterresourcequotasResource, name, opts), &v1alpha1.ApplicationAwareClusterResourceQuota{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeApplicationAwareClusterResourceQuotas) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(applicationawareclusterresourcequotasResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ApplicationAwareClusterResourceQuotaList{})
	return err
}

// Patch applies the patch and returns the patched applicationAwareClusterResourceQuota.
func (c *FakeApplicationAwareClusterResourceQuotas) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ApplicationAwareClusterResourceQuota, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(applicationawareclusterresourcequotasResource, name, pt, data, subresources...), &v1alpha1.ApplicationAwareClusterResourceQuota{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ApplicationAwareClusterResourceQuota), err
}
