package encryption

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/configuration"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/log"
	"github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/internal/types"
	"maunium.net/go/maulogger/v2"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/id"
	"maunium.net/go/mautrix/util/dbutil"

	_ "github.com/mattn/go-sqlite3" // driver for sqlite3
)

// GetCryptoStore initializes a sql crypto store
//lint:ignore SA4009 Try to move to MySQL later
func GetCryptoStore(debug bool, db *sql.DB, config *configuration.Matrix) (crypto.Store, id.DeviceID, error) {
	var deviceID id.DeviceID
	username := fmt.Sprintf("@%s:%s", config.Username, strings.ReplaceAll(strings.ReplaceAll(config.Homeserver, "https://", ""), "http://", ""))

	err := os.MkdirAll("data", 0755)
	if err != nil {
		return nil, deviceID, err
	}

	// Currently the library does not support MySQL
	db, err = sql.Open("sqlite3", "data/olm.db")
	if err != nil {
		log.Warn(err.Error())
		return nil, deviceID, err
	}

	// Use device ID from database if available otherwise fallback to settings
	err2 := db.QueryRow("SELECT device_id FROM crypto_account WHERE account_id=$1", username).Scan(&deviceID)
	if err2 != nil && err2 != sql.ErrNoRows {
		log.Warn("Failed to scan device ID: " + err2.Error())
		deviceID = id.DeviceID(config.DeviceID)
	}

	cryptoDB, err := dbutil.NewWithDB(db, "sqlite3")
	if err != nil {
		return nil, deviceID, err
	}

	cryptoStore := crypto.NewSQLCryptoStore(cryptoDB, dbutil.MauLogger(maulogger.Create()), username, deviceID, []byte(config.DeviceKey))

	err = cryptoStore.Upgrade()
	if err != nil {
		return nil, deviceID, err
	}

	return cryptoStore, deviceID, err
}

// GetOlmMachine initializes a new olm machine
func GetOlmMachine(debug bool, client *mautrix.Client, store crypto.Store, database types.Database, stateStore *StateStore) *crypto.OlmMachine {
	if client == nil {
		log.Warn("client nil")
		panic("client nil")
	}
	if store == nil {
		log.Warn("store nil")
		panic("store nil")
	}

	logger, err := newCryptoLogger(debug)
	if err != nil {
		panic(err)
	}
	machine := crypto.NewOlmMachine(client, logger, store, stateStore)

	return machine
}
