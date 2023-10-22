package gin

import (
	"controller/k8s"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func getDeployment(c *gin.Context) {
	list, err := k8s.GetDepoyments("")
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(200, list.Items)
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
	return
}

func createDeployment(c *gin.Context) {
	owner := c.PostForm("owner")
	deployment, err := k8s.CreateDeployment(owner, k8s.DefaultDeployment)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.String(http.StatusOK, deployment.GetObjectMeta().GetName())
}

func deleteDeployment(c *gin.Context) {
	name := c.Params.ByName("name")
	err := k8s.DeleteDeployment("default", name)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.String(http.StatusOK, "deleted "+name)
}
