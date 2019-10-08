package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	"codeberg.org/momar/logg"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

var db = initDB()

func initDB() *sql.DB {
	goose.SetDialect("sqlite3")
	db, err := sql.Open("sqlite3", "./data/storage.db")
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db is nil")
	}

	if err := goose.Up(db, "./sql"); err != nil {
		logg.Error("Error in Database migration: %s", err)
		os.Exit(1)
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
	Mail     string `json:"mail"`
	Password string `json:"-"`
	FürReise int    `json:"fürReise"`
}

type Beleg struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Datum    string   `json:"datum"`
	Betrag   int      `json:"betrag"`
	ReiseID  int      `json:"-"`
	Von      Nutzer   `json:"von"`
	An       []Nutzer `json:"an"`
	Gelöscht bool     `json:"gelöscht"`
}

// All diese Funktionen führen nur stumpf ihre Aktion aus und testen nicht, ob Berechtigungen vorhanden sind.

// getNutzer gibt den Nutzer mit der id oder mail zurück (ohne seine Reisen)
func getNutzer(n Nutzer) (Nutzer, error) {
	if n.ID == 0 && n.Mail == "" {
		return Nutzer{}, errors.New("Nutzer nicht spezifiziert")
	}
	var (
		reise    sql.NullInt64
		password sql.NullString
	)
	if n.ID != 0 {
		err := db.QueryRow("Select name, mail, password, reise FROM Nutzer WHERE id = ?", n.ID).Scan(&n.Name, &n.Mail, &password, &reise)
		if err != nil {
			logg.Error(err.Error())
			return Nutzer{}, err
		}
		if password.Valid {
			n.Password = password.String
		}
		if reise.Valid {
			n.FürReise = int(reise.Int64)
		}
		return n, nil
	}
	err := db.QueryRow("SELECT name, id, password, reise FROM Nutzer WHERE mail = ?", n.Mail).Scan(&n.Name, &n.ID, &password, &reise)
	if err != nil {
		return Nutzer{}, err
	}
	if password.Valid {
		n.Password = password.String
	}
	if reise.Valid {
		n.FürReise = int(reise.Int64)
	}
	return n, nil
}

// createNutzer erstellt einen neuen Nutzer mit den Angaben
// Das Passwort muss schon gehashed sein
func createNutzer(n Nutzer) (Nutzer, error) {
	if n.Mail == "" || n.Name == "" || n.Password == "" {
		return Nutzer{}, errors.New("Nutzer nicht genau spezifiziert")
	}
	_, err := db.Exec("Insert into Nutzer (name, mail, password) Values (?, ?, ?)", n.Name, n.Mail, n.Password)
	if err != nil {
		return Nutzer{}, err
	}
	db.QueryRow("Select id From Nutzer Where mail = ?", n.Mail).Scan(&n.ID)
	return n, nil
}

// createAnzeigeNutzer erstellt einen Nutzer für eine Reise id, sodass sich nicht alle sich einen Account erstellen müssen,
// Er wird auch gleich zu der Reise hinzu gefügt.
func createAnzeigeNutzer(n Nutzer) (Nutzer, error) {
	if n.Name == "" || n.FürReise == 0 {
		return Nutzer{}, errors.New("Nutzer nicht genau spezifiziert " + n.Name + " - " + strconv.Itoa(n.FürReise))
	}
	n.Mail = generateDummyMail()
	res, err := db.Exec("Insert into Nutzer (name, mail, reise) Values (?, ?, ?)", n.Name, n.Mail, n.FürReise)
	if err != nil {
		return Nutzer{}, err
	}
	id64, _ := res.LastInsertId()
	n.ID = int(id64)
	return n, nil
}

// updateNutzer ändert die nicht leeren Felder
func updateNutzer(n Nutzer) (Nutzer, error) {
	return n, errors.New("Noch nicht implementiert")
}

// deleteNutzer löscht den Nutzer
func deleteNutzer(n Nutzer) error {
	if n.ID == 0 && n.Mail == "" {
		return errors.New("Nutzer nicht spezifiziert")
	}
	if n.ID != 0 {
		err := db.QueryRow("Delete FROM Nutzer WHERE id = ?", n.ID).Scan(&n.Name, &n.Mail, &n.Password)
		if err != nil {
			return err
		}
		return nil
	}
	err := db.QueryRow("Delete FROM Nutzer WHERE mail = ?", n.ID).Scan(&n.Name, &n.ID, &n.Password)
	if err != nil {
		return err
	}
	return nil
}

