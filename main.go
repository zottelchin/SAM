package main

import (
	"github.com/gin-gonic/gin"
	parcelServe "github.com/moqmar/parcel-serve"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

var session *r.Session

func main() {
	loadConfig()

	router := gin.Default()
	router.GET("/api/config", func(c *gin.Context) { c.JSON(200, Config.Key_required) })

	router.POST("/api/login", login)
	router.POST("/api/logout", logout)
	router.POST("/api/register", register)
	router.GET("/api/me", me)

	router.GET("/api/reisen", meineReisen) //alle reisen
	router.PUT("/api/reisen", neueReise)   // Reise Anlegen
	router.GET("/api/reisen/:id", reiseNummer)
	router.PUT("/api/reisen/:id", reiseBearbeiten)

	router.PUT("/api/reisen/:id/mitreisende", reisePersonHinzufügen)
	router.DELETE("/api/reisen/:id/mitreisende", reisePersonEntfernen)

	router.PUT("/api/reisen/:id/beleg", neuerBeleg)
	router.PUT("/api/reisen/:id/beleg/:bid", belegBearbeiten)
	router.DELETE("/api/reisen/:id/beleg/:bid", belegLöschen)

	parcelServe.Serve("frontend", router, AssetNames(), MustAsset)
	router.Run(":2222")

}
