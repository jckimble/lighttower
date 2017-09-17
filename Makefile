VERSION=$(shell git describe --tags --abbrev=0 --dirty="-dev")

all: clean lighttower

lighttower:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$(VERSION)" -a -installsuffix cgo -o lighttower-linux-amd64 .
clean:
	rm -f lighttower-linux-amd64

release: lighttower
	sha256sum lighttower-linux-amd64 > sha256sum.txt
	openssl dgst -sha256 -sign lighttower.key -out lighttower-linux-amd64.sig lighttower-linux-amd64
	go get github.com/itchio/gothub
	$(GOPATH)/bin/gothub release --user jckimble --repo lighttower --tag $(VERSION) --name "Automatic Release $(VERSION)" --description ""
	$(GOPATH)/bin/gothub upload --user jckimble --repo lighttower --tag $(VERSION) --name "lighttower-linux-amd64" --file ./lighttower-linux-amd64
	$(GOPATH)/bin/gothub upload --user jckimble --repo lighttower --tag $(VERSION) --name "sha256sum.txt" --file ./sha256sum.txt
	$(GOPATH)/bin/gothub upload --user jckimble --repo lighttower --tag $(VERSION) --name "lighttower-linux-amd64.sig" --file ./lighttower-linux-amd64.sig
