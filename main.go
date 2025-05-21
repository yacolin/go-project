package main

import (
	"net/http"

	"go-project/configs"
	"go-project/controllers"
	"go-project/middlewares"
	"go-project/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.StaticFile("/favicon.ico", "./static/favicon.ico")
	r.Static("/static", "./static")
	r.Use(middlewares.ErrorHandler(), middlewares.ValidatorMiddleware())

	configs.ConnectMysql()

	// 初始化OSS配置
	// ossConfig := configs.OSSConfig{
	// 	Endpoint:        os.Getenv("OSS_ENDPOINT"),
	// 	AccessKeyID:     os.Getenv("OSS_ACCESS_KEY_ID"),
	// 	AccessKeySecret: os.Getenv("OSS_ACCESS_KEY_SECRET"),
	// 	BucketName:      os.Getenv("OSS_BUCKET_NAME"),
	// }

	// if err := configs.InitOSS(ossConfig); err != nil {
	// 	log.Fatalf("Failed to initialize OSS: %v", err)
	// }

	r.POST("/api/v1/login", controllers.Login)
	r.POST("/api/v1/register", controllers.Register)
	r.POST("/api/v1/refresh", controllers.Refresh)

	v1 := r.Group("/api/v1")
	// v1.Use(middlewares.JWTAuthMiddleware())
	{
		// v1.GET("/files", controllers.GetAllFiles)
		// v1.POST("/files", controllers.UploadFile)
		// v1.GET("/files/:id", controllers.GetFileByID)
		// v1.DELETE("/files/:id", controllers.DeleteFile)

		v1.GET("/teams", controllers.GetAllTeams)
		v1.GET("/team/:id", controllers.GetTeamByID)

		v1.GET("/books", controllers.GetAllBooks)
		v1.GET("/books/search", controllers.SearchBooks)
		v1.GET("/book/:id", controllers.GetBookByID)
		v1.POST("/book", controllers.CreateBook)
		v1.PUT("book/:id", controllers.UpdateBook)
		v1.DELETE("/book/:id", controllers.DeleteBook)

		v1.GET("/pets", controllers.GetAllPets)

		v1.GET("/albums", controllers.GetAllAlbums)
		v1.GET("albums/search", controllers.SearchAlbums)
		v1.GET("/album/:id", controllers.GetAlbumByID)
		v1.POST("/album", controllers.CreateAlbum)
		v1.PUT("/album/:id", controllers.UpdateAlbum)
		v1.DELETE("/album/:id", controllers.DeleteAlbum)

		v1.GET("/songs", controllers.GetAllSongs)

		// 评论相关
		v1.POST("/comments", controllers.CreateComment)
		v1.PUT("/comments/:id", controllers.UpdateComment)
		v1.DELETE("/comments/:id", controllers.DeleteComment)
		v1.GET("/photos/:id/comments", controllers.GetCommentsByPhotoID)
	}

	r.NoRoute(func(c *gin.Context) {
		utils.ErrorResponse(
			c,
			http.StatusNotFound,
			utils.ErrorNotFound,
			gin.H{"path": c.Request.URL.Path},
		)
	})

	r.Run()
}
