# Go Store

## Key-value memory database with option to persist

![image](https://i.imgur.com/g0OVt3o.png)

## Ever needed a database in a hurry? Go store could be right for you.

<br>
Easy and intuitive command line tool allows you to spin up a database avaliable from web or locally in a few seconds.  
Server can be run over a custom TCP protocol or over HTTP.  
Database can be be kept in memory only or persisted to a json file upon exit.  
You can also spin up a server from an existing file with or without modifying it.  
This is a hobby project and not meant for production use, but could be useful for testing/development phase.

<br>
<br>

#### The CLI tool is powered by [cobra](https://github.com/spf13/cobra)

<br>

## Installation

Go 1.15 and newer

```
go install github.com/maracko/go-store@latest
```

Go versions older then 1.15

```
go get -u github.com/maracko/go-store
```

## Usage

```
go-store [command] [argument] --flags
```

`go-store help` can be used to get help on commands

## Config

### All the options can be set inside a config file or using environment variables.

#### Environment variables have the prefix `GOSTORE_`. Like `GOSTORE_PORT`

#### Config file in format `.gostore-config.{extension}` file in the working directory

#### The file can be saved in any language viper supports (json, yaml, toml, etc.)

#### To learn about the flags you can use, see below.

## Examples

<br>

## HTTP

<br>

```
go-store server HTTP -p 8888 -l /home/mario/database.json
```

This command will open and read the json file provided under `-l` flag and start a HTTP server on port `-p`  
You can then send HTTP requests to the server  
Logs will be outputed to stdout in format `Time` `Method` `URL` `Host`

```
2021/05/02 18:54:02 HTTP server started
2021/05/02 18:54:22 GET /testKey localhost:8888
```

### Server flags

- **--location -l** => Location of database file to read from (if blank or not provided a empty database in memory only will be initialised)
- **--port -p** => Port on which to start serving http
- **--tls-port** => Port on which to start serving https
- **--memory -m** => If present database won't be saved upon exit (even if read from a file first)
- **--private-key -k** => Used for HTTPS. Put a path to key
- **--certificate -c** => Used for HTTPS. Put a path to certificate
- **--token -t** => Used for auth. Send in `Authorization` header
- **--continous-write -c** => If you want to keep saving the DB to the disks
- **--write-interval -i** => How many minutes to wait between writes. Default is 1 minute, if 0 will always write
  <br>

### **HTTP Requests**

<br>

**http://{host}:{port}/{key}**
<br>

- GET => returns key/keys
- POST => creates a new key
- PATCH => update existing key
- DELETE => delete key/keys

For retrieving operations just add key/s in the URI path. To retrieve multiple values set multiple keys split with a comma.
<br>

### Data

When you want to create/update keys you must send data inside request **body in JSON format**.

### Example:

**GET**  
 `http://localhost:8888/myKey`  
 or  
 `http://localhost:8888/myKey,myOtherKey,anotherKey`
<br/>

**POST**  
 `http://localhost:8888`  
 _BODY_ =

```json
{
  "key": "myKey",
  "value": "myValue"
}
```

<br/>

**PATCH**  
 `http://localhost:8888`  
 _BODY_ =

```json
{
  "key": "myKey",
  "value": "myNewValue"
}
```

<br/>

**DELETE**  
 `http://localhost:8888/myKey`  
 or  
 `http://localhost:8888/myKey,myOtherKey,anotherKey`

<br/>

## TCP

<br>

```
go-store server TCP -p 8888 -l /home/mario/database.json

2021/05/02 18:57:33 TCP server started
2021/05/02 18:58:09 Accepted connection from [::1]:34126
```

### Server flags

- **--location -l** => location of database file to read from (if blank or not provided a empty database in memory only will be initialised)
- **--port -p** => port on which to start serving http
- **--memory -m** => if present database won't be saved upon exit (even if read from a file first)
- **--continous-write -c** => if you want to keep saving the DB to the disks
- **--write-interval -i** => how many minutes to wait between writes. Default is 1 minute
  <br>

```
go-store client -s localhost -p 9999

Welcome to go-store server!
2021/05/02 18:58:01 Connected to localhost:9999
$:set foo bar
Created new key foo
$:get foo
bar
$:
```

**TCP currently only supports strings for both key and value, and will do no encoding on them (so no complex types)**  
<br>

## TCP Client

```
go-store client -s localhost -p 8888
```

### Client flags

- **--command -c** => put a list of commands to execute and you will get output to stdout. Optionally chain them with `**` to execute them in sequence.

<br>

- **get [key]** => returns a single key
- **set [key] [value]** => set a new key
- **upd [key] [value]** => update existing key
- **del [key]** => deletes key
  <br>
