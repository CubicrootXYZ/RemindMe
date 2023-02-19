package encryption

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/CubicrootXYZ/gologger"
	_ "github.com/mattn/go-sqlite3" // driver
	"maunium.net/go/maulogger/v2"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/id"
	"maunium.net/go/mautrix/util/dbutil"
)

// NewCryptoStore sets up a new crypto store.
// The crypto store saves data into a SQLite file, it will be created at data/olm.db.
func NewCryptoStore(username, deviceKey, homeserver, confDeviceID string, logger gologger.Logger) (crypto.Store, id.DeviceID, error) {
	var deviceID id.DeviceID
	usernameFull := fmt.Sprintf("@%s:%s", username, strings.ReplaceAll(strings.ReplaceAll(homeserver, "https://", ""), "http://", ""))

	err := os.MkdirAll("data", 0755)
	if err != nil {
		return nil, deviceID, err
	}

	// Currently the library does not support MySQL
	db, err := sql.Open("sqlite3", "data/olm.db")
	if err != nil {
		logger.Err(err)
		return nil, deviceID, err
	}

	cryptoDB, err := dbutil.NewWithDB(db, "sqlite3")
	if err != nil {
		return nil, deviceID, err
	}

	// Use device ID from database if available otherwise fallback to settings
	err2 := db.QueryRow("SELECT device_id FROM crypto_account WHERE account_id=$1", usernameFull).Scan(&deviceID)
	if err2 != nil && err2 != sql.ErrNoRows {
		logger.Errorf("Failed to scan device ID: " + err2.Error())
		deviceID = id.DeviceID(confDeviceID)
	}

	cryptoStore := crypto.NewSQLCryptoStore(cryptoDB, dbutil.MauLogger(maulogger.Create()), usernameFull, deviceID, []byte(deviceKey))

	err = cryptoStore.DB.Upgrade()
	if err != nil {
		return nil, deviceID, err
	}

	return cryptoStore, deviceID, err
}
