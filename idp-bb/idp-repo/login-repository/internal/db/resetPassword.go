package db

import (
	"context"
	"idp-repository/login-repository/internal/utils"
	pb "idp-repository/protos/login"
	"log"

	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c DBConfig) ResetPassword(ctx context.Context, req *pb.Request) (*pb.Status, error) {
	userId := req.UserId

	sqlStatement := `
		UPDATE users SET updatedAt = (select current_timestamp at time zone ('utc')), passwordHash = $1
		WHERE userId = $2
		RETURNING userId
	`

	hashedPassword, err := utils.HashPassword(req.Password)
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
