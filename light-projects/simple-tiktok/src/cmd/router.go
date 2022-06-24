package main

import (
	"github.com/gin-gonic/gin"
	"simple-tiktok/src/controller"
)

func handle(r *gin.Engine) {

	// 测试
	r.Any("/ping", controller.Ping)

	basic := r.Group("/douyin")

	// 视频流
	feed := basic.Group("/feed")
	feed.GET("/")

	// 用户相关
	userGroup := basic.Group("/user")
	{
		// 获取用户登录信息
		userGroup.GET("/")

		// 新用户注册
		userGroup.POST("/register/")

		// 用户登录
		userGroup.POST("/login/")
	}

	// 视频投稿相关
	publishGroup := basic.Group("/publish")
	{
		// 用户上传视频
		publishGroup.POST("/action/")

		// 直接列出用户投稿过的所有视频
		publishGroup.GET("/list/")
	}

	// 点赞相关
	favoriteGroup := basic.Group("/favorite")
	{
		// 点赞 取消点赞
		favoriteGroup.POST("/action/")

		// 获取点赞列表
		favoriteGroup.GET("/list/")
	}

	// 评论相关
	commentGroup := basic.Group("/comment")
	{
		// 评论
		commentGroup.POST("/action/")

		// 倒叙查看评论
		commentGroup.GET("/list/")
	}

	// 用户间关系操作 如关注 获取关注列表
	relationGroup := basic.Group("/relation")
	{
		// 对指定用户 关注 取关
		relationGroup.POST("/action/")

		// 获取用户的关注列表
		relationGroup.GET("/follow/list/")

		// 获取用户的粉丝列表
		relationGroup.GET("/follower/list/")
	}
}
