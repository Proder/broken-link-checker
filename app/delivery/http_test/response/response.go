package response

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context) {
	html := "<html>" +
		"<a href='' >link</a>" +
		"<a href='/' /a>" +
		"<a href='http://localhost:8080/' /a>" +
		"<a href='http://localhost:8080/success/' /a>" +
		"<a href='/success/' /a>" +
		"<a href='/success/1' /a>" +
		"<a href='/success/2' /a>" +
		"<a href='/success/3/' /a>" +
		getRandomLink("/success/", 500) +
		"<a href='http://localhost:8080/error/' /a>" +
		"<a href='http://localhost:8080/error/url/' /a>" +
		"<a href='/error/url/' /a>" +
		"<a href='/error/url/1' /a>" +
		"<a href='/error/url/2' /a>" +
		"<a href='/error/url/3/' /a>" +
		getRandomLink("/error/url/", 500) +
		"<a href='/error/server/' /a>" +
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

func getRandomLink(href string, count int) string {
	var (
		prefix  = "<a href='"
		postfix = "/' /a>"
	)

	buffer := bytes.Buffer{}
	for i := 0; i < count; i++ {
		randStr := strconv.FormatInt(int64(i+100), 10) // strconv.FormatInt(int64(rand.Intn(10000)+10), 10)
		buffer.WriteString(prefix)
		buffer.WriteString(href)
		buffer.WriteString(randStr)
		buffer.WriteString(postfix)
	}

	return buffer.String()
}
