package k8s

import (
	"context"
	"errors"
	apiv1 "k8s.io/api/core/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func CreateService(namespace string) (*coreV1.Service, error) {
	if namespace == "" {
		return nil, errors.New("No namespace defined")
	}

	service2Create := &apiv1.Service{
		metav1.TypeMeta{Kind: "Service", APIVersion: "v1"},
		metav1.ObjectMeta{
			Name: namespace,
		},
		coreV1.ServiceSpec{
			Selector: map[string]string{
				"app":   "cs2server",
				"owner": namespace,
			},
			Ports: []coreV1.ServicePort{
				{
					Name:       "cs2server",
					Port:       8080,
					Protocol:   apiv1.ProtocolTCP,
					TargetPort: intstr.FromInt32(8080),
				},
			},
		}, coreV1.ServiceStatus{},
	}

	service, err := k8sClient.CoreV1().Services(namespace).Create(
		context.TODO(), service2Create, metav1.CreateOptions{},
	)
	if err != nil {
		return nil, err
	}

	return service, nil
}
