package server

import (
	"context"
	"net/http"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/steinfletcher/apitest"
)

var DefaultCfg = &ServerCfg{
	SecretKey: "f8e03b9275d79e802b593409e6073a1ff31be30c0dc72566870bb7d7d992e630", // Base64: +OA7knXXnoArWTQJ5gc6H/Mb4wwNxyVmhwu319mS5jA=
	DBUrl:     "root:secret7@(localhost:3306)/test?multiStatements=true&parseTime=true",
}

var srv *Server
var logger *zap.Logger

func TestMain(m *testing.M) {
	logger, _ = zap.NewDevelopment()
	srv = New(DefaultCfg, logger)
	srv.Init(context.Background())
	code := m.Run()
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
