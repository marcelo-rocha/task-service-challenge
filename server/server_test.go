package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/golang-jwt/jwt/v4"
	"github.com/marcelo-rocha/task-service-challenge/domain/entities"
	"github.com/marcelo-rocha/task-service-challenge/persistence"
	"github.com/steinfletcher/apitest"
)

var DefaultCfg = &ServerCfg{
	SecretKey: "f8e03b9275d79e802b593409e6073a1ff31be30c0dc72566870bb7d7d992e630", // Base64: +OA7knXXnoArWTQJ5gc6H/Mb4wwNxyVmhwu319mS5jA=
	DBUrl:     "root:secret7@(localhost:3306)/test?multiStatements=true&parseTime=true",
	NATSUrl:   "localhost:4222",
}

var srv *Server
var logger *zap.Logger
var tasks *persistence.Tasks

func TestMain(m *testing.M) {
	logger, _ = zap.NewDevelopment()
	srv = New(DefaultCfg, logger)
	ctx := context.Background()
	srv.Init(ctx)
	users := persistence.NewUsers(srv.dbConnection, srv.logger)
	tasks = persistence.NewTasks(srv.dbConnection, srv.logger)
	insertTestUsers(ctx, users)
	code := m.Run()
	users.RestoreInitialSetup(ctx)

	os.Exit(code)
}

func TestGetRoot(t *testing.T) {
	apitest.New().
		Handler(srv.Router).
		Get("/").
		Expect(t).
		Status(http.StatusOK).
		End()
}

const (
	DemoUserId     = 2
	OperatorUserId = 3
)

func insertTestUsers(ctx context.Context, users *persistence.Users) error {
	id, err := users.InsertUser(ctx, "demo", "demonstration", entities.Technician, true, &persistence.DefaultAdminUserId)
	if err != nil {
		return err
	}
	if id != DemoUserId {
		return errors.New("unexpected id")
	}

	id, err = users.InsertUser(ctx, "operator", "assistent operator", entities.Technician, true, &persistence.DefaultAdminUserId)
	if err != nil {
		return err
	}
	if id != OperatorUserId {
		return errors.New("unexpected id")
	}
	return nil
}

func getAdminToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         strconv.Itoa(int(persistence.DefaultAdminUserId)),
		"iat":         time.Now().UTC().Unix(),
		UserKindClaim: string(entities.Manager),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(AuthencationSecretKey)
	return tokenString, err
}

func getDemoToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         strconv.Itoa(int(DemoUserId)),
		"iat":         time.Now().UTC().Unix(),
		UserKindClaim: string(entities.Technician),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(AuthencationSecretKey)
	return tokenString, err
}

func getOperatorToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":         strconv.Itoa(int(OperatorUserId)),
		"iat":         time.Now().UTC().Unix(),
		UserKindClaim: string(entities.Technician),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(AuthencationSecretKey)
	return tokenString, err
}
