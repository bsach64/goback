-- STORES CLIENT info
CREATE TABLE "clients" (
	"id" INTEGER,
	"ip" STRING,
	"username" STRING NOT NULL UNIQUE,
	"alive" BOOLEAN,
	PRIMARY KEY("id")
);

-- Keep track of files for each of the clients
CREATE TABLE "files" (
	"id" INTEGER,
	"client_id" INTEGER,
	"name" STRING,
	"size" INTEGER,
	"created_at" INTEGER,
	PRIMARY KEY("id"),
	FOREIGN KEY("client_id") REFERENCES "clients"("id")
);
