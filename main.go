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

type user struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	Email                string   `json:"email"`
	Password             string   `json:"password"`
	FollowedChannelsByID []string `json:"followed_channels_by_id"`
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type channel struct {
	ChannelId   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	ChannelType string `json:"channel_type"`
	Description string `json:"description"`
}

type createUserInput struct {
	// binding adalah tag yang digunakan untuk menentukan aturan validasi pada field struct saat melakukan binding data dari request body. Dalam kasus ini, kita menggunakan binding:"required" untuk menandai bahwa field Name pada struct createUserInput wajib diisi saat melakukan binding data JSON dari request body. Jika field Name tidak diisi atau kosong, maka proses binding akan gagal dan menghasilkan error. Dengan menggunakan tag binding:"required", kita dapat memastikan bahwa data yang diterima dari request body memiliki field Name yang valid dan tidak kosong sebelum melanjutkan ke proses selanjutnya.
	Name                 string   `json:"name"`
	Email                string   `json:"email"`
	Password             string   `json:"password"`
	FollowedChannelsByID []string `json:"followed_channels_by_id"`
}
type updateUser struct {
	Name                 *string   `json:"name"`
	Email                *string   `json:"email"`
	Password             *string   `json:"password"`
	FollowedChannelsByID *[]string `json:"followed_channels_by_id"`
}

type createChannelInput struct {
	ChannelName string `json:"channel_name" binding:"required"`
	ChannelType string `json:"channel_type" binding:"required"`
	Description string `json:"description"`
}

type message struct {
	MessageID string `json:"message_id"`
	SenderID  string `json:"sender_id"`
	Content   string `json:"content"`
	Timestamp int    `json:"timestamp"`
}

type chat struct {
	ChatID   string    `json:"chat_id"`
	UserID   [2]string `json:"user_id"`
	Messages []message `json:"messages"`
}

var lastUserID int
var lastChannelID int
var userJSON []byte
var channelJSON []byte
var users []user
var channels []channel
var chats []chat
var chatsJSON []byte
var lastchatID int

func main() {

	var err error
	userJSON, err = os.ReadFile("data/users.json")
	channelJSON, err = os.ReadFile("data/channels.json")
	chatsJSON, err := os.ReadFile("data/datachat.json")

	// jika terjadi error saat membaca file users.json, maka program akan panic dan menampilkan pesan error. Hal ini dilakukan untuk memastikan bahwa program tidak melanjutkan eksekusi jika file tidak dapat dibaca, sehingga mencegah terjadinya kesalahan lebih lanjut yang mungkin terjadi akibat data yang tidak tersedia.
	if err != nil {
		panic(err)
	}

	// json.Unmarshal() merupakan fungsi untuk mengkonversi data JSON menjadi struct atau slice dalam bahasa Go. Fungsi ini menerima dua parameter, yaitu data JSON yang akan dikonversi dan variabel yang akan menampung hasil konversi. Dalam kasus ini, kita mengkonversi data JSON yang dibaca dari file users.json menjadi slice of user dan menyimpannya dalam variabel users.
	json.Unmarshal(userJSON, &users)
	json.Unmarshal(channelJSON, &channels)
	json.Unmarshal(chatsJSON, &chats)
	lastUserID = len(users)
	lastChannelID = len(channels)
	lastchatID = len(chats)

	router := gin.Default()
	router.GET("api/public/users", getUser)
	router.GET("api/public/users/:id", getUserByID)
	router.POST("api/public/channels", addChannel)
	router.POST("api/public/login", handlerLogin)
	router.POST("api/public/users", addUser)

	protected := router.Group("api/private")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.GET("/channels", getChannel)
	protected.POST("/users", editUserByID)
	protected.POST("/chats", addChat)
	protected.GET("/chats", getChat)

	router.Run(":8080")
}

