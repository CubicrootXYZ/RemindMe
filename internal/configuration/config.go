package configuration

// Config holds all settings and credentials the application needs
type Config struct {
	Debug            bool `default:"false"`
	MatrixBotAccount Matrix
	MatrixUser       string `required:"true"`
	Database         Database
}

// Matrix holds the information for accessing the bots account
type Matrix struct {
	Username   string `required:"true"`
	Password   string `required:"true"`
	Homeserver string `required:"true"`
}

// Database holds all data for connection to the database
type Database struct {
	Connection string `required:"true"`
}
