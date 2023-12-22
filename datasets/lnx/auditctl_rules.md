# AuditCTL Rules

1. Track Execution of bash
  > auditctl -a always,exit -F arch=b64 -S execve -F exe=/bin/bash

2. Monitor File Access for /etc directory
  > auditctl -a always,exit -F path=/etc/ -F perm=rwa

3. Monitor Failed login attempts
  > auditctl -a always,exit -F arch=b64 -S login -F success=0

  *This works on 64bit architectures,change your os architectures type*

4. Monitor process Execution by a specific user(this one ir root)
> auditctl -a always,exit -F arch=b64 -S execve -F auid=root

5. Monitor changes to critical system file like shadow/passwd
  > auditctl -a always,exit -F arch=b64 -S open -F path=/etc/passwd -F perm=wa
  > auditctl -a always,exit -F arch=b64 -S open -F path=/etc/shadow -F perm=wa

6. Monitor changes to network configuration attempts
  > auditctl -a always,exit -F arch=b64 -S open -F path=/etc/hosts -F perm=wa
  > auditctl -a always,exit -F arch=b64 -S open -F path=/etc/resolv.conf -F perm=wa

7. Monitor privilege escalation attempts.
  > auditctl -a always,exit -F arch=b64 -S execve -F exe=/usr/bin/sudo

8. To watch a file for changes (2 ways to express):
  > auditctl -w /etc/shadow -p waauditctl -a exit,always -F path=/etc/shadow -F perm=wa

9. To recursively watch a directory for changes (2 ways to express):
  > auditctl -w /etc/ -p waauditctl -a exit,always -F dir=/etc/ -F perm=wa

Reference [https://linux.die.net/man/8/auditctl] auditctl
Rules are stored at /etc/audit/audit.rules

## Requirements
This are requirements to have relevant linux configurations to log system events.
1. Enable auditd and install necessary plugins.
  > sudo apt-get install auditd audispd-plugins
  > sudo systemctl enable auditd
