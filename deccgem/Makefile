deccgem: $(shell find . -name "*.go")
	GOPATH=`pwd` go build src/cmd/deccgem.go
	tar -c deccgem | gzip -c > deccgem.tgz
	mv deccgem.tgz ../roles/decc.server/files
