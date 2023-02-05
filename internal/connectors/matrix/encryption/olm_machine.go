package encryption

import (
	"github.com/CubicrootXYZ/gologger"
	"github.com/CubicrootXYZ/gologger/olmlogger"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
)

// NewOlmMachine sets up a new olm machine.
func NewOlmMachine(client *mautrix.Client, cryptoStore crypto.Store, stateStore crypto.StateStore, logger gologger.Logger) (*crypto.OlmMachine, error) {
	olmLogger := olmlogger.New(logger)

	olm := crypto.NewOlmMachine(client, olmLogger, cryptoStore, stateStore)
	err := olm.Load()

	return olm, err
}
