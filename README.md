# ODIN

## What is Odin
Odin is an all in one pentesting suite intended to automate the pentesting stages.  Created with the Mitre Attack framework in mind. This entails from Reconnaissance to impact. To start of, I built the recon stage(working on adding shodan and IOT Capabilities) and a C2 that will help manage compromised networks/forests.
Odin is also made with Detection in mind in that it will provide detection capabilities. This is to GIVE guys getting into cyber security an all in one suite to introduce them to the tasks performed all round.
In the long run it comes in handy as a purple teaming toolkit and suitable for small organizations/SME's with a zero budget on cyber security. Hope you get to enjoy it or tweak it as per your needs and it helps you out.
All sessions,connectors are stored into the db, runtime data is called through api calls to your currently running c2(default at 55677,55678)
At the very least, Odin is meant for an organization to run it by themselves single handedly.


## Getting Started

> go to installation https://github.com/alphamystic/odin/setup

## How Odin,Wheagle and Loki work

k
### Start odin
Odin by itself is a stand alone enumeration tooolkit that takes in a target does basic recon and store the data into a mysql/maria db for you.
On start up, run the cli with ./odin while at thie root direcctory
Type help or:
  1. cli ask for target
  2. cli ask for multiple targets file

Enumerate target: Handlers
  If .com find ip and multiple Ips,
  Scan Ip,multiple IPs for services running
  Enumerate the services
#### Future implementations
  1. Scan with shodan api key's
  2. Return the potential vulnerabilities or An AD domain
  3. Can also take in a vulnerability and compromise a target

#### Descriptors
*skip this if ypu're interested only in running it out, check out setup*
  1. Kowalski
We should have a manage engine that injests logs per second and determines what step the agent should take.


## FEATURES
### Odin Vulnerability Scanner (OVS)
Automated Recon
Add your own vulnerability scanners and load them as per your need or use the default minimal version
CVE Checker (Still working on how to implement a CVE Database that will act as a reference point)
IT Asset Management
Keep track of APT's,threats and various viruses

### Wheagle (Command line for post Exploitation anm,d managing of multiple targets)
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



## PIPELINE
1. Add a command line for all plugins
2. Loot and persistence
3. Pivoting
4. AD Pentesting
5. Syscall unhooking and EDR Detection and evasion techniques
