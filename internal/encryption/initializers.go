package encryption

import (
	"database/sql"
	"fmt"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/id"

	_ "github.com/mattn/go-sqlite3"
)

// GetCryptoStore initializes a sql crypto store
func GetCryptoStore(db *sql.DB, config *configuration.Matrix) (crypto.Store, error) {
	account := fmt.Sprintf("%s/%s", config.Username, config.DeviceID)

	// TODO get propper sql support to work?
	db, err := sql.Open("sqlite3", "olm.db")
	if err != nil {
		log.Warn(err.Error())
		panic(err)
	}

	cryptoStore := crypto.NewSQLCryptoStore(db, "sqlite3", account, id.DeviceID(config.DeviceID), []byte(config.DeviceKey), cryptoLogger{"Crypto"})

	err = cryptoStore.CreateTables()

	return cryptoStore, err
}

// GetOlmMachine initializes a new olm machine
func GetOlmMachine(client *mautrix.Client, store crypto.Store, database types.Database) *crypto.OlmMachine {
	if client == nil {
		log.Warn("client nil")
		panic("client nil")
	}
	if store == nil {
		log.Warn("store nil")
		panic("store nil")
	}

	machine := crypto.NewOlmMachine(client, cryptoLogger{"Crypto"}, store, NewStateStore(database))

	return machine
}
