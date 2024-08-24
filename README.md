# enpass

## Introduction
I wrote some Python scripts to manipulate my Enpass database, but I wanted to do more. I wanted a portable binary that can be used to read and parse my Enpass data. There is already an amazing project [here](https://github.com/hazcod/enpass-cli.git), but I wanted to make the CLI a little more intuitive and add some features to it. I borrowed the logic from hazcod's project and used Cobra to build the CLI.

## Requirements
* golang 1.22+
* darwin or linux OS

## Installation
Clone the repository and from within the repository directory, type `make build`. This will create a directory with the given value of `GOOS` and install the binary there. It will also create a tarball which will eventually be used for Homebrew formulae.

## Features
* List all entries from the database
* List all entries from the database and show the passwords
* Display the password for a given item to STDOUT
* Copy the password for a given item to the clipboard
* Output in JSON, YAML, list, or table format
* Show trashed items
* Try to auto-detect the location of the Enpass vault
* Specify columns to sort by for list and show operations
* Filter by multiple logins (subtitles), titles, categories, or uuids using wildcards
* Colorize output for JSON, YAML, list, and default views
* Colors are defined centrally so all colorized outputs use the same color scheme
* And more....

## Usage
```
enpass is a command line interface for the Enpass password manager

Usage:
  enpass [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  copy        Copy the password of a vault entry to the clipboard
  help        Help about any command
  list        List vault entries without displaying the password
  pass        Print the password of a vault entry to STDOUT
  show        List vault entries, displaying the password
  version     Print the current enpass version

Flags:
  -c, --category stringArray   Filter based on record category. Wildcards (%) are allowed. Can be used multiple times.
  -h, --help                   help for enpass
  -k, --keyfile string         Path to your Enpass vault keyfile.
      --log string             The log level, one of: debug, error, fatal, info, panic, trace, warn (default "info")
  -l, --login stringArray      Filter based on record login. Wildcards (%) are allowed. Can be used multiple times.
      --nocolor                Disable colorized output and logging.
  -n, --non-interactive        Disable prompts and fail instead.
  -p, --pin                    Enable PIN.
      --sensitive              Force category and title searches to be case-sensitive.
  -t, --title stringArray      Filter based on record title. Wildcards (%) are allowed. Can be used multiple times.
      --type string            The type of your card. (password, ...) (default "password")
  -u, --uuid stringArray       Filter based on record uuid. Can be used multiple times.
  -v, --vault string           Path to your Enpass vault.

Use "enpass [command] --help" for more information about a command.
```

## List Usage
```
List vault entries without displaying the password

Usage:
  enpass list [flags]

Flags:
  -h, --help                  help for list
      --json                  Output the data as JSON.
      --list                  Output the data as list, similar to SQLite line mode.
  -o, --orderby stringArray   Specify fields to sort by. Can be used multiple times. (default [title])
      --table                 Output the data as a table.
      --trashed               Show trashed items.
      --yaml                  Output the data as YAML.

Global Flags:
  -c, --category stringArray   Filter based on record category. Wildcards (%) are allowed. Can be used multiple times.
  -k, --keyfile string         Path to your Enpass vault keyfile.
      --log string             The log level, one of: debug, error, fatal, info, panic, trace, warn (default "info")
  -l, --login stringArray      Filter based on record login. Wildcards (%) are allowed. Can be used multiple times.
      --nocolor                Disable colorized output and logging.
  -n, --non-interactive        Disable prompts and fail instead.
  -p, --pin                    Enable PIN.
      --sensitive              Force category and title searches to be case-sensitive.
  -t, --title stringArray      Filter based on record title. Wildcards (%) are allowed. Can be used multiple times.
      --type string            The type of your card. (password, ...) (default "password")
  -u, --uuid stringArray       Filter based on record uuid. Can be used multiple times.
  -v, --vault string           Path to your Enpass vault.
```

## Examples
List the `Discord` record and output to JSON format
```
$ enpass list --title Discord --json
Enter vault password:
[
    {
        "uuid": "xxxxxxx",
        "created": 1722787097,
        "card_type": "password",
        "updated": 1722787097,
        "title": "Discord",
        "subtitle": "user@example.com",
        "category": "login",
        "label": "Password",
        "sensitive": true,
        "icon": "{\"fav\":\"discord.com\",\"image\":{\"file\":\"misc/login\"},\"type\":1,\"uuid\":\"\"}",
        "raw_value": "xxxxxxx"
    }
]
```

Copy the `Foo` record's password to the clipboard
```
$ enpass copy --title Foo
Enter vault password:
The password for "Foo" was copied to the clipboard
```

List all records with the login user@example.com and output in to table format
```
enpass list --login user@example.com --table
Enter vault password:
title                 login                 category
--------------------- --------------------- --------
Discord               user@example.com      login
Playstation           user@example.com      login
Xbox                  user@example.com      login
Twitch                user@example.com      login
```

List all records containing GitHub, forcing case-sensitivity, and output to a table format
```
enpass list --title %GitHub% --sensitive --table
Enter vault password:
title                       login                        category
--------------------------- ---------------------------- --------
GitHub                      gdanko@example.com           login
GitHub                      https://github.com           computer
Work GitHub (gdanko-work)   https://github.workplace.com computer
```

## To Do
* Allow the user to specify fields to display
* Make sure all the other stuff works
