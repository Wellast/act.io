package gin

import (
	"controller/k8s"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func getDeployment(c *gin.Context) {
	namespace := c.Query("namespace")
	if namespace == "" {
		c.String(http.StatusInternalServerError, errors.New("?namespace= is not defined").Error())
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
	var progress []string
	namespace := c.PostForm("namespace")
	if namespace == "" {
		c.String(http.StatusInternalServerError, errors.New("namespace is not defined").Error())
		return
	}

	//	NAMESPACE
	ns, err := k8s.CreateNamespace(namespace)
	if err != nil {
		if err.Error() == "namespaces \""+namespace+"\" already exists" {
			ns, err = k8s.GetNamespace(namespace)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}
	fmt.Println("namespace: ", ns.GetObjectMeta().GetName())
	progress = append(progress, "+namespace")

	//	DEPLOYMENT
	customDeployment := k8s.DefaultDeployment
	customDeployment.ObjectMeta.Name = namespace
	customDeployment.Spec.Template.ObjectMeta.Labels["owner"] = namespace
	deployment, err := k8s.CreateDeployment(namespace, customDeployment)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("deployment", deployment.GetObjectMeta().GetName())
	progress = append(progress, "+deployment")

	//	SERVICE
	service, err := k8s.CreateService(namespace)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println("service:", service.GetObjectMeta().GetName())
	progress = append(progress, "+service")

	//	INGRESS
	ingress, err := k8s.CreateIngress(namespace)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Println(ingress.GetObjectMeta().GetName())
	progress = append(progress, "+ingress")

	c.String(http.StatusOK, strings.Join(progress, " "))
	return
}

func deleteDeployment(c *gin.Context) {
	namespace := c.Params.ByName("namespace")
	err := k8s.DeleteNamespace(namespace)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "namespace "+namespace+" deleted")
}
