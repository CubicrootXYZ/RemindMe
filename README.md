# Matrix reminder and calendar bot - RemindMe

[![GitHub license](https://img.shields.io/github/license/CubicrootXYZ/matrix-reminder-and-calendar-bot)](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/CubicrootXYZ/matrix-reminder-and-calendar-bot)](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/issues)
[![Actions Status](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/workflows/Main/badge.svg?branch=main)](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/actions)

![Logo](media/Logo.png)

A matrix bot that handles reminders and knows your agenda.

⚠️ The main branch currently contains the `v2` which is not yet considered stable. `v1` is moved to the similiar named branch. ⚠️

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
* Import reminders from an Ical link (see API documentation about third party resources)

## 👩‍🔧 Contribute

See our [contribution guidelines](https://github.com/CubicrootXYZ/RemindMe/blob/main/CONTRIBUTING.md).

## 🔍 How to use the bot

After you have installed the bot it will invite every user in the config in a channel. Accept the invite and you are ready to interact with it.

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

## ⚙️ Installation

See our [installation guides](https://github.com/CubicrootXYZ/RemindMe/wiki/Installation). We provide docker container images or you can build the binary yourself. 

## 📚 Further documentation 

Take a look into our [wiki](https://github.com/CubicrootXYZ/matrix-reminder-and-calendar-bot/wiki). It provides you with further information and troubleshooting guides.

### API

⚠️ The documentation currently is only available for `v1` ⚠️

The bot offers an API. It needs to be enabled in the settings where the api key for the "Admin-Authentication" needs to be set. 

The documentation can be found at [cubicrootxyz.github.io/RemindMe/](https://cubicrootxyz.github.io/RemindMe/).

## 🎁 Related projects

Any project missing? Open a pull request!

* [RemindMe-Web](https://github.com/CubicrootXYZ/RemindMe-Web) - Web UI for controlling the bot (only works with `v1`)

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
* [Golang ICAL](https://github.com/arran4/golang-ical)
* [Mock](https://github.com/golang/mock)
* [Go SQLite 3](https://github.com/mattn/go-sqlite3)
* [RRULE Go](https://github.com/teambition/rrule-go )
