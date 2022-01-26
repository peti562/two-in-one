package controller

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FibonacciTestSuite struct {
	suite.Suite
	ctrl       *gomock.Controller
	Recorder   *httptest.ResponseRecorder
	Context    echo.Context
	controller *FibonacciController
}

func TestFibonacciSuite(t *testing.T) {
	suite.Run(t, new(FibonacciTestSuite))
}

// SetupAllSuite has a SetupSuiteForCard method, which will run before the tests in the suite are run.
func (suite *FibonacciTestSuite) SetupSuite() {
	suite.ctrl = gomock.NewController(suite.T())
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	// Track the response payloads
	suite.Recorder = httptest.NewRecorder()
	suite.Context = echo.New().NewContext(request, suite.Recorder)
	suite.controller = NewFibonacciController()
}

type Test struct {
	Input          uint
	ExpectedOutput *big.Int
}

func (suite *FibonacciTestSuite) Test_calc() {

	for _, test := range GetTests() {
		actualOutput := suite.controller.calc(test.Input)
		suite.Equal(test.ExpectedOutput.Int64(), actualOutput.Int64())
	}
}

func GetTests() []Test {
	tests := []Test{
		{
			Input:          0,
			ExpectedOutput: big.NewInt(int64(0)),
		},
		{
			Input:          1,
			ExpectedOutput: big.NewInt(int64(1)),
		},
		{
			Input:          2,
			ExpectedOutput: big.NewInt(int64(1)),
		},
		{
			Input:          3,
			ExpectedOutput: big.NewInt(int64(2)),
		},
		{
			Input:          4,
			ExpectedOutput: big.NewInt(int64(3)),
		},
		{
			Input:          5,
			ExpectedOutput: big.NewInt(int64(5)),
		},
		{
			Input:          6,
			ExpectedOutput: big.NewInt(int64(8)),
		},
		{
			Input:          7,
			ExpectedOutput: big.NewInt(int64(13)),
		},
		{
			Input:          8,
			ExpectedOutput: big.NewInt(int64(21)),
		},
		{
			Input:          9,
			ExpectedOutput: big.NewInt(int64(34)),
		},
		{
			Input:          10,
			ExpectedOutput: big.NewInt(int64(55)),
		},
		{
			Input:          11,
			ExpectedOutput: big.NewInt(int64(89)),
		},
		{
			Input:          12,
			ExpectedOutput: big.NewInt(int64(144)),
		},
		{
			Input:          13,
			ExpectedOutput: big.NewInt(int64(233)),
		},
		{
			Input:          14,
			ExpectedOutput: big.NewInt(int64(377)),
		},
		{
			Input:          15,
			ExpectedOutput: big.NewInt(int64(610)),
		},
		{
			Input:          16,
			ExpectedOutput: big.NewInt(int64(987)),
		},
		{
			Input:          17,
			ExpectedOutput: big.NewInt(int64(1597)),
		},
		{
			Input:          18,
			ExpectedOutput: big.NewInt(int64(2584)),
		},
		{
			Input:          19,
			ExpectedOutput: big.NewInt(int64(4181)),
		},
		{
			Input:          20,
			ExpectedOutput: big.NewInt(int64(6765)),
		},
		{
			Input:          21,
			ExpectedOutput: big.NewInt(int64(10946)),
		},
		{
			Input:          22,
			ExpectedOutput: big.NewInt(int64(17711)),
		},
		{
			Input:          23,
			ExpectedOutput: big.NewInt(int64(28657)),
		},
		{
			Input:          24,
			ExpectedOutput: big.NewInt(int64(46368)),
		},
		{
			Input:          25,
			ExpectedOutput: big.NewInt(int64(75025)),
		},
		{
			Input:          26,
			ExpectedOutput: big.NewInt(int64(121393)),
		},
		{
			Input:          27,
			ExpectedOutput: big.NewInt(int64(196418)),
		},
		{
			Input:          28,
			ExpectedOutput: big.NewInt(int64(317811)),
		},
		{
			Input:          29,
			ExpectedOutput: big.NewInt(int64(514229)),
		},
		{
			Input:          30,
			ExpectedOutput: big.NewInt(int64(832040)),
		},
		{
			Input:          31,
			ExpectedOutput: big.NewInt(int64(1346269)),
		},
		{
			Input:          32,
			ExpectedOutput: big.NewInt(int64(2178309)),
		},
		{
			Input:          33,
			ExpectedOutput: big.NewInt(int64(3524578)),
		},
		{
			Input:          34,
			ExpectedOutput: big.NewInt(int64(5702887)),
		},
		{
			Input:          35,
			ExpectedOutput: big.NewInt(int64(9227465)),
		},
		{
			Input:          36,
			ExpectedOutput: big.NewInt(int64(14930352)),
		},
		{
			Input:          37,
			ExpectedOutput: big.NewInt(int64(24157817)),
		},
		{
			Input:          38,
			ExpectedOutput: big.NewInt(int64(39088169)),
		},
		{
			Input:          39,
			ExpectedOutput: big.NewInt(int64(63245986)),
		},
		{
			Input:          40,
			ExpectedOutput: big.NewInt(int64(102334155)),
		},
		{
			Input:          41,
			ExpectedOutput: big.NewInt(int64(165580141)),
		},
		{
			Input:          42,
			ExpectedOutput: big.NewInt(int64(267914296)),
		},
		{
			Input:          43,
			ExpectedOutput: big.NewInt(int64(433494437)),
		},
		{
			Input:          44,
			ExpectedOutput: big.NewInt(int64(701408733)),
		},
		{
			Input:          45,
			ExpectedOutput: big.NewInt(int64(1134903170)),
		},
		{
			Input:          46,
			ExpectedOutput: big.NewInt(int64(1836311903)),
		},
		{
			Input:          47,
			ExpectedOutput: big.NewInt(int64(2971215073)),
		},
		{
			Input:          48,
			ExpectedOutput: big.NewInt(int64(4807526976)),
		},
		{
			Input:          49,
			ExpectedOutput: big.NewInt(int64(7778742049)),
		},
		{
			Input:          50,
			ExpectedOutput: big.NewInt(int64(12586269025)),
		},
	}

	return tests
}
