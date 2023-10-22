package k8s

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DefaultDeployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name: "demo-deployment",
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: int32Ptr(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				"app": "cs2server",
			},
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": "cs2server",
				},
			},
			Spec: apiv1.PodSpec{
				Containers: []apiv1.Container{
					{
						Name:  "cs2server",
						Image: "kicbase/echo-server:1.0",
						Ports: []apiv1.ContainerPort{
							{
								Name:          "http",
								Protocol:      apiv1.ProtocolTCP,
								ContainerPort: 80,
							},
						},
					},
				},
			},
		},
	},
}

/*	GET DEPLOYMENTS example:
	list, err := k8s.GetDepoyments("default")
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
*/
func GetDepoyments(namespace string) (*v1.DeploymentList, error) {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}

	deploymentsClient := k8sClient.AppsV1().Deployments(namespace)

	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return list, nil
}

/*	CREATE DEPLOYMENT example:
result, err := k8s.CreateDeployment("default", k8s.DefaultDeployment)
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
*/
func CreateDeployment(namespace string, deployment *appsv1.Deployment) (*v1.Deployment, error) {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}

	deployments := k8sClient.AppsV1().Deployments(namespace)
	result, err := deployments.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

/*	DELETE DEPLOYMENT example:
err = k8s.DeleteDeployment("default", "demo-deployment")
*/
func DeleteDeployment(namespace string, name string) error {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}

	deletePolicy := metav1.DeletePropagationForeground
	err := k8sClient.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err
	}
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