// getNutzersReisen gibt die Reisen zurück, in denen ein Nutzer spezifiziert ist
func getNutzersReisen(n Nutzer) ([]Reise, error) {
	if n.ID == 0 {
		return nil, errors.New("Nutzer nicht spezifiziert")
	}
	rows, err := db.Query("Select id, name, archiviert From Reisen inner join Nutzer_Reisen on id = reise_id where nutzer_id = ?", n.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	e := ""
	reisen := []Reise{}
	for rows.Next() {
		tmp := Reise{}
		if err := rows.Scan(&tmp.ID, &tmp.Name, &tmp.Archiviert); err != nil {
			e += fmt.Sprintln(err)
		}

		rowsU, _ := db.Query("Select nutzer_id from Nutzer_Reisen Where reise_id = ?", tmp.ID)
		for rowsU.Next() {
			userID := 0
			rowsU.Scan(&userID)
			u, _ := getNutzer(Nutzer{ID: userID})
			tmp.Mitreisende = append(tmp.Mitreisende, u)
		}

		reisen = append(reisen, tmp)
	}
	if e != "" {
		return reisen, errors.New(e)
	}
	return reisen, nil
}

// createReise erstellt eine neue Reise mit den gegegebenen Angaben
func createReise(r Reise) (Reise, error) {
	res, err := db.Exec("Insert into Reisen (name) Values (?)", r.Name)
	if err != nil {
		return Reise{}, err
	}
	id, _ := res.LastInsertId()
	r.ID = int(id)
	return r, nil
}

// gibt die Reise mit der ID zurück
func getReise(id int) (Reise, error) {
	if id == 0 {
		return Reise{}, errors.New("Reise ID nicht spezifiziert")
	}
	r := Reise{}
	err := db.QueryRow("Select id, name, archiviert FROM Reisen Where id = ?", id).Scan(&r.ID, &r.Name, &r.Archiviert)
	if err != nil {
		return Reise{}, err
	}
	rows, _ := db.Query("Select nutzer_id from Nutzer_Reisen Where reise_id = ?", r.ID)
	for rows.Next() {
		userID := 0
		rows.Scan(&userID)
		u, _ := getNutzer(Nutzer{ID: userID})
		r.Mitreisende = append(r.Mitreisende, u)
	}
	rows, _ = db.Query("Select id, name, datum, betrag, von from Belege Where reise_id = ? and archiviert = false", r.ID)
	for rows.Next() {
		tmp := Beleg{}
		rows.Scan(&tmp.ID, &tmp.Name, &tmp.Datum, &tmp.Betrag, &tmp.Von.ID)
		tmp.Von, _ = getNutzer(Nutzer{ID: tmp.Von.ID})

		userRows, _ := db.Query("Select nutzer_id from bezahlt_fuer Where beleg_id = ?", tmp.ID)
		for userRows.Next() {
			userID := 0
			userRows.Scan(&userID)
			u, _ := getNutzer(Nutzer{ID: userID})
			tmp.An = append(tmp.An, u)
		}
		r.Ausgaben = append(r.Ausgaben, tmp)
	}
	return r, nil
}

// istDabei gibt WAHR zurück, wenn jemand bei einer Reise dabei ist
func istDabei(n Nutzer, r Reise) bool {
	res := 0
	db.QueryRow("Select Exists(Select 1 From Nutzer_Reisen Where nutzer_id = ? and reise_id = ?)", n.ID, r.ID).Scan(&res)
	if res == 1 {
		return true
	}
	return false
}

// updateReise ändert die gesetzten Parameter
func updateReise(r Reise) (Reise, error) {
	if r.Name != "" {
		_, err := db.Exec("Update Reisen Set name = ? Where id = ?", r.Name, r.ID)
		if err != nil {
			r, _ = getReise(r.ID)
			return r, err
		}
	}

	if r.Archiviert {
		_, err := archivReise(r)
		if err != nil {
			r, _ = getReise(r.ID)
			return r, err
		}
	}
	r, _ = getReise(r.ID)
	return r, nil
}

// archivReise setzt das Archivierungs Bit auf Wahr, sodass die Reise als abgeschlossen angezeigt wird
func archivReise(r Reise) (Reise, error) {
	_, err := db.Exec("Update Reisen Set archiviert = true Where id = ?", r.ID)
	if err != nil {
		return r, err
	}
	r.Archiviert = true
	return r, err
}

// addNutzerReise fügt einen Nutzer als Mitreisenden hinzu
func addNutzerReise(n Nutzer, r Reise) error {
	_, err := db.Exec("Insert Into Nutzer_Reisen (nutzer_id, reise_id) Values (?, ?)", n.ID, r.ID)
	return err
}

// removeNutzerReise entfernt einen Nutzer aus der Reihe der Mitreisenden
func removeNutzerReise(n Nutzer, r Reise) error {
	n, _ = getNutzer(n)
	res, err := db.Exec("Delete From Nutzer_Reisen Where nutzer_id = ? and reise_id = ?", n.ID, r.ID)
	rows, _ := res.RowsAffected()
	if err == nil && int(rows) != 1 {
		err = errors.New("Nutzer nicht beteiligt")
	}
	return err
}

// erstellt einen neun Beleg für eine Reise
func createBeleg(b Beleg, r Reise) (Beleg, error) {
	userID, _ := getNutzer(b.Von)
	b.Von = userID
	res, err := db.Exec("Insert into Belege (name, datum, betrag, reise_id, von) Values (?, ?, ?, ?, ?)", b.Name, b.Datum, b.Betrag, r.ID, userID.ID)
	if err != nil {
		return Beleg{}, err
	}
	id, _ := res.LastInsertId()
	b.ID = int(id)
	an := []Nutzer{}
	for _, x := range b.An {
		id, _ := getNutzer(x)
		an = append(an, id)
		db.Exec("Insert Into bezahlt_fuer (beleg_id, nutzer_id) Values (?, ?)", b.ID, id.ID)
	}
	b.An = an
	return b, nil
}

// ändert einen bestehenden Beleg
func updateBeleg(b Beleg) (Beleg, error) {
	old := Beleg{}
	err := db.QueryRow("Select name, datum, betrag, von from Belege Where id = ? and archiviert = false", b.ID).Scan(&old.Name, &old.Datum, &old.Betrag, &old.Von.ID)
	if err != nil {
		logg.Error(err.Error())
		return Beleg{}, err
	}
	if old.Name != b.Name && b.Name != "" {
		_, err := db.Exec("Update Belege Set name = ? Where id = ?", b.Name, b.ID)
		if err != nil {
			return b, err
		}
	}
	if old.Betrag != b.Betrag && b.Betrag != 0 {
		_, err := db.Exec("Update Belege Set betrag = ? Where id = ?", b.Betrag, b.ID)
		if err != nil {
			return b, err
		}
	}
	if old.Datum != b.Datum && b.Datum != "" {
		_, err := db.Exec("Update Belege Set datum = ? Where id = ?", b.Datum, b.ID)
		if err != nil {
			return b, err
		}
	}

	von, _ := getNutzer(Nutzer{Mail: b.Von.Mail})
	if old.Von.ID != von.ID && von.ID != 0 {
		_, err := db.Exec("Update Belege Set von = ? Where id = ?", von.ID, b.ID)
		if err != nil {
			return b, err
		}
	}

	db.Exec("DELETE FROM bezahlt_fuer WHERE beleg_id = ?", b.ID)
	for _, n := range b.An {
		n, _ = getNutzer(n)
		db.Exec("Insert Into bezahlt_fuer (beleg_id, nutzer_id) Values (?, ?)", b.ID, n.ID)
	}
	return Beleg{}, nil
}

// setzt das Archivierungsbit für einen Beleg, sodass dieser nicht mehr angezeigt wird
func archivBeleg(b Beleg) (Beleg, error) {
	_, err := db.Exec("Update Belege Set archiviert = true Where id = ?", b.ID)
	if err != nil {
		return b, err
	}
	b.Gelöscht = true
	return b, err
}

func generateDummyMail() string {
	var mail string
	for {
		mail = randomString(11) + "@sam"
		_, err := getNutzer(Nutzer{Mail: mail})
		if err == sql.ErrNoRows {
			break
		}
	}
	return mail
}
