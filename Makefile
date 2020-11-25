PROJECT="fileWatcher"
BIN_TARGET="/usr/local/bin/file-watcher"
CONF_TARGET="/etc/"

default :
	echo ${PROJECT}

install :
	go get github.com/fsnotify/fsnotify
	go get github.com/spf13/viper
	go build -o ${BIN_TARGET} -v
	chmod +x ${BIN_TARGET}
	cp ./file_watcher.yaml ${CONF_TARGET}
	nohup ${BIN_TARGET} >> /dev/null &

test : install
	go test ./..

.PHONY : default install test
