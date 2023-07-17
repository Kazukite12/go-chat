package userController

import (
	"net/http"
	"time"

	"github.com/Kazukite12/go-chat/models"
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("my_secret_key")

func Auth(c *gin.Context) {
	// get cookie off req

	tokenString, err := c.Cookie("jwt")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unathorized"})
	}

	//decode it

	//parse
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Youre not logged in"})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)

		}

		var user models.User
		models.DB.First(&user, claims["id"])

		if user.Id == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("user", user)

		c.Next()
	}

	//find the user wth token sub
	//attach to req

	//continue

}

func Validate(c *gin.Context) {

	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{"message": user})

}

func Register(c *gin.Context) {
	var data map[string]string

	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)

	if err != nil {
		panic("cannot accept the password")
	}

	user := models.User{
		Username: data["username"],
		Password: password,
	}

	models.DB.Create(&user)

	c.JSON(http.StatusOK, gin.H{"user": user})

}

func Login(c *gin.Context) {
	var data map[string]string

	if err := c.ShouldBindJSON(&data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	var user models.User

	models.DB.Where("username = ?", data["username"]).First(&user)

	if user.Id == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "username or password incorrect"})
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = user.Id
	claims["username"] = user.Username
	claims["classes"] = user.Classroom
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "could not generate token"})
	}
	c.SetCookie("jwt", tokenString, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "sucsses"})
}

func Logout(c *gin.Context) {

	_, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you are not login"})

		return
	}

	expiredCookie := &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, expiredCookie)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out succsesfully"})

}
