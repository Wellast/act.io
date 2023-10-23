package gin

import (
	"controller/k8s"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getDeployment(c *gin.Context) {
	namespace := c.Query("namespace")
	_ = c.Query("name")
	if namespace == "" {
		c.String(http.StatusInternalServerError, errors.New("namespace or name is not defined").Error())
		return
	}

	list, err := k8s.GetDepoyments(namespace)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, list.Items)
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
	return
}

func createDeployment(c *gin.Context) {
	namespace := c.PostForm("namespace")
	name := c.PostForm("name")

	if namespace == "" || name == "" {
		c.String(http.StatusInternalServerError, errors.New("namespace or name is not defined").Error())
		return
	}

	progress := gin.H{
		"namespace": namespace,
		"name":      name,
		"error":     false,
		"deleted":   gin.H{},
	}

	//	NAMESPACE
	_, err := k8s.CreateNamespace(namespace)
	if err != nil {
		if err.Error() != "namespaces \""+namespace+"\" already exists" {
			progress["deleted"].(gin.H)["namespace"] = err.Error()
			progress["error"] = true
			c.JSON(http.StatusForbidden, progress)
			return
		}
	}

	//	DEPLOYMENT
	customDeployment := k8s.DefaultDeployment
	customDeployment.ObjectMeta.Name = name
	customDeployment.Spec.Selector.MatchLabels["app"] = name
	customDeployment.Spec.Selector.MatchLabels["namespace"] = namespace
	customDeployment.Spec.Template.ObjectMeta.Labels["app"] = name
	customDeployment.Spec.Template.ObjectMeta.Labels["namespace"] = namespace
	_, err = k8s.CreateDeployment(namespace, customDeployment)
	if err != nil {
		progress["deleted"].(gin.H)["deployment"] = err.Error()
		progress["error"] = true
	} else {
		progress["deleted"].(gin.H)["deployment"] = "success"
	}

	//	SERVICE
	_, err = k8s.CreateService(namespace, name)
	if err != nil {
		progress["deleted"].(gin.H)["service"] = err.Error()
		progress["error"] = true
	} else {
		progress["deleted"].(gin.H)["service"] = "success"
	}

	//	INGRESS
	_, err = k8s.CreateIngress(namespace, name)
	if err != nil {
		progress["deleted"].(gin.H)["ingress"] = err.Error()
		progress["error"] = true
	} else {
		progress["deleted"].(gin.H)["ingress"] = "success"
	}

	if progress["error"] == false {
		c.JSON(http.StatusOK, progress)
	} else {
		c.JSON(http.StatusForbidden, progress)
	}
	return
}

func deleteDeployment(c *gin.Context) {
	namespace := c.Query("namespace")
	name := c.Query("name")

	if namespace == "" || name == "" {
		c.String(http.StatusInternalServerError, errors.New("namespace or name is not defined").Error())
		return
	}

	progress := gin.H{
		"namespace": namespace,
		"name":      name,
		"error":     false,
		"deleted":   gin.H{},
	}

	err := k8s.DeleteDeployment(namespace, name)
	if err != nil {
		progress["deleted"].(gin.H)["deployment"] = err.Error()
		progress["error"] = true
	} else {
		progress["deleted"].(gin.H)["deployment"] = "success"
	}

	err = k8s.DeleteService(namespace, name)
	if err != nil {
		progress["deleted"].(gin.H)["service"] = err.Error()
		progress["error"] = true
	} else {
		progress["deleted"].(gin.H)["service"] = "success"
	}

	err = k8s.DeleteIngress(namespace, name)
	if err != nil {
		progress["deleted"].(gin.H)["ingress"] = err.Error()
		progress["error"] = true
	} else {
		progress["deleted"].(gin.H)["ingress"] = "success"
	}

	if progress["error"] == false {
		c.JSON(http.StatusOK, progress)
	} else {
		c.JSON(http.StatusForbidden, progress)
	}
	return
}
