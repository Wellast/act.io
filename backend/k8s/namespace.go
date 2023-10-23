package k8s

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1 "k8s.io/api/core/v1"
)

func GetNamespace(namespace string) (*coreV1.Namespace, error) {
	ns, err := k8sClient.CoreV1().Namespaces().Get(context.TODO(), namespace, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func CreateNamespace(namespace string) (*coreV1.Namespace, error) {
	ns, err := k8sClient.CoreV1().Namespaces().Create(
		context.TODO(),
		&coreV1.Namespace{
			v1.TypeMeta{},
			v1.ObjectMeta{Name: namespace},
			coreV1.NamespaceSpec{},
			coreV1.NamespaceStatus{},
		},
		v1.CreateOptions{},
	)
	if err != nil {
		return nil, err
	}

	return ns, nil
}

func DeleteNamespace(namespace string) error {
	err := k8sClient.CoreV1().Namespaces().Delete(context.TODO(), namespace, v1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
