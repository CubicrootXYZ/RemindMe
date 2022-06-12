package encryption

import (
	"database/sql"
	"os"
	"testing"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/stretchr/testify/assert"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
)

var testDb *sql.DB

func TestMain(m *testing.M) {
	cleanUp()
	db, err := sql.Open("sqlite3", "data/olm.db")
	if err != nil {
		log.Warn(err.Error())
		panic(err)
	}
	testDb = db

	exitCode := m.Run()

	cleanUp()

	os.Exit(exitCode)
}

func cleanUp() {
	os.Remove("data/olm.db")
}

func TestEncryption_GetCryptoStoreOnSuccess(t *testing.T) {
	store, _, err := GetCryptoStore(false, testDb, &configuration.Matrix{
		DeviceID: "1234",
		Username: "admin",
	})

	assert.NoError(t, err)
	assert.NotNil(t, store)
}

func TestEncryption_GetOlmMachineOnSuccess(t *testing.T) {
	mach := GetOlmMachine(false, getTestClient(), getTestStore(), nil, nil)

	assert.NotNil(t, mach)
}

func getTestClient() *mautrix.Client {
	client, err := mautrix.NewClient("https://mydomain.tld", "", "")
	if err != nil {
		panic(err)
	}

	return client
}

func getTestStore() crypto.Store {
	store, _, err := GetCryptoStore(false, testDb, &configuration.Matrix{
		DeviceID: "1234",
		Username: "admin",
	})

	if err != nil {
		panic(err)
	}

	return store
}
