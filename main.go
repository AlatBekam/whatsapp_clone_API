package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	JsonWebToken "whatsapp-clone-api/JWT"
	"whatsapp-clone-api/middleware"

	// github.com/gin-gonic/gin merupakan package yang digunakan untuk membuat web framework di Go. Package ini menyediakan berbagai fitur untuk memudahkan pengembangan aplikasi web, seperti routing, middleware, dan rendering template. Dalam kasus ini, kita menggunakan package gin untuk membuat API yang dapat menangani request dan response dalam format JSON.
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Name string `json:"name"`
	Password string `json:"password"`
}

var lastUserID int
var lastChannelID int
var userJSON []byte
var channelJSON []byte
var users []user
var channels []channel


func main() {

	var err error
	userJSON, err = os.ReadFile("data/users.json")
	channelJSON, err = os.ReadFile("data/channels.json")

	// jika terjadi error saat membaca file users.json, maka program akan panic dan menampilkan pesan error. Hal ini dilakukan untuk memastikan bahwa program tidak melanjutkan eksekusi jika file tidak dapat dibaca, sehingga mencegah terjadinya kesalahan lebih lanjut yang mungkin terjadi akibat data yang tidak tersedia.
	if err != nil {
		panic(err)
	}

	// json.Unmarshal() merupakan fungsi untuk mengkonversi data JSON menjadi struct atau slice dalam bahasa Go. Fungsi ini menerima dua parameter, yaitu data JSON yang akan dikonversi dan variabel yang akan menampung hasil konversi. Dalam kasus ini, kita mengkonversi data JSON yang dibaca dari file users.json menjadi slice of user dan menyimpannya dalam variabel users.
	json.Unmarshal(userJSON, &users)
	json.Unmarshal(channelJSON, &channels)
	lastUserID = len(users)
	lastChannelID = len(channels)

	router := gin.Default()
	router.GET("api/public/users", getUser)
	router.GET("api/public/users/:id", getUserByID)
	router.POST("api/public/channels", addChannel)
	router.POST("api/private/login", handlerLogin)
	router.POST("api/public/users", addUser)

	protected := router.Group("api/private")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.GET("/channels", getChannel)

	router.Run("localhost:8080")
}

func handlerLogin(c *gin.Context) {
	var req LoginRequest
	userFind := false

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
				"error":err.Error(),
		})
	return
	}

	// fmt.Println("===== REQUEST MASUK =====")
	// fmt.Println("Nama:", req.Name)
	// fmt.Println("Pass:", req.Password)
	// fmt.Println("=========================")

	for _, a := range users {
		if req.Name == a.Name && req.Password == a.Password {
			userFind = true
			break
		}
	}

	if !userFind {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"Invalid User or Password"})
		return
	}

	token, err := JsonWebToken.GenerateJWT(req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"Failed generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success":true, "token":token})

// 	response := gin.H{
// 	"success": true,
// 	"token":   token,
// }

// fmt.Println("===== RESPONSE =====")
// fmt.Println(response)
// fmt.Println("====================")

// c.JSON(http.StatusOK, response)

}





func generateUserID() string {
	lastUserID++
	return fmt.Sprintf("%d", lastUserID)
}

func generateChannelID() string {
	lastChannelID++
	return fmt.Sprintf("ch_%d", lastChannelID)
}



