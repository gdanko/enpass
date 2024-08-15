# enpass

## Introduction
I wrote some Python scripts to manipulate my Enpass database, but I wanted to do more. I wanted a portable binary that can be used to read and parse my Enpass data. There is already an amazing project [here](https://github.com/hazcod/enpass-cli.git), but I wanted to make the CLI a little more intuitive and add some features to it. I borrowed the logic from hazcod's project and used Cobra to build the CLI.

## Requirements
* golang 1.22+
* darwin or linux OS

## Installation
Clone the repository and from within the repository directory, type `make build`. This will create a directory with the given value of `GOOS` and install the binary there. It will also create a tarball which will eventually be used for Homebrew formulae.

## Features
* List all items from the database
* List all items from the database and show the passwords
* Display the password for a given item to STDOUT
* Copy the password for a given item to the clipboard
* Output in JSON, YAML, list, or a tabular format
* Show trashed items
* Try to auto-detect the location of the Enpass vault
* And more....

## Usage
```
enpass is a command line interface for the Enpass password manager

Usage:
  enpass [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  copy        Copy the password of a vault entry matching FILTER to the clipboard
  help        Help about any command
  list        List vault entries matching FILTER without password
  pass        Print the password of a vault entry matching FILTER to stdout
  show        List vault entries matching FILTER with password
  version     Print the current enpass version

Flags:
      --and                    Combines filters with AND instead of default OR.
  -c, --category stringArray   The category of your card. Can be used multiple times.
  -h, --help                   help for enpass
  -k, --keyfile string         Path to your Enpass vault keyfile.
      --log string             The log level from debug (5) to panic (1). (default "4")
      --nonInteractive         Disable prompts and fail instead.
      --pin                    Enable PIN.
      --type string            The type of your card. (password, ...) (default "password")
  -v, --vault string           Path to your Enpass vault

Use "enpass [command] --help" for more information about a command.
```

## Examples
List the `Foo` record and output to JSON format
```
$ enpass list Foo --json
Enter vault password:
[
    {
        "uuid": "xxxxxxx",
        "created": 1722787097,
        "card_type": "password",
        "updated": 1722787097,
        "title": "Foo",
        "subtitle": "user@example.com",
        "category": "login",
        "label": "Password",
        "sensitive": true,
        "icon": "{\"fav\":\"foo.com\",\"image\":{\"file\":\"misc/login\"},\"type\":1,\"uuid\":\"\"}",
        "raw_value": "xxxxxxx"
    }
]
```

Copy the `Foo` record's password to the clipboard
```
$ enpass copy Foo
Enter vault password:
The password for "Foo" was copied to the clipboard
```

List entries containing `user@example.com`
```
$ enpass list user@example.com
Enter vault password:
title                     login                category
------------------------- -------------------- --------
Amazon                    user@example.com     login
AT&T                      user@example.com     login
AWS                       user@example.com     login
Bandcamp                  user@example.com     login
Best Buy                  user@example.com     login
```

## To Do
* Allow the user to specify fields to display in the tabular view
* Allow better flexibility with the card types
* Make sure all the other stuff works
