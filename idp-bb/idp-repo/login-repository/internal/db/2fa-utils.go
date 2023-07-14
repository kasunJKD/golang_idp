package db

import (
	"log"

	_ "github.com/lib/pq"
	pb "idp-repository/protos/login"
)

func (c DBConfig) Check2faEnabled(req *pb.OtpRequest) (int32, error) {
	var checker int32
	sqlStatement := `SELECT OTP FROM users WHERE userId = $1`
	err := c.roDb.QueryRow(sqlStatement, req.UserId).Scan(&checker)

	if err != nil {
		log.Fatal(err)
	}

	return checker, err
}

func (c DBConfig) Update2faStatus (req *pb.OtpRequest) (int32, error) {
	
	sqlStatement := `UPDATE users SET OTP = $1, updatedAt = (select current_timestamp at time zone ('utc'))
		WHERE userId = $2
		RETURNING 2fa`

	var (
		status int32
	)
	
	err := c.rwDb.QueryRow(sqlStatement,
		req.GetOtpEnabled(),
		req.GetUserId(),
	).Scan(&status)

	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	return status, err
}
