package k8s

import (
	"context"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetIngress(namespace string, name string) (*networkingv1.Ingress, *networkingv1.IngressList, error) {
	ingress, err := k8sClient.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}
	list, _ := k8sClient.NetworkingV1().Ingresses(namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: "name=" + name + ",namespace=" + namespace,
	})
	if err != nil {
		return nil, nil, err
	}
	return ingress, list, nil
}

func CreateIngress(namespace, name string) (*networkingv1.Ingress, error) {
	ptp := networkingv1.PathTypePrefix
	ingress2Create := networkingv1.Ingress{
		v1.TypeMeta{APIVersion: "networking.k8s.io/v1", Kind: "Ingress"},
		v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"namespace": namespace,
				"name":      name,
			},
		},
		networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				networkingv1.IngressRule{
					namespace + "." + name + ".actio.live", networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								networkingv1.HTTPIngressPath{
									PathType: &ptp,
									Path:     "/",
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: name,
											Port: networkingv1.ServiceBackendPort{Number: 8080},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		networkingv1.IngressStatus{},
	}
	ingress, err := k8sClient.NetworkingV1().Ingresses(namespace).Create(
		context.TODO(),
		&ingress2Create,
		v1.CreateOptions{},
	)

	if err != nil {
		return nil, err
	}

	return ingress, nil
}

func DeleteIngress(namespace, name string) error {
	err := k8sClient.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}
