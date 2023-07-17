package main

import (
	"github.com/Kazukite12/go-chat/models"
	"github.com/Kazukite12/go-chat/routes"
)

func main() {

	models.ConnectDB()
	routes.Routes()
}
