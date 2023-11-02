package gin

import (
	"controller/k8s"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func createDeployment(c *gin.Context) {
	namespace := c.PostForm("namespace")
	name := c.PostForm("name")
	steamName := c.PostForm("steam_name")
	steamPass := c.PostForm("steam_pass")
	steamGuard := c.PostForm("steam_guard")

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

	//	PERSISTENT VOLUME CLAIM - STEAM
	_, err = k8s.CreatePersistentVolumeClaim(namespace, name+"-steam", "1Gi")
	if err != nil {
		progress["status"].(gin.H)["pvc_steam"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["pvc_steam"] = "success"
	}
	//	PERSISTENT VOLUME CLAIM - GAMESERVER
	_, err = k8s.CreatePersistentVolumeClaim(namespace, name+"-gameserver", "45Gi")
	if err != nil {
		progress["status"].(gin.H)["pvc_gameserver"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["pvc_gameserver"] = "success"
	}

	//	JOB - Steam init
	customJob := k8s.DefaultSteamJob
	customJob.Spec.Template.Spec.Volumes[0].Name = name + "-steam"
	customJob.Spec.Template.Spec.Volumes[0].VolumeSource.PersistentVolumeClaim.ClaimName = name + "-steam"
	customJob.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name = name + "-steam"
	customJob.Labels = map[string]string{"name": name}
	customJob.Spec.Template.Spec.Containers[0].Args = []string{
		"bash", "/home/steam/steamcmd/steamcmd.sh",
	}
	if steamGuard != "" {
		customJob.Spec.Template.Spec.Containers[0].Args = append(customJob.Spec.Template.Spec.Containers[0].Args,
			"+set_steam_guard_code", steamGuard,
		)
	}
	customJob.Spec.Template.Spec.Containers[0].Args = append(customJob.Spec.Template.Spec.Containers[0].Args,
		"+login", steamName, steamPass, "+quit",
	)
	_, err = k8s.CreateJob(namespace, customJob)
	if err != nil {
		progress["status"].(gin.H)["job"] = err.Error()
		progress["error"] = true
	} else {
		progress["status"].(gin.H)["job"] = "success"
	}

	//	DEPLOYMENT
	customDeployment := k8s.DefaultServerDeployment
	customDeployment.ObjectMeta.Name = name
	customDeployment.Spec.Selector.MatchLabels["app"] = name
	customDeployment.Spec.Selector.MatchLabels["namespace"] = namespace
	customDeployment.Spec.Template.ObjectMeta.Labels["app"] = name
	customDeployment.Spec.Template.ObjectMeta.Labels["namespace"] = namespace
	customDeployment.Spec.Template.Spec.Volumes[0].Name = name + "-gameserver"
	customDeployment.Spec.Template.Spec.Volumes[0].VolumeSource.PersistentVolumeClaim.ClaimName = name + "-gameserver"
	customDeployment.Spec.Template.Spec.Containers[0].Name = name
	customDeployment.Spec.Template.Spec.Containers[0].VolumeMounts[0].Name = name + "-gameserver"
	pod2conf := map[string]string{
		"STEAMGUARD": steamGuard,
		"STEAMUSER":  steamName,
		"STEAMPASS":  steamPass,
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
	//	_, err = k8s.CreateIngress(namespace, name)
	//	if err != nil {
	//		progress["status"].(gin.H)["ingress"] = err.Error()
	//		progress["error"] = true
	//	} else {
	//		progress["status"].(gin.H)["ingress"] = "success"
	//	}

	if progress["error"] == false {
		c.JSON(http.StatusOK, progress)
	} else {
		c.JSON(http.StatusForbidden, progress)
	}
	return
}
