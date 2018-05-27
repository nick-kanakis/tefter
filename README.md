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

7. Export notes all notes from "lists" and "expenses" notebooks to a json file
```
tefter export -n lists,expenses
```

8. Import notes from a json file at path "documents/notes.json"
```
tefter import documents/notes.json
```

9. Print available notebooks & notes (titles only)
```
tefter overview -d
```

10. Initiate the graphic interface.
```
tefter print -a
```

11. Search notes for "2018" keyword
```
tefter search 2018
```

12. Initiate rest API endpoint on port 8081
```
tefter serve -p 8081
```

13. Update note with id 42, remove tag "2018" and add tag "2019" also set the title to "Bali 2019"
```
tefter update 42 -t "Bali 2019" --tags -2018,2019
```

14. Update notebook with title "lists", to title "2018 lists"
```
tefter updateNotebook "lists" "2018 lists"
```
