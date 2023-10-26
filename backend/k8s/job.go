package k8s

import (
	"context"
	v1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

var restartPolicyNever = apiv1.ContainerRestartPolicy("Never")
var backoffLimit = int32(0)
var DefaultSteamJob = v1.Job{
	TypeMeta: metav1.TypeMeta{
		APIVersion: "batch/v1",
		Kind:       "Job",
	},
	ObjectMeta: metav1.ObjectMeta{
		Name: "steam-init",
	},
	Spec: v1.JobSpec{
		BackoffLimit: &backoffLimit,
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
						Name:          "steam-init",
						Image:         "cm2network/steamcmd",
						Args:          []string{},
						RestartPolicy: &restartPolicyNever,
						VolumeMounts: []apiv1.VolumeMount{
							{
								Name:      "REPLACE",
								MountPath: "/home/steam/Steam",
							},
						},
					},
				},
			},
		},
	},
}

func GetJob(namespace, name string) (*apiv1.Pod, error) {
	pod, err := k8sClient.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return pod, nil
}

func CreateJob(namespace string, pod *apiv1.Pod) (*apiv1.Pod, error) {
	result, err := k8sClient.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteJob(namespace string, name string) error {
	deletePolicy := metav1.DeletePropagationForeground
	err := k8sClient.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err
	}
	return nil
}

func WatchJob(namespace string, name string) (<-chan watch.Event, error) {
	w, err := k8sClient.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{
		LabelSelector: name,
	})
	if err != nil {
		return nil, err
	}
	return w.ResultChan(), nil
}
