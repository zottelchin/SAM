package main

import (
	"math/rand"
	"strings"
	"time"

	"codeberg.org/momar/ternary"

	log "codeberg.org/momar/logg"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	if isMailInUse(u.Mail) {
		log.Warn("Diese Email ist schon in der Datenbank vorhanden.")
		c.String(400, "Diese Email ist schon in der Datenbank vorhanden.")
		c.Abort()
		return
	}
	u.Password = hashPassword(u.Password)
	v := Nutzer{Password: u.Password, Name: u.Name, Mail: u.Mail}
	if e := db.Create(&v).Error; e != nil {
		log.Error("Registrierung fehlgeschlagen: %s", e)
		c.Status(500)
		c.Abort()
		return
	}

	c.Status(200)
	return

}

func login(c *gin.Context) {
	l := Login{}
	c.BindJSON(&l)
	u := Nutzer{}
	if e := db.First(&u, "mail = ?", l.Mail).Error; e != nil {
		if e == gorm.ErrRecordNotFound {
			log.Error("Nutzer ist nicht in der Datenbank. %s", l.Mail)
			c.String(401, "Login fehlgeschlagen. Bitte überprüfen Sie Ihre Logindaten.")
			c.Abort()
			return
		}
		log.Error("%s", e)
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

func isMailInUse(mail string) bool {
	var c int
	if e := db.Model(&Nutzer{}).Where(Nutzer{Mail: mail}).Count(&c).Error; e != nil {
		log.Error("Prüfen ob die Mail in der Datenbank ist, hat einen Fehler zurückgegeben: %s", e)
	}
	return ternary.If(c == 0, false, true).(bool)
}

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

var hashes = map[string]int{"abc": 4}

func addHash(u Nutzer) string {
	log.Info("Füge neuen Hash für %s hinzu", u.Name)
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
	log.Info("Lösche Hash %s für %s", hash, userByHash(c))
	hash = strings.TrimPrefix(hash, "Bearer ")
	delete(hashes, hash)
}

func userByHash(c *gin.Context) Nutzer {
	hash := c.GetHeader("Authorization")
	hash = strings.TrimPrefix(hash, "Bearer ")
	id := hashes[hash]
	u := Nutzer{}
	db.First(&u, id)
	return u
}

func isLoggedIn(c *gin.Context) bool {
	log.Info("Teste ob Nutzer eingelogged ist")
	hash := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	log.Info("Hash ist: %s", hash)
	if hash != "" && hashes[hash] != 0 {
		log.Info("Nutzer ist angemeldet")
		return true
	}
	log.Info("Nutzer ist nicht angemeldet")
	return false
}
