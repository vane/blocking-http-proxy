blocking-http-proxy
====

## Description

Just blocking stuff that is in ```block.yaml```  
This code is ugly as hell and I know it but works with http/https traffic and is 
keeping me away from unwanted websites that save data on my computer.  

## Install
```bash
go build main.go
./main
```
## Usage
Point your computer http/https traffic to ```localhost``` port ```11666``` open browser and cry with blocked internet 

## Tips
Wanna go to facebook - remove last line in block.yaml
