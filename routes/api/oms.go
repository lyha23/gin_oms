package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/wejectchen/ginblog/api/v1"
)

func RegisterOMSRouter(r *gin.Engine) {
	/*
		后台管理路由接口
	*/
	auth := r.Group("/php")
	// auth.Use(middleware.JwtToken())
	{
		// 用户模块的路由接口
		auth.GET("/ata_selection", v1.GetATASelection)
	}
}
