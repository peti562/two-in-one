package controller

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	mocketHelper "two-in-one/helper/mocket"
	structHelper "two-in-one/helper/struct"
	"two-in-one/model"
)

type CommentTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	Context      echo.Context
	Recorder     *httptest.ResponseRecorder
	MocketDb     *gorm.DB
	MocketClient *mocketHelper.Helper
	controller   *CommentController
}

// SetupAllSuite has a SetupSuiteForCard method, which will run before the tests in the suite are run.
func (suite *CommentTestSuite) SetupSuite() {

	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	mocket.Catcher.PanicOnEmptyResponse = true

	mocketDriver := mocketHelper.Open("mocket")
	suite.ctrl = gomock.NewController(suite.T())
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Content-Type", "application/json")

	// Track the response payloads
	suite.Recorder = httptest.NewRecorder()
	suite.Context = echo.New().NewContext(request, suite.Recorder)
	suite.MocketDb, _ = gorm.Open(mocketDriver, &gorm.Config{})
	suite.MocketClient = mocketHelper.New(suite.MocketDb)
	suite.controller = NewCommentController(suite.MocketDb)
}

func (suite *CommentTestSuite) Test_GetCommentById_Success() {
	suite.Context.SetParamNames("commentId")
	suite.Context.SetParamValues("1")

	suite.MocketClient.Select(&mocketHelper.Data{
		Model: &model.Comment{},
		Response: []map[string]interface{}{
			structHelper.MapAsGorm(&model.Comment{
				Id:      1,
				Body:    "",
				Deleted: false,
				UserId:  123,
			}),
		},
	})

	suite.NoError(suite.controller.GetCommentById(suite.Context))

}

func (suite *CommentTestSuite) Test_GetCommentByUserId_Success() {
	suite.Context.SetParamNames("userId")
	suite.Context.SetParamValues("1")

	suite.MocketClient.Select(&mocketHelper.Data{
		Model: &model.Comment{},
		Response: []map[string]interface{}{
			structHelper.MapAsGorm(&model.Comment{
				Id:      1,
				Body:    "",
				Deleted: false,
				UserId:  123,
			}),
		},
	})

	suite.NoError(suite.controller.GetCommentByUserId(suite.Context))
}

func (suite *CommentTestSuite) Test_UpdateComment_Success() {
	suite.Context.SetParamNames("commentId")
	suite.Context.SetParamValues("1")

	suite.MocketClient.Update(&mocketHelper.Data{
		Model: &model.Comment{Id: 1},
	})

	suite.NoError(suite.controller.UpdateComment(suite.Context))
}

func (suite *CommentTestSuite) Test_DeleteComment_Success() {
	suite.Context.SetParamNames("commentId")
	suite.Context.SetParamValues("1")

	suite.MocketClient.Update(&mocketHelper.Data{
		Model: &model.Comment{Id: 1},
	})

	suite.NoError(suite.controller.DeleteComment(suite.Context))
}

func (suite *CommentTestSuite) Test_CreateComment() {

	suite.MocketClient.Insert(&mocketHelper.Data{
		Model: &model.Comment{Id: 1},
	})

	suite.NoError(suite.controller.CreateComment(suite.Context))
}

func TestCommentSuite(t *testing.T) {
	suite.Run(t, new(CommentTestSuite))
}

func (suite *CommentTestSuite) TearDownTest() {
	suite.ctrl.Finish()
}
