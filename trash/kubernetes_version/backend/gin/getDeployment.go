package gin

import (
	"controller/k8s"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/apps/v1"
	jobV1 "k8s.io/api/batch/v1"
	coreV1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	"net/http"
	"strings"
)

func getDeployment(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")

	if namespace == "" || name == "" {
		c.String(http.StatusBadRequest, errors.New("namespace or name is not defined").Error())
		return
	}

	progress := gin.H{
		"namespace": namespace,
		"name":      name,
		"error":     false,
		"status":    gin.H{},
	}

	ns, err := k8s.GetNamespace(namespace)
	if err != nil {
		progress["error"] = true
		progress["status"].(gin.H)["namespace"] = err.Error()

		c.JSON(http.StatusForbidden, progress)
		return
	}

	// NAMESPACE - DONE
	progress["status"].(gin.H)["namespace"] = gin.H{
		"name":               ns.GetObjectMeta().GetName(),
		"creation_timestamp": ns.GetCreationTimestamp(),
	}

	// PERSISTENT VOLUME CLAIM - STEAM - DONE
	pvcSteam, err := k8s.GetPersistentVolumeClaim(namespace, name+"-steam")
	if err != nil {
		progress["error"] = true
		progress["status"].(gin.H)["pvc_steam"] = err.Error()
	} else {
		progress["status"].(gin.H)["pvc_steam"] = getPVCObject(pvcSteam)
	}

	// PERSISTENT VOLUME CLAIM - GAMESERVER - DONE
	pvcGameserver, err := k8s.GetPersistentVolumeClaim(namespace, name+"-gameserver")
	if err != nil {
		progress["error"] = true
		progress["status"].(gin.H)["pvc_gameserver"] = err.Error()
	} else {
		progress["status"].(gin.H)["pvc_gameserver"] = getPVCObject(pvcGameserver)
	}

	// JOB
	job, jobList, err := k8s.GetJob(namespace, "steam-init")
	if err != nil {
		progress["error"] = true
		progress["status"].(gin.H)["job"] = err.Error()
	} else {
		progress["status"].(gin.H)["job"] = getJobObject(job, jobList)
	}

	// DEPLOYMENT- DONE
	deployment, err := k8s.GetDepoyments(namespace, name)
	if err != nil {
		progress["error"] = true
		progress["status"].(gin.H)["deployment"] = err.Error()
	} else {
		progress["status"].(gin.H)["deployment"] = getDeploymentObject(deployment)
	}

	// SERVICE - DONE
	service, err := k8s.GetService(namespace, name)
	if err != nil {
		progress["error"] = true
		progress["status"].(gin.H)["service"] = err.Error()
	} else {
		progress["status"].(gin.H)["service"] = getServiceObject(service)
	}

	// INGRESS - DONE
	//	ingress, ingressList, err := k8s.GetIngress(namespace, name)
	//	if err != nil {
	//		progress["error"] = true
	//		progress["status"].(gin.H)["ingress"] = err.Error()
	//	} else {
	//		progress["status"].(gin.H)["ingress"] = getIngressObject(ingress, ingressList)
	//	}

	c.JSON(http.StatusOK, progress)
	return

}

func getPVCObject(pvc *coreV1.PersistentVolumeClaim) gin.H {
	var accessModeList []string
	for _, accessMode := range pvc.Status.AccessModes {
		accessModeList = append(accessModeList, string(accessMode))
	}
	return gin.H{
		"name":               pvc.GetName(),
		"creation_timestamp": pvc.GetCreationTimestamp(),
		"access_mode":        strings.Join(accessModeList, ","),
		"capacity":           pvc.Status.Capacity.Storage(),
		"phase":              pvc.Status.Phase,
		"mode":               pvc.Spec.VolumeMode,
		"class_name":         pvc.Spec.StorageClassName,
	}
}

func getJobObject(job *jobV1.Job, jobList *jobV1.JobList) gin.H {
	var jobItems []gin.H
	for _, item := range jobList.Items {
		var volumes []gin.H
		for _, volume := range item.Spec.Template.Spec.Volumes {
			volumes = append(volumes, gin.H{
				"name": volume.Name,
				"pvc":  volume.PersistentVolumeClaim.ClaimName,
			})
		}
		var containers []gin.H
		for _, container := range item.Spec.Template.Spec.Containers {
			containers = append(containers, gin.H{
				"name":         container.Name,
				"image":        container.Image,
				"volumeMounts": container.VolumeMounts,
			})
		}
		jobItems = append(jobItems, gin.H{
			"volumes":    volumes,
			"containers": containers,
		})
	}
	return gin.H{
		"name":               job.Name,
		"creation_timestamp": job.CreationTimestamp,
		"completion":         *job.Spec.Completions,
		"status_succeeded":   job.Status.Succeeded,
		"list":               jobItems,
	}
}

func getServiceObject(service *coreV1.Service) gin.H {
	var ports []string
	for _, port := range service.Spec.Ports {
		ports = append(ports, fmt.Sprint(port.Protocol)+" "+fmt.Sprint(port.Port)+":"+fmt.Sprint(port.TargetPort.IntVal))
	}
	return gin.H{
		"name":               service.GetName(),
		"creation_timestamp": service.GetCreationTimestamp(),
		"type":               service.Spec.Type,
		"cluster_ip":         service.Spec.ClusterIP,
		"ports":              ports,
		"selector":           service.Spec.Selector,
	}
}

func getDeploymentObject(deployment *v1.Deployment) gin.H {
	var containers []gin.H
	for _, container := range deployment.Spec.Template.Spec.Containers {
		var ports []string
		for _, port := range container.Ports {
			ports = append(ports, fmt.Sprint(port.ContainerPort)+"/"+fmt.Sprint(port.Protocol))
		}
		containers = append(containers, gin.H{
			"image": container.Image,
			"name":  container.Name,
			"ports": ports,
		})
	}
	return gin.H{
		"name":               deployment.GetName(),
		"creation_timestamp": deployment.GetCreationTimestamp(),
		"ready":              fmt.Sprint(deployment.Status.ReadyReplicas) + "/" + fmt.Sprint(deployment.Status.Replicas),
		"containers":         containers,
		//		"image":              deployment.Spec.Template.Spec.Containers,
	}
}

func getIngressObject(ingress *networkingV1.Ingress, ingressList *networkingV1.IngressList) gin.H {
	var ig []gin.H
	for _, it := range ingressList.Items {
		var rules []gin.H
		for _, itRule := range it.Spec.Rules {
			var paths []gin.H
			for _, itPath := range itRule.HTTP.Paths {
				paths = append(paths, gin.H{
					"path":         itPath.Path,
					"path_prefix":  itPath.PathType,
					"service_port": itPath.Backend.Service.Port.Number,
					"service_name": itPath.Backend.Service.Name,
				})
			}
			rules = append(rules, gin.H{
				"rules": paths,
			})
		}

		var ip []string
		for _, itIP := range it.Status.LoadBalancer.Ingress {
			ip = append(ip, itIP.IP)
		}
		ig = append(ig, gin.H{
			"class": it.Spec.IngressClassName,
			"hosts": rules,
			"ip":    ip,
		})
	}

	return gin.H{
		"name":               ingress.GetObjectMeta().GetName(),
		"creation_timestamp": ingress.GetCreationTimestamp(),
		"list":               ig,
	}
}
