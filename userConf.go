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
	Password string `json:"password"`
	FollowedChannelsByID []string `json:"followed_channels_by_id"`
}

type createUserInput struct {
	// binding adalah tag yang digunakan untuk menentukan aturan validasi pada field struct saat melakukan binding data dari request body. Dalam kasus ini, kita menggunakan binding:"required" untuk menandai bahwa field Name pada struct createUserInput wajib diisi saat melakukan binding data JSON dari request body. Jika field Name tidak diisi atau kosong, maka proses binding akan gagal dan menghasilkan error. Dengan menggunakan tag binding:"required", kita dapat memastikan bahwa data yang diterima dari request body memiliki field Name yang valid dan tidak kosong sebelum melanjutkan ke proses selanjutnya.
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	FollowedChannelsByID []string `json:"followed_channels_by_id"`
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

	newUser := user {
		ID : generateUserID(),
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
		FollowedChannelsByID: req.FollowedChannelsByID,
	}

	// kita melakukan iterasi pada slice users untuk memeriksa apakah ada user dengan ID yang sama dengan newUser.ID. Jika ditemukan, maka kita mengirimkan response dengan status code http.StatusConflict (409) dan pesan "user with the same ID already exists". Jika tidak ditemukan, maka kita melanjutkan proses untuk menambahkan data user baru ke dalam slice users. Namun, dalam kasus ini, kita tidak perlu melakukan pengecekan untuk ID yang sama karena ID pada struct user dihasilkan secara otomatis menggunakan fungsi generateUserID(), sehingga tidak mungkin ada dua user dengan ID yang sama. Oleh karena itu, kita dapat langsung menambahkan data user baru ke dalam slice users tanpa perlu melakukan pengecekan untuk ID yang sama.
	for _, a := range users {
		if newUser.Email == a.Email {
			c.IndentedJSON(http.StatusConflict, gin.H{"success":false, "message": "Email Already Use"})
			return
		}
	}


	// users = append(users, newUser) merupakan fungsi untuk menambahkan data user baru yang telah diterima dari request body ke dalam slice users. Fungsi append() digunakan untuk menambahkan elemen baru ke dalam slice. Dalam kasus ini, kita menambahkan newUser ke dalam slice users, sehingga data user baru tersebut akan disimpan dalam slice users dan dapat diakses melalui endpoint GET /users. Setelah menambahkan data user baru ke dalam slice users, kita mengirimkan response dengan status code http.StatusCreated (201) dan data user baru yang telah ditambahkan dalam format JSON.
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, gin.H{"success":true, "data":newUser})
	// json.MarshalIndent() merupakan fungsi untuk mengkonversi data dalam format struct atau slice menjadi format JSON dengan indentasi yang rapi. Fungsi ini menerima tiga parameter, yaitu data yang akan dikonversi, prefix untuk setiap baris (dalam kasus ini kita menggunakan string kosong), dan indentasi yang digunakan untuk setiap level (dalam kasus ini kita menggunakan dua spasi). Fungsi ini akan mengembalikan data dalam format JSON yang telah diindentasikan dengan rapi. Dalam kasus ini, kita mengkonversi slice users yang telah diperbarui dengan data user baru menjadi format JSON dengan indentasi yang rapi dan menyimpannya dalam variabel data.
	// perbedaan antara json.Marshal() dan json.MarshalIndent() adalah bahwa json.Marshal() akan menghasilkan output JSON dalam format yang lebih ringkas tanpa indentasi, sedangkan json.MarshalIndent() akan menghasilkan output JSON dengan indentasi yang rapi untuk meningkatkan keterbacaan. Dalam kasus ini, kita menggunakan json.MarshalIndent() untuk menghasilkan output JSON yang lebih mudah dibaca saat menyimpan data ke dalam file users.json.
	data, _ := json.MarshalIndent(users, "", "  ")
	// os.WriteFile() merupakan fungsi untuk menulis data ke dalam file. Fungsi ini menerima tiga parameter, yaitu nama file yang akan ditulis, data yang akan ditulis, dan permission untuk file tersebut. Dalam kasus ini, kita menulis data JSON yang telah diindentasikan dengan rapi ke dalam file "data/users.json" dengan permission 0644 (read and write untuk owner, read untuk group dan others).
	os.WriteFile("data/users.json", data, 0644)
}
