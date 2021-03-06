package main

import (
	"strconv"

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
	reisen, err := getNutzersReisen(u)
	if err != nil {
		logg.Error("Fehler beim Laden der Reisen für %s: %s", u.Mail, err)
		c.Status(500)
		c.Abort()
		return
	}
	c.JSON(200, reisen)
}

func neueReise(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}

	r := Reise{}
	c.BindJSON(&r)
	r.Mitreisende = append(r.Mitreisende, userByHash(c))
	r, err := createReise(r)
	if err != nil {
		logg.Error("Fehler beim erstellen einer neuen Reise %s", err)
		c.Status(500)
		c.Abort()
		return
	}
	err = addNutzerReise(userByHash(c), r)
	if err != nil {
		logg.Error("Fehler beim hinzugügen des erstellers zur Reise: %s", err)
		c.Status(500)
		c.Abort()
		return
	}
	c.JSON(200, r)
}

func reiseBearbeiten(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	änderungen := Reise{ID: id}
	c.BindJSON(&änderungen)
	reise, err := updateReise(änderungen)
	if err != nil {
		logg.Error("Fehler beim ändern der Reise %d: %s", id, err)
		c.Status(500)
		c.Abort()
		return
	}
	c.JSON(200, reise)
}

func reiseNummer(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	if !istDabei(userByHash(c), Reise{ID: id}) {
		c.Status(403)
		c.Abort()
		return
	}
	reise, err := getReise(id)
	if err != nil {
		logg.Error("Fehler beim Abrufen der Reise %d: %s", id, err)
	}
	c.JSON(200, reise)
}

func reisePersonHinzufügen(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	u := Nutzer{}
	err := c.BindJSON(&u)
	if err != nil {
		logg.Error("Fehler beim Binden des Nutzers %s, um eine Person zu einer Reise hinzu zu fügen: %s", u.Mail, err)
		c.String(500, err.Error())
		c.Abort()
		return
	}
	u, err = getNutzer(u)
	if err != nil {
		logg.Error("Fehler beim Abrufen des NUtzers %d, um eine Person zu einer Reise hinzu zu fügen: %s", id, err)
		c.String(500, err.Error())
		c.Abort()
		return
	}
	r, err := getReise(id)
	if err != nil {
		logg.Error("Fehler beim Abrufen der Reise %d, um eine Person hinzu zu fügen: %s", id, err)
		c.String(500, err.Error())
		c.Abort()
		return
	}
	err = addNutzerReise(u, r)
	if err == nil {
		r.Mitreisende = append(r.Mitreisende, u)
		c.JSON(200, r)
		return
	}
	logg.Error("Nutzer konnte nicht zu Reise hinzugefügt werden: %s", err)
	c.String(500, err.Error())
	c.Abort()
}

func reisePersonEntfernen(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	u := Nutzer{Mail: c.Param("mail")}
	if u.Mail == "" {
		logg.Error("Nutzer nicht Spezifiziert")
		c.String(400, "Nutzer nicht spezifiziert")
		c.Abort()
		return
	}
	r, err := getReise(id)
	if err != nil {
		logg.Error("Fehler beim abrufen der Reise %d um einen Nutzer zu entfernen: %s", id, err)
		c.String(400, "%s", err)
		c.Abort()
		return
	}
	err = removeNutzerReise(u, r)
	if err != nil {
		logg.Error("Fehler beim entfernen einer Peson aus der Reise %d: %s", id, err)
		c.String(500, "%s", err)
		c.Abort()
		return
	}
	r, _ = getReise(id)
	c.JSON(200, r)

}

func reiseAnzeigeNutzerhinzufügen(c *gin.Context) {
	if !isLoggedIn(c) {
		c.Status(401)
		c.Abort()
		return
	}
	logg.Info("Nutzer hat Rechte")
	id, _ := strconv.Atoi(c.Param("id"))
	r, err := getReise(id)
	if err != nil {
		logg.Error("Fehler beim Abrufen der Reise %d, um eine Person hinzu zu fügen: %s", id, err)
		c.String(500, err.Error())
		c.Abort()
		return
	}
	logg.Info("Reise %d konnte abgerufen werden", id)
	user := Nutzer{}
	c.BindJSON(&user)
	user.FürReise, _ = strconv.Atoi(c.Param("id"))
	user, err = createAnzeigeNutzer(user)
	if err != nil {
		logg.Error("Fehler beim erstellen eines Dummy Nutzers: %s", err)
		c.String(500, "%s", err)
		c.Abort()
		return
	}
	logg.Info("Nutzer wurde erstellt: %v+", user)
	err = addNutzerReise(user, r)
	if err != nil {
		logg.Error("Fehler beim erstellen eines Dummy Nutzers: %s", err)
		c.String(500, "%s", err)
		c.Abort()
		return
	}
	r, _ = getReise(user.FürReise)
	c.JSON(200, r)
}