func addChat(c *gin.Context) {
	myIDAny, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID Doesnt exist"})
		return
	}
	myID := myIDAny.(string)

	var req struct {
		ReceiverID string `json:"receiver_id"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, a := range chats {
		u := a.UserID
		if (u[0] == myID && u[1] == req.ReceiverID) || (u[0] == req.ReceiverID && u[1] == myID) {
			c.JSON(http.StatusOK, gin.H{"chat_id": a.ChatID})
			return
		}
	}

	// Jika chat tidak ditemukan, buat chat baru
	lastchatID++
	newChatID := fmt.Sprintf("chat_%d", lastchatID)
	newChat := chat{
		ChatID:   newChatID,
		UserID:   [2]string{myID, req.ReceiverID},
		Messages: []message{},
	}
	chats = append(chats, newChat)

	// Simpan ke file
	data, _ := json.MarshalIndent(chats, "", "  ")
	os.WriteFile("data/datachat.json", data, 0644)

	c.JSON(http.StatusCreated, gin.H{"chat_id": newChatID})
}

// getChat returns all chats for the currently logged-in user
func getChat(c *gin.Context) {
	myIDAny, exist := c.Get("userID")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID Doesnt exist"})
		return
	}
	myID := myIDAny.(string)

	// Find all chats that the user is part of
	var userChats []chat
	for _, a := range chats {
		u := a.UserID
		if u[0] == myID || u[1] == myID {
			userChats = append(userChats, a)
		}
	}

	c.JSON(http.StatusOK, gin.H{"chats": userChats})
}

func handlerLogin(c *gin.Context) {
	var req LoginRequest
	userFind := false
	var userID string

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
			userID = a.ID
			break
		}
	}

	if !userFind {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User or Password"})
		return
	}

	token, err := JsonWebToken.GenerateJWT(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "token": token})

	// 	response := gin.H{
	// 	"success": true,
	// 	"token":   token,
	// }

	// fmt.Println("===== RESPONSE =====")
	// fmt.Println(response)
	// fmt.Println("====================")

	// c.JSON(http.StatusOK, response)

}

func getUser(c *gin.Context) {
	// c.indentedJSON() merupakan fungsi untuk mengirimkan response dalam format JSON dengan indentasi yang rapi. Fungsi ini menerima dua parameter, yaitu status code HTTP dan data yang akan dikirimkan dalam format JSON. Dalam kasus ini, kita mengirimkan status code http.StatusOK (200) dan data userJSON yang berisi informasi tentang semua user.
	c.IndentedJSON(http.StatusOK, users)
	// fmt.Println(lastUserID)
	// fmt.Println(lastChannelID)

	// c.Data() merupakan fungsi untuk mengirimkan response dengan format data tertentu. Fungsi ini menerima tiga parameter, yaitu status code HTTP, content type, dan data yang akan dikirimkan. Dalam kasus ini, kita mengirimkan status code http.StatusOK (200), content type "application/json", dan data userJSON yang berisi informasi tentang semua user.
	// c.Data(http.StatusOK, "application/json", userJSON)

	// kenapa menggunakan c.Data() bukan c.IndentedJSON() karena c.IndentedJSON() akan mengirimkan data dalam format JSON dengan indentasi yang rapi, sedangkan c.Data() akan mengirimkan data dalam format yang sesuai dengan content type yang ditentukan. Dalam kasus ini, kita ingin mengirimkan data dalam format JSON tanpa indentasi, sehingga menggunakan c.Data() lebih tepat.
}

func getChannel(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, channels)

	// fmt.Println("All Headers:", c.Request.Header)

	response := gin.H{"response": channels}

	fmt.Println("===== RESPONSE =====")
	fmt.Println(response)
	fmt.Println("====================")
}

func getUserByID(c *gin.Context) {
	// c.Param() merupakan fungsi untuk mengambil parameter dari URL. Dalam kasus ini, kita mengambil parameter "id" yang didefinisikan dalam route "/users/:id". Parameter ini akan digunakan untuk mencari user dengan ID yang sesuai dalam slice users.
	idParam := c.Param("id")
	// strconv.Atoi() merupakan fungsi untuk mengkonversi string menjadi integer. Namun, dalam kasus ini, ID pada struct user didefinisikan sebagai string, sehingga tidak perlu melakukan konversi ke integer. Oleh karena itu, kita dapat langsung membandingkan idParam dengan ID pada struct user tanpa perlu menggunakan strconv.Atoi().
	// id, _ := strconv.Atoi(idParam)

	// kita melakukan iterasi pada slice users untuk mencari user dengan ID yang sesuai dengan idParam. Jika ditemukan, maka kita mengirimkan response dengan status code http.StatusOK (200) dan data user yang ditemukan dalam format JSON. Jika tidak ditemukan, maka kita mengirimkan response dengan status code http.StatusNotFound (404) dan pesan "user not found".
	// pada golang, _ digunakan untuk mengabaikan nilai yang dikembalikan oleh fungsi. Dalam kasus ini, kita mengabaikan nilai error yang dikembalikan oleh strconv.Atoi() karena kita tidak perlu melakukan konversi ke integer. Namun, jika kita ingin menangani error tersebut, kita dapat menggunakan variabel lain untuk menyimpan nilai error dan melakukan pengecekan sebelum melanjutkan eksekusi.
	for _, a := range users {
		if a.ID == idParam {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "user not found"})
}

func saveUserToJson(users []user) error {
	file, err := os.Create("data/users.json")

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(users)
}

func editUserByID(c *gin.Context) {
	var req updateUser
	found := false

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	idParam, exist := c.Get("userID")
	IDParam := idParam.(string)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ID Doesnt exist"})
		return
	}

	for a := range users {
		if users[a].ID == IDParam {
			if req.Name != nil {
				users[a].Name = *req.Name
			}
			if req.Email != nil {
				users[a].Email = *req.Email
			}
			if req.Password != nil {
				users[a].Password = *req.Password
			}
			if req.FollowedChannelsByID != nil {
				users[a].FollowedChannelsByID = *req.FollowedChannelsByID
			}
			found = true
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	}

	err := saveUserToJson(users)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"success": true,
	// 	"user": users,
	// })

	// fmt.Println("User ID dari token:", IDParam)
	// fmt.Println("Update name jadi:", *req.Name)
	// fmt.Println("Update Password jadi:", *req.Password)
	// fmt.Println("Update email jadi:", req.Email)
	// fmt.Println("tambah channel:", req.FollowedChannelsByID)

	// for _, a := range users {
	// 	if a.ID == idParam {
	// 		c.IndentedJSON(http.StatusOK, )
	// 	}
	// }
}

func generateUserID() string {
	lastUserID++
	return fmt.Sprintf("%d", lastUserID)
}
func generateChannelID() string {
	lastChannelID++
	return fmt.Sprintf("ch_%d", lastChannelID)
}

func addUser(c *gin.Context) {
	var req createUserInput

	// var newUser user merupakan deklarasi variabel newUser dengan tipe data user. Variabel ini akan digunakan untuk menyimpan data user baru yang diterima dari request body dalam format JSON.
	// var newUser user
	// c.BindJSON() merupakan fungsi untuk mengikat data JSON yang diterima dari request body ke dalam variabel yang telah dideklarasikan. Fungsi ini menerima satu parameter, yaitu pointer ke variabel yang akan menampung data JSON. Dalam kasus ini, kita mengikat data JSON yang diterima dari request body ke dalam variabel newUser. Jika terjadi error saat proses pengikatan, maka kita mengirimkan response dengan status code http.StatusBadRequest (400) dan pesan error yang terjadi.
	if err := c.BindJSON(&req); err != nil {
		// c.JSON() merupakan fungsi untuk mengirimkan response dalam format JSON. Fungsi ini menerima dua parameter, yaitu status code HTTP dan data yang akan dikirimkan dalam format JSON. Dalam kasus ini, kita mengirimkan status code http.StatusBadRequest (400) dan data berupa objek JSON yang berisi pesan error yang terjadi saat proses pengikatan data JSON ke dalam variabel newUser.
		// gin.H merupakan tipe data yang digunakan untuk membuat objek JSON dalam format key-value. Dalam kasus ini, kita membuat objek JSON dengan key "error" dan value berupa pesan error yang terjadi saat proses pengikatan data JSON ke dalam variabel newUser.
		c.JSON(http.StatusBadRequest, gin.H{
			// err.Error() merupakan metode untuk mendapatkan pesan error dalam bentuk string dari variabel err yang berisi informasi tentang error yang terjadi saat proses pengikatan data JSON ke dalam variabel newUser. Pesan error ini akan dikirimkan sebagai value dari key "error" dalam objek JSON yang dikirimkan sebagai response.
			"error": err.Error(),
		})
		return
	}

	newUser := user{
		ID:                   generateUserID(),
		Name:                 req.Name,
		Email:                req.Email,
		Password:             req.Password,
		FollowedChannelsByID: req.FollowedChannelsByID,
	}

	// kita melakukan iterasi pada slice users untuk memeriksa apakah ada user dengan ID yang sama dengan newUser.ID. Jika ditemukan, maka kita mengirimkan response dengan status code http.StatusConflict (409) dan pesan "user with the same ID already exists". Jika tidak ditemukan, maka kita melanjutkan proses untuk menambahkan data user baru ke dalam slice users. Namun, dalam kasus ini, kita tidak perlu melakukan pengecekan untuk ID yang sama karena ID pada struct user dihasilkan secara otomatis menggunakan fungsi generateUserID(), sehingga tidak mungkin ada dua user dengan ID yang sama. Oleh karena itu, kita dapat langsung menambahkan data user baru ke dalam slice users tanpa perlu melakukan pengecekan untuk ID yang sama.
	for _, a := range users {
		if newUser.Email == a.Email {
			c.IndentedJSON(http.StatusConflict, gin.H{"success": false, "message": "Email Already Use"})
			return
		}
	}

	// users = append(users, newUser) merupakan fungsi untuk menambahkan data user baru yang telah diterima dari request body ke dalam slice users. Fungsi append() digunakan untuk menambahkan elemen baru ke dalam slice. Dalam kasus ini, kita menambahkan newUser ke dalam slice users, sehingga data user baru tersebut akan disimpan dalam slice users dan dapat diakses melalui endpoint GET /users. Setelah menambahkan data user baru ke dalam slice users, kita mengirimkan response dengan status code http.StatusCreated (201) dan data user baru yang telah ditambahkan dalam format JSON.
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, gin.H{"success": true, "data": newUser})
	// json.MarshalIndent() merupakan fungsi untuk mengkonversi data dalam format struct atau slice menjadi format JSON dengan indentasi yang rapi. Fungsi ini menerima tiga parameter, yaitu data yang akan dikonversi, prefix untuk setiap baris (dalam kasus ini kita menggunakan string kosong), dan indentasi yang digunakan untuk setiap level (dalam kasus ini kita menggunakan dua spasi). Fungsi ini akan mengembalikan data dalam format JSON yang telah diindentasikan dengan rapi. Dalam kasus ini, kita mengkonversi slice users yang telah diperbarui dengan data user baru menjadi format JSON dengan indentasi yang rapi dan menyimpannya dalam variabel data.
	// perbedaan antara json.Marshal() dan json.MarshalIndent() adalah bahwa json.Marshal() akan menghasilkan output JSON dalam format yang lebih ringkas tanpa indentasi, sedangkan json.MarshalIndent() akan menghasilkan output JSON dengan indentasi yang rapi untuk meningkatkan keterbacaan. Dalam kasus ini, kita menggunakan json.MarshalIndent() untuk menghasilkan output JSON yang lebih mudah dibaca saat menyimpan data ke dalam file users.json.
	data, _ := json.MarshalIndent(users, "", "  ")
	// os.WriteFile() merupakan fungsi untuk menulis data ke dalam file. Fungsi ini menerima tiga parameter, yaitu nama file yang akan ditulis, data yang akan ditulis, dan permission untuk file tersebut. Dalam kasus ini, kita menulis data JSON yang telah diindentasikan dengan rapi ke dalam file "data/users.json" dengan permission 0644 (read and write untuk owner, read untuk group dan others).
	os.WriteFile("data/users.json", data, 0644)
}

func addChannel(c *gin.Context) {
	var req createChannelInput

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newChannel := channel{
		ChannelId:   generateChannelID(),
		ChannelName: req.ChannelName,
		ChannelType: req.ChannelType,
		Description: req.Description,
	}

	channels = append(channels, newChannel)
	c.IndentedJSON(http.StatusCreated, newChannel)
	data, _ := json.MarshalIndent(channels, "", "  ")
	os.WriteFile("data/channels.json", data, 0644)
}
