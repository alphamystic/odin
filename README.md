# ODIN

## Getting Started
>> go to installation https://github.com/alphamystic
### Start odin
> cli ask for target
> cli ask for multiple targets file

> Enumerate target: Handlers
    If .com find ip and multiple Ips,
    Scan Ip,multiple IPs for services running
    Enumerate the services
    Return the potential vulnerabilities or An AD domain
    Can also take in a vulnerability and compromise a target

> Attack the target: Skipper
    Use the specified vulnerabilities to try and get foothold into the network

> Pivotting: Skipper
    Enumerate host and find potential priv escalation points
    Find a way to pivot through the network and infect other machines
    Store CnC incide brain: (AdminC2) This can be a particular forest


### Start Wheagle
> Initialize current sessions and connectors
> generate adminc2
> generate implant as per the adminC2
> interact with adminC2 or an implant

## FEATURES
### Odin Vulnerability Scanner (OVS)
Automated Recon
Add your own vulnerability scanners and load them as per your need or use the default minimal version
CVE Checker (Still working on how to implement a CVE Database that will act as a reference point)
IT Asset Management
Keep track of APT's,threats and various viruses

### Wheagle (Command line for post Exploitation and managing of multiple targets)
Generation of bot c2 and implant
Session Management for both c2 and implant (locally and db persistence)
Multiple C2 administration
Shellcode generation
Run shellcode
Generate a library (.dll or .so) as both minion and admin fro persistence
Process injection
Taking screenshots of multiple screens
Encrypted commmunication over grpc's protobuf
Add your Own dropper or C2 in form of a plugin
Tunnelling over http,https/tcp or DNS
    Generate your own DNS server binary and run it anywhere then have your implant communicate to it.
    Generate a DNS implant or run a DNS Client by the default implant
      *Personallly I would use a minimal implant that runs a shellcode of the DNS client*
Staged and Stageless payload

## Directory
1. .brain
    Stores all "brains" rather sessions connections and configurations
2. lib
      A library of all functions helping to run odin
3. payloads
      List of payloads to be used for brute-forcing or enumerations.



### License

### Feedback

## METHODOLOGY
So I want to automate my hacking/pen-testing process. I have for some reason convinced myself that I can do it.

I'm using penguins of madagascar to sought of co-ordinate the attack
  1. Skiper: Get's a shell and pivots or escalate privs and so on.
  2. Kowalski: Create exploits from vulnerabilities
  3. Rico: Brute-force for vulnerabilities
  4. Private: Scans for vulnerabilities

  Private and Rico generate Vulnerabilities and feed them into Kowalski to create exploits.
  They are created in plugin form exporting a function Scan(reconData)  taking recon data to return vulnerabilities

Classified,Mort,KingJulien,Maurice

## PIPELINE
1. Add a command line for all plugins
2. Loot and persistence
3. Pivoting
4. AD Pentesting
5. Syscall unhooking and EDR Detection and evasion techniques
