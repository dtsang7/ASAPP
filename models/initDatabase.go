package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rubenv/sql-migrate"
	"log"
)

type DAO struct {
	db         *sql.DB
	driverName string
}

func CreateDAO(driverName string, dataSource string) *DAO {
	//connect to database
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		log.Fatal("Unable to open DB", err.Error())
	}
	//database to struct
	return &DAO{db, driverName}
}

func (dao *DAO) RunMigrations() {

	migrations := migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	n, err := migrate.Exec(dao.db, dao.driverName, migrations, migrate.Up)
	if err != nil {
		log.Fatal("Unable to migrate", err.Error())
	}
	log.Printf("Applied %d migrations.\n", n)
}
