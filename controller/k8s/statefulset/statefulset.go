package statefulset

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pigs/controller"
	"pigs/controller/response"
	"pigs/models/k8s"
	"pigs/pkg/k8s/Init"
	"pigs/pkg/k8s/parser"
	"pigs/pkg/k8s/statefulset"
)

func GetStatefulSetListController(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	dataSelect := parser.ParseDataSelectPathParameter(c)
	nsQuery := parser.ParseNamespacePathParameter(c)

	data, err := statefulset.GetStatefulSetList(client, nsQuery, dataSelect)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}

	response.OkWithData(data, c)
	return
}

func DeleteCollectionStatefulSetController(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	var statefulSetList []k8s.StatefulSetData

	err = controller.CheckParams(c, &statefulSetList)
	if err != nil {
		response.FailWithMessage(http.StatusNotFound, err.Error(), c)
		return
	}

	err = statefulset.DeleteCollectionStatefulSet(client, statefulSetList)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}

	response.Ok(c)
	return
}

func DeleteStatefulSetController(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}
	var statefulSet k8s.StatefulSetData

	err = controller.CheckParams(c, &statefulSet)
	if err != nil {
		response.FailWithMessage(http.StatusNotFound, err.Error(), c)
		return
	}

	err = statefulset.DeleteStatefulSet(client, statefulSet.Namespace, statefulSet.Name)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}

	response.Ok(c)
	return
}

func RestartStatefulSetController(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}
	var statefulSetData k8s.StatefulSetData
	err2 := controller.CheckParams(c, &statefulSetData)
	if err2 != nil {
		response.FailWithMessage(response.ParamError, err2.Error(), c)
		return
	}
	err3 := statefulset.RestartStatefulSet(client, statefulSetData.Name, statefulSetData.Namespace)
	if err3 != nil {
		response.FailWithMessage(response.InternalServerError, err3.Error(), c)
		return
	}
	response.Ok(c)
	return
}

func ScaleStatefulSetController(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}
	var scaleData k8s.ScaleStatefulSet

	err2 := controller.CheckParams(c, &scaleData)
	if err2 != nil {
		response.FailWithMessage(response.ParamError, err2.Error(), c)
		return
	}

	err = statefulset.ScaleStatefulSet(client, scaleData.Namespace, scaleData.Name, *scaleData.ScaleNumber)
	if err != nil {
		response.FailWithMessage(response.InternalServerError, err.Error(), c)
		return
	}

	response.Ok(c)
	return
}

func DetailStatefulSetController(c *gin.Context) {
	client, err := Init.ClusterID(c)
	if err != nil {
		response.FailWithMessage(response.ParamError, err.Error(), c)
		return
	}
	namespace := c.Query("namespace")
	name := c.Query("name")
	if name == "" || namespace == "" {
		response.FailWithMessage(response.ParamError, "缺少必要的参数", c)
		return
	}
	data, err := statefulset.GetStatefulSetDetail(client, namespace, name)

	if err != nil {
		response.FailWithMessage(response.ERROR, err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}
