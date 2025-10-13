BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "Users" (
	"Username"	TEXT,
	PRIMARY KEY("Username")
);
INSERT INTO "Users" ("Username") VALUES ('Alle');

CREATE TABLE IF NOT EXISTS Wishes (
	"id"	INTEGER,
	"from"	TEXT NOT NULL,
	"to"	TEXT NOT NULL,
	"wish"	TEXT NOT NULL,
	"timestamp"	datetime,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("from") REFERENCES "Users"("Username"),
	FOREIGN KEY("to") REFERENCES "Users"("Username")
);

-- Trigger, der bei jedem INSERT die lokale Zeit setzt
CREATE TRIGGER IF NOT EXISTS set_local_timestamp
AFTER INSERT ON Wishes
FOR EACH ROW
WHEN NEW.timestamp IS NULL
BEGIN
    UPDATE Wishes
    SET timestamp = datetime('now', 'localtime')
    WHERE id = NEW.id;
END;
COMMIT;