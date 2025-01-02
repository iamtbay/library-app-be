package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func initRoutes(r *gin.Engine) {
	r.Use(CORSMw())
	r.Static("/uploads", "./public/uploads")
	handlers := initHandlers()
	route := r.Group("/api/v1")

	route.POST("/auth/login", isNotAuthenticatedMW(), handlers.login)
	route.POST("/auth/register", isNotAuthenticatedMW(), handlers.register)

	//
	route.Use(IsAuthenticatedMW())
	route.GET("", handlers.getAllBooks)
	route.GET("/auth/check-user", handlers.checkUser)
	route.POST("/auth/logout", handlers.logout)
	route.POST("/add-book-file", handlers.addDst)
	route.POST("/add-book", handlers.addABook)
	route.GET("/book/:id", handlers.getABook)
	route.DELETE("/book/:id", handlers.deleteABook)
}

func CORSMw() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://kugulupark.netlify.app")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No content
			return
		}

		c.Next()
	}
}



func IsAuthenticatedMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("accessToken")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthenticated, please login!",
			})
			c.Abort()
			return
		}

		_, err = parseJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func isNotAuthenticatedMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("accessToken")
		if err != nil {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "User already logged in",
		})
		c.Abort()

	}
}
