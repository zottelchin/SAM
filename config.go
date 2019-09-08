package main

import (
	"fmt"

	"codeberg.org/momar/logg"
	"codeberg.org/momar/ternary"
	"github.com/jinzhu/configor"
)

func loadConfig() {
	if err := configor.Load(&Config, "./data/config.yml"); err != nil {
		logg.Error("Laden der Konfigurationsdatei ist fehlgeschlagen %s", err)
	}
	result := ""
	result += fmt.Sprintln("Die Konfiguration wurde geladen:")
	result += fmt.Sprintln(ternary.If(Config.Key_required, fmt.Sprintf(" - Es wird ein Registrierungsschlüssel benötigt und der lautet: %s", Config.Register_Key), " - Es wird kein Registrierungsschlüssel benötigt").(string))
	logg.Info(result)
}

var Config = struct {
	Register_Key string
	Key_required bool

	Mail struct {
		Host     string
		Port     int
		Login    string
		Password string
		Sender   string
	}
}{}
