package main

import (
	"strconv"
	"strings"

	"codeberg.org/momar/logg"
	"github.com/gin-gonic/gin"
)

func geteilteReise(c *gin.Context) {
	key := c.Param("key")
	pass := strings.TrimPrefix(c.Param("pass"), "/")
	logg.Info("key: %s \n pass: %s", key, pass)
	exists, passwort, protected, id := getShare(key)
	if !exists {
		c.Status(404)
		return
	}
	if protected && pass == "" {
		c.Status(402)
		return
	}
	if protected && pass != passwort {
		c.Status(401)
		return
	}
	reise, err := getReise(id)
	if err != nil {
		logg.Error(err.Error())
		c.Status(500)
		return
	}
	c.JSON(200, reise)
}

func reiseTeilen(c *gin.Context) {
	Body := struct {
		Key string
	}{}
	c.BindJSON(&Body)
	id, _ := strconv.Atoi(c.Param("id"))
	erstelleFreigabe(id,Body.Key)
}
