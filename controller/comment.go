package controller

import (
	"net/http"
	"strconv"
	"two-in-one/model"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CommentController controller object
type CommentController struct {
	gormDb *gorm.DB
}

func NewCommentController(
	gormDb *gorm.DB,
) *CommentController {

	// Create the base controller instance
	newInstance := &CommentController{}

	newInstance.gormDb = gormDb

	return newInstance
}

func (tc *CommentController) GetCommentById(c echo.Context) error {

	commentId, exception := strconv.Atoi(c.Param("commentId"))
	if exception != nil {
		// should be proper error handling here
		return exception
	}
	var comment model.Comment

	if exception := comment.FindById(tc.gormDb, uint32(commentId)); exception != nil {
		// should be proper error handling here
		return exception
	}

	return c.JSON(http.StatusOK, comment)
}

func (tc *CommentController) GetCommentByUserId(c echo.Context) error {
	userId, exception := strconv.Atoi(c.Param("userId"))
	if exception != nil {
		// should be proper error handling here
		return exception
	}

	var comment model.Comment

	comments, exception := comment.GetByUserId(tc.gormDb, uint32(userId))
	if exception != nil {
		// should be proper error handling here
		return exception
	}

	return c.JSON(http.StatusOK, comments)
}

func (tc *CommentController) UpdateComment(c echo.Context) error {

	commentId, exception := strconv.Atoi(c.Param("commentId"))

	if exception != nil {
		// should be proper error handling here
		return exception
	}

	var comment model.Comment

	comment.Body = "This is a comment"
	comment.UserId = 5

	if exception := comment.UpdateBody(tc.gormDb, uint32(commentId), comment.Body); exception != nil {
		// should be proper error handling here
		return exception
	}

	return c.JSON(http.StatusOK, comment)
}

func (tc *CommentController) CreateComment(c echo.Context) error {

	tx := tc.gormDb.Begin()

	var comment model.Comment

	comment.Body = "This is a comment"
	comment.UserId = 5

	tx.Save(&comment)

	if exception := tx.Commit().Error; exception != nil {
		tx.Rollback()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":   true,
		"commentId": comment.Id,
	})
}

func (tc *CommentController) DeleteComment(c echo.Context) error {

	commentId, exception := strconv.Atoi(c.Param("commentId"))

	if exception != nil {
		// should be proper error handling here
		return exception
	}

	var comment model.Comment

	if exception := comment.Delete(tc.gormDb, uint32(commentId)); exception != nil {
		// should be proper error handling here
		return exception
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
	})
}
