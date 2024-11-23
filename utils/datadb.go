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

type FileInfo struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

const CreateClientTable = `
	CREATE TABLE "clients" (
		"ip" TEXT,
		"alive" BOOLEAN,
		PRIMARY KEY("ip")
	);
`

const CreateFilesTable = `
	CREATE TABLE "files" (
		"id" INTEGER,
		"client_ip" TEXT,
		"name" TEXT,
		"size" INTEGER,
		"created_at" DEFAULT CURRENT_TIMESTAMP,
		"status" TEXT CHECK("status" IN ('inprogress', 'done')),
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
	row := db.DB.QueryRow("SELECT alive FROM clients WHERE ip = ?;", clientInfo.IP)
	var alive bool
	err := row.Scan(&alive)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.DB.Exec("INSERT INTO clients(ip, alive) VALUES(?, ?);", clientInfo.IP, clientInfo.Alive)
		if err != nil {
			return err
		}
	}
	return err
}

func (db *DBConn) StartFileUpload(ClientIP string, filename string, size int64) error {
	_, err := db.DB.Exec(`INSERT INTO files(client_ip, name, size, status) VALUES(?, ?, ?, 'inprogress');`, ClientIP, filename, size)
	if err != nil {
		return err
	}
	return nil
}

func (db *DBConn) FinishFileUpload(ClientIP string, filename string) error {
	_, err := db.DB.Exec(`UPDATE files SET status = 'done' WHERE client_ip = ? AND name = ?;`, ClientIP, filename)
	if err != nil {
		return err
	}
	return nil
}
