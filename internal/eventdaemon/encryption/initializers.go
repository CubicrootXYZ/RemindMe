package encryption

import (
	"database/sql"
	"fmt"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/id"
)

// GetCryptoStore initializes a sql crypto store
func GetCryptoStore(db *sql.DB, config *configuration.Matrix) (crypto.Store, error) {
	account := fmt.Sprintf("%s/%s", config.Username, config.DeviceID)

	cryptoStore := crypto.NewSQLCryptoStore(db, "sql", account, id.DeviceID(config.DeviceID), []byte(config.DeviceKey), cryptoLogger{"Crypto"})

	err := cryptoStore.CreateTables()

	return cryptoStore, err
}

// GetOlmMachine initializes a new olm machine
func GetOlmMachine(client *mautrix.Client, store crypto.Store) *crypto.OlmMachine {
	// TODO crypto state
	machine := crypto.NewOlmMachine(client, cryptoLogger{"Crypto"}, store, nil)

	return machine
}
