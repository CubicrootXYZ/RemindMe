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
	db, err := sql.Open("sqlite3", "testdb.db")
	if err != nil {
		log.Warn(err.Error())
		panic(err)
	}
	testDb = db

	m.Run()

	cleanUp()
}

func cleanUp() {
	os.Remove("testdb.db")
}

func TestEncryption_GetCryptoStoreOnSuccess(t *testing.T) {
	store, err := GetCryptoStore(testDb, &configuration.Matrix{
		DeviceID: "1234",
		Username: "admin",
	})

	assert.NoError(t, err)
	assert.NotNil(t, store)
}

func TestEncryption_GetOlmMachineOnSuccess(t *testing.T) {
	olm := GetOlmMachine(getTestClient(), getTestStore())

	assert.NotNil(t, olm)
}

func getTestClient() *mautrix.Client {
	client, err := mautrix.NewClient("https://mydomain.tld", "", "")
	if err != nil {
		panic(err)
	}

	return client
}

func getTestStore() crypto.Store {
	store, err := GetCryptoStore(testDb, &configuration.Matrix{
		DeviceID: "1234",
		Username: "admin",
	})

	if err != nil {
		panic(err)
	}

	return store
}
