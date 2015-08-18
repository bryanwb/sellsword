build:
	godep go build github.com/bryanwb/sellsword/cmd/sellsword

install: build
	cp -f sellsword ssw /usr/local/bin/
	chmod +x /usr/local/bin/sellsword /usr/local/bin/ssw
