BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "Users" (
	"Username"	TEXT,
	PRIMARY KEY("Username")
);
CREATE TABLE IF NOT EXISTS "Wishes" (
	"id"	INTEGER,
	"from"	TEXT NOT NULL,
	"to"	TEXT NOT NULL,
	"wish"	TEXT NOT NULL,
	"timestamp"	TEXT NOT NULL DEFAULT current_timestamp,
	PRIMARY KEY("id" AUTOINCREMENT),
	FOREIGN KEY("from") REFERENCES "Users"("Username"),
	FOREIGN KEY("to") REFERENCES "Users"("Username")
);
COMMIT;