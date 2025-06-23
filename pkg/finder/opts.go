package finder

import (
	"log"
	"sync"

	l "finder/pkg/util/log"
)

type Finder struct {
	logger  *l.Log
	address string
	port    uint

	m sync.RWMutex

	mid sync.Map //uint,string

	name sync.Map //string,uint
	nick sync.Map //string,uint

	genre  map[string][]MusicDataInfo //string,[]MInfo
	artist map[string][]MusicDataInfo //string,[]MInfo
	SDVXManager
}

type Options func(*Finder)

// WithLog 自定义配置
func WithLog(path string, age, size uint) Options {
	return func(f *Finder) {
		var err error
		f.logger, err = l.New(
			path,
			age,
			size)
		if err != nil {
			f.panic(err)
		}
	}
}

// WithServer 自定义配置
func WithServer(address string, port uint) Options {
	return func(f *Finder) {
		f.address = address
		f.port = port
	}
}

// logf 打印日志(如果没有启用则打到控制台)
func (f *Finder) logf(format string, v ...interface{}) {
	if f.logger != nil {
		f.logger.Printf(format, v...)
	}
	log.Printf(format, v...)
}

// logln 打印日志(如果没有启用则打到控制台)
func (f *Finder) logln(v ...any) {
	if f.logger != nil {
		f.logger.Println(v...)
	}
	log.Println(v...)
}

// panic 崩溃
func (f *Finder) panic(v any) {
	f.logf("%v", v)
	panic(v)
}
