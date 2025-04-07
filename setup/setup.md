# SETUP

This script is work in progress just git clone and go build for now:

## Requirements
*This installation guide assumes a kali linux install as it leverages built in tools. If only interested in the C2 then you can just ignore this*
1. Nmap
2. TheHarvester
3. Subfinder

## Getting it.
> git clone https://github.com/alphamystic/odin
> cd odin
> go mod tidy

*If on kali/debian concide adding this, The rest from arch you probably can figure things out incase of failure*
Install this for kali/debian builds
> sudo apt install -y mingw-w64
> sudo apt-get install gcc-multilib && apt-get install gcc-mingw-w64

> cd odin

### Install the DB
Create a database odin
Import the mysql file on odin/lib/db/odin.sql

@TODO

### Running Odin
*Always remember to run odin with sudo privileges as it runs nmap.
At the root folder,  build odin
>  go build -ldflags="-s -w"  --race  .
> sudo ./odin

#### Start scanning with odin


### Running Wheagle
cd into wheagle and build it out:
> cd wheagle
>  go build -ldflags="-s -w"  --race  .
> /wheagle -ban=true or /wheagle -ban=false

### Running Loki
By default loki runs on port 9000
cd into loki:
> go build -ldflags=" -s -w" --race
> ./loki
