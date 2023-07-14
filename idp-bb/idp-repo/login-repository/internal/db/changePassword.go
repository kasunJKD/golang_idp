package db

import (
	"idp-repository/login-repository/internal/utils"
	pb "idp-repository/protos/login"
	"log"

	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c DBConfig) ChangePassword(req *pb.Passwordreq) (*pb.Status, error) {
	userId := req.GetUserId()
	newPassword := req.GetNewPassword()

	//check name exists
	b, _ := c.checkPassword(req)
	if b  == false{
		return nil, status.Errorf(codes.Unauthenticated, "old password incorrect")
	}

	sqlStatement := `
		UPDATE users SET updatedAt = (select current_timestamp at time zone ('utc')), passwordHash = $1
		WHERE userId = $2
		RETURNING userId
	`

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password cannot be hashed")
	}

	err = c.rwDb.QueryRow(sqlStatement,
		hashedPassword,
		userId,
	).Scan(&userId)

	if err != nil {
		log.Fatalln(err)
		return nil, status.Errorf(codes.Internal, "database update error")
	}

	res := &pb.Status{
		Code: 0,
		Message: "Password updated",
	}

	return res, nil


}

func (c DBConfig) checkPassword (req *pb.Passwordreq) (bool, error) {
	var storedHash string
	sqlStatement := `SELECT passwordHash FROM users C WHERE C.userId = $1`
	err := c.roDb.QueryRow(sqlStatement, req.UserId).Scan(&storedHash)

	if err != nil {
		log.Fatal(err)
	}

	//check password is correct (oldPassword)
	err = utils.CheckPassword(req.OldPassword, storedHash)

	if err != nil {
		return false, err 
	}
	
	return true, err
}