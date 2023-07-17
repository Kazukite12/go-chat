package routes

import (
	"net/http"

	userController "github.com/Kazukite12/go-chat/controllers/UserController"
	"github.com/gin-gonic/gin"
)

func Routes() {
	router := gin.Default()

	router.GET("/api/user/validate", userController.Auth, userController.Validate)
	router.POST("api/user/register", userController.Register)
	router.POST("api/user/login", userController.Login)
	router.POST("api/user/logout", userController.Logout)

	http.ListenAndServe("localhost:8080", router)

}
