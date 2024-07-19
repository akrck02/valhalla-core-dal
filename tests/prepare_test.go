package tests

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/akrck02/valhalla-core-dal/configuration"
	"github.com/akrck02/valhalla-core-dal/database"
	"github.com/akrck02/valhalla-core-sdk/log"
)

var pathSetupDone bool = false

func TestMain(m *testing.M) {
	setupPathIfNecessary()
	satupDatabase()
	os.Exit(m.Run())
}

func setupPathIfNecessary() {
	if pathSetupDone {
		return
	}

	var _, current_execution_dir, _, _ = runtime.Caller(0)
	var BASE_PATH = current_execution_dir
	var _ = configuration.SetBasePath(BASE_PATH)

	// substract the last 2 directories
	BASE_PATH = BASE_PATH[:strings.LastIndex(BASE_PATH, "/")]
	BASE_PATH = BASE_PATH[:strings.LastIndex(BASE_PATH, "/")] + "/"

	configuration.SetBasePath(BASE_PATH)
	configuration.LoadConfiguration(BASE_PATH + ".env")

}

func satupDatabase() {

	log.Jump()
	log.Info("Setting up test environment...")
	database.SetupTest()
	pathSetupDone = true
	log.Jump()

}
