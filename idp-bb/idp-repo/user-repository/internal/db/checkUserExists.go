package db

import (
	"fmt"
	_ "github.com/lib/pq"
	pb "idp-repository/protos/user"
)

func (c DBConfig) CheckAuthUserExists(req *pb.Request) (bool, error) {
	var isAuthenticated bool
	c.log.Infolog.Println("checking if user exists")
	sqlStatement := `SELECT CASE WHEN
	 (SELECT COUNT(*) FROM users US WHERE US.email = $1) = 0 THEN 0 ELSE 1 END`
	err := c.roDb.QueryRow(sqlStatement, req.Email).Scan(&isAuthenticated)

	if err != nil {
		c.log.Errorlog.Fatal(err)
	}
	c.log.Infolog.Println("User exists = " + fmt.Sprint(isAuthenticated))
	return isAuthenticated, err
}
