package gin

import (
	"controller/k8s"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func deleteDeployment(c *gin.Context) {
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

	err := k8s.DeleteJob(namespace, "steam-init")
	if err != nil {
		progress["status"].(gin.H)["job"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["job"] = "success"
	}

	err = k8s.DeleteDeployment(namespace, name)
	if err != nil {
		progress["status"].(gin.H)["deployment"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["deployment"] = "success"
	}

	err = k8s.DeleteService(namespace, name)
	if err != nil {
		progress["status"].(gin.H)["service"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["service"] = "success"
	}

	err = k8s.DeleteIngress(namespace, name)
	if err != nil {
		progress["status"].(gin.H)["ingress"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["ingress"] = "success"
	}

	if progress["error"] == false {
		c.JSON(http.StatusOK, progress)
	} else {
		c.JSON(http.StatusForbidden, progress)
	}
	return
}
