build:
	godep go build -o sellsword main.go appset.go app.go env.go

install: build
	cp -f sellsword ssw /usr/local/bin/
	chmod +x /usr/local/bin/sellsword /usr/local/bin/ssw
