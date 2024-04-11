package api

import (
	"os"
	"testing"

	"github.com/aseerkt/go-simple-bank/pkg/db"
	"github.com/aseerkt/go-simple-bank/pkg/utils"
	"github.com/gin-gonic/gin"
)

func newTestServer(store db.Store) *Server {
	config := utils.LoadConfig("../..")
	gin.SetMode(gin.TestMode)
	return NewServer(store, &config)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
