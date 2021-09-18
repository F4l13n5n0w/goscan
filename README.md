# Goscan is a golang version fast port scanner
## It can be used in Red Teaming for low priv windows user without install nmap

## Guide, help manual
goscan.exe -h


```root@vscode-server001:~/codes/goscan# ./goscan -h

    ________        _________                       ._. ._.
   /  _____/  ____ /   _____/ ____ _____    ____    | | | |
  /   \  ___ /  _ \\_____  \_/ ___\\__  \  /    \   | | | |
  \    \_\  (  <_> )        \  \___ / __ \|   |  \   \|  \|
   \______  /\____/_______  /\___  >____  /___|  /   __  __
          \/              \/     \/     \/     \/    \/  \/
  
                                             by  F4l13n5n0w
                                             version: 0.0.2
Usage of ./goscan:
  -iL string
        Input the file path of target IP list
  -ip string
        target IP address
  -p string
        ports range, split by ',' (default is top1k)
  -rt int
        TCP message reader timeout (seconds), this is for service detection (default 5)
  -st int
        TCP scan connection timeout (seconds) (default 5)
  -thread int
        set thread number, make sure not too high (default 100)
root@vscode-server001:~/codes/goscan# ```


Generate top-x ports from nmap-services file:
```/* Command to generate top-x ports from nmap-services file:
   sort -r -k3 /usr/share/nmap/nmap-services | grep tcp | head -n 100 | cut -f 2 | cut -d '/' -f 1 | sed -n 'H;${x;s/\n/,/g;s/^,//;p;}'

   ┌──(root💀kali)-[/usr/share/nmap]
   └─# sort -r -k3 /usr/share/nmap/nmap-services | grep tcp | head -n 10 | cut -f 2 | cut -d '/' -f 1 | sed -n 'H;${x;s/\n/,/g;s/^,//;p;}'
   80,23,443,21,22,25,3389,110,445,139
*/```