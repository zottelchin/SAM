package main

import (
	"codeberg.org/momar/logg"
	"github.com/gin-gonic/gin"
	parcelServe "github.com/moqmar/parcel-serve"
)

func main() {
	logg.Info("SAM startet jetzt... \n Builddate: %s \n Commit: %s\n Version: %s", buildDate, gitHash, version)
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
	router.POST("/api/reisen/:id", reiseBearbeiten)

	router.POST("/api/reisen/:id/mitreisende", reisePersonHinzufügen)
	router.PUT("/api/reisen/:id/mitreisende", reiseAnzeigeNutzerhinzufügen)
	router.DELETE("/api/reisen/:id/mitreisende/:mail", reisePersonEntfernen)

	router.PUT("/api/reisen/:id/beleg", neuerBeleg)
	router.PUT("/api/reisen/:id/beleg/:bid", belegBearbeiten)
	router.DELETE("/api/reisen/:id/beleg/:bid", belegLöschen)

	router.GET("/api/share/:key/*pass", geteilteReise)
	router.PUT("/api/reisen/:id/share", reiseTeilen)

	router.GET("/api/version", func(c *gin.Context) {
		c.JSON(418, gin.H{
			"version": version,
			"build":   buildDate,
			"commit":  gitHash,
		})
	})

	parcelServe.Serve("frontend", router, AssetNames(), MustAsset)
	router.Run(":2222")

}
