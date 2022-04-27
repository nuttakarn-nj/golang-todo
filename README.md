# Build with parameter --ldflags
go build -o main --ldflags "-X main.buildcommit=`git rev-parse --short HEAD` -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`"
main