# timeIt
Contraption to measure run time of standalone binaries. 

```bash
#!/bin/bash 
cmd="uptime"
start=$(date +%s)
$cmd 
end=$(date +%s)
echo "$end took $((end - start)) second(s) to run $cmd"
```

```
$ ./measure.sh 
 19:18:50 up 10:13,  1 user,  load average: 0.39, 1.04, 2.14
1658283530 took 0 second(s) to run uptime
```

While this works well for bash, it is a <a href="https://stackoverflow.com/questions/673523/how-do-i-measure-execution-time-of-a-command-on-the-windows-command-line">bit more</a> difficult when it comes to measuring execution time in other environments / operating systems.

# How does timeIt work?
Consider this example:
```bash
$ ls
timeIt-linux-amd64
$ cp timeIt-linux-amd64 uptime
$ ./uptime 
Fatal error: open /tmp/timeIt/uptime: no such file or directory
```
`timeIt` is copied as `uptime`. When we try to run `./uptime`, `timeIt` will try to execute the named binary in the pickup folder `/tmp/timeIt`. 

```
$ mkdir /tmp/timeIt && cp -v `which uptime` /tmp/timeIt/uptime
'/usr/bin/uptime' -> '/tmp/timeIt/uptime'
$ ./uptime 
 11:17:44 up 1 day, 19:49,  5 users,  load average: 0.25, 0.79, 0.99
```
`timeIt` executes `/tmp/timeIt/uptime` in this example.
This also creates the logs we are ultimately interested in:
```bash
$ cat /tmp/timeIt/timeIt.log 
1658513865005659 timeIt/1.03c took 1.160092ms to run cmd: /tmp/timeIt/uptime
```

## Pickup Folder
`timeIt` uses `C:\Temp\timeIt` as "pickup" folder on Windows and `/tmp/timeIt/` everywhere else. 

Post-execution logs also are stored in this folder.


# Use case
https://github.com/Azure/kubelogin/issues/102 is a real-life example.

You want to run `kubectl` and see how much time it takes for it to finish. Behind the scenes, `kubectl` evokes yet another binary, `kubelogin`, as an exec plugin. 

You would like to know how much time `kubectl` and `kubelogin` each take to complete. 

Meet `timeIt`: a simple attempt to solve this sort of challenges.
## Linux
```bash
$ kubectl get nodes
NAME                                STATUS   ROLES   AGE     VERSION
aks-nodepool1-30293584-vmss000000   Ready    agent   7d21h   v1.24.0
aks-nodepool1-30293584-vmss000001   Ready    agent   7d21h   v1.24.0
aks-nodepool1-30293584-vmss000002   Ready    agent   7d21h   v1.24.0

```


Create the pickup folder and copy `kubelogin` and `kubectl` binaries over:

```bash 
$ mkdir /tmp/timeIt/`

$ which kubectl kubelogin 
/usr/local/bin/kubectl
/usr/local/bin/kubelogin

$ cp /usr/local/bin/kubelogin /usr/local/bin/kubectl /tmp/timeIt/ -v 
'/usr/local/bin/kubelogin' -> '/tmp/timeIt/kubelogin'
'/usr/local/bin/kubectl' -> '/tmp/timeIt/kubectl'
```
Save and keep aside the original binaries:
```bash
$ sudo cp /usr/local/bin/kubelogin /usr/local/bin/kubelogin.save -v 
'/usr/local/bin/kubelogin' -> '/usr/local/bin/kubelogin.save'
$ sudo cp /usr/local/bin/kubelogin /usr/local/bin/kubelogin.save -v
'/usr/local/bin/kubelogin' -> '/usr/local/bin/kubelogin.save'
```
Replace with `timeIt`:
```bash
~/git/timeIt$ sudo cp bin/timeIt-linux-amd64 /usr/local/bin/kubelogin -v 
'bin/timeIt-linux-amd64' -> '/usr/local/bin/kubelogin'
~/git/timeIt$ sudo cp bin/timeIt-linux-amd64 /usr/local/bin/kubectl -v 
'bin/timeIt-linux-amd64' -> '/usr/local/bin/kubectl'
```

At this point, run `kubectl` command:
```bash
$ kubectl get nodes
NAME                                STATUS   ROLES   AGE     VERSION
aks-nodepool1-30293584-vmss000000   Ready    agent   7d22h   v1.24.0
aks-nodepool1-30293584-vmss000001   Ready    agent   7d22h   v1.24.0
aks-nodepool1-30293584-vmss000002   Ready    agent   7d22h   v1.24.0
```
`timeIt` will evoke the binaries in the `pickup folder` and will log the results in `timeIt.log`:
```bash
$ cat /tmp/timeIt/timeIt.log 
1658285443001959 timeIt/1.03a took 35.263568ms to run cmd: /tmp/timeIt/kubelogin get-token --login azurecli --server-id redacted
1658285443434396 timeIt/1.03a took 515.303212ms to run cmd: /tmp/timeIt/kubectl get nodes
```

### Another example:
```bash
$ kubelogin remove-tokens
$ kubectl get jobs
No resources found in default namespace.
```
Inspecting the timeIt log:
```bash
$ tail -2  /tmp/timeIt/timeIt.log 
1658286406623108 timeIt/1.03a took 289.759517ms to run cmd: /tmp/timeIt/kubelogin get-token --login azurecli --server-id redacted
1658286406958111 timeIt/1.03a took 670.823205ms to run cmd: /tmp/timeIt/kubectl get jobs
```
## Undoing the changes
Revert the binaries you saved to their original form:
```bash 
$ sudo cp /usr/local/bin/kubelogin.save /usr/local/bin/kubelogin
$ sudo cp /usr/local/bin/kubectl.save /usr/local/bin/kubectl
```

## Windows 

```powershell
PS C:\Users\user1\Downloads\kStuff> dir


    Directory: C:\Users\user1\Downloads\kStuff


