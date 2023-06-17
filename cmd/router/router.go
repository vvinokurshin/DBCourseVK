package router

import (
	"github.com/labstack/echo/v4"
	forumDelivery "github.com/vvinokurshin/DBCourseVK/internal/forum/delivery"
	postDelivery "github.com/vvinokurshin/DBCourseVK/internal/post/delivery"
	serviceDelivery "github.com/vvinokurshin/DBCourseVK/internal/service/delivery"
	threadDelivery "github.com/vvinokurshin/DBCourseVK/internal/thread/delivery"
	userDelivery "github.com/vvinokurshin/DBCourseVK/internal/user/delivery"
)

func AddRoutes(e *echo.Echo, userH userDelivery.DeliveryI, forumH forumDelivery.DeliveryI, threadH threadDelivery.DeliveryI,
	postH postDelivery.DeliveryI, serviceH serviceDelivery.DeliveryI) {
	// user
	e.POST("/api/user/:nickname/create", userH.CreateUser)
	e.GET("/api/user/:nickname/profile", userH.GetUser)
	e.POST("/api/user/:nickname/profile", userH.UpdateUser)

	// forum
	e.POST("/api/forum/create", forumH.CreateForum)
	e.GET("/api/forum/:slug/details", forumH.GetForum)
	e.GET("/api/forum/:slug/users", forumH.GetUsersByForum)

	// thread
	e.POST("/api/forum/:slug/create", threadH.CreateThread)
	e.GET("/api/forum/:slug/threads", threadH.GetThreadsByForum)
	e.POST("/api/thread/:slug_or_id/vote", threadH.CreateVote)
	e.GET("/api/thread/:slug_or_id/details", threadH.GetThread)
	e.POST("/api/thread/:slug_or_id/details", threadH.UpdateThread)

	// post
	e.POST("/api/thread/:slug_or_id/create", postH.CreatePosts)
	e.GET("/api/thread/:slug_or_id/posts", postH.GetPostsByThread)
	e.GET("/api/post/:id/details", postH.GetPost)
	e.POST("/api/post/:id/details", postH.UpdatePost)

	// service
	e.POST("/api/service/clear", serviceH.ClearAll)
	e.GET("api/service/status", serviceH.GetStatus)
}
