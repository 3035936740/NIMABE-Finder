package main

import (
	"finder/pkg/finder"
	"finder/pkg/util/config"
	"flag"
	"log"
	"os"
	"os/signal"
)

const (
	version = "0.0.1-alpha"
)
const (
	flagConfig = "c" //-c %{path} 自定义配置文件路径
)

var (
	configFile = flag.String(flagConfig, "config.toml", "set a config file")
)

// waitQuit 阻塞等待应用退出(ctrl+c / kill)
func waitQuit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Printf("finder exit")
}

func main() {
	flag.Parse()

	conf := &finder.Config{}
	if err := config.Load(*configFile, conf); err != nil {
		log.Panic(err)
	}

	log.Println("load config success:", conf)

	srv := finder.New(
		finder.WithLog(conf.Log.FilePath, conf.Log.MaxAgeHours, conf.Log.MaxRotationMegabytes),
		finder.WithServer(conf.Server.Address, conf.Server.Port),
	)

	log.Printf("finder %s running...%v", version, srv.Start())

	waitQuit()
}
