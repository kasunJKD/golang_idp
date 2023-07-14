package db

import (
	"fmt"
	_ "github.com/lib/pq"
	pb "idp-repository/protos/user"
)

func (c DBConfig) CheckIdpProviderLinked(req *pb.Request) (bool, error) {
	var isAuthenticated bool
	c.log.Infolog.Println("checking if provider linked")
	sqlStatement := `SELECT CASE WHEN
	 (SELECT COUNT(*) FROM linkedAccounts l WHERE l.userId = $1 AND l.providerId = $2) = 0 THEN 0 ELSE 1 END`
	err := c.roDb.QueryRow(sqlStatement, req.UserId, req.ProviderId).Scan(&isAuthenticated)

	if err != nil {
		c.log.Errorlog.Fatal(err)
	}
	c.log.Infolog.Println("Provider linked = " + fmt.Sprint(isAuthenticated))
	return isAuthenticated, err
}