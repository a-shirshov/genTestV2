package llm

import (
	"github.com/gin-gonic/gin"
)

type Delivery interface {
	GenerateTest(c *gin.Context)
}