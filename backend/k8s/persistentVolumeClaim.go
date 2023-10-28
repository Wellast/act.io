package k8s

import (
	"context"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPersistentVolumeClaim(namespace string, name string) (*coreV1.PersistentVolumeClaim, error) {
	ns, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return ns, nil
}

func CreatePersistentVolumeClaim(namespace string, name string, quantity string) (*coreV1.PersistentVolumeClaim, error) {
	qv := resource.QuantityValue{}
	_ = qv.Set(quantity)
	ns, err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Create(
		context.TODO(),
		&coreV1.PersistentVolumeClaim{
			TypeMeta: v1.TypeMeta{
				APIVersion: "v1",
				Kind:       "PersistentVolumeClaim",
			},
			ObjectMeta: v1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					"name":      name,
					"namespace": namespace,
				},
			},
			Spec: coreV1.PersistentVolumeClaimSpec{
				AccessModes: []coreV1.PersistentVolumeAccessMode{coreV1.ReadWriteOnce},
				Resources: coreV1.ResourceRequirements{
					Requests: coreV1.ResourceList{
						"storage": qv.Quantity,
					},
				},
			},
			Status: coreV1.PersistentVolumeClaimStatus{},
		},
		v1.CreateOptions{},
	)
	if err != nil {
		return nil, err
	}

	return ns, nil
}

func DeletePersistentVolumeClaim(namespace string, name string) error {
	err := k8sClient.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
