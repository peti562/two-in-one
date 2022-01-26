package main

import (
	"two-in-one/controller"

	dic "github.com/DrBenton/minidic"
	"github.com/labstack/echo/v4"
)

func createEndpoints(e *echo.Echo, container dic.Container) {

	commentController := container.Get("Controller.Comment").(*controller.CommentController)
	fibonacciController := container.Get("Controller.Fibonacci").(*controller.FibonacciController)

	// todo 	ideally a middleware here would check get the userId from the
	// todo 	auth token and just call comments/, but now I simplified it to prevent overcomplicating
	e.GET("comments/:userId", commentController.GetCommentByUserId)

	commentGroup := e.Group("/comment")
	commentGroup.GET("/:commentId", commentController.GetCommentById)
	commentGroup.GET("/:commentId", commentController.UpdateComment)
	commentGroup.GET("/:commentId", commentController.DeleteComment)

	commentGroup.GET("/create", commentController.CreateComment)

	e.GET("fibonacci/:n", fibonacciController.Get)
}
