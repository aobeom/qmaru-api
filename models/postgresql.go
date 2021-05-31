package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"reflect"

	"qmaru-api/configs"

	_ "github.com/lib/pq"
)

type PostgreSQL struct{}

var Psql *PostgreSQL
var pdb *sql.DB

func init() {
	Psql := new(PostgreSQL)
	pdb = Psql.Connect()
}

func (p *PostgreSQL) Connect() *sql.DB {
	cfg := configs.DBCfg()

	host := cfg["host"].(string)
	port := int(cfg["port"].(float64))
	user := cfg["user"].(string)
	password := cfg["password"].(string)
	dbname := cfg["dbname"].(string)

	info := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbname)
	db, err := sql.Open("postgres", info)
	if err != nil {
		log.Panic("DB Connect Failed", err)
	}
	return db
}

func (p *PostgreSQL) Exec(sql string, args ...interface{}) {
	stmt, err := pdb.Prepare(sql)
	if err != nil {
		log.Panic(err)
	}
	defer stmt.Close()
	stmt.Exec(args...)
}

func (p *PostgreSQL) Query(sql string, args ...interface{}) (rows *sql.Rows, err error) {
	stmt, err := pdb.Prepare(sql)
	if err != nil {
		log.Panic(err)
	}

	defer stmt.Close()
	rows, err = stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (p *PostgreSQL) QueryOne(sql string, args ...interface{}) *sql.Row {
	stmt, err := pdb.Prepare(sql)
	if err != nil {
		log.Panic(err)
	}

	defer stmt.Close()
	row := stmt.QueryRow(args...)
	return row
}

func InitTable() {
	tables := []interface{}{
		CrondTime{},
		DramaInfo{},
		MediaInfo{},
		RadikoInfo{},
		StInfo{},
		StToken{},
	}
	defer pdb.Close()
	for _, table := range tables {
		var buffer bytes.Buffer
		rType := reflect.TypeOf(table)
		rName := DBName(rType.Name())
		DBFiled(rType, &buffer)
		rFiled := buffer.Bytes()[0 : len(buffer.Bytes())-1]

		sql := fmt.Sprintf("CREATE TABLE %s (%s)", rName, rFiled)
		_, err := pdb.Exec(sql)
		if err != nil {
			log.Println("Create Error:", err)
		}
	}
}

func DBPing() {
	err := pdb.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
