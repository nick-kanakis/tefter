# tefter
[![Build Status](https://travis-ci.org/nicolasmanic/tefter.svg?branch=master)](https://travis-ci.org/nicolasmanic/tefter)
[![Go Report Card](https://goreportcard.com/badge/github.com/nicolasmanic/tefter)](https://goreportcard.com/report/github.com/nicolasmanic/tefter)
[![Coverage Status](https://coveralls.io/repos/github/nicolasmanic/tefter/badge.svg?branch=master)](https://coveralls.io/github/nicolasmanic/tefter?branch=master)

Tefter is a simple note manager written in Go, inspired by [Shiori](https://github.com/RadhiFadlillah/shiori).
New notes can be created (with vim), updated and viewed without exiting the terminal.
Notes can also be collected to notebooks and flagged with tags.

## Table of Content
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Examples](#examples)

## Features
- Use it without ever leaving your terminal
- Organize notes per notebook and tags
- Search notes based on notebooks, tags, or by a keyword
- Import/Export from/to a json file.
- All package into one executable file
- Rest API thor 3rd party integration

## Installation

You can download the latest version from the release page [TODO]. To build from source, you should have `go` installed and then run:
```
go get github.com/nicolasmanic/tefter
```
Finally put the binary in your `PATH`.

## Usage

```
Tefter is a simple memo book application

Usage:
  tefter [command]

Available Commands:
  account        Add/Delete/Print account
  add            Create a new note
  delete         Delete one or more notes based on ID(s)
  deleteNotebook Delete one or more notebooks based on title
  export         Exports notes to json format
  help           Help about any command
  import         Import notes from json file
  overview       Take a quick glance at the available notebooks and notes
  print          Print notes
  search         Search notes given a keyword
  serve          Initiate rest API interface
  update         Update existing note
  updateNotebook Set new title to an existing notebook

Flags:
  -h, --help   help for tefter

Use "tefter [command] --help" for more information about a command.
```

## Examples

1. Create a new account for rest API.
```
tefter account add
```

2. Delete a account .
```
tefter account delete
```

3. Print active usernames .
```
tefter account print
```

4. Create a new note, with title: "Bali 2018", tags: "vacation" & "summer" and insert it into notebook with title: "lists".
```
tefter add -t "Bali 2018" --tags vacation,summer -n lists
```

5. Delete notes, with id 42 & 23
```
tefter delete 42,23
```
6. Delete notebooks with title: "lists" and "expenses".

```
tefter deleteNotebook lists,expenses
```

