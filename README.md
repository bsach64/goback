# Goback

**Goback** is a distributed local backup system designed to facilitate secure and efficient file backups. The system includes both a client and server component to manage and perform backups across different machines.

## Getting Started

### Prerequisites

- Go 1.22 or higher
- SSH key for server operations (generate one for yourself)

### Install golanci-lint
It is necessary to install golanci-lint to run the project.

for linux users:
```bash
yay golangci-lint
```
Other users:
```bash
https://golangci-lint.run/welcome/install/
```

### Generating SSH Key

To run the server, you will need a SSH key. Generate it using:

```bash
ssh-keygen -t rsa -b 4096 -f id_rsa
```
> Place the generated id_rsa file inside the `private` directory in the root folder of the project.

## Building the Project
You can build the project using `make`. Available commands in the `Makefile` are:

* `clean`: Remove build artifacts
* `test`: Lint the code and Run tests
* `build`: Lint the code and Build the binary
* `format`: Format and tidy the code 
* `lint`: Lint, format and tidy the code
* `tidy`: Tidy up dependencies

To build the project, run:
```bash
make build
```
OR
```bash
make
```

## Running the Binary

After building the project, you can run the binary using:

Client: To open up the client prompt, run:

```bash
./goback client
```
Server: To open the server prompt, which shows listen logs and an exit command, run:

```bash
./goback server
```
If you want to run using `goback` only you can export a path variable instead run inside the root folder1: 

```bash
export PATH=$PATH:$(pwd)
```


## Usage
### Client
The client command opens a prompt allowing you to interact with the backup system. It includes options for uploading files, listing directories, and exiting.

![client-usage](https://i.imgur.com/Ocu6ZqM.png)
### Server
The server command starts a server that listens for incoming client connections. You will see logs of the server's activity and have the option to exit the server prompt.

![alt text](https://i.imgur.com/crutHwp.png)

## Contributing
Contributions are welcome! Please ensure that your code adheres to the project's style guidelines and includes tests.
