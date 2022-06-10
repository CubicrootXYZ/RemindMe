# Matrix reminder and calendar bot - RemindMe

[![GitHub license](https://img.shields.io/github/license/CubicrootXYZ/matrix-reminder-and-calendar-bot)](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/CubicrootXYZ/matrix-reminder-and-calendar-bot)](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/issues)
[![Actions Status](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/workflows/Main/badge.svg?branch=main)](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/actions)

![Logo](media/Logo.png)

A matrix bot that handles reminders and knows your agenda.

## Example

![Example chat interaction](media/Chat_Example.png)

## 📋 Features

* Schedule reminders
* Edit and delete reminders
* Timezone support
* Natural language understanding
* Quick actions via reactions
* Daily message with open reminders for the day
* Repeatable reminders
* iCal export of all reminders _(via API)_
* Block users _(via API)_
* Allow bot to be invited _(enable in settings)_
* End to end encrypted channels _(enable in settings)_

The following features are seen as **experimental**, we do not recommend them for use in production. Data losses or data leaks might happen.

* Multi-User channels

## 👩‍🔧 Contribute

I really enjoy help in making this bot even better. So we all can enjoy the work on this project please follow the rules. 

You can contribute in many ways to this project:

* Report issues and improvement ideas
* Test new features
* Improve the code basis (open a pull request)
* Add new features (open a pull request)

### Issues, ideas and more

Please submit your issues or specific feature requests as "Issues". 

General ideas and concepts can be discussed in the "Discussions" section.

### Contributing code

Fork this repository and add your changes. Open a pull request to merge them in the main branch of this repository.

## 🔍 How to use the bot

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

### Setup a multi user channel

1. Create a new matrix group channel
2. Invite the bot - you are the admin of the channel now
3. Invite any user you want to participate
4. Add user that should be able to interact with the bot with `add user @username`


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

You should persist the data stored in `data/` (`/run/data/` in the docker image) and frequently back it up. There is crucial information for the end to end encryption stored in.

### Plain

Plain installation requires knowledge about building go binaries and installing arbitrary packages. We recommend using the prebuild docker containers.

1. Install the dependencies
    1. You need `libolm-dev` with at least version 3 (e.g. for debian buster run `apt install libolm-dev/buster-backports -y`)
    2. Install `gcc` which is required by cgo (e.g. for debian buster run `apt install gcc -y`)
2. Download the code
3. Run `go build -o /app/bin /app/cmd/remindme/main.go` to build the binary in `/app/bin`
4. Setup your config file
5. Run the binary
6. Make sure to persists the `data` folder as it contains important data for the service 

### Docker

Different versions are available on docker hub:

[Docker Hub](https://hub.docker.com/r/cubicrootxyz/remindme)

You are missing a docker container for your architecture? We'd love to see you contributing to this project by opening a pull request with the build instructions for it.

## 📚 Further documentation 

Take a look into our [wiki](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/wiki). It provides you with further information and troubleshooting guides.

### API

The bot offers an API. It needs to be enabled in the settings where the api key for the "Admin-Authentication" needs to be set. 

The documentation can be found at [cubicrootxyz.github.io/matrix-reminder-and-calendar-bot/](https://cubicrootxyz.github.io/matrix-reminder-and-calendar-bot/).

## 🎁 Related projects

Any project missing? Open a pull request!

* [RemindMe-Web](https://github.com/CubicrootXYZ/RemindMe-Web) - Web UI for controlling the bot

## ❤️ Attribution

Great thanks to the libraries used in this project:

* [Mautrix](https://github.com/tulir/mautrix-go)
* [Gorm](https://gorm.io/)
* [Gorm MySQL](https://github.com/go-gorm/mysql)
* [Naturaldate](https://github.com/tj/go-naturaldate)
* [Configor](https://github.com/jinzhu/configor)
* [Uniuri](https://github.com/dchest/uniuri)
* [Go-Naturalduration](https://github.com/CubicrootXYZ/gonaturalduration)
* [Gorm](https://github.com/go-gorm/gorm)
* [Stretchr/Testify](https://github.com/stretchr/testify)
* [Gin](https://github.com/gin-gonic/gin)
* [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
* [zap](https://github.com/uber-go/zap)
* [CubicrootXYZ/gormlogger](https://github.com/CubicrootXYZ/gormlogger)
