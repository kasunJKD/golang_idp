package db

import (
	"fmt"
	_ "github.com/lib/pq"
	pb "idp-repository/protos/user"
)

func (c DBConfig) CheckIdpAccountLinked(req *pb.Request) (bool, error) {
	var isAuthenticated bool
	c.log.Infolog.Println("checking if account linked")
	sqlStatement := `SELECT CASE WHEN
	 (SELECT COUNT(*) FROM linkedAccounts l WHERE l.providerId = $1 AND l.federatedId = $2) = 0 THEN 0 ELSE 1 END`
	err := c.roDb.QueryRow(sqlStatement, req.ProviderId, req.FederatedId).Scan(&isAuthenticated)

	if err != nil {
		c.log.Errorlog.Fatal(err)
	}
	c.log.Infolog.Println("Account linked = " + fmt.Sprint(isAuthenticated))
	return isAuthenticated, err
}