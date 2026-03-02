package main

import (
	"encoding/json"
	"fmt"
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
	userJSON, err = os.ReadFile("data/users.json")

	// jika terjadi error saat membaca file users.json, maka program akan panic dan menampilkan pesan error. Hal ini dilakukan untuk memastikan bahwa program tidak melanjutkan eksekusi jika file tidak dapat dibaca, sehingga mencegah terjadinya kesalahan lebih lanjut yang mungkin terjadi akibat data yang tidak tersedia.
	if err != nil {
		panic(err)
	}

	// json.Unmarshal() merupakan fungsi untuk mengkonversi data JSON menjadi struct atau slice dalam bahasa Go. Fungsi ini menerima dua parameter, yaitu data JSON yang akan dikonversi dan variabel yang akan menampung hasil konversi. Dalam kasus ini, kita mengkonversi data JSON yang dibaca dari file users.json menjadi slice of user dan menyimpannya dalam variabel users.
	json.Unmarshal(userJSON, &users)

	router := gin.Default()
	router.GET("/users", getUser)
	router.GET("/users/:id", getUserByID)
	router.POST("/users", addUser)
	router.Run("localhost:8080")
}


func getUser(c *gin.Context) {
	// c.indentedJSON() merupakan fungsi untuk mengirimkan response dalam format JSON dengan indentasi yang rapi. Fungsi ini menerima dua parameter, yaitu status code HTTP dan data yang akan dikirimkan dalam format JSON. Dalam kasus ini, kita mengirimkan status code http.StatusOK (200) dan data userJSON yang berisi informasi tentang semua user.
	c.IndentedJSON(http.StatusOK, users)

	// c.Data() merupakan fungsi untuk mengirimkan response dengan format data tertentu. Fungsi ini menerima tiga parameter, yaitu status code HTTP, content type, dan data yang akan dikirimkan. Dalam kasus ini, kita mengirimkan status code http.StatusOK (200), content type "application/json", dan data userJSON yang berisi informasi tentang semua user.
	// c.Data(http.StatusOK, "application/json", userJSON)

	// kenapa menggunakan c.Data() bukan c.IndentedJSON() karena c.IndentedJSON() akan mengirimkan data dalam format JSON dengan indentasi yang rapi, sedangkan c.Data() akan mengirimkan data dalam format yang sesuai dengan content type yang ditentukan. Dalam kasus ini, kita ingin mengirimkan data dalam format JSON tanpa indentasi, sehingga menggunakan c.Data() lebih tepat.
}

func getUserByID(c *gin.Context) {
	// c.Param() merupakan fungsi untuk mengambil parameter dari URL. Dalam kasus ini, kita mengambil parameter "id" yang didefinisikan dalam route "/users/:id". Parameter ini akan digunakan untuk mencari user dengan ID yang sesuai dalam slice users.
	idParam := c.Param("id")
	// strconv.Atoi() merupakan fungsi untuk mengkonversi string menjadi integer. Namun, dalam kasus ini, ID pada struct user didefinisikan sebagai string, sehingga tidak perlu melakukan konversi ke integer. Oleh karena itu, kita dapat langsung membandingkan idParam dengan ID pada struct user tanpa perlu menggunakan strconv.Atoi().
	// id, _ := strconv.Atoi(idParam)

	// kita melakukan iterasi pada slice users untuk mencari user dengan ID yang sesuai dengan idParam. Jika ditemukan, maka kita mengirimkan response dengan status code http.StatusOK (200) dan data user yang ditemukan dalam format JSON. Jika tidak ditemukan, maka kita mengirimkan response dengan status code http.StatusNotFound (404) dan pesan "user not found".
	// pada golang, _ digunakan untuk mengabaikan nilai yang dikembalikan oleh fungsi. Dalam kasus ini, kita mengabaikan nilai error yang dikembalikan oleh strconv.Atoi() karena kita tidak perlu melakukan konversi ke integer. Namun, jika kita ingin menangani error tersebut, kita dapat menggunakan variabel lain untuk menyimpan nilai error dan melakukan pengecekan sebelum melanjutkan eksekusi.
	for _, a := range users {
		if a.ID == 	idParam {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func addUser(c *gin.Context) {
	var newUser user
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
	data, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile("users.json", data, 0644)
	fmt.Println("POST MASUK")
}


