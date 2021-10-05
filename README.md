# How to use
First fill your server detail first in **.env.example** file, and then rename it into **.env** and then copy **.env** to same location bin you want to use, ex : **bin/windows**

```bash
cp .env.example bin/windows/.env
```

Let say you want forward your local port 3306 to your server so it can be accessed by public using port 6689, you can use

```bash
tobira.exe -local=3306 -remote=6689
```

or you if you want random remote port you can just use

```
tobira.exe -local=3306
```

![Tobira usage](/screenshot/1.png)

# Tobira Server
You can use your ssh server on host server, by editing **/etc/ssh/sshd_config** and enable reverse tunneling. or you can use docker image that I've built and modified from [LinuxServer/DockerOpenSSHServer](https://github.com/linuxserver/docker-openssh-server). all you need is just run docker-compose that have already been in this repo see **docker-compose.yml**
```yaml
---
version: "2.1"
services:
  openssh-server:
    image: scz10/tobira:v0.1
    network_mode: "host"
    container_name: tobira-server
    hostname: tobira-server #optional
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Asia/Jakarta
      - PASSWORD_ACCESS=true
      - USER_PASSWORD=passwordhere
      - USER_NAME=tobira
    restart: unless-stopped
```
## Parameters

Container images are configured using parameters passed at runtime (such as those above). These parameters are separated by a colon and indicate `<external>:<internal>` respectively. For example, `-p 8080:80` would expose port `80` from inside the container to be accessible from the host's IP on port `8080` outside the container.

| Parameter | Function |
| :----: | --- |
| `--hostname=` | Optionally the hostname can be defined. |
| `-p 2222` | ssh port |
| `-e PUID=1000` | for UserID - see below for explanation |
| `-e PGID=1000` | for GroupID - see below for explanation |
| `-e TZ=Europe/London` | Specify a timezone to use EG Europe/London |
| `-e PUBLIC_KEY=yourpublickey` | Optional ssh public key, which will automatically be added to authorized_keys. |
| `-e PUBLIC_KEY_FILE=/path/to/file` | Optionally specify a file containing the public key (works with docker secrets). |
| `-e PUBLIC_KEY_DIR=/path/to/directory/containing/_only_/pubkeys` | Optionally specify a directory containing the public keys (works with docker secrets). |
| `-e SUDO_ACCESS=false` | Set to `true` to allow `linuxserver.io`, the ssh user, sudo access. Without `USER_PASSWORD` set, this will allow passwordless sudo access. |
| `-e PASSWORD_ACCESS=false` | Set to `true` to allow user/password ssh access. You will want to set `USER_PASSWORD` or `USER_PASSWORD_FILE` as well. |
| `-e USER_PASSWORD=password` | Optionally set a sudo password for `linuxserver.io`, the ssh user. If this or `USER_PASSWORD_FILE` are not set but `SUDO_ACCESS` is set to true, the user will have passwordless sudo access. |
| `-e USER_PASSWORD_FILE=/path/to/file` | Optionally specify a file that contains the password. This setting supersedes the `USER_PASSWORD` option (works with docker secrets). |
| `-e USER_NAME=linuxserver.io` | Optionally specify a user name (Default:`linuxserver.io`) |
| `-v /config` | Contains all relevant configuration files. |


# Build Tobira by yourself 
To build for your cpu architecture, you can build this project like this : 
```bash
# build for linux/amd64
GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/tobira main.go 
# build for windows/amd64
GOOS=windows GOARCH=amd64 go build -o bin/windows/tobira.exe main.go
# build for linux/arm64
GOOS=linux GOARCH=arm64 go build -o bin/linux/arm64/tobira main.go
# build for linux/arm/v7
GOOS=linux GOARCH=arm GOARM=7 go build -o bin/linux/armv7/tobira main.go
```