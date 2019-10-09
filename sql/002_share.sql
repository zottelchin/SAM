-- +goose Up

CREATE TABLE IF NOT EXISTS Freigaben (
    schluessel VARCHAR PRIMARY KEY,
    passwort VARCHAR, 
    reise_id INTEGER NOT NULL,
    FOREIGN KEY (reise_id) REFERENCES Reisen (id) ON DELETE CASCADE
);