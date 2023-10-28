package gin

import (
	"controller/k8s"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getDeployment(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")

	if namespace == "" || name == "" {
		c.String(http.StatusForbidden, errors.New("namespace or name is not defined").Error())
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
	} else {
		progress["status"].(gin.H)["namespace"] = gin.H{
			"name":              ns.GetObjectMeta().GetName(),
			"creationTimestamp": ns.GetCreationTimestamp(),
		}

		pvcSteam, err := k8s.GetPersistentVolumeClaim(namespace, name+"-steam")
		if err != nil {
			progress["error"] = true
			progress["status"].(gin.H)["pvc_steam"] = err.Error()
		} else {
			progress["status"].(gin.H)["pvc_steam"] = pvcSteam
		}
		pvcGameserver, err := k8s.GetPersistentVolumeClaim(namespace, name+"-gameserver")
		if err != nil {
			progress["error"] = true
			progress["status"].(gin.H)["pvc_gameserver"] = err.Error()
		} else {
			progress["status"].(gin.H)["pvc_gameserver"] = pvcGameserver
		}

		deployment, err := k8s.GetDepoyments(namespace, name)
		if err != nil {
			progress["error"] = true
			progress["status"].(gin.H)["deployment"] = err.Error()
		} else {
			progress["status"].(gin.H)["deployment"] = deployment.GetObjectMeta()
		}

		service, err := k8s.GetService(namespace, name)
		if err != nil {
			progress["error"] = true
			progress["status"].(gin.H)["service"] = err.Error()
		} else {
			progress["status"].(gin.H)["service"] = gin.H{
				"name":              service.GetName(),
				"creationTimestamp": service.GetCreationTimestamp(),
			}
		}

		ingress, ingressList, err := k8s.GetIngress(namespace, name)
		if err != nil {
			progress["error"] = true
			progress["status"].(gin.H)["ingress"] = err.Error()
		} else {
			progress["status"].(gin.H)["ingress"] = ingress
			progress["status"].(gin.H)["ingress_list"] = ingressList
		}

		c.JSON(http.StatusOK, progress)
		return
	}
}
