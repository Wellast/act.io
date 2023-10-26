package k8s

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DefaultServerDeployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name: "REPLACE",
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: int32Ptr(1),
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{},
		},
		Template: apiv1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{},
			},
			Spec: apiv1.PodSpec{
				Volumes: []apiv1.Volume{
					{
						Name: "REPLACE",
						VolumeSource: apiv1.VolumeSource{
							PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
								ClaimName: "REPLACE",
							},
						},
					},
				},
				Containers: []apiv1.Container{
					{
						Name:  "REPLACE",
						Image: "joedwards32/cs2",
						VolumeMounts: []apiv1.VolumeMount{
							{
								Name:      "REPLACE",
								MountPath: "/home/steam/cs2-dedicated/",
							},
						},
						Ports: []apiv1.ContainerPort{
							{
								Name:          "http",
								Protocol:      apiv1.ProtocolTCP,
								ContainerPort: 8080,
							},
						},
						Env: []apiv1.EnvVar{
							{Name: "STEAMGUARD", Value: "REPLACE"},
							{Name: "STEAMUSER", Value: "REPLACE"},
							{Name: "STEAMPASS", Value: "REPLACE"},
						},
					},
				},
			},
		},
	},
}

/*	GET DEPLOYMENTS example:
	list, err := k8s.GetDepoyments("default", "cs2server")
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
*/
func GetDepoyments(namespace, name string) (*v1.Deployment, error) {
	deployment, err := k8sClient.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

/*	CREATE DEPLOYMENT example:
result, err := k8s.CreateDeployment("default", k8s.DefaultServerDeployment)
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
*/
func CreateDeployment(namespace string, deployment *appsv1.Deployment) (*v1.Deployment, error) {
	result, err := k8sClient.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

/*	DELETE DEPLOYMENT example:
err = k8s.DeleteDeployment("default", "demo-deployment")
*/
func DeleteDeployment(namespace string, name string) error {
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
