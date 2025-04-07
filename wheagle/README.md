# WHEAGLE
Wheagle is a c2 framework that allows for easier control of multiple C2's and implant sessions from compromised machines.


Golang Tensor flow https://www.youtube.com/watch?v=Al_T3O27DO8

## Features
1. C2 generation
  Generate exe/elf/dll os so files to run as C2 on the fly.
2. Implant generation
  Generate exe/elf/dll or so files to use as implants and assign each to a Mother ship
3. Screenshot
  Use the command "screenshot" to take a screenshot of a c2 or implant
4. Persistance
  WINDOWS: Persist through by installing implant or minion as service,startup file or through the registry(All this should be easily detectable by most/any EDR/ good AV's)
  LINUX: Persist using a cronjob or a service
5. Get to root/System by injecting into a process
### 6. HUNTER
  1. Generate a dropper of your choice(stager or stageless)
  2. Hide payloads into zip or iso files
  3. Spike your payload by injecting it into another legitimate file.
  4. Encode your payloads with xor,hex,base64 (Nothing major)
  5. Generate shellcode from your payload. (Biggest disadvantage, you have to create a C2 or minion the get shellcode from it)
          *Use any payload doesn't matter whether it' from wheagle or not.*
  6. Generate a bootkit/rootkit from a payload. curiosity

### Pipeline Features
Tunnelling and Pivoting
C2 for http/https/udp and tcp
Minions surpporting both protobuf and other protocols or separate(we'll know when it happens)
Plugins for hunter
Multiple Operator support.

# TUNNEL WORKABILITY

If tunnel address is present:
  1. AdminTunnelAddress: This is a tunneled address exposing the adminC2 address allowing operator to connect to.
  2. ImplantTunnelAddress: This is a tunneled address for implants to connect to the admin.


#TODO
Brush up on all possible generators
Finish on working rather adding work to pool. Just get all commands working properly and zero timeout/waiting
Create screenshots
Cretae RevShell (streaming or whichever the way possible)
Ensure droppers work perfectly(mutants and all self/url down-loaders) Basically Persistence and
Create encrypted Communication
Brush up on all commands and ensure they are working
Try out our lsass dumper too
Create DNS Tunnelling (might take some time)

Git Push
Take a brake from codding and all this then come back and start all the extra/additional features

Create DB to save/load connectors and sessions (try deleting n starting new)
Establish connections well:
  Authenticate
  Run Command
    InC2 and In Implant
    Run Shell on both
    Take screenshot
    Transfer Files
    Enable various default privilege escalation (read from itself and inject into memory)
Create more minions/persistence methods:
  BootKit
  RootKit


Brush up on the server and get the templates ready
Create The API server

go get gopkg.in/vmihailenco/msgpack.v2
