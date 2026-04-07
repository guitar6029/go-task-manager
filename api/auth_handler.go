package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var body AuthRequest

		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid body"})
			return
		}

		if body.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email missing"})
			return
		}

		if body.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password missing"})
			return
		}

		//check if user exists in DB
		var userID int
		var email string
		var password string

		err := db.QueryRow(`SELECT id, email, password FROM users WHERE email = ?`, body.Email).Scan(&userID, &email, &password)
		if err == sql.ErrNoRows {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "server error"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(password), []byte(body.Password))
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}

		// jwt token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			c.JSON(500, gin.H{"error": "could not create token"})
			return
		}

		c.JSON(200, gin.H{
			"token": tokenString,
		})
	}
}

func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body AuthRequest

		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid body"})
			return
		}

		if body.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email missing"})
			return
		}

		if body.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password missing"})
			return
		}

		var existingID int
		//check if user exists
		err := db.QueryRow(`SELECT id FROM users WHERE email = ?`, body.Email).Scan(&existingID)
		if err == nil {
			c.JSON(409, gin.H{
				"error": "user already exists",
			})
			return
		}
		if err != sql.ErrNoRows {
			// real DB error
			c.JSON(500, gin.H{
				"error": "server error",
			})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "could not hash password"})
			return
		}

		_, err = db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", body.Email, string(hashedPassword))
		if err != nil {
			c.JSON(500, gin.H{
				"error": "could not create user",
			})
			return
		}

		c.JSON(201, gin.H{"message": "user created"})

	}
}
