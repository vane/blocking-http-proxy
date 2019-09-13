blocking-http-proxy
====

## Description

Just blocking stuff that is in ```block.yaml```  
This code might be a bit ugy cause it's my first code in go and I know it 
but works with http/https traffic and is keeping me away from unwanted 
websites that save data on my computer.  

## Install
```bash
go build main.go
./main
```
## Usage
By default point your computer http/https traffic to ```localhost``` port ```11666``` open browser and cry with blocked internet  
You can type optional arguments like ```host```  ```port``` or custom yaml file ```block.yaml```
All blocked requests are logged to file ```block.log``` and allowed ones to ```allow.log```

## Tips
Wanna go to facebook - remove last line in block.yaml
