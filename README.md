# enpass

## Introduction
I wrote some Python scripts to manipulate my Enpass database, but I wanted to do more. I wanted a portable binary that can be used to read and parse my Enpass data. There is already an amazing project [here](https://github.com/hazcod/enpass-cli.git), but I wanted to make the CLI a little more intuitive and add some features to it. I borrowed the logic from hazcod's project and used Cobra to build the CLI.

## Requirements
* golang 1.22+
* darwin or linux OS

## Installation
* Clone the repository and from within the repository directory
* Type `make build`. This will create a bin directory and install the binary there. It will also create a tarball which will eventually be used for Homebrew formulae.
* Copy <repo_root>/enpass.yml.SAMPLE to ~/.enpass.yml

## Installation (Homebrew)
* `brew tap gdanko/homebrew`
* `brew install gdanko/homebrew/enpass`

## Features
* List all entries from the database
* List all entries from the database and show the passwords
* Display the password for a given item to STDOUT
* Copy the password for a given item to the clipboard
* Output in YAML, list, or table format
* Show trashed items
* Try to auto-detect the location of the Enpass vault
* Specify columns to sort by for list and show operations
* Filter by multiple logins (subtitles), titles, categories, or uuids using wildcards
* Colorize output for YAML, list, and default views
* Colors are defined centrally so all colorized outputs use the same color scheme
* Basic options can be set in ~/.enpass.yml, please see enpass.yml.SAMPLE

## The `~/.enpass.yml` file
This file currently supports the following options
* `vault_path` - The absolute path to your vault file
* `colors` - Configure colors for output
    * `alias_color`
    * `anchor_color`
    * `bool_color`
    * `key_color`
    * `null_color`
    * `number_color`
    * `string_color`
* `output_style` - One of `list`, `table`, or `yaml`
* `default_labels` - A YAML array of labels, you will need to parse your database file to find all available values.
* `orderby` - A YAML array of fields to sort the output by.

## Usage
```
$ enpass
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
  -y, --label stringArray      Filter based on record field label. Can be used multiple times
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
$ enpass list --help
List vault entries without displaying the password

Usage:
  enpass list [flags]

Flags:
  -h, --help                  help for list
      --list                  Output the data as list, similar to SQLite line mode.
  -o, --orderby stringArray   Specify fields to sort by. Can be used multiple times. Valid: card_type, category, created, label, last_used, subtitle, title, updated
      --table                 Output the data as a table.
      --trashed               Show trashed items.
      --yaml                  Output the data as YAML.

Global Flags:
  -c, --category stringArray   Filter based on record category. Wildcards (%) are allowed. Can be used multiple times.
  -k, --keyfile string         Path to your Enpass vault keyfile.
  -y, --label stringArray      Filter based on record field label. Can be used multiple times
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
List the `Discord` record and output to YAML format
```
$ enpass list --title Discord --yaml
Enter vault password:
- uuid: xxxxxxx
  created: 1710107066
  card_type: password
  title: Discord
  subtitle: user@example.com
  category: login
  label: Password
  last_used: 1730565183
  sensitive: true
  icon: "{\"fav\":\"discord.com\",\"image\":{\"file\":\"misc/login\"},\"type\":1,\"uuid\":\"\"}"
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

## Troubleshooting
You need to get the value of the hex-encoded key
* In `openEncryptedDatabase()` you need to add a line to print the key to the console, `fmt.Println(hex.EncodeToString(dbKey)[:masterKeyLength])`
* Once you have the key, do the following
    * Open the database file using DB Browser for SQLite
    * You will see a `SQLCipher encryption` dialog
        * Check `SQLCipher3 Defaults`
        * In the `Password` field, put `x'<encoded_hext_string>'`
        * Click `OK`
* You can now query the database to look

## To Do
* Allow the user to specify fields to display
* Make sure all the other stuff works
