# gotalk

## Simple Multi-user ad-hoc communication program.
**The communication is secured using tls over tcp. The program can be started in server mode or in client mode (see below). In client mode the program starts a graphical user interface to accomodate both conversations and status messages. The client GUI is built using `fyne` (https://fyne.io/), a portable graphical toolkit.**

&NewLine; 
**Build the software:**
- install golang
- install fyne (see also https://developer.fyne.io/index.html)
- clone / download this repo
- rename `model.go.example` to `model.go`
- run `openssl ecparam -genkey -name prime256v1 -out server.key`
- replace `serverKey` constant dummy content with content of `server.key` file
- run `openssl req -new -x509 -key server.key -out server.pem -days 3650`
- replace `rootCert` constant dummy content with content of `server.pem` file
- run `go build`


&NewLine;  
&NewLine;  

**Run the software in server mode:**

	gotalk server [<port>] 

**Examples:**

    ./gotalk server
    ./gotalk server 8089

Server termination by SIGHUP (for the time being)

**Run the software in client mode:**

	gotalk client [<nickname> [<address>] [<port>]]

**Examples:**

    ./gotalk client
    ./gotalk client MyNick
    ./gotalk client MyNick 127.0.0.1
    ./gotalk client MyNick 127.0.0.1 8089

![Client example](https://github.com/ulritter/gotalk/blob/main/example.png)

Client commands:
- `/exit` - terminate connection and exit
- `/list` - displays active users in room
- `/nick <nickname>` - change nickname

&NewLine;   
&NewLine;   

In all cases \<address\> defaults to `localhost` and \<port\> defaults to `8089` and \<nickname\> defaults to `J_Doe`

