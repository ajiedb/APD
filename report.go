package routes

import (
	"github.com/gin-gonic/gin"

	c "ujicoba-go/controllers"
	mid "ujicoba-go/middlewares"
	"ujicoba-go/utils"
)

func ReportRouter(r *gin.RouterGroup) {
	reportController := c.ReportController{}

	r.POST("", mid.VerifyUser([]int{utils.USER_ROLE_ADMIN}), reportController.CreateReportRequest)
	r.GET("", mid.VerifyUser([]int{utils.USER_ROLE_ADMIN}), reportController.List)
	r.GET("/download/:uuid/:fileName", mid.VerifyUser([]int{utils.USER_ROLE_ADMIN}), utils.GetFileFromMinio)
	r.GET("/outcomes/:uuid", mid.VerifyUser([]int{utils.USER_ROLE_ADMIN}), reportController.Outcomes)
	r.GET("/score-option/:uuid", mid.VerifyUser([]int{}), reportController.ScoreOption)
	r.GET("/score-question/:uuid", mid.VerifyUser([]int{}), reportController.ScoreQuestion)
	r.GET("/raw-non-essay-short/:uuid", mid.VerifyUser([]int{}), reportController.RawNonEssayShort)
	r.GET("/raw-essay-short/:uuid", mid.VerifyUser([]int{}), reportController.RawEssayShort)
	r.GET("/response-jawaban/:uuid", mid.VerifyUser([]int{}), reportController.ResponseJawaban)
}
