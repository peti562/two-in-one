package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"math/big"
	"net/http"
	"strconv"
)

// FibonacciController controller object
type FibonacciController struct {
}

func NewFibonacciController() *FibonacciController {
	return &FibonacciController{}
}

func (fc *FibonacciController) Get(c echo.Context) error {
	n, exception := strconv.Atoi(c.Param("n"))
	if exception != nil {
		// should be proper error handling here
		fmt.Println(exception)
		return exception
	}

	result := fc.calc(uint(n))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"input":  n,
		"output": result.String(),
	})
}
func (fc *FibonacciController) calc(n uint) *big.Int {
	if n <= 1 {
		return big.NewInt(int64(n))
	}
	var second, first = big.NewInt(0), big.NewInt(1)

	for i := uint(1); i < n; i++ {
		second.Add(second, first)
		first, second = second, first
	}

	return first
}
