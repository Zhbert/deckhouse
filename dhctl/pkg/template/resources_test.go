// Copyright 2021 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package template

import (
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	nsKind = schema.GroupVersionKind{Version: "v1", Kind: "Namespace"}
	cmKind = schema.GroupVersionKind{Version: "v1", Kind: "ConfigMap"}
)

func fromUnstructured(unstructuredObj unstructured.Unstructured, obj interface{}) {
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObj.UnstructuredContent(), obj)
	if err != nil {
		panic(err)
	}
}

func TestResourcesWithoutTemplateData(t *testing.T) {
	t.Run("parse resources from multidocument without template data", func(t *testing.T) {
		resources, err := ParseResources("testdata/resources/without_tmp.yaml", nil)
		require.NoError(t, err)

		require.Contains(t, resources.Items, nsKind)
		require.Len(t, resources.Items[nsKind].Items, 2)

		ns := v1.Namespace{}
		fromUnstructured(resources.Items[nsKind].Items[0], &ns)

		require.Equal(t, ns.Name, "test-ns")

		fromUnstructured(resources.Items[nsKind].Items[1], &ns)
		require.Equal(t, ns.Name, "another-ns")

		require.Contains(t, resources.Items, cmKind)
		require.Len(t, resources.Items[cmKind].Items, 1)

		cm := v1.ConfigMap{}
		fromUnstructured(resources.Items[cmKind].Items[0], &cm)

		require.Equal(t, cm.Namespace, "test-ns")
		require.Equal(t, cm.Name, "some-cm")

		require.Contains(t, cm.Data["key"], "value")
	})
}

func TestResourcesWithTemplateData(t *testing.T) {
	const expectedValueFromCloudData = "id1"
	t.Run("parses template resources and put data in manifests", func(t *testing.T) {
		resources, err := ParseResources("testdata/resources/with_tmp.yaml", map[string]interface{}{
			"cloudDiscovery": map[string]interface{}{
				"networkId": map[string]interface{}{
					"ru-central1-a": expectedValueFromCloudData,
					"ru-central1-b": expectedValueFromCloudData + "1",
					"ru-central1-c": expectedValueFromCloudData + "2",
				},

				"anotherKey": "anotherValue",
			},
		})
		require.NoError(t, err)

		require.Contains(t, resources.Items, cmKind)
		require.Len(t, resources.Items[cmKind].Items, 1)

		cm := v1.ConfigMap{}
		fromUnstructured(resources.Items[cmKind].Items[0], &cm)

		require.Equal(t, cm.Namespace, "test-ns")
		require.Equal(t, cm.Name, "some-cm")

		require.Equal(t, cm.Data["key"], "value")
		require.Equal(t, cm.Data["fromCloudDiscovery"], expectedValueFromCloudData)
		require.Equal(t, cm.Data["sprigFuncAvailable"], "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b")
	})
}

func TestResourcesNotExistsTemplateDataReturnError(t *testing.T) {
	t.Run("returns error if value not found in data", func(t *testing.T) {
		resources, err := ParseResources("testdata/resources/with_tmp.yaml", map[string]interface{}{
			"cloudDiscovery": map[string]interface{}{
				"anotherKey": "anotherValue",
			},
		})

		require.Error(t, err)
		require.Nil(t, resources)
	})
}
