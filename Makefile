all:
	@go build -ldflags "-X main.date=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.commit=`git rev-parse --short HEAD` -X main.version=`git describe --tags --long`"
