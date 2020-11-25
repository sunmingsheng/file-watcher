## **fileWatcher**

fileWatcher是一个使用 golang 实现的文件变动监听器，支持通过配置文件配置的方式添加监听任务，支持动态的加载配置文件，代码量极少，使用起来也比较简单

一:使用方式

1:下载原代码

`
go get github.com/sunmingsheng/file-watcher
`

2:编译安装

`
make install
`

3:添加监听任务，编辑/ect/file_watcher.yaml，添加以下配置(example)

```
mapping:
    {PATH}/nginx.conf: nginx -s reload
```

4:修改{PATH}/nginx.conf，观察nginx是否完成reload操作(example)


二:实现依赖:

1:使用fsnotify库实现文件变化监听

2:使用viper库实现配置的动态加载
