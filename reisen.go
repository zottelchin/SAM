package main

import (
	"codeberg.org/momar/logg"
	"github.com/gin-gonic/gin"
)

func meineReisen(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}

	u := userByHash(c)
	reisen := []Reise{}
	logg.Info("%v", u)
	db.Model(&u).Preload("Mitreisende").Preload("Ausgaben").Related(&reisen, "MeineReisen")

	c.JSON(200, reisen)
}

func neueReise(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}

	u := userByHash(c)
	r := Reise{}
	r.Mitreisende = append(r.Mitreisende, u)
	c.BindJSON(&r)
	if db.Create(&r).Error != nil {
		logg.Error("Fehler beim erstellen der neuen Reise")
		c.Status(500)
		c.Abort()
		return
	}
	c.JSON(200, r)

}
