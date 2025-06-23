package finder

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	//Log 日志
	Log struct {
		FilePath             string //日志存储路径
		MaxAgeHours          uint   //日志轮转最大生命周期(小时)
		MaxRotationMegabytes uint   //日志最大文件大小(MB)
	}

	//Server 服务器
	Server struct {
		Address string //服务器地址
		Port    uint   //端口
	}
}

func FullPath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Panic(err)
	}

	path, err := filepath.Abs(file)
	if err != nil {
		log.Panic(err)
	}

	return path[0:strings.LastIndex(path, string(os.PathSeparator))] + string(os.PathSeparator)
}