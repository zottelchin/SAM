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
	router.POST("/api/logout", logout)
	router.POST("/api/register", register)
	router.GET("/api/user", listUser)

	router.GET("/api/reisen", meineReisen) //alle reisen
	router.PUT("/api/reisen", neueReise)   // Reise Anlegen
	router.GET("/api/reisen/:id", reiseNummer)
	router.PUT("/api/reisen/:id", reiseBearbeiten)

	router.PUT("/api/reisen/:id/mitreisende", reisePersonHinzufügen)
	router.DELETE("/api/reisen/:id/mitreisende", reisePersonEntfernen)

	router.PUT("/api/reisen/:id/beleg", neuerBeleg)
	router.PUT("/api/reisen/:id/beleg/:bid", belegBearbeiten)
	router.DELETE("/api/reisen/:id/beleg/:bid", belegLöschen)

	parcelServe.Serve("frontend2", router, nil, nil)
	router.Run(":2222")

}
