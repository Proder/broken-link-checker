package response

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context) {
	rand := strconv.FormatInt(time.Now().UnixNano(), 10)
	html := "<html>" +
		"<a href='' >link</a>" +
		"<a href='/' /a>" +
		"<a href='http://localhost:8080/' /a>" +
		"<a href='http://localhost:8080/success' /a>" +
		"<a href='/success/' /a>" +
		"<a href='/success/1' /a>" +
		"<a href='/success/2' /a>" +
		"<a href='/success/3/' /a>" +
		"<a href='/success/" + rand + "' /a>" +
		"<a href='http://localhost:8080/error/' /a>" +
		"<a href='http://localhost:8080/error/url' /a>" +
		"<a href='/error/url' /a>" +
		"<a href='/error/url/1' /a>" +
		"<a href='/error/url/2' /a>" +
		"<a href='/error/url/3/' /a>" +
		"<a href='/error/url/" + rand + "' /a>" +
		"<a href='/error/server' /a>" +
		"<a href='/error/server/1' /a>" +
		"<a href='/error/server/2' /a>" +
		"<a href='/error/server/3/' /a>" +
		"</html>"

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func Redirect(c *gin.Context) {
	c.JSON(301, nil)
}

func ErrorUrl(c *gin.Context) {
	c.JSON(400, nil)
}

func ErrorServer(c *gin.Context) {
	c.JSON(500, nil)
}
