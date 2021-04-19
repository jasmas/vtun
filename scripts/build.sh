#!bin/bash
export GO111MODULE=on
export GOPROXY=https://goproxy.cn

UNAME=$(uname)

if [ "$UNAME" == "Linux" ] ; then
    GOOS=linux GOARCH=amd64 go build -o ./bin/vtun ./main.go
elif [ "$UNAME" == "Darwin" ] ; then
    GOOS=darwin GOARCH=amd64 go build -o ./bin/vtun ./main.go
elif [[ "$UNAME" == CYGWIN* || "$UNAME" == MINGW* ]] ; then
    GOOS=windows GOARCH=amd64 go build -o ./bin/vtun.exe ./main.go
fi

echo "done!"
