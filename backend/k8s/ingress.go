package k8s

import (
	"context"
	"fmt"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
)

func CreateIngress(namespace string) (*networkingv1.Ingress, error) {
	ptp := networkingv1.PathTypePrefix
	ingress2Create := networkingv1.Ingress{
		v1.TypeMeta{APIVersion: "networking.k8s.io/v1", Kind: "Ingress"},
		v1.ObjectMeta{
			Name:        namespace,
			Annotations: map[string]string{},
		},
		networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				networkingv1.IngressRule{
					"", networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								networkingv1.HTTPIngressPath{
									PathType: &ptp,
									Path:     "/",
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: namespace,
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
	b, _ := json.Marshal(ingress2Create)
	fmt.Println(string(b))
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
