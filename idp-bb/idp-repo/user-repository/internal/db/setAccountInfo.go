package db

import (
	"context"
	pb "idp-repository/protos/user"

	_ "github.com/lib/pq"
)

func (c DBConfig) SetAccountInfo(ctx context.Context, req *pb.Request) (*pb.Response, error) {

	sqlStatement := `
	WITH inUserInfo AS (
		UPDATE userinfo SET displayName = coalesce(NULLIF($1, ''), displayName), firstName = coalesce(NULLIF($2, ''), firstName), lastName = coalesce(NULLIF($3, ''), lastName), photoUrl = coalesce(NULLIF($4, ''), photoUrl)
		WHERE userId = $5
		RETURNING userId
	)
	UPDATE users SET updatedAt = (select current_timestamp at time zone ('utc'))
	WHERE userId = $5
	RETURNING (SELECT userId FROM inUserInfo)`

	var (
		userId string
	)

	err := c.rwDb.QueryRow(sqlStatement,
		req.GetDisplayName(),
		req.GetFirstName(),
		req.GetLastName(),
		req.GetPhotoUrl(),
		req.GetUserId(),
	).Scan(&userId)

	if err != nil {
		c.log.Errorlog.Fatalln(err)
	}

	c.log.Infolog.Println("User's Info updated: userId = " + userId)

	return nil, err

}
