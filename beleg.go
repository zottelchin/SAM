package main

import (
	"strconv"

	"codeberg.org/momar/logg"
	"github.com/gin-gonic/gin"
)

func neuerBeleg(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}
	beleg := Beleg{}
	err := c.BindJSON(&beleg)
	logg.Error("%s", err)
	id, _ := strconv.Atoi(c.Param("id"))
	r, err := createBeleg(beleg, Reise{ID: id})
	logg.Info("%v", r)
	if err != nil {
		logg.Error("Fehler beim erstellen des neuen Belegs: %s", err)
		c.Status(500)
		c.Abort()
		return
	}
	c.JSON(200, r)
}

func belegBearbeiten(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}
	id, _ := strconv.Atoi(c.Param("bid"))
	beleg := Beleg{ID: id}
	c.BindJSON(&beleg)
	r, err := updateBeleg(beleg)
	if err != nil {
		logg.Error("Fehler beim Beleg bearbeiten: %s", err)
		c.Status(500)
		c.Abort()
		return
	}
	c.JSON(200, r)
}

func belegLöschen(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}
	id, _ := strconv.Atoi(c.Param("bid"))
	r, err := archivBeleg(Beleg{ID: id})
	if err != nil {
		logg.Error("Fehler beim Beleg löschen: %s", err)
		c.Status(500)
		c.Abort()
		return
	}
	c.JSON(200, r)
}
