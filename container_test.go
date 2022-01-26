package main

import (
	"testing"

	"two-in-one/controller"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Test_buildContainer(t *testing.T) {
	_ = godotenv.Load()

	// Create a mock DB
	gormDb, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// Close our DB
	defer closeConnection(gormDb)

	// Build our container
	container := buildContainer(gormDb)

	// Get the workers
	commentController := container.Get("Controller.Comment")
	// Controllers or workers
	assert.IsType(t, &controller.CommentController{}, commentController)
}
