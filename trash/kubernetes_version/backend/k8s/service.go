package k8s

import (
	"context"
	"errors"
	apiv1 "k8s.io/api/core/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func GetService(namespace string, name string) (*coreV1.Service, error) {
	ns, err := k8sClient.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func CreateService(namespace, name string) (*coreV1.Service, error) {
	if namespace == "" || name == "" {
		return nil, errors.New("no namespace or name defined")
	}

	service2Create := &apiv1.Service{
		metav1.TypeMeta{Kind: "Service", APIVersion: "v1"},
		metav1.ObjectMeta{
			Name: name,
		},
		coreV1.ServiceSpec{
			Type: "NodePort",
			Selector: map[string]string{
				"app":       name,
				"namespace": namespace,
			},
			Ports: []coreV1.ServicePort{
				{
					Name:       "tcp",
					Port:       27015,
					Protocol:   apiv1.ProtocolTCP,
					TargetPort: intstr.FromInt32(27015),
				},
				{
					Name:       "udp",
					Port:       27015,
					Protocol:   apiv1.ProtocolUDP,
					TargetPort: intstr.FromInt32(27015),
				},
			},
		}, coreV1.ServiceStatus{},
	}

	service, err := k8sClient.CoreV1().Services(namespace).Create(context.TODO(), service2Create, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return service, nil
}

func DeleteService(namespace, name string) error {
	if namespace == "" || name == "" {
		return errors.New("no namespace or name defined")
	}

	err := k8sClient.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
