package db

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"idp-repository/pkg/logger"
	pb "idp-repository/protos/user"
	"log"
	"testing"
	"time"
)

func TestCheckAuthUserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}

	email := "test@gmail.com"

	//Test case #1 (check user not exist)
	log.Println("Test case #1: check user not exist")
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(0)
	mock.ExpectQuery(`SELECT CASE WHEN \(SELECT COUNT\(\*\) FROM users US WHERE US.email = \$1\) = 0 THEN 0 ELSE 1 END`).WithArgs(email).WillReturnRows(rows)

	exists, err := c.CheckAuthUserExists(&pb.Request{Email: email})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := false
	if exists != expected {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, exists)
	}
	log.Println("test case #1 is passed")

	//Test case #2 (check user exist)
	log.Println("Test case #2: check user exist")
	rows = sqlmock.NewRows([]string{"exists"}).AddRow(1)
	mock.ExpectQuery(`SELECT CASE WHEN \(SELECT COUNT\(\*\) FROM users US WHERE US.email = \$1\) = 0 THEN 0 ELSE 1 END`).WithArgs(email).WillReturnRows(rows)

	exists, err = c.CheckAuthUserExists(&pb.Request{Email: email})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected = true
	if exists != expected {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, exists)
	}
	log.Println("test case #2 is passed")

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

}

func TestCheckIdpAccountLinked(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}

	providerId := "google"
	federatedId := "123"

	//Test case #1 (check user not exist)
	log.Println("Test case #1: check idp account not linked")
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(0)
	expectedQuery := `SELECT CASE WHEN \(SELECT COUNT\(\*\) FROM linkedAccounts l WHERE l.providerId = \$1 AND l.federatedId = \$2\) = 0 THEN 0 ELSE 1 END`
	mock.ExpectQuery(expectedQuery).WithArgs(providerId, federatedId).WillReturnRows(rows)

	exists, err := c.CheckIdpAccountLinked(&pb.Request{ProviderId: providerId, FederatedId: federatedId})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := false
	if exists != expected {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, exists)
	}
	log.Println("test case #1 is passed")

	//Test case #2 (check user exist)
	log.Println("Test case #2: check idp account already linked")
	rows = sqlmock.NewRows([]string{"exists"}).AddRow(1)
	mock.ExpectQuery(expectedQuery).WithArgs(providerId, federatedId).WillReturnRows(rows)

	exists, err = c.CheckIdpAccountLinked(&pb.Request{ProviderId: providerId, FederatedId: federatedId})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected = true
	if exists != expected {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, exists)
	}
	log.Println("test case #2 is passed")

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

}

func TestCheckIdpProviderLinked(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}

	userId := "google"
	providerId := "google"

	//Test case #1 (check user not exist)
	log.Println("Test case #1: check idp provider not linked")
	rows := sqlmock.NewRows([]string{"exists"}).AddRow(0)
	expectedQuery := `SELECT CASE WHEN \(SELECT COUNT\(\*\) FROM linkedAccounts l WHERE l.userId = \$1 AND l.providerId = \$2\) = 0 THEN 0 ELSE 1 END`
	mock.ExpectQuery(expectedQuery).WithArgs(userId, providerId).WillReturnRows(rows)

	exists, err := c.CheckIdpProviderLinked(&pb.Request{UserId: userId, ProviderId: providerId})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected := false
	if exists != expected {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, exists)
	}
	log.Println("test case #1 is passed")

	//Test case #2 (check user exist)
	log.Println("Test case #2: check idp provider already linked")
	rows = sqlmock.NewRows([]string{"exists"}).AddRow(1)
	mock.ExpectQuery(expectedQuery).WithArgs(userId, providerId).WillReturnRows(rows)

	exists, err = c.CheckIdpProviderLinked(&pb.Request{UserId: userId, ProviderId: providerId})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expected = true
	if exists != expected {
		t.Errorf("Unexpected result. Expected: %v, Got: %v", expected, exists)
	}
	log.Println("test case #2 is passed")

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}

}

func TestGetAccountInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}


	email := "test@gmail.com"
	req := &pb.Request{
		Email: email,
	}

	//Test case #1 (get account info by email)
	log.Println("Test case #1: get account info by email")

	expectedQuery := `select u.userId, u.createdAt, u.updatedAt, u.emailVerified, u.passwordHash, u.OTP, i.displayName, i.firstName, i.lastName, i.photoUrl, i.gender, i.address, i.age, i.experience, i.playingTime, i.preferredPlatforms from users u JOIN userInfo i ON u.userId = i.userId where u.email= \$1`
	expectedRows := sqlmock.NewRows([]string{"userId", "createdAt", "updatedAt", "emailVerified", "passwordHash", "OTP", "displayName", "firstName", "lastName", "photoUrl", "gender", "address", "age", "experience", "playingTime", "preferredPlatforms"}).
		AddRow("TestUserId", time.Now(), time.Now(), true, "passwordHash", 0, "TestDisplayName", "TestFirstName", "TestLastName", "photo-url", "male", "21 Street", 25, "3 years", 4, "PC, Console")
	mock.ExpectQuery(expectedQuery).WithArgs(email).WillReturnRows(expectedRows)

	res, err := c.GetAccountInfo(context.Background(), req)
	if err != nil {
		t.Fatalf("Function under test returned an error: %v", err)
	}

	log.Println("account info: ", res)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetAccountInfoById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}


	userId := "123"
	req := &pb.Request{
		UserId: userId,
	}

	//Test case #1 (get account info by userId)
	log.Println("Test case #1: get account info by userId")

	expectedQuery := `select u.userId, u.createdAt, u.updatedAt, u.emailVerified, u.OTP, i.displayName, i.firstName, i.lastName, i.photoUrl, i.gender, i.address, i.age, i.experience, i.playingTime, i.preferredPlatforms, u.email from users u JOIN userInfo i ON u.userId = i.userId where u.userId= \$1`
	expectedRows := sqlmock.NewRows([]string{"userId", "createdAt", "updatedAt", "emailVerified", "OTP", "displayName", "firstName", "lastName", "photoUrl", "gender", "address", "age", "experience", "playingTime", "preferredPlatforms", "email"}).
		AddRow("TestUserId", time.Now(), time.Now(), true, 0, "TestDisplayName", "TestFirstName", "TestLastName", "photo-url", "male", "21 Street", 25, "3 years", 4, "PC, Console", "test@gmail.com")
	mock.ExpectQuery(expectedQuery).WithArgs(userId).WillReturnRows(expectedRows)

	res, err := c.GetAccountInfoById(context.Background(), req)
	if err != nil {
		t.Fatalf("Function under test returned an error: %v", err)
	}

	log.Println("account info: ", res)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetLinkedAccountInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database connection: %v", err)
	}
	defer db.Close()

	customlogger := logger.NewCustomLogger()
	c := DBConfig{roDb: db, rwDb: db, log: customlogger}


	userId := "123"
	providerId := "google"
	req := &pb.Request{
		UserId: userId,
		ProviderId: providerId,
	}

	//Test case #1 (get account info by userId)
	log.Println("Test case #1: get account info by userId")

	expectedQuery := `select l.federatedId, l.email, l.linkedUserId from linkedAccounts l where l.userId = \$1 AND l.providerId = \$2`
	expectedRows := sqlmock.NewRows([]string{"federatedId", "email", "linkedUserId"}).
		AddRow("456", "test@gmail.com", "123")
	mock.ExpectQuery(expectedQuery).WithArgs(userId, providerId).WillReturnRows(expectedRows)

	res, err := c.GetLinkedAccountInfo(context.Background(), req)
	if err != nil {
		t.Fatalf("Function under test returned an error: %v", err)
	}

	log.Println("linked account info: ", res)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

//this one keep throwing "panic: Query: could not match actual sql"
//func TestCreateNewUser(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("Failed to create mock database connection: %v", err)
//	}
//	defer db.Close()
//
//	customlogger := logger.NewCustomLogger()
//	c := DBConfig{roDb: db, rwDb: db, log: customlogger}
//
//	expectedQuery := `WITH inUser AS \(INSERT INTO users\(userId, email, emailVerified, createdAt, updatedAt, passwordHash\) VALUES \(gen_random_uuid\(\), \$1, \$2, \(select current_timestamp at time zone \('utc'\)\), \(select current_timestamp at time zone \('utc'\)\), \$3\) RETURNING userId, createdAt, updatedAt\) INSERT INTO userinfo \(userId, displayName, firstName, lastName, photoUrl, gender, address, age, experience, playingTime, preferredPlatforms\) SELECT IU.userId, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12, \$13 FROM inUser IU RETURNING \(SELECT inUser.userId FROM inUser\), \(SELECT inUser.createdAt FROM inUser\), \(SELECT inUser.updatedAt FROM inUser\)`
//
//	rows := sqlmock.NewRows([]string{"userId", "createdAt", "updatedAt"}).
//		AddRow("someUserId", time.Now(), time.Now())
//	//mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
//	mock.ExpectQuery(expectedQuery).WithArgs("test@gmail.com", true, "password", "TestDisplayName", "TestFirstName", "TestLastName", "photo-url", "Male", "21 Street", 25, "3 years", 4, "PC, Console").WillReturnRows(rows)
//
//	err = mock.ExpectationsWereMet()
//	if err != nil {
//		//this throw: "Unfulfilled expectations: there is a remaining expectation which was not matched: ExpectedQuery => expecting Query, QueryContext or QueryRow"
//		t.Errorf("Unfulfilled expectations: %v", err)
//	}
//
//	req := &pb.Request{
//		Email:              "test@gmail.com",
//		Password:           "password",
//		DisplayName:        "TestDisplayName",
//	}
//
//	res, err := c.CreateNewUser(req)
//	if err != nil {
//		t.Fatalf("Function under test returned an error: %v", err)
//	}
//
//	print(res)
//
//}


