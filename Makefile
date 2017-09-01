VERSION=$(shell git describe --tags --abbrev=0 --dirty="-dev")

all: clean lighttower

lighttower:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -a -installsuffix cgo -o lighttower-linux-amd64 .
clean:
	rm -f lighttower

release: lighttower
	go get github.com/itchio/gothub
	~/go/bin/gothub release --user jckimble --repo lighttower --tag $(VERSION) --name "Automatic Release $(VERSION)" --description ""
	$(GOPATH)/bin/gothub upload --user jckimble --repo lighttower --tag $(VERSION) --name "lighttower-linux-amd64" --file ./lighttower-linux-amd64
