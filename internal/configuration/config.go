package configuration

// Config holds all settings and credentials the application needs
type Config struct {
	Debug            bool `default:"false"`
	MatrixBotAccount Matrix
	MatrixUsers      []string `required:"true"`
	Database         Database
	Webserver        Webserver
	BotSettings      BotSettings
}

// Matrix holds the information for accessing the bots account
type Matrix struct {
	Username   string `required:"true"`
	Password   string `required:"true"`
	Homeserver string `required:"true"`
	DeviceID   string `default:"123456"`
	DeviceKey  string `required:"true"`
}

// BotSettings holds information about the bot itself
type BotSettings struct {
	AllowInvites bool
	MaxUser      int64 `default:"-1"`
}

// Database holds all data for connection to the database
type Database struct {
	Connection string `required:"true"`
}

// Webserver holds all data for the webserver
type Webserver struct {
	Enabled bool
	APIkey  string
	BaseURL string
}
