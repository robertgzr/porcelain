all:
	@go build -ldflags "-X main.buildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.buildRev=`git rev-parse --short HEAD` -X main.buildTag=`git describe --tags --long`"
