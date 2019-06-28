package main

import (
	"database/sql"

	"codeberg.org/momar/logg"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pressly/goose"
)

var db = initDB()

func initDB() *sql.DB {
	goose.SetDialect("sqlite3")
	if db, err := sql.Open("sqlite3", "./data/storage.db"); err != nil {
		panic(err)
	}
	if db == nil {
		panic("db is nil")
	}

	if err := goose.Up(db, "./sql"); err != nil {
		logg.Error("Error in Database migratoin: %s", err)
	}

	return db
}

type Reise struct {
	ID          int
	Name        string   `json:"name"`
	Mitreisende []Nutzer `json:"mitreisende"`
	Ausgaben    []Beleg  `json:"belege"`
	Archiviert  bool     `json:"archiviert"`
}

type Nutzer struct {
	ID       int    `json:"-"`
	Name     string `json:"name"`
	Mail     string `json:"mail" gorm:"unique;not null"`
	Password string `json:"-"`
}

type Beleg struct {
	ID       int      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	Name     string   `json:"name"`
	Datum    string   `json:"datum"`
	Betrag   float64  `json:"betrag"`
	ReiseID  int      `json:"-"`
	Von      Nutzer   `json:"von" gorm:"foreignkey:ID"`
	An       []Nutzer `gorm:"many2many:bezahlt_fuer" json:"an"`
	Gelöscht bool     `json:"gelöscht"`
}

// All diese Funktionen führen nur stumpf ihre Aktion aus und testen nicht, ob Berechtigungen vorhanden sind.

// getNutzer gibt den Nutzer mit der id oder mail zurück (ohne seine Reisen)
func getNutzer(n Nutzer) (Nutzer, error) {}

// createNutzer erstellt einen neuen Nutzer mit den Angaben
func createNutzer(n Nutzer) (Nutzer, error) {}

// createAnzeigeNutzer erstellt einen Nutzer für eine Reise id, sodass sich nicht alle sich einen Account erstellen müssen,
// Er wird auch gleich zu der Reise hinzu gefügt.
func createAnzeigeNutzer(n Nutzer) (Nutzer, error) {}

// updateNutzer ändert die nicht leeren Felder
func updateNutzer(n Nutzer) (Nutzer, error) {}

// deleteNutzer löscht den Nutzer
func deleteNutzer(n Nutzer) error {}

// getNutzersReisen gibt die Reisen zurück, in denen ein Nutzer spezifiziert ist
func getNutzersReisen(n Nutzer) ([]Reise, error) {}

// createReise erstellt eine neue Reise mit den gegegebenen Angaben
func createReise(r Reise) (Reise, error) {}

// istDabei gibt WAHR zurück, wenn jemand bei einer Reise dabei ist
func istDabei(n Nutzer, r Reise) bool {}

// updateReise ändert die gesetzten Parameter
func updateReise(r Reise) (Reise, error) {}

// archivReise setzt das Archivierungs Bit auf Wahr, sodass die Reise als abgeschlossen angezeigt wird
func archivReise(r Reise) (Reise, error) {}

// addNutzerReise fügt einen Nutzer als Mitreisenden hinzu
func addNutzerReise(n Nutzer, r Reise) error {}

// removeNutzerReise entfernt einen Nutzer aus der Reihe der Mitreisenden
func removeNutzerReise(n Nutzer, r Reise) error {}

// erstellt einen neun Beleg für eine Reise
func createBeleg(b Beleg, r Reise) (Reise, error) {}

// ändert einen bestehenden Beleg
func updateBeleg(b Beleg) (Beleg, error) {}

// setzt das Archivierungsbit für einen Beleg, sodass dieser nicht mehr angezeigt wird
func archivBeleg(b Beleg) (Beleg, error) {}
