# Matrix reminder and calendar bot - RemindMe

[![GitHub license](https://img.shields.io/github/license/CubicrootXYZ/matrix-reminder-and-calendar-bot)](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/CubicrootXYZ/matrix-reminder-and-calendar-bot)](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/issues)
[![Actions Status](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/workflows/Main/badge.svg?branch=main)](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/workflows/actions)


A matrix bot that handles reminders and knows your agenda.

![Example list of reminders](Screenshots/reminders.png)
![Example list of commands](Screenshots/commands.png)

## 📋 Features

* Schedule reminders
* Edit and delete reminders
* Timezone support
* Natural language understanding
* Quick actions via reactions
* Daily message with open reminders for the day
* Repeatable reminders
* iCal export of all reminders

## 👥 Contribute

I really enjoy help in making this bot even better. So we all can enjoy the work on this project please follow the rules. 

### Issues, ideas and more

Please submit your issues or specific feature requests as "Issues". 

General ideas and concepts can be discussed in the "Discussions" section.

### Contributing code

Fork this repository and add your changes. Open a pull request to merge them in the master branch of this repository.

## ℹ️ How to use it

After you have installed the bot he will invite every user in the config in a channel. Accept the invite and you are ready to interact with him.

### New Reminder

To make a new reminder talk to the bot like this: 
* `make laundry at sunday 16:00`
* `walking with the dog 6am`
* `brunch with alan at sunday`

It tries to understand your natural language as best as it can. 

### List all available commands 

To get all commands just type one of these lines:
* `commands`
* `list all commands`
* `show all commands`
* `help`

### Set timezone

* `set timezone Europe/Berlin`

### Daily reminder overview

To activate a daily message with the reminders of the day:

* `set daily reminder at 10:00`
* `change daily reminder to 10am`

To deactivate it:

* `delete daily reminder`

## ⚙️ Installation

In any case you need a config file with your preferences and credentials.

1. Copy the `config.example.yml` from the repository
2. Rename it to `config.yml`
3. Fill in the settings
    1. `debug` do not touch unless you know what you are doing
    2. `matrixbotaccount` those are the credentials and the homeserver (url of the matrix instance) of the bots account. You need to create one yourself.
    3. `matrixusers` those users will be able to interact with the bot. Enter a username in the format `@username:instance.tld` so for the user "test123" at the instance "matrix.org" this would be `@test123:matrix.org`
    4. `database` enter a database connection here. You need to enter a MySQL [connection string](https://github.com/go-sql-driver/mysql#dsn-data-source-name). If your database-server is running at the domain "mydatabase.org" at port 3306 and your credentials to log in are "root" and "12345" and the database you created is named "remindme_database" then your connection string would look like this: `root:12345@tcp(mydatabase.org:3306)/remindme_database`.
4. Now you need to copy the file to the folder where the binary is executed.
    * Using the "plain" method: put the binary you build and the config file in the same folder. Execute them from there.
    * Using the pre-build docker image you need to mount the file to `/run/config.yml`

### Plain

1. Download the code
2. Run `go build -o /app/bin /app/cmd/remindme/main.go` to build the binary in `/app/bin`
3. Setup your config file
4. Run the binary

### Docker

Different versions are available on docker hub:

[Docker Hub](https://hub.docker.com/r/cubicrootxyz/remindme)

## API

The bot offers an API. 

### Calendar

**[GET] /calendar**

Returns a list of available calendars (one per user and channel). The returned token is needed to access further information.

Header-Parameters:

* Authorization: String _your api key_

**[GET] /calendar/{id}/ical**

Returns an `ics` file with all reminders of that calendar.

URL-Parameters:

* ID: uint _the calendars id_

Query-Parameters:

* Token: string _the calendars token_

**[PATCH] /calendar/{id}**

Generates a new secret/token for the calendar and removes the old one.

URL-Parameters:

* ID: uint _the calendars id_

Header-Parameters:

* Authorization: String _your api key_


## ❤️ Attribution

Great thanks to the libraries used in this project:

* [Mautrix](https://github.com/tulir/mautrix-go)
* [Gorm](https://gorm.io/)
* [Gorm MySQL](https://github.com/go-gorm/mysql)
* [Naturaldate](https://github.com/tj/go-naturaldate)
* [Configor](https://github.com/jinzhu/configor)
* [Uniuri](https://github.com/dchest/uniuri)
* [Go-Naturalduration](https://github.com/CubicrootXYZ/gonaturalduration)
