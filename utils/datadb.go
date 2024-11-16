package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type DBConn struct {
	Filename string
	DB       *sql.DB
}

type ClientInfo struct {
	IP    string
	Alive bool
}

const CreateClientTable = `
	CREATE TABLE "clients" (
		"ip" STRING,
		"alive" BOOLEAN,
		PRIMARY KEY("ip")
	);
`

const CreateFilesTable = `
	CREATE TABLE "files" (
		"id" INTEGER,
		"client_ip" STRING,
		"name" STRING,
		"size" INTEGER,
		"created_at" INTEGER,
		PRIMARY KEY("id"),
		FOREIGN KEY("client_ip") REFERENCES "clients"("ip")
	);
`

// This function creates a fresh db every time it is called
func CreateDBConn(filename string) (DBConn, error) {
	err := os.Remove(filename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return DBConn{}, err
		}
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", filename))
	if err != nil {
		db.Close()
		return DBConn{}, err
	}

	_, err = db.Exec(CreateClientTable)
	if err != nil {
		db.Close()
		return DBConn{}, err
	}

	_, err = db.Exec(CreateFilesTable)
	if err != nil {
		db.Close()
		return DBConn{}, err
	}

	return DBConn{Filename: filename, DB: db}, nil
}

func (db *DBConn) WriteClientInfo(clientInfo ClientInfo) error {
	var tmpClient ClientInfo
	rows, err := db.DB.Query("SELECT * FROM clients;")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&tmpClient.IP, &tmpClient.Alive)
		if err != nil {
			return err
		}

		if tmpClient.IP == clientInfo.IP {
			_, err = db.DB.Exec("UPDATE clients SET alive = 1 WHERE ip = ?", clientInfo.IP)
			if err != nil {
				return err
			}
			return nil
		}
	}

	_, err = db.DB.Exec("INSERT INTO clients(ip, alive) VALUES(?, ?);", clientInfo.IP, clientInfo.Alive)
	if err != nil {
		return err
	}
	return nil
}
