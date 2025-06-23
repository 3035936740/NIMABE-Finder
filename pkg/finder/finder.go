package finder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type MusicInfo struct {
	Data []map[uint]MusicDataInfo `json:"data"`
}

type MusicDataInfo struct {
	Title   string `json:"title"`
	Version uint   `json:"version"`

	AsciiTitle string `json:"asciiTitle"`
	Genre      string `json:"genre"`
	Artist     string `json:"artist"`

	MID uint `json:"entryId"`

	Difficult map[string]MusicDifficult `json:"difficulties"`
}

type MusicDifficult struct {
	Beginner   uint `json:"beginner"`
	Normal     uint `json:"normal"`
	Hyper      uint `json:"hyper"`
	Another    uint `json:"another"`
	Legendaria uint `json:"legendaria"`
}

type NickInfo struct {
	Nick map[string]uint `json:"nicks"`
}

func New(opts ...Options) (f *Finder) {
	f = &Finder{}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (f *Finder) Start() error {
	if e := f.reload(); e != nil {
		return e
	}

	if e := f.sdvxLoadUni(); e != nil {
		return e
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(context *gin.Context) {
		data, _ := context.GetRawData()
		context.Set("data", data)

		clientIp := context.ClientIP()
		if ip := context.GetHeader("X-Real-IP"); len(ip) > 0 {
			clientIp = ip
		}
		context.Set("clientIp", clientIp)

		f.logln(clientIp, context.Request.URL)

		context.Next()
	})

	f.routers(router)

	f.logln("server listen on", f.address, f.port)

	return router.Run(fmt.Sprintf("%s:%d", f.address, f.port))
}

// 加载歌库
func (f *Finder) loadMusicDB(path string) error {
	f.m.RLock()
	defer f.m.RUnlock()

	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	musics := &MusicInfo{}
	err = json.Unmarshal(jsonBytes, musics)
	if err != nil {
		return err
	}

	counts := 0
	for _, data := range musics.Data {
		for mid, music := range data {
			f.mid.Store(mid, music.Title)
			f.name.Store(music.Title, mid)
			f.genre[music.Genre] = append(f.genre[music.Genre], music)
			f.artist[music.Artist] = append(f.genre[music.Artist], music)
			counts++
		}
	}

	f.logln("load total db musics:", counts)
	f.logln("load total db artist:", len(f.artist))
	f.logln("load total db genre:", len(f.genre))

	return nil
}

// 加载外号
func (f *Finder) loadNickName(path string) error {
	f.m.RLock()
	defer f.m.RUnlock()

	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	nick := make(map[string]uint)
	err = json.Unmarshal(jsonBytes, &nick)
	if err != nil {
		return err
	}

	counts := 0
	for nick, mid := range nick {
		f.nick.Store(nick, mid)
		counts++
	}

	f.logln("load total nicks:", counts)

	return nil
}

// 写入外号
func (f *Finder) saveNickName(path string) error {
	f.m.Lock()
	defer f.m.Unlock()

	m := make(map[string]uint)
	f.nick.Range(func(nick, mid any) bool {
		m[nick.(string)] = mid.(uint)
		return true
	})

	bytes, err := json.Marshal(m)

	if err != nil {
		return err
	}

	f.logln("save total nicks:", len(m))

	return os.WriteFile(path, bytes, 0755)
}

func (f *Finder) reload() error {
	f.nick = sync.Map{}
	f.name = sync.Map{}
	f.mid = sync.Map{}
	f.genre = make(map[string][]MusicDataInfo)
	f.artist = make(map[string][]MusicDataInfo)

	if err := f.loadMusicDB(filepath.Join(FullPath(), "music_data.json")); err != nil {
		return err
	}

	return f.loadNickName(filepath.Join(FullPath(), "music_nick.json"))
}

func (f *Finder) sdvxLoadUni() error {
	// SDVXLoad
	if e := f.SDVXManager.LoadData("music_db.xml"); e != nil {
		return e
	}

	if e := f.SDVXManager.LoadAliases("aliases.json"); e != nil {
		return e
	}

	return nil
}
