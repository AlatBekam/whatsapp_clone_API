package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type user struct {
	ID	string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var userJSON []byte
var users []user

func main() {
	var err error
	userJSON, err = os.ReadFile("users.json")
	if err != nil {
		panic(err)
	}

	json.Unmarshal(userJSON, &users)

	router := gin.Default()
	router.GET("/users", getUser)
	router.GET("/users/:id", getUserByID)
	router.Run("localhost:8080")
}


func getUser(c *gin.Context) {
	// c.IndentedJSON(http.StatusOK, userJSON)
	c.Data(http.StatusOK, "application/json", userJSON)
}

func getUserByID(c *gin.Context) {
	idParam := c.Param("id")

	for _, a := range users {
		if a.ID == 	idParam {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}