// old unit testing using actual db connection
//func TestUserDbFunctions(t *testing.T) {
//	os.Setenv("DB_PORT", "45432")
//	os.Setenv("DB_USER", "postgres")
//	os.Setenv("DB_PASS", "9221")
//
//	customlogger := logger.NewCustomLogger()
//	env , err := env.GetEnv()
//	if err != nil {
//		customlogger.Errorlog.Fatalf("Loading env: %v", err)
//	}
//	rwDB, err := postgres.Connect(env)
//	if err != nil {
//		customlogger.Errorlog.Fatal("RW Postgresql init: ", err)
//	}
//	defer rwDB.Close()
//	roDB, err := postgres.Connect(env)
//	if err != nil {
//		customlogger.Errorlog.Fatal("RO Postgresql init: ", err)
//	}
//	defer roDB.Close()
//
//	c := DBConfig{roDb: roDB, rwDb: roDB, log: customlogger}
//
//	log.Println("Testing user db functions ------------>")
//
//	TestEmail := "test_" + fmt.Sprint(time.Now().Nanosecond())[:6] + "@gmail.com"
//	//Test case #1 (check user not exist)
//	log.Println("Test case #1: check user not exist ------------>")
//	req := &pb.Request{
//		Email : TestEmail,
//	}
//
//	res, err := c.CheckAuthUserExists(req)
//
//	if err != nil {
//		t.Fatalf("error test case #1: %v", err)
//	} else if res {
//		t.Fatalf("test case #1 is failed")
//	} else {
//		log.Println("test case #1 is passed")
//	}
//
//	//Test case #2 (create the user)
//	log.Println("Test case #2: create the user ------------>")
//	req = &pb.Request{
//		Email : TestEmail,
//		DisplayName : "TestCase1",
//		Password: "pass",
//	}
//
//	res2, err := c.CreateNewUser(req)
//
//	if err != nil {
//		t.Fatalf("error test case #2: %v", err)
//	} else {
//		log.Println("test case #2 is passed")
//		log.Println(res2)
//	}
//
//	//Test case #3 (check user exist)
//	log.Println("Test case #3: check user exist ------------>")
//	req = &pb.Request{
//		Email : TestEmail,
//	}
//
//	res, err = c.CheckAuthUserExists(req)
//
//	if err != nil {
//		t.Fatalf("error test case #3: %v", err)
//	} else if res {
//		log.Println("test case #3 is passed")
//	} else {
//		t.Fatalf("test case #3 is failed")
//	}
//
//	//Test case #4 (get account info by email)
//	log.Println("Test case #4: get account info by email ------------>")
//	req = &pb.Request{
//		Email : TestEmail,
//	}
//
//	userInfo, err := c.GetAccountInfo(context.Background(), req)
//
//	if err != nil {
//		t.Fatalf("error test case #4: %v", err)
//	} else {
//		log.Println("test case #4 is passed")
//	}
//
//	//Test case #5 (set account info)
//	log.Println("Test case #5: set account info ------------>")
//	req = &pb.Request{
//		UserId: userInfo.Users.UserId,
//		DisplayName : "TestName",
//	}
//
//	_, err = c.SetAccountInfo(context.Background(), req)
//
//	if err != nil {
//		t.Fatalf("error test case #5: %v", err)
//	} else {
//		log.Println("test case #5 is passed")
//	}
//
//	//Test case #6 (get account info by userId)
//	log.Println("Test case #6: get account info by userId ------------>")
//	req = &pb.Request{
//		UserId: userInfo.Users.UserId,
//	}
//
//	userInfo, err = c.GetAccountInfoById(context.Background(), req)
//
//	if err != nil {
//		t.Fatalf("error test case #6: %v", err)
//	} else {
//		log.Println("test case #6 is passed")
//	}
//
//	log.Println("All test cases are finished.")
//
//}
