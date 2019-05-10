package main

import (
	"github.com/gin-gonic/gin"
	parcelServe "github.com/moqmar/parcel-serve"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

var session *r.Session

func main() {
	loadConfig()
	migrate()

	router := gin.Default()
	router.POST("/api/login", login)
	router.GET("/api/logout", logout)
	router.POST("/api/register", register)
	router.GET("/api/reise", meineReisen)
	router.POST("/api/reise", neueReise)
	parcelServe.Serve("frontend", router, nil, nil)
	router.Run(":2222")

}
