package controllers

import (
	"log"
	"net/http"
	"time"

	"login_jwt/db"
	"login_jwt/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func Login(c *gin.Context) {
	con := db.CreateCon()      //Prepare DB Connection
	users := new(models.Users) //Prepare Variable for Payloads
	hash_password := ""

	type Users_Login struct {
		Username     string `json:"Username" binding:"required"`
		Password     string `json:"Password" binding:"required"`
		Nama_lengkap string `json:"Nama_lengkap" binding:"-"`
	}

	var users_login Users_Login
	//Get The Payloads
	if err := c.BindJSON(&users_login); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	//Get The Record from DB
	query := "SELECT password,nama_lengkap FROM  " + users.TableName() + "  WHERE username='" + users_login.Username + "' LIMIT 1"
	rows, err := con.Query(query)
	defer con.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	for rows.Next() {

		if err := rows.Scan(&hash_password, &users_login.Nama_lengkap); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		} else {
			if comparePasswords(hash_password, []byte(users_login.Password)) {
				sign := jwt.New(jwt.GetSigningMethod("HS256"))
				token, err := sign.SignedString([]byte("secret"))
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": err.Error(),
					})
					c.Abort()
				}
				c.JSON(http.StatusOK, gin.H{
					"username": users_login.Username,
					"token":    token,
				})
				return
			}
		}
	}
	defer con.Close()

	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "error",
	})
	return
}

func Register(c *gin.Context) {
	con := db.CreateCon()      //Prepare DB Connection
	users := new(models.Users) //Prepare Variable for Payloads
	var users_payload models.Users

	//Get The Payloads
	if err := c.BindJSON(&users_payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	//Prepare the Data
	users_payload.Id = uuid.Must(uuid.NewRandom())
	users_payload.Created_at = time.Now()
	users_payload.Password = hashAndSalt([]byte(users_payload.Password))

	//Check is there same Username or not on DB
	query := "SELECT username FROM  " + users.TableName() + "  WHERE username='" + users_payload.Username + "' LIMIT 1"
	rows, err := con.Query(query)
	defer con.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	counter := 0
	for rows.Next() {
		counter += 1
	}
	if counter > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Username already registered",
		})
		return
	}

	//Store the Data to DB
	insForm, err := con.Prepare("INSERT INTO " + users.TableName() + " (id,username,password,nama_lengkap,created_at) VALUES(?,?,?,?,?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	insForm.Exec(users_payload.Id, users_payload.Username, users_payload.Password, users_payload.Nama_lengkap, users_payload.Created_at)
	defer con.Close()

	//Make Response
	c.JSON(http.StatusOK, gin.H{
		"message":  "done",
		"nama":     users_payload.Nama_lengkap,
		"username": users_payload.Username,
	})
	return
}
