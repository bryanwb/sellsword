build:
	go build sellsword.go

install: build
	cp -f sellsword ssw /usr/local/bin/
	chmod +x /usr/local/bin/sellsword /usr/local/bin/ssw
