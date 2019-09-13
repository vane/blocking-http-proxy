blocking-http-proxy
====

## Description

Just blocking stuff that is in ```block.yaml```  
This code might be a bit ugy cause it's my first code in golang and I know it 
but it works with http/https traffic and is keeping me away from unwanted 
websites that try to save data to my computer.  

## Install
```bash
go build main.go
```
## Usage
```bash
./main
```
By default point your computer http/https traffic to ```localhost``` port ```11666``` open browser and cry
You can type optional arguments:  
```
Arguments:

  -h  --help   Print help information
      --host   host
      --port   port
      --block  YAML block file
```
So it can be used like this 
```bash
./main --host 0.0.0.0 00 --port 6666 --block block.yaml
```
Requests are logged to file ```block.log``` (for blocked requests) and ```allow.log``` (for allowed ones)  

## Tips
Wanna go to facebook - remove last line in block.yaml
