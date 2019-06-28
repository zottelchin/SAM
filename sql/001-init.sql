-- +goose up

CREATE TABLE IF NOT EXISTS Reisen (
    id INT PRIMARY KEY AUTOINCREMENT, 
    name VARCHAR, 
    archiviert bool
);

CREATE TABLE IF NOT EXISTS Nutzer (
    id INT PRIMARY KEY AUTOINCREMENT, 
    name VARCHAR, 
    mail VARCHAR, 
    password VARCHAR,
    reise INT --Hier wird eine Reise ID eingefügt, wenn der Nutzer nur exitiert um angezeigt zu werden, aber kein echter Nutzer ist 
);

CREATE TABLE IF NOT EXISTS Belege (
    id INT PRIMARY KEY AUTOINCREMENT,
    name VARCHAR,
    datum VARCHAR,
    betrag REAL,
    reise_id INT,
    von INT,
    FOREIGN KEY (reise_id) REFERENCES Reisen (id) ON DELETE CASCADE,
    FOREIGN KEY (von) REFERENCES Nutzer (id)
);

CREATE TABLE IF NOT EXISTS Nutzer_Reisen (
    nutzer_id INT,
    reisen_id INT,
    PRIMARY KEY (nutzer_id, reisen_id),
    FOREIGN KEY (nutzer_id) REFERENCES Nutzer (id),
    FOREIGN KEY (reisen_id) REFERENCES Reisen (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bezahlt_fuer (
    beleg_id INT,
    nutzer_id INT,
    PRIMARY KEY (beleg_id, nutzer_id),
    FOREIGN KEY (beleg_id) REFERENCES Belege (id),
    FOREIGN KEY (nutzer_id) REFERENCES Nutzer (id)
);