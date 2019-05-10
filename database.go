package main

import (
	"time"

	"github.com/jinzhu/gorm"

	log "codeberg.org/momar/logg"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func migrate() {
	var err error
	if db, err = gorm.Open("sqlite3", "./storage.db"); err != nil {
		log.Error("%s", err)
	}
	db.SingularTable(true)

	if e := db.AutoMigrate(&Nutzer{}, &Beleg{}, &Reise{}).Error; e != nil {
		log.Error("%s", e)
	} else {
		log.Info("Tabellen Nutzer, Belege und Reisen wurde erstellt oder war schon vorhanden.")
	}
}

type Reise struct {
	ID          int
	Name        string   `json:"name"`
	Mitreisende []Nutzer `json:"mitreisende" gorm:"many2many:Nutzer_Reisen;"`
	Ausgaben    []Beleg  `json:"belege" gorm:"foreignkey:ReiseID"`
}

type Nutzer struct {
	ID          int     `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Name        string  `json:"name"`
	Mail        string  `json:"mail"`
	Password    string  `json:"-"`
	MeineReisen []Reise `gorm:"many2many:Nutzer_Reisen;"`
}

type Beleg struct {
	ID      int       `gorm:"primary_key" json:"id"`
	Name    string    `json:"name"`
	Datum   time.Time `json:"datum"`
	Betrag  float32   `json:"betrag"`
	ReiseID uint
}

func (Beleg) TableName() string {
	return "Belege"
}

func (Nutzer) TableName() string {
	return "Nutzer"
}

func (Reise) TableName() string {
	return "Reisen"
}
