package Test

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webapptrials/Classroom/Shared"
)

func TestData(c *gin.Context) {
	Shared.SendGinGenericApiResponse(c, http.StatusOK, "Test api functional!", "Success!", "")
}
