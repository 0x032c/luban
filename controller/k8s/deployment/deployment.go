package deployment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"pigs/common"
	"pigs/controller"
	"pigs/controller/response"
	"pigs/models/k8s"
	"pigs/pkg/k8s/Init"
	k8scommon "pigs/pkg/k8s/common"
	"pigs/pkg/k8s/deployment"
	"pigs/pkg/k8s/parser"
	"pigs/pkg/k8s/service"
	"strings"
)

func GetDeploymentList(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	dataSelect := parser.ParseDataSelectPathParameter(c)
	nameSpace := c.Query("namespace")
	var p = &k8scommon.NamespaceQuery{
		Namespaces: strings.Split(nameSpace, ","),
	}

	data, err := deployment.GetDeploymentList(client, p, dataSelect)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}

	response.OkWithData(data, c)
	return
}

func DeleteCollectionDeployment(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	var deploymentData []k8s.RemoveDeploymentData

	err = controller.CheckParams(c, &deploymentData)
	if err != nil {
		response.FailWithMessage(http.StatusNotFound, err.Error(), c)
		return
	}

	err = deployment.DeleteCollectionDeployment(client, deploymentData)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}

	response.Ok(c)
	return
}

func DeleteDeployment(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}
	var deploymentData k8s.RemoveDeploymentToServiceData

	err2 := controller.CheckParams(c, &deploymentData)
	if err2 != nil {
		response.FailWithMessage(http.StatusNotFound, err2.Error(), c)
		return
	}

	err = deployment.DeleteDeployment(client, deploymentData.Namespace, deploymentData.DeploymentName)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	common.LOG.Info(fmt.Sprintf("deployment：%v, 已删除", deploymentData.DeploymentName))

	if deploymentData.IsDeleteService {
		serviceErr := service.DeleteService(client, deploymentData.Namespace, deploymentData.ServiceName)

		if serviceErr != nil {
			common.LOG.Error("删除相关Service出错", zap.Any("err: ", serviceErr))
			response.FailWithMessage(response.InternalServerError, err.Error(), c)
			return
		}
	}
	response.Ok(c)
	return
}

func ScaleDeployment(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}

	var scaleData k8s.ScaleDeployment

	err2 := controller.CheckParams(c, &scaleData)
	if err2 != nil {
		response.FailWithMessage(http.StatusNotFound, err2.Error(), c)
		return
	}

	err = deployment.ScaleDeployment(client, scaleData.Namespace, scaleData.DeploymentName, scaleData.ScaleNumber)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}

	response.Ok(c)
	return
}

func RestartDeploymentController(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}
	var restartDeployment k8s.RestartDeployment
	err2 := controller.CheckParams(c, &restartDeployment)
	if err2 != nil {
		response.FailWithMessage(response.ParamError, err2.Error(), c)
		return
	}
	err3 := deployment.RestartDeployment(client, restartDeployment.DeploymentName, restartDeployment.Namespace)
	if err3 != nil {
		response.FailWithMessage(response.InternalServerError, err3.Error(), c)
		return
	}
	response.Ok(c)
	return

}

func GetDeploymentToServiceController(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}

	var Deployment k8s.RestartDeployment
	err2 := controller.CheckParams(c, &Deployment)
	if err2 != nil {
		response.FailWithMessage(response.ParamError, err2.Error(), c)
		return
	}

	data, err := service.GetDeploymentToService(client, Deployment.Namespace, Deployment.DeploymentName)
	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), c)
		return
	}
	response.OkWithData(data, c)
	return
}

func DetailDeploymentController(c *gin.Context) {

	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}
	namespace := c.Query("namespace")
	name := c.Query("name")

	if name == "" || namespace == "" {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}

	data, err := deployment.GetDeploymentDetail(client, namespace, name)

	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}
