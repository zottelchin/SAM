package main

import (
	"math/rand"
	"strings"
	"time"

	"codeberg.org/momar/logg"

	log "codeberg.org/momar/logg"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

//// Gin Handler

func register(c *gin.Context) {
	log.Info("Ein neuer Nutzer möchte sich registrieren")
	u := Register{}
	c.BindJSON(&u)
	if Config.Key_required && u.Key != Config.Register_Key {
		log.Info("Es wird ein Registrierungsschlüssel benötigt, aber der mitgesendete stimmt nicht mit unserem überein")
		c.Status(403)
		c.Abort()
		return
	}
	u.Password = hashPassword(u.Password)
	v := Nutzer{Password: u.Password, Name: u.Name, Mail: u.Mail}
	v, e := createNutzer(v)
	if e != nil {
		log.Error("Registrierung fehlgeschlagen: %s", e)
		c.Status(500)
		c.Abort()
		return
	}
	log.Info("Neuer Nutzer mit der ID %d hat sich registriert", v.ID)
	c.Status(200)
}

func login(c *gin.Context) {
	l := Login{}
	c.BindJSON(&l)
	u, err := getNutzer(Nutzer{Mail: l.Mail})
	if err != nil {
		log.Error("Fehler beim Abrufen des Nutzers aus der Datenbank: %s", err)
		c.Status(500)
		c.Abort()
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(l.Passwort)); err == nil {
		c.JSON(200,
			gin.H{
				"hash": addHash(u),
			})
		return
	}
	log.Warn("Falsches Password für %s", l.Mail)
	c.String(401, "Login fehlgeschlagen. Bitte überprüfen Sie Ihre Logindaten.")
}

func logout(c *gin.Context) {
	removeHash(c)
	c.Status(200)
}

//// Helperfunctions

func randomString(l int) string {
	pool := "qwertzuiopasdfghjklyxcvbnmQWERTZUIOPASDFGHJKLYXCVBNM<>.,;-_:()0123456789"
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < l; i++ {
		s += string(pool[rand.Intn(len(pool))])
	}
	return s
}

type Login struct {
	Mail     string `json:"mail"`
	Passwort string `json:"password"`
}

type Register struct {
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
	Key      string `json:"key"`
}

//// Hash Password

func hashPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Error("%s", err)
	}
	return string(hash)
}

//// Hash Map and Hash Functions

var hashes = map[string]int{"abc": 1}

func addHash(u Nutzer) string {
	for {
		newHash := randomString(55)
		if hashes[newHash] == 0 {
			hashes[newHash] = u.ID
			return newHash
		}
	}
}

func removeHash(c *gin.Context) {
	hash := c.GetHeader("Authorization")
	hash = strings.TrimPrefix(hash, "Bearer ")
	delete(hashes, hash)
}

func userByHash(c *gin.Context) Nutzer {
	hash := c.GetHeader("Authorization")
	hash = strings.TrimPrefix(hash, "Bearer ")
	id := hashes[hash]
	u, err := getNutzer(Nutzer{ID: id})
	if err != nil {
		logg.Error("Fehler beim abrufen des Nutzers %d: %s", id, err)
	}
	return u
}

func isLoggedIn(c *gin.Context) bool {
	hash := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if hash != "" && hashes[hash] != 0 {
		return true
	}
	return false
}
