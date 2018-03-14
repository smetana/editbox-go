## Development Notes


### Go Installation
```
# Choose preferred version from https://golang.org/dl/
wget https://redirector.gvt1.com/edgedl/go/go1.9.2.linux-amd64.tar.gz
tar -C ~/go1.9 -xzf go1.9.2.linux-amd64.tar.gz

# run or add to ~/.bashrc
export GOROOT=$HOME/go1.9
export PATH=$PATH:$GOROOT/bin
```
### Workspace setup
```
export GOPATH=$HOME/projects/go
export PATH=$PATH:$(go env GOPATH)/bin
```

### Run Project
```
cd ~/projects/editbox
go get -v -d ./

go run _examples/example.go
```

### Tests
```
go test
```

### Before commit
```
gofmt -w editbox.go
```
