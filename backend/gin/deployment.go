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
	name := c.Query("name")
	if namespace == "" || name == "" {
		c.String(http.StatusForbidden, errors.New("namespace or name is not defined").Error())
		return
	}

	deployment, err := k8s.GetDepoyments(namespace, name)
	if err != nil {
		c.String(http.StatusForbidden, err.Error())
		return
	}
	c.JSON(http.StatusOK, deployment)
	return
}

func createDeployment(c *gin.Context) {
	namespace := c.PostForm("namespace")
	name := c.PostForm("name")

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

	//	NAMESPACE
	_, err := k8s.CreateNamespace(namespace)
	if err != nil {
		if err.Error() != "namespaces \""+namespace+"\" already exists" {
			progress["status"].(gin.H)["namespace"] = err.Error()
			progress["error"] = true
			c.JSON(http.StatusForbidden, progress)
			return
		}
	}

	//	PERSISTENT VOLUME CLAIM
	_, err = k8s.CreatePersistentVolumeClaim(namespace, name)
	if err != nil {
		progress["status"].(gin.H)["pvc"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["pvc"] = "success"
	}

	//	POD - Steam init
	customPod := k8s.DefaultSteamJob
	customPod.Spec.Volumes[0].Name = name
	customPod.Spec.Volumes[0].VolumeSource.PersistentVolumeClaim.ClaimName = name
	customPod.Spec.Containers[0].VolumeMounts[0].Name = name
	customPod.Labels = map[string]string{"name": name}
	customPod.Spec.Containers[0].Args = []string{
		"bash", "/home/steam/steamcmd/steamcmd.sh", "+login",
		appConfig.Steam.Username, appConfig.Steam.Password,
		"+quit",
	}
	_, err = k8s.CreatePod(namespace, customPod)
	if err != nil {
		progress["status"].(gin.H)["pod"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["pod"] = "success"
	}
	watchChan, err := k8s.WatchPod(namespace, name)
	if err != nil {
		panic(err)
	}
	for event := range watchChan {
		fmt.Println(event)
	}

	//	DEPLOYMENT
	customDeployment := k8s.DefaultServerDeployment
	customDeployment.ObjectMeta.Name = name
	customDeployment.Spec.Selector.MatchLabels["app"] = name
	customDeployment.Spec.Selector.MatchLabels["namespace"] = namespace
	customDeployment.Spec.Template.ObjectMeta.Labels["app"] = name
	customDeployment.Spec.Template.ObjectMeta.Labels["namespace"] = namespace
	customDeployment.Spec.Template.Spec.Volumes[0].Name = name
	customDeployment.Spec.Template.Spec.Volumes[0].VolumeSource.PersistentVolumeClaim.ClaimName = name
	customDeployment.Spec.Template.Spec.Containers[0].Name = name
	customDeployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name = name
	pod2conf := map[string]string{
		"STEAMGUARD": appConfig.Steam.Guardcode,
		"STEAMUSER":  appConfig.Steam.Username,
		"STEAMPASS":  appConfig.Steam.Password,
	}
	for idx, _ := range customDeployment.Spec.Template.Spec.Containers[0].Env {
		customDeployment.Spec.Template.Spec.Containers[0].Env[idx].Value = pod2conf[customDeployment.Spec.Template.Spec.Containers[0].Env[idx].Name]
	}
	_, err = k8s.CreateDeployment(namespace, customDeployment)
	if err != nil {
		progress["status"].(gin.H)["deployment"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["deployment"] = "success"
	}

	//	SERVICE
	_, err = k8s.CreateService(namespace, name)
	if err != nil {
		progress["status"].(gin.H)["service"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["service"] = "success"
	}

	//	INGRESS
	_, err = k8s.CreateIngress(namespace, name)
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

	err := k8s.DeletePod(namespace, "steam-init")
	if err != nil {
		progress["status"].(gin.H)["pod"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["pod"] = "success"
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