Mode                 LastWriteTime         Length Name
----                 -------------         ------ ----
-a----         7/19/2022   8:44 PM       46279680 kubectl.exe
-a----         7/19/2022   8:45 PM       43229184 kubelogin.exe
-a----         7/19/2022   8:39 PM        2271232 timeIt-win-amd64.exe
```
Temporarily add this folder to the `$env:PATH`:
```
PS C:\Users\user1>  $env:path += ";C:\Users\user1\Downloads\kStuff"
```
Verify `$env:path`:
```
PS C:\Users\user1> kubelogin --version
kubelogin version
git hash: v0.0.16/b61a6ef2b83d81e57b5ecc15d384b9da367f365e
Go version: go1.17.12
Build time: 2022-07-18T20:32:10Z
```
Create the pickup folder and copy kubelogin and kubectl binaries over:
```PowerShell
PS C:\Users\user1\Downloads\kStuff> mkdir "C:\Temp\timeIt\"
PS C:\Users\user1\Downloads\kStuff> cp .\kubectl.exe  C:\Temp\timeIt\
PS C:\Users\user1\Downloads\kStuff> cp .\kubelogin.exe C:\Temp\timeIt\
```
Save the original binaries first:

```PowerShell
PS C:\Users\user1\Downloads\kStuff> cp .\kubectl.exe .\kubectl.exe.save
PS C:\Users\user1\Downloads\kStuff> cp .\kubelogin.exe .\kubelogin.exe.save
```
 Replace with `timeIt`:
```Powershell
PS C:\Users\user1\Downloads\kStuff> cp .\timeIt-win-amd64.exe .\kubelogin.exe
PS C:\Users\user1\Downloads\kStuff> cp .\timeIt-win-amd64.exe .\kubectl.exe
```
Run `kubectl`: 
```
PS C:\Users\user1> kubectl get nodes -A
NAME                                STATUS   ROLES   AGE     VERSION
aks-nodepool1-30293584-vmss000000   Ready    agent   7d23h   v1.24.0
aks-nodepool1-30293584-vmss000001   Ready    agent   7d23h   v1.24.0
aks-nodepool1-30293584-vmss000002   Ready    agent   7d23h   v1.24.0
```

Inspect the timeIt logs:
```powershell
PS C:\Users\user1> cat C:\Temp\timeIt\timeIt.log
1658289550088029 timeIt/1.03a took 2.8090766s to run cmd: C:\Temp\timeIt\kubelogin.exe get-token --login azurecli --server-id Redacted
1658289552063280 timeIt/1.03a took 5.7506374s to run cmd: C:\Temp\timeIt\kubectl.exe get nodes -A
```
## Undoing the changes
Revert the binaries you saved to their original form:
```PowerShell 
PS C:\Users\user1\Downloads\kStuff> cp .\kubectl.exe.save .\kubectl.exe
PS C:\Users\user1\Downloads\kStuff> cp .\kubelogin.exe.save .\kubelogin.exe
```
# Priming 
Environments with AntiVirus/Malware scanners may take a tad bit longer for the first run - the OS sees a never-before-seen binary so it waits for an AV scan to finish before running it. A simple workaround would be to "prime" the binary once by passing it a random flag:
```bash
$ cp timeIt-linux-amd64 uname 
$ cp `which uname` /tmp/timeIt/ 
$ ./uname -h
/tmp/timeIt/uname: invalid option -- 'h'
Try '/tmp/timeIt/uname --help' for more information.
Fatal error: exit status 1
$ 
```
This would make the binaries "known" to the ecosystem, so subsequent invocations should not be blocked on AV scans. 

# Optional: Build 
If you prefer to build the binaries yourself, you will need the [Go Programming Language](https://golang.org/dl/) installed on your System. 

Clone this repo and build as: 
``` 
go build -o timeIt
```

# Caveats
- Execution time of binaries that background themselves (example: `calc.exe`, `notepad.exe`) and return immediately cannot be measured accurately. 

