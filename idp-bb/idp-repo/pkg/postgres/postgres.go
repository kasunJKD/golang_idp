package postgres

import (
	"database/sql"
	"fmt"
	"idp-repository/env"
	"log"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	DB *sql.DB
}

//Connect database
func Connect(env *env.Env) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		env.Postgres.Host,
		env.Postgres.Port,
		env.Postgres.User,
		env.Postgres.Password,
		env.Postgres.Dbname)
	log.Println(psqlInfo)
	dbConn, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("Error Connecting to database", err)
		return nil, err
	}
	//defer dbConn.Close()
	err = dbConn.Ping()
	if err != nil {
		log.Fatal("Error Pinging to database", err)
		return nil, err
	}

	log.Println("Successfully Connected to database")

	return dbConn, nil
}

//Initialize database ::TODO
//db.InitDB(dbConn)
//db function definitions::TODO