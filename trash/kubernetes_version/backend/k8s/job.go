package k8s

import (
	"context"
	v1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

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
				Labels: map[string]string{
					"name": "steam-init",
				},
			},
			Spec: apiv1.PodSpec{
				RestartPolicy: "Never",
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
						Name:  "steam-init",
						Image: "cm2network/steamcmd",
						Args:  []string{},
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

func GetJob(namespace, name string) (*v1.Job, *v1.JobList, error) {
	job, err := k8sClient.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, err
	}

	jobList, err := k8sClient.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}

	return job, jobList, nil
}

func CreateJob(namespace string, job v1.Job) (*v1.Job, error) {
	result, err := k8sClient.BatchV1().Jobs(namespace).Create(context.TODO(), &job, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DeleteJob(namespace string, name string) error {
	err := k8sClient.BatchV1().Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func WatchJob(namespace string, name string) (<-chan watch.Event, error) {
	w, err := k8sClient.BatchV1().Jobs(namespace).Watch(context.TODO(), metav1.ListOptions{
		LabelSelector: name,
	})
	if err != nil {
		return nil, err
	}
	return w.ResultChan(), nil
}
