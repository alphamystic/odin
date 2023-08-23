# MY TODO's
https://www.youtube.com/watch?v=dXG2SIB1XMM
https://www.youtube.com/watch?v=k3oOlBIW2hk
https://www.youtube.com/watch?v=JPuA9TdC-2M

## BULDING PIPELINE
0. Always initialize a new feature and place into the Pipeline
1. Recon (should receive a target) classify recon output and return a vulnerability @Done
2. Attack (should receive vulnerability with a payload that works)
3. Pivoting and priv escalation

@TODO Go back to db and session manager and add a windows file system directory
Fix windows generator zoo.ShellcodeRunner thingy
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w"   test.go


First we finish wheagle to create c2 and agent @DONE
Create a shell-code runner where possible
Add a downloader from c2 and from operator
Look into other features slowly as per the need
Add a GetAgentData for both admin and agent

Finish on getting recon data HalfDone(Finish on NMAP)
Persistence of recon data to DB @DONE
create at-least  (2) vulnerability scanners, one for rico and another for private

Add login functionalities
Persist to MySQL DB functions
Find a way to stream something like rdp or any desk from an internal network

So anything backdoor or mothership will be handled by wheagle and odin will only load it as a plugin
With that in mind let odin load the scanners during a test and set a mode to prevent restarting/rebuilding
odin will manage all the other BS (users,contacts,appointments,reports,api users and keys,issues,projects~plugins)

AD https://learn.microsoft.com/en-us/azure/active-directory/authentication/concept-authentication-web-browser-cookies
dievus ad-generator
rubeus,mimickatz
https://www.youtube.com/watch?v=gH9qyHVc9-M&list=RDCMUCJ2U9Dq9NckqHMbcUupgF0A&start_radio=1&rv=gH9qyHVc9-M&t=1971
cache poisoning: https://portswigger.net/research/practical-web-cache-poisoning
MSF Reference: https://docs.rapid7.com/metasploit/
Blue teaming ref: https://github.com/A-poc/BlueTeam-Tools
Mimickatz: https://github.com/gentilkiwi
SafeKatz: https://github.com/GhostPack/SafetyKatz

Vulnerable laravel:
git clone https://github.com/appelsiini/vulnerable-laravel-app.git

ad mindmap
https://github.com/Orange-Cyberdefense/ocd-mindmaps


https://mega.nz/file/3XJCyD5C#qAda14pWUjd5u4wjOYmzCI52UMa1rUFulh7V0kBGZk8

https://gist.github.com/jhaddix/78cece26c91c6263653f31ba453e273b.
https://hackerone.com/reports/961046
https://nuxtjs.org/docs/2.x/concepts/server-side-rendering
https://hackerone.com/reports/335330
https://hackerone.com/reports/1180697
https://hackerone.com/reports/352869
https://hackerone.com/reports/178152
https://hackerone.com/reports/827052
https://hackerone.com/reports/1062888
https://hackerone.com/reports/303744
https://hackerone.com/reports/1154542
https://hackerone.com/reports/1125425
https://securitylab.github.com/research/now-you-c-me-part-two/
https://hackerone.com/reports/584603
https://hackerone.com/reports/591295
https://hackerone.com/reports/547630
https://securitylab.github.com/research/libssh2-integer-overflow/
https://securitylab.github.com/research/libssh2-integer-overflow-CVE-2019-17498/
https://github.com/swisskyrepo/PayloadsAllTheThings
https://github.com/EdOverflow/bugbounty-cheatsheet
https://hackerone.com/reports/300305
https://hackerone.com/reports/759247
https://hackerone.com/reports/974892

vulerable laravel
https://github.com/appelsiini/vulnerable-laravel-app

https://twitter.com/Shubham_pen/status/1639126141373153282
https://twitter.com/AccentInvesting/status/1638873298359713793



go build -buildmode=c-shared -ldflags="-w -s -H=windowsgui" -o Updater.dll

sudo apt install -y mingw-w64
Install this for kali/debian builds
sudo apt-get install gcc-multilib && apt-get install gcc-mingw-w64
GOOS=windows GOARCH=386 CGO_ENABLED=1 CXX_FOR_TARGET=i686-w64-mingw32-g++ CC_FOR_TARGET=i686-w64-mingw32-gcc go build test.go
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared -ldflags="-w -s -H=windowsgui" test.go

https://medium.com/@warrenbutterworth/finding-initial-access-on-a-real-life-penetration-test-86ed5503ae48
10.81..19.146
net income perpendicularism

// GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared -ldflags="-w -s -H=windowsgui" -o updater.dll
// Windows: go build -buildmode=c-shared -ldflags="-w -s -H=windowsgui" -o updater.dll
// Run: rundll32.exe ./updater.dll,Updater
