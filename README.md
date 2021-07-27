# Matrix reminder and calendar bot
A matrix bot that handles reminders and knows your agenda.

**This little boy needs to grow a lot!**

## How to use it

After you have installed the bot he will invite you in a channel accept the invite and you are ready to interact with him. 

### New Reminder

To make a new reminder talk to the bot like this: 
```
Make laundry at Sunday 16:00
```

It tries to understand your natural language as best as it can. 

### List all available commands 

To get all commands just type one of these lines:
```
commands
list all commands
show all commands
```

### List all reminders

You can use one of those commands to list all pending reminders in a channel:
```
list
list all reminders
list reminders
list my reminders
show 
show all reminders
show reminders
reminders 
reminder
```

## Installation

### Plain

1. Download the code
2. Run `go build -o /app/bin /app/cmd/remindme/main.go` to build the binary in `/app/bin`
3. Run the binary

### Docker

Needs to be done :).
