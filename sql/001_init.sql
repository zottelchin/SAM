-- +goose Up

CREATE TABLE IF NOT EXISTS Reisen (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    name VARCHAR, 
    archiviert bool DEFAULT 0
);

CREATE TABLE IF NOT EXISTS Nutzer (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    name VARCHAR, 
    mail VARCHAR UNIQUE, 
    password VARCHAR,
    reise INTEGER --Hier wird eine Reise ID eingefügt, wenn der Nutzer nur exitiert um angezeigt zu werden, aber kein echter Nutzer ist 
);

CREATE TABLE IF NOT EXISTS Belege (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR,
    datum VARCHAR,
    betrag INTEGER,
    reise_id INTEGER,
    von INTEGER,
    archiviert BOOL DEFAULT 0,
    FOREIGN KEY (reise_id) REFERENCES Reisen (id) ON DELETE CASCADE,
    FOREIGN KEY (von) REFERENCES Nutzer (id)
);

CREATE TABLE IF NOT EXISTS Nutzer_Reisen (
    nutzer_id INTEGER,
    reise_id INTEGER,
    PRIMARY KEY (nutzer_id, reise_id),
    FOREIGN KEY (nutzer_id) REFERENCES Nutzer (id),
    FOREIGN KEY (reise_id) REFERENCES Reisen (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bezahlt_fuer (
    beleg_id INTEGER,
    nutzer_id INTEGER,
    PRIMARY KEY (beleg_id, nutzer_id),
    FOREIGN KEY (beleg_id) REFERENCES Belege (id),
    FOREIGN KEY (nutzer_id) REFERENCES Nutzer (id)
);

INSERT OR IGNORE INTO sqlite_sequence (name, seq) VALUES ('Reisen', 0);
INSERT OR IGNORE INTO sqlite_sequence (name, seq) VALUES ('Nutzer', 0);
INSERT OR IGNORE INTO sqlite_sequence (name, seq) VALUES ('Belege', 0);