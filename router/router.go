package router

import (
	"codeup.aliyun.com/xhey/server/point-manage/controller"
	"codeup.aliyun.com/xhey/server/point-manage/middle"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//定义router 结构体
type router struct{}

//实例化route对象,可使用该对象点出首字母大写的方法(挎包调用)
var Router router

//初始化路由
func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/ping", controller.Ping.Ping)
	//router.POST("/next/point/login", controller.Ldap.Login)
	//临时获取token接口,测试使用
	router.POST("/next/point/token", controller.Token.GetToken)
	router.POST("/next/point/parsetoken", controller.Token.ParseToken)
	//回调接口,不做jwt认证
	router.POST("/next/point/oss/callback", middle.Oplogs(), controller.Oss.Callback)
	//prometheus
	router.GET("/metrics", middle.PromHandler(promhttp.Handler()))
	apigroup := router.Group("/next/point")
	{
		apigroup.GET("/oss/token", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Oss.GetOssToken)

		apigroup.POST("/version/create", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Version.Create)
		apigroup.GET("/version/list", controller.Version.List)
		apigroup.POST("/version/info", controller.Version.VersionInfo)
		apigroup.POST("/version/update", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Version.Update)
		apigroup.POST("/version/release", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Version.UpdateToRelease)

		apigroup.GET("/class/list", controller.Classification.List)
		apigroup.POST("/class/create", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Classification.Create)
		apigroup.POST("/class/update", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Classification.Update)

		//apigroup.POST("/event/create", middle.Roleif(), controller.Event.Create)
		apigroup.POST("/event/update", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Event.Update)
		apigroup.POST("/event/removeDemand", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Event.Remove)
		apigroup.POST("/event/listByDemand", controller.Event.ListByDemand)
		apigroup.POST("/event/listByKeyword", controller.Event.SearchByKeyword)
		apigroup.POST("/event/listByVersionCode", controller.Event.ListByVersionCode)
		apigroup.POST("/event/listByCategoryId", controller.Event.ListByCategoryId)
		//简明的事件列表
		apigroup.POST("/event/conciseByCategoryId", controller.Event.ConciseListByClassId)
		apigroup.POST("/event/addversion", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Event.EventAddVersion)
		apigroup.POST("/event/delete", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Event.Delete)           //软删除---update
		apigroup.DELETE("/event/realdelete", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Event.RealDelete) //硬删除---delete
		apigroup.POST("/event/info", controller.Event.Info)

		apigroup.POST("/tags/create", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Tags.Create)
		apigroup.DELETE("/tags/delete", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Tags.DeleteById)

		apigroup.POST("/role/adduser", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Ldap.Role)
		apigroup.GET("/role/info", middle.JWTAuth(), middle.Oplogs(), controller.Ldap.Info)

		apigroup.POST("/attribute/delete", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Attribute.Delete)
		apigroup.POST("/attribute/update", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Attribute.Update)
		apigroup.DELETE("/attribute/realdelete", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Attribute.RealDelete)
		apigroup.POST("/attribute/tell", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Attribute.Tell)
		apigroup.POST("/attribute/recommend", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Attribute.Recommend)

		apigroup.POST("/value/delete", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Value.Delete)
		apigroup.POST("/value/update", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Value.Update)
		apigroup.POST("/value/listbykey", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Value.ByAttKeyGetValues)
		apigroup.DELETE("/value/realdelete", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Value.RealDelete)

		apigroup.POST("/demand/update", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Demand.Update)
		apigroup.POST("/create", controller.Point.Create) //客户端埋点上报
		apigroup.POST("/select", controller.Point.Select)

		apigroup.GET("/kernel/event", controller.Kernel.ListKernel)
		apigroup.POST("/kernel/tagevent", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Kernel.TagEventKernalCreate)
		apigroup.POST("/kernel/tageventremove", middle.JWTAuth(), middle.Oplogs(), middle.Roleif(), controller.Kernel.TagEventKernalRemove)

	}
}
