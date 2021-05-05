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

```
go install github.com/maracko/go-store
```

## Usage

```
go-store [command] [argument] --flags
```

`go-store help` can be used to get help on commands

## Examples
<br>

## HTTP
<br> 

```
go-store server HTTP -p 8888 -l /home/mario/database.json
```
This command will open and read the json file provided under `-l` flag and start a HTTP server on `-p`  
You can then send HTTP requests to the server  
Logs will be outputed to stdout
```
2021/05/02 18:54:02 HTTP server started
2021/05/02 18:54:22 GET /go-store?key=aa localhost:8888
```

### Server flags

- **--location -l** => location of database file to read from (if blank or not provided a empty database in memory only will be initialised)
- **--port -p** => port on which to start the server
- **--memory -m** => if present database won't be saved upon exit (even if read from a file first)
<br>

### HTTP Requests
<br>

 **http://{host}:{port}/go-store?key=key**
<br>

- GET => returns key/keys
- POST => creates a new key
- PATCH => update existing key
- DELETE => deletes a key

For retrieving operations a query param `key` must be set. To retrieve multiple values set multiple keys split with a comma
<br>

### Example:  
 **GET**  `http://localhost:8888/go-store?key=myKey,myOtherKey,anotherKey`
<br> 

### Data

When you want to create/update keys you must send data inside request **body in JSON format**.
<br>

## TCP
<br>

```
go-store server TCP -p 8888 -l /home/mario/database.json

2021/05/02 18:57:33 TCP server started
2021/05/02 18:58:09 Accepted connection from [::1]:34126
```
TCP server supports the same flags as HTTP. To interact with it use the `go-store client`

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
**TCP currently only supports strings for both key and value, a.k.a will do no encoding on them (so no complex types)**  
<br>



## TCP Client commands
<br>

- **get [key]** => returns a single key  
- **set [key] [value]** => set a new key  
- **upd [key] [value]** => update existing key
- **del [key]** => deletes key
<br>

## TODO

- Make TCP client/server support complex types and commands
- Support HTTPS for HTTP server
- Add authorization support
- TODO in comment
- Other things I will think of

