&{$env:GOARCH="386"; go build -a -o .\build\ecalc32.exe .}
&{$env:GOARCH="amd64"; go build -a -o .\build\ecalc.exe .}
