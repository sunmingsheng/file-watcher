package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os/exec"
	"strings"
	"sync"
)

type Configs struct {
	mapping map[string]string
	mu sync.Mutex
}

var configs = Configs{
	mapping: make(map[string]string),
}

var watcher *fsnotify.Watcher
var err error

const mappingKey = "mapping"

func main() {

	watcher, err = fsnotify.NewWatcher()
	defer watcher.Close()

	viper.SetConfigName("file_watcher")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	loadConf()
    watchConf()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if hookCommand, ok := configs.mapping[event.Name]; ok {
						hook(hookCommand)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()
	<-done
}


func loadConf() {
	configs.mu.Lock()
	defer configs.mu.Unlock()
	if !viper.IsSet(mappingKey) {
		panic("config error")
	}
	configs.mapping = viper.GetStringMapString(mappingKey)
	for fileName, _ := range configs.mapping {
		_ = watcher.Remove(fileName)
		_ = watcher.Add(fileName)
	}
}

func watchConf() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		loadConf()
	})
}

func hook(hookParam string) {
	if len(hookParam) == 0 {
		return
	}
	command := strings.Split(hookParam, " ")
	cmd := exec.Command(command[0], command[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	fmt.Println("command: " + hookParam)
	if err != nil {
		fmt.Println(fmt.Sprintf("cmd.Run() failed with %s\n", err))
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Println(fmt.Sprintf("out:\n%s\nerr:\n%s\n", outStr, errStr))
}