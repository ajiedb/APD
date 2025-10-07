package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	apmgin "go.elastic.co/apm/module/apmgin/v2"

	"ujicoba-go/config"
	"ujicoba-go/controllers"
	"ujicoba-go/middlewares"
	r "ujicoba-go/routes"
	"ujicoba-go/utils"
	"time"
	"ujicoba-go/models"
	"github.com/robfig/cron"
)

var version = "No Version Provided"

func genReportTask() {
	fmt.Println("===================================================")
	fmt.Println("Running cron task:", time.Now().Format(time.RFC3339))
	outcomes := models.ReportOutcomes{}
	incompleteReports, err := outcomes.GetInCompleteReport()
	if err != nil {
		fmt.Println("Error getting incomplete reports:", err)
		return
	}

	if len(incompleteReports) == 0 {
		fmt.Println("No incomplete reports found.")
		return
	}
	
	for _, report := range incompleteReports {
		if(report.Type == 1){
			fmt.Printf("Processing report UUID: %s (type: %d)", report.UUID, report.Type)
			result, err := controllers.GenerateRawEssayShortReport(report.UUID)
			if err != nil {
				fmt.Println("Failed to generate report for UUID %s: %v\n", report.UUID, err)
				continue
			}
			fmt.Println("Successfully generated report for UUID %s: %s\n", result.UUID, result.FileName)
		}else if(report.Type == 2){
			fmt.Printf("Processing report UUID: %s (type: %d)", report.UUID, report.Type)
			result, err := controllers.GenerateRawNonEssayShortReport(report.UUID)
			if err != nil {
				fmt.Println("Failed to generate report for UUID %s: %v\n", report.UUID, err)
				continue
			}
			fmt.Println("Successfully generated report for UUID %s: %s\n", result.UUID, result.FileName)
		}else if(report.Type == 3){
			fmt.Printf("Processing report UUID: %s (type: %d)", report.UUID, report.Type)
			result, err := controllers.GenerateResponseJawabanReport(report.UUID)
			if err != nil {
				fmt.Println("Failed to generate report for UUID %s: %v\n", report.UUID, err)
				continue
			}
			fmt.Println("Successfully generated report for UUID %s: %s\n", result.UUID, result.FileName)
		}else if(report.Type == 4){
			fmt.Printf("Processing report UUID: %s (type: %d)", report.UUID, report.Type)
			result, err := controllers.GenerateScoreReport(report.UUID)
			if err != nil {
				fmt.Println("Failed to generate report for UUID %s: %v\n", report.UUID, err)
				continue
			}
			fmt.Println("Successfully generated report for UUID %s: %s\n", result.UUID, result.FileName)
		}else if(report.Type == 5){
			fmt.Printf("Processing report UUID: %s (type: %d)", report.UUID, report.Type)
			result, err := controllers.GenerateQuestionReport(report.UUID)
			if err != nil {
				fmt.Println("Failed to generate report for UUID %s: %v\n", report.UUID, err)
				continue
			}
			fmt.Println("Successfully generated report for UUID %s: %s\n", result.UUID, result.FileName)
		}
		
	}
}

func main() {
	c := cron.New()
	c.AddFunc("0 * * * *", genReportTask) // Cron expression for every hours
	c.Start()
	
	gin.SetMode(gin.ReleaseMode)
	port := "8080"
	if config.EnvVariable("APP_PORT") != "" {
		port = config.EnvVariable("APP_PORT")
	}

	configCORS := "*"
	if config.EnvVariable("ALLOWED_ORIGINS") != "" {
		configCORS = config.EnvVariable("ALLOWED_ORIGINS")
	}

	router := gin.Default()

	router.Use(middlewares.ErrorHandler())
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", configCORS)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Transfer-Encoding, Connection, X-Powered-By, Cache-Control, Date, Server")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		utils.ApiService{}.SetContextAPI(c)

		c.Next()
	})

	router.Use(apmgin.Middleware(router))

	router.Use(func(c *gin.Context) {
		// meng set fullpath current server host
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		fullPath := scheme + "://" + c.Request.Host + c.Request.URL.Path
		c.Set("fullPath", fullPath)
		c.Next()
	})

	router.GET("/", func(c *gin.Context) {
		commitSHA := config.EnvVariable("COMMIT_SHA")
		c.JSON(200, gin.H{
			"message": "Ujicoba Go API: " + version,
			"sha":     commitSHA,
		})
	})

	v1 := router.Group("/api/v1")
	r.AuthRouter(v1.Group("/auth"))
	r.UserRouter(v1.Group("/user"))
	r.EventRouter(v1.Group("/events"))
	r.GatewayRouter(v1.Group("/gateway"))
	r.ExamRouter(v1.Group("/exam"))
	r.TestParticipantsRouter(v1.Group("/testparticipants"))
	r.TestAttemptsRouter(v1.Group("/testattempts"))
	r.TestAnswersRouter(v1.Group("/testanswers"))
	r.SyncRouter(v1.Group("/sync"))
	r.IntegrationRouter(v1.Group("/integration"))
	r.AnnouncementsRouter(v1.Group("/announcements"))
	r.RecapRouter(v1.Group("/recap"))
	r.FileRouter(v1.Group("/file"))
	r.LogRouter(v1.Group("/log"))
	r.ReportRouter(v1.Group("/reports"))

	resultController := controllers.ResultController{}
	v1.GET("export-jawaban", middlewares.VerifyUser([]int{utils.USER_ROLE_ADMIN}), resultController.ExportTestPackageAnswerResults)

	subDataExamController := controllers.SubDataExamController{}
	v1.GET("import-school-template", middlewares.VerifyUser([]int{utils.USER_ROLE_ADMIN}), subDataExamController.DownloadTemplateSchoolExam)
	updateController := controllers.UpdateController{}
	v1.POST("update/app", middlewares.VerifyOfflineMode(), middlewares.VerifyUser([]int{utils.USER_ROLE_PENGAWAS}), updateController.UpdateApp)
	exportAnswersController := controllers.ExportAnswersController{}
	v1.POST("export/answer", middlewares.VerifyOfflineMode(), middlewares.VerifyUser([]int{}), exportAnswersController.Download)
	uploadJawabanController := controllers.UploadJawabanController{}
	v1.POST("upload/jawaban", middlewares.VerifyUser([]int{}), uploadJawabanController.Upload)
	syncUrlController := controllers.SyncUrlController{}
	v1.POST("sync/pusmenjar-pusmendik", middlewares.VerifyUser([]int{}), syncUrlController.SyncPusmenjarPusmendik)

	fmt.Println("running on port :" + port)
	err := router.Run(":" + port)
	if err != nil {
		utils.WriteLog(fmt.Sprint("Error : ", err.Error()), "error_run")
		log.Fatal(err)
	}

}
