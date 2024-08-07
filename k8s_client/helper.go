package k8s_client

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/rest"
)

func NewRestClient(restConfig *rest.Config, gv schema.GroupVersion) (rest.Interface, error) {
	restConfig.ContentConfig = resource.UnstructuredPlusDefaultContentConfig()
	restConfig.GroupVersion = &gv

	if len(gv.Group) == 0 {
		restConfig.APIPath = "/api"
	} else {
		restConfig.APIPath = "/apis"
	}

	return rest.RESTClientFor(restConfig)
}

func Apply(jsonData []byte, restConfig *rest.Config, mapper meta.RESTMapper) error {

	return nil
}

func Describe(restConfig *rest.Config, gvk schema.GroupVersionKind, namespace, name string) (string, error) {

	return "", nil
}

func Get(jsonData string, restConfig *rest.Config, mapper meta.RESTMapper) ([]*metav1.Table, error) {

	return nil, nil

}

func Delete(jsonData string, restConfig *rest.Config, mapper meta.RESTMapper) error {

	return nil
}

func setDefaultNamespaceIfScopedAndNoneSet(unstructured *unstructured.Unstructured, helper *resource.Helper) {
	namespace := unstructured.GetNamespace()

	if helper.NamespaceScoped && namespace == "" {

		namespace = "default"
		unstructured.SetNamespace(namespace)
	}
}

// 处理空的 metav1.table
func handlerUndefinedTable() *metav1.Table {
	emptyColumn := []metav1.TableColumnDefinition{
		{Name: "Name", Type: "string", Format: "name", Description: ""},
		{Name: "Description", Type: "string", Format: "description", Description: ""},
	}
	table := &metav1.Table{
		ColumnDefinitions: emptyColumn,
	}

	table.Kind = "Table"
	table.APIVersion = "meta.k8s.io/v1"

	emptyRow := metav1.TableRow{
		//Object: runtime.RawExtension{Object: &runtime.Unknown{
		//	TypeMeta: runtime.TypeMeta{APIVersion: "unknown/v1", Kind: "Unknown"},
		//}},
	}
	emptyRow.Cells = []interface{}{
		"undefined", "该资源可能已经删除",
	}

	table.Rows = []metav1.TableRow{emptyRow}
	return table
}
