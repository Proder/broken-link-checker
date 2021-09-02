package checker

import (
	response2 "broken-link-checker/app/internal/delivery/http/api/v1/response"
	models2 "broken-link-checker/app/internal/models"
	linkChecker2 "broken-link-checker/app/internal/service/linkChecker"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type responseData struct {
	Duration   string
	BreakLinks []string
}

func SearchBrokenLinks(c *gin.Context) {
	var data models2.CheckerRequestData

	if err := c.BindJSON(&data); err != nil {
		log.Println("checker -> SearchBrokenLinks: error parse data. reason: ", err.Error())
		response2.Error(c, "Error in request data: "+err.Error())
		return
	}

	checker := linkChecker2.Checker{}
	if err := checker.Run(data.Link, data.Depth); err != nil {
		log.Println("checker -> SearchBrokenLinks: linkChecker error. reason: ", err.Error())
		response2.Error(c, "Server error: "+err.Error())
		return
	}

	fmt.Printf("Broken links found: %d\n", len(checker.GetBreakLinks()))
	fmt.Printf("Time spent: %s\n", checker.GetDuration())

	response2.Success(c, responseData{
		Duration:   checker.GetDuration(),
		BreakLinks: checker.GetBreakLinks(),
	})
}
