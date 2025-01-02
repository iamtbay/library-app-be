package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handlers struct{}

func initHandlers() *Handlers {
	return &Handlers{}
}

var services = initServices()

// GET ALL BOOKS
func (x *Handlers) getAllBooks(c *gin.Context) {
	fmt.Println(c.ClientIP())
	books, err := services.getAllBooks()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"books": books})
}

// ADD A BOOK
func (x *Handlers) addABook(c *gin.Context) {
	var bookInfo NewBookInfo
	err := c.BindJSON(&bookInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	err = services.addABook(bookInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"message": "book added"})
}

func (x *Handlers) getABook(c *gin.Context) {
	id := c.Param("id")
	book, err := services.getABook(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"book": book,
	})
}

func (x *Handlers) deleteABook(c *gin.Context) {
	id := c.Param("id")

	bookInfo, err := services.getABook(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = services.deleteABook(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	filepath := filepath.Join("./public/", bookInfo.Destination)
	fmt.Println(filepath)
	err = os.Remove(filepath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"message": "book deleted"})
}

func (x *Handlers) login(c *gin.Context) {
	var authInfo AuthInfo
	err := c.BindJSON(&authInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	userInfo, ss, err := services.login(authInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//set cookie
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("accessToken", ss, 3600, "/", "netlify.app", true, true)


	c.JSON(http.StatusOK, gin.H{
		"message": "Logged in successfully",
		"data":    userInfo,
	})
}

func (x *Handlers) register(c *gin.Context) {
	var authInfo NewAuthInfo
	err := c.BindJSON(&authInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	err = services.register(authInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registered successfully",
	})
}

func (x *Handlers) checkUser(c *gin.Context) {
	tokenString, _ := c.Cookie("accessToken")
	authInfo, err := parseJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "user not authorized!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "user verified",
		"user":    authInfo,
	})
}

func (x *Handlers) logout(c *gin.Context) {
	c.SetCookie("accessToken", "", -1, "/", "netlify.app", true, true)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "User successfully logged out.",
	})
}

// update dst
func (x *Handlers) addDst(c *gin.Context) {

	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	fileName := strings.ReplaceAll(file.Filename, " ", "_")
	file.Filename = fileName

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to open file",
		})
		return
	}
	defer openedFile.Close()

	//server
	dst := fmt.Sprintf("/mnt/data/uploads/%v", file.Filename)
	//local
	//	dst := fmt.Sprintf("./public/uploads/%v", file.Filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{"message": "book saved", "destination": fmt.Sprintf("uploads/%v", file.Filename)})
}
