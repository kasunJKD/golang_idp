package db

import (
	"context"
	"database/sql"
	"idp-repository/pkg/logger"
	pb "idp-repository/protos/login"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//note:: columns should be what is being scanned in the original queries

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestCheck2faEnabled(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}
    
    testRequest := &pb.OtpRequest{
        UserId: "321432",
    }

    log.Println("Test case #1: 2fa enabled")
    rows := sqlmock.NewRows([]string{"OTP"}).AddRow(1)
    mock.ExpectQuery(`SELECT OTP FROM users WHERE userId = \$1`).WithArgs(testRequest.UserId).WillReturnRows(rows)
  
    checker, err := c.Check2faEnabled(testRequest)

    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

	var expected int32 = 1
    if checker != expected {
		t.Errorf("Unexpected result. Expected: %v, Got %v", expected, checker)
	}
	log.Println("Test case #1: passed")

	log.Println("Test case #2: 2fa not enabled")
	query := "SELECT OTP FROM users WHERE userId = \\$1"
    rows = sqlmock.NewRows([]string{"OTP"}).AddRow(0)
    mock.ExpectQuery(query).WithArgs(testRequest.UserId).WillReturnRows(rows)
  
    checker, err = c.Check2faEnabled(testRequest)

    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }

	expected = 0
    if checker != expected {
		t.Errorf("Unexpected result. Expected: %v, Got %v", expected, checker)
	}

	log.Println("Test case #2: passed")

}

func TestCheckClientidExists (t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}

	request := &pb.ClientReq{
		ClientId: "test321432",
	}

	log.Println("Test case #1: Client id exists")
	query := `SELECT CASE WHEN \(SELECT COUNT\(\*\) FROM clients C WHERE C.clientId = \$1\) = 0 THEN 0 ELSE 1 END`
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(0)
	mock.ExpectQuery(query).WithArgs(request.ClientId).WillReturnRows(rows)

	b, err := c.CheckClientidExists(request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := false 
	if b != expected {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, b)
	}
	log.Println("test case #1 passed")

	log.Println("Test case #2: Client id not exists")
	query = `SELECT CASE WHEN \(SELECT COUNT\(\*\) FROM clients C WHERE C.clientId = \$1\) = 0 THEN 0 ELSE 1 END`
	rows = sqlmock.NewRows([]string{"exists"}).AddRow(1)
	mock.ExpectQuery(query).WithArgs(request.ClientId).WillReturnRows(rows)

	b, err = c.CheckClientidExists(request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected = true 
	if b != expected {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, b)
	}
	log.Println("test case #2 passed")


}

func TestCheckClientNameExists (t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}

	request := &pb.ClientReq{
		ClientName: "testname",
	}

	log.Println("Test case #1: Client name exists")
	query := `SELECT CASE WHEN \(SELECT COUNT\(\*\) FROM clients C WHERE C.clientName = \$1\) = 0 THEN 0 ELSE 1 END`
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(0)
	mock.ExpectQuery(query).WithArgs(request.ClientName).WillReturnRows(rows)

	b, err := c.CheckClientNameExists(request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := false 
	if b != expected {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, b)
	}
	log.Println("test case #1 passed")

	log.Println("Test case #2: Client name not exists")
	query = `SELECT CASE WHEN \(SELECT COUNT\(\*\) FROM clients C WHERE C.clientName = \$1\) = 0 THEN 0 ELSE 1 END`
	rows = sqlmock.NewRows([]string{"exists"}).AddRow(1)
	mock.ExpectQuery(query).WithArgs(request.ClientName).WillReturnRows(rows)

	b, err = c.CheckClientNameExists(request)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected = true 
	if b != expected {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, b)
	}
	log.Println("test case #2 passed")


}

func TestGetClientInfo (t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}

	request := &pb.ClientReq{
		ClientId: "testId321",
	}

	tt, _ := time.Parse(time.UnixDate,"2020-04-09 11:24:14.785868848 +0000 UTC m=+0.000187421")

	log.Println("Test case #1: Get Client Info")
	query := `select distinct clientId, clientName, clientSecret, projectId, userId, redirectUrl, createdAt, updatedAt, active
				from clients 
				where clientId= \$1`
	rows := sqlmock.NewRows([]string{"clientId", "clientName", "clientSecret",
						"projectId", "userId", "redirectUrl", "createdAt", "updatedAt", "active"}).
					AddRow("testId321", "testname", "testSecret", "testprojectId", "testUser", "www.redirectUrltest",
								tt, tt, true)
	mock.ExpectQuery(query).WithArgs(request.ClientId).WillReturnRows(rows)

	res, err := c.GetClientInfo(context.Background(), request)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := &pb.ClientRes{
		ClientId: "testId321",
		ClientName: "testname",
		ClientSecret: "testSecret",
		ProjectId: "testprojectId",
		UserId: "testUser",
		RedirectUrl: "www.redirectUrltest",
		CreatedAt:     timestamppb.New(tt),
		UpdatedAt:     timestamppb.New(tt),
		Active: true,
	}

	if reflect.DeepEqual(expected, res) == false {
		t.Fatalf("Unexpected result. Expected: %v, Got: %v", expected, res)
	}
	log.Println("test case #1 passed")	

}

func TestGetClientsByUserId (t *testing.T){
	db, mock := NewMock()
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}

	request := &pb.ClientReq{
		UserId: "user123",
	}

	tt, _ := time.Parse(time.UnixDate,"2020-04-09 11:24:14.785868848 +0000 UTC m=+0.000187421")

	log.Println("Test case #1: Get Clients by userid")
	query := `select C.clientId, C.clientName, C.clientSecret, C.redirectUrl, C.createdAt, C.updatedAt, C.active
				from Clients C
				where C.userId = \$1`
	rows := sqlmock.NewRows([]string{"clientId", "clientName", "clientSecret",
						"redirectUrl", "createdAt", "updatedAt", "active"}).
					AddRow("testId321", "testname", "testSecret", "www.redirectUrltest",
								tt, tt, true)
	mock.ExpectQuery(query).WithArgs(request.UserId).WillReturnRows(rows)

	res, err := c.GetClientsByUserId(context.Background(), request)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}


	cc := &pb.ClientRes {
			ClientId: "testId321",
			ClientName: "testname",
			ClientSecret: "testSecret",
			RedirectUrl: "www.redirectUrltest",
			CreatedAt: timestamppb.New(tt),
			UpdatedAt: timestamppb.New(tt),
			Active: true,
		}
	clients := make([]*pb.ClientRes, 0)
	clients = append(clients, cc)
	expected := &pb.GetClientsByUserIdResponse{
		Clients: clients,
	}

	if reflect.DeepEqual(expected, res) == false {
		t.Fatalf("Unexpected result. Expected: %v, Got: %v", expected, res)
	}
	log.Println("test case #1 passed")	

	
}



