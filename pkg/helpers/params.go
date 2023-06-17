package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetParamId(c *gin.Context, key string) (int, error) {
	value := c.Param(key)

	id, err := strconv.Atoi(value)

	if err != nil {
		return 0, err
	}

	return id, nil
}
