package db

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	pb "idp-repository/protos/login"
)

func (c DBConfig) CheckClientidExists(req *pb.ClientReq) (bool, error) {
	var isAuthenticated bool
	log.Println("checking if clientId exists")
	sqlStatement := `SELECT CASE WHEN
	 (SELECT COUNT(*) FROM clients C WHERE C.clientId = $1) = 0 THEN 0 ELSE 1 END`
	err := c.roDb.QueryRow(sqlStatement, req.ClientId).Scan(&isAuthenticated)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("User exists = " + fmt.Sprint(isAuthenticated))
	return isAuthenticated, err
}

func (c DBConfig) CheckClientNameExists(req *pb.ClientReq) (bool, error) {
	var isAuthenticated bool
	log.Println("checking if clientId exists")
	sqlStatement := `SELECT CASE WHEN
	 (SELECT COUNT(*) FROM clients C WHERE C.clientName = $1) = 0 THEN 0 ELSE 1 END`
	err := c.roDb.QueryRow(sqlStatement, req.ClientName).Scan(&isAuthenticated)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client Name exists = " + fmt.Sprint(isAuthenticated))
	return isAuthenticated, err
}
