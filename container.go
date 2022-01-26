package main

import (
	"two-in-one/controller"

	dic "github.com/DrBenton/minidic"
	"gorm.io/gorm"
)

func buildContainer(gormDb *gorm.DB) dic.Container {

	// Create our container
	container := dic.NewContainer()

	container.Add(dic.NewInjection("Controller.Comment", func(c dic.Container) *controller.CommentController {
		return controller.NewCommentController(gormDb)
	}))
	container.Add(dic.NewInjection("Controller.Fibonacci", func(c dic.Container) *controller.FibonacciController {
		return controller.NewFibonacciController()
	}))

	return container
}
