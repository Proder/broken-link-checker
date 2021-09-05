package response

import (
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	var response struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}
	response.Status = "Success"
	response.Data = data

	c.JSON(200, response)
}

func Error(c *gin.Context, msg string) {
	var response struct {
		Status string `json:"status"`
		Msg    string `json:"smg"`
	}
	response.Status = "Error"
	response.Msg = msg

	c.JSON(500, response)
}
