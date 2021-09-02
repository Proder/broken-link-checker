package checker

import (
	"broken-link-checker/app/internal/delivery/http/api/v1/response"
	"broken-link-checker/app/internal/models"
	"broken-link-checker/app/internal/service/linkChecker"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type responseData struct {
	Duration   string
	BreakLinks []string
}

func SearchBrokenLinks(c *gin.Context) {
	var data models.CheckerRequestData

	if err := c.BindJSON(&data); err != nil {
		log.Println("checker -> SearchBrokenLinks: error parse data. reason: ", err.Error())
		response.Error(c, "Error in request data: "+err.Error())
		return
	}

	checker := linkChecker.Checker{}
	if err := checker.Run(data.Link, data.Depth); err != nil {
		log.Println("checker -> SearchBrokenLinks: linkChecker error. reason: ", err.Error())
		response.Error(c, "Server error: "+err.Error())
		return
	}

	fmt.Printf("Broken links found: %d\n", len(checker.GetBreakLinks()))
	fmt.Printf("Time spent: %s\n", checker.GetDuration())

	response.Success(c, responseData{
		Duration:   checker.GetDuration(),
		BreakLinks: checker.GetBreakLinks(),
	})
}
