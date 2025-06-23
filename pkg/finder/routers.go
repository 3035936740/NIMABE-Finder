package finder

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (f *Finder) routers(r *gin.Engine) {
	f.logln("add router GET /")
	r.GET("/", f.getIndex)

	f.logln("add router GET /set")
	r.GET("/set", f.getSet)

	f.logln("add router GET /get")
	r.GET("/get", f.getGet)

	f.logln("add router GET /del")
	r.GET("/del", f.getDel)

	f.logln("add router GET /nicks")
	r.GET("/nicks", f.getNicks)

	f.logln("add router GET /songs")
	r.GET("/songs", f.getSongs)

	f.logln("add router GET /reload")
	r.GET("/reload", f.getReload)

	f.logln("add router Get /sdvx/get")
	r.GET("/sdvx/get", f.getSDVXGet)

	f.logln("add router Get /sdvx/reload")
	r.GET("/sdvx/reload", f.getSDVXReload)

	f.logln("add router Get /sdvx/aliases")
	r.GET("/sdvx/aliases", f.getSDVXAliasList)

	f.logln("add router Get /sdvx/matchid")
	r.GET("/sdvx/matchid", f.getSDVXMatchId)

	f.logln("add router Get /sdvx/existid")
	r.GET("/sdvx/existid", f.getSDVXIdExist)

	f.logln("add router Get /sdvx/addali")
	r.GET("/sdvx/addali", f.addSDVXAlias)

	f.logln("add router Get /sdvx/delali")
	r.GET("/sdvx/delali", f.delSDVXAlias)
}

// getIndex 服务是不是活着
func (f *Finder) getIndex(c *gin.Context) {
	c.String(http.StatusOK, "MaoMaNi - Finder Service Living...")
}

// getSet 设置外号的外号名和值
func (f *Finder) getSet(c *gin.Context) {
	id, _ := c.GetQuery("id")

	nick, _ := c.GetQuery("nick")

	//KV是否都有值
	if id == "" {
		c.String(http.StatusBadRequest, "id was nil")
		return
	}

	if nick == "" {
		c.String(http.StatusBadRequest, "nick was nil")
		return
	}

	//判断K是否存在
	if _, exists := f.nick.Load(nick); exists {
		c.String(http.StatusBadRequest, "id was exists")
		return
	}

	ids, _ := strconv.Atoi(id)

	if _, exists := f.mid.Load(uint(ids)); !exists {
		c.String(http.StatusBadRequest, "id was not exists")
		return
	}

	f.nick.Store(nick, uint(ids))

	f.logln("save nicks:", id, nick)

	c.String(http.StatusOK, fmt.Sprintf("id: %s, nick: %s writed: %v", id, nick, f.saveNickName(filepath.Join(FullPath(), "music_nick.json"))))
}

// getGet 根据外号名获取外号的值
func (f *Finder) getGet(c *gin.Context) {
	m := make(map[string]uint)

	nick, _ := c.GetQuery("nick")

	max, _ := c.GetQuery("max")

	nickId, _ := strconv.Atoi(nick)

	maxCount, _ := strconv.Atoi(max)

	if maxCount == 0 {
		maxCount = 5
	}

	if nick == "" {
		c.String(http.StatusBadRequest, "nick was nil")
		return
	}

	if id, ok := f.nick.Load(nick); ok {
		if name, ok := f.mid.Load(id); ok {
			m[name.(string)] = id.(uint)
			c.JSON(http.StatusOK, m)
			return
		}
	}

	if nickId > 0 {
		if name, ok := f.mid.Load(uint(nickId)); ok {
			m[name.(string)] = uint(nickId)
			c.JSON(http.StatusOK, m)
			return
		}
	}

	if mid, ok := f.name.Load(nick); ok {
		m[nick] = mid.(uint)
		c.JSON(http.StatusOK, m)
		return
	}

	f.nick.Range(func(key, value any) bool {
		search := key.(string)

		if maxCount <= len(m) {
			return false
		}

		if strings.Contains(search, nick) {
			if name, ok := f.mid.Load(value.(uint)); ok {
				m[name.(string)] = value.(uint)
			}
		}

		return true
	})

	f.nick.Range(func(key, value any) bool {
		search := key.(string)

		if maxCount <= len(m) {
			return false
		}

		if strings.Contains(strings.ToLower(search), nick) {
			if name, ok := f.mid.Load(value.(uint)); ok {
				m[name.(string)] = value.(uint)
			}
		}

		return true
	})

	f.nick.Range(func(key, value any) bool {
		search := key.(string)

		if maxCount <= len(m) {
			return false
		}

		if strings.Contains(strings.ToLower(search), strings.ToLower(nick)) {
			if name, ok := f.mid.Load(value.(uint)); ok {
				m[name.(string)] = value.(uint)
			}
		}

		return true
	})

	if musics, ok := f.artist[nick]; ok {
		for _, music := range musics {
			if maxCount <= len(m) {
				break
			}
			m[music.Title] = music.MID
		}
	}

	if musics, ok := f.genre[nick]; ok {
		for _, music := range musics {
			if maxCount <= len(m) {
				break
			}
			m[music.Title] = music.MID
		}
	}

	/*
		for art, musics := range f.artist {
			if maxCount <= len(m) {
				break
			}

			if strings.Contains(art, nick) {
				for _, music := range musics {
					if maxCount <= len(m) {
						break
					}
					m[music.Title] = music.MID
				}
			}
		}

		for art, musics := range f.artist {
			if maxCount <= len(m) {
				break
			}

			if strings.Contains(strings.ToLower(art), nick) {
				for _, music := range musics {
					if maxCount <= len(m) {
						break
					}
					m[music.Title] = music.MID
				}
			}
		}

		for art, musics := range f.artist {
			if maxCount <= len(m) {
				break
			}

			if strings.Contains(strings.ToLower(art), strings.ToLower(nick)) {
				for _, music := range musics {
					if maxCount <= len(m) {
						break
					}
					m[music.Title] = music.MID
				}
			}
		}

		for gen, musics := range f.genre {
			if maxCount <= len(m) {
				break
			}

			if strings.Contains(gen, nick) {
				for _, music := range musics {
					if maxCount <= len(m) {
						break
					}
					m[music.Title] = music.MID
				}
			}
		}

		for gen, musics := range f.genre {
			if maxCount <= len(m) {
				break
			}

			if strings.Contains(strings.ToLower(gen), nick) {
				for _, music := range musics {
					if maxCount <= len(m) {
						break
					}
					m[music.Title] = music.MID
				}
			}
		}

		for gen, musics := range f.genre {
			if maxCount <= len(m) {
				break
			}

			if strings.Contains(strings.ToLower(gen), strings.ToLower(nick)) {
				for _, music := range musics {
					if maxCount <= len(m) {
						break
					}
					m[music.Title] = music.MID
				}
			}
		}
	*/

	f.name.Range(func(key, value any) bool {
		search := key.(string)

		if maxCount <= len(m) {
			return false
		}

		if strings.Contains(search, nick) {
			m[key.(string)] = value.(uint)
		}

		return true
	})

	f.name.Range(func(key, value any) bool {
		search := key.(string)

		if maxCount <= len(m) {
			return false
		}

		if strings.Contains(strings.ToLower(search), nick) {
			m[key.(string)] = value.(uint)
		}

		return true
	})

	f.name.Range(func(key, value any) bool {
		search := key.(string)

		if maxCount <= len(m) {
			return false
		}

		if strings.Contains(strings.ToLower(search), strings.ToLower(nick)) {
			m[key.(string)] = value.(uint)
		}

		return true
	})

	c.JSON(http.StatusOK, m)
}

// getDel 删外号
func (f *Finder) getDel(c *gin.Context) {
	nick, _ := c.GetQuery("nick")

	if nick == "" {
		c.String(http.StatusBadRequest, "nick was nil")
		return
	}

	f.logln("delete nicks:", nick)

	f.nick.Delete(nick)
	f.saveNickName(filepath.Join(FullPath(), "music_nick.json"))
	c.String(http.StatusOK, "")
}

// getNicks 外号列表
func (f *Finder) getNicks(c *gin.Context) {
	m := make(map[string]uint)
	f.nick.Range(func(key, value any) bool {
		m[key.(string)] = value.(uint)
		return true
	})
	c.JSON(http.StatusOK, m)
}

// getSongs 歌单
func (f *Finder) getSongs(c *gin.Context) {
	m := make(map[string]uint)
	f.name.Range(func(key, value any) bool {
		m[key.(string)] = value.(uint)
		return true
	})
	c.JSON(http.StatusOK, m)
}

// getSongs 歌单
func (f *Finder) getReload(c *gin.Context) {
	c.JSON(http.StatusOK, f.reload())
}

// getSDVXGet 搜歌
func (f *Finder) getSDVXGet(c *gin.Context) {
	// id找歌
	idMatch, isIdMatch := c.GetQuery("id")
	// 名称匹配曲目
	queryMatch, isQueryMatch := c.GetQuery("query")

	var result any
	var err error

	if isIdMatch {
		result, err = f.SDVXManager.Get(idMatch)
		if err != nil {
			result = nil
		}
	} else if isQueryMatch {
		resultList := make([]SDVXMusicInfo, 0)
		ids := f.SDVXManager.SimpleMatch(queryMatch)
		for _, id := range ids {
			info, err := f.SDVXManager.Get(id)
			if err != nil || info == nil {
				continue // 跳过无效的 `info`
			}
			resultList = append(resultList, *info)
		}
		result = resultList
	} else {
		result = f.SDVXManager.GetAll()
	}

	c.JSON(http.StatusOK, result)
}

// getSDVXReload 加载sdvx数据库和别名
func (f *Finder) getSDVXReload(c *gin.Context) {
	if e := f.sdvxLoadUni(); e != nil {
		c.JSON(http.StatusInternalServerError, "failure")
		return
	}

	c.JSON(http.StatusOK, "ok")
}

// getSDVXAliasList 获取别名列表
func (f *Finder) getSDVXAliasList(c *gin.Context) {
	var result any
	idMatch, isIdMatch := c.GetQuery("id")
	if isIdMatch {
		aliases, status, err := f.SDVXManager.GetAlias(idMatch)

		var msg string

		if err != nil {
			msg = err.Error()
		}

		result = map[string]any{
			"aliases": aliases,
			"status":  status,
			"msg":     msg,
		}
	} else {
		result = f.SDVXManager.GetAliases()
	}
	c.JSON(http.StatusOK, result)
}

// getSDVXMatchId 通过别名或者曲目获匹配到曲目id
func (f *Finder) getSDVXMatchId(c *gin.Context) {
	query, isQuery := c.GetQuery("query")

	result := map[string]any{
		"msg":      "",
		"status":   Success,
		"contents": nil,
	}

	if !isQuery {
		result["msg"] = "missing query parameters"
		result["status"] = MissingParameters
		c.JSON(http.StatusBadRequest, result)
		return
	}

	isNoCase, hasIsNoCase := c.GetQuery("isnocase")
	if !hasIsNoCase {
		isNoCase = "0"
	}

	isFuzzy, hasIsFuzzy := c.GetQuery("isfuzzy")
	if !hasIsFuzzy {
		isFuzzy = "0"
	}

	isAlias, hasIsAlias := c.GetQuery("isalias")
	if !hasIsAlias {
		isAlias = "0"
	}

	useNoCase := false
	if isNoCase != "0" {
		useNoCase = true
	}

	useFuzzy := false
	if isFuzzy != "0" {
		useFuzzy = true
	}

	if isAlias == "0" {
		result["contents"] = f.SDVXManager.Match(query, useNoCase, useFuzzy)
		// 曲目名称查找id
	} else {
		// 曲目别名查找id
		result["contents"] = f.SDVXManager.MatchAlias(query, useNoCase, useFuzzy)
	}
	c.JSON(http.StatusOK, result)
}

// getSDVXIdExist 判断id是否存在
func (f *Finder) getSDVXIdExist(c *gin.Context) {
	id, isId := c.GetQuery("id")

	result := map[string]any{
		"msg":    "",
		"status": Success,
		"exist":  false,
	}

	if !isId {
		result["msg"] = "missing id parameters"
		result["status"] = MissingParameters
		c.JSON(http.StatusBadRequest, result)
		return
	}

	exist, err := f.SDVXManager.Exist(id)
	result["exist"] = exist
	if err != nil {
		result["status"] = UnknownError
		result["msg"] = err.Error()
		c.JSON(http.StatusBadRequest, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// addSDVXAlias 添加别名
func (f *Finder) addSDVXAlias(c *gin.Context) {
	id, isId := c.GetQuery("id")
	alias, isAlias := c.GetQuery("alias")

	result := map[string]any{
		"msg":    "",
		"status": Success,
	}

	if !(isId || isAlias) {
		result["status"] = MissingParameters
		result["msg"] = "missing 'id' or 'alias' parameters"
		c.JSON(http.StatusBadRequest, result)
		return
	}

	alias = strings.TrimSpace(alias)

	if alias == "" {
		result["status"] = EmptyString
		result["msg"] = "alias cannot be an empty string"
		c.JSON(http.StatusBadRequest, result)
		return
	}

	status, err := f.SDVXManager.AddAlias(id, alias)
	result["status"] = status
	if err != nil {
		result["msg"] = err.Error()
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

// delSDVXAlias 删除别名
func (f *Finder) delSDVXAlias(c *gin.Context) {
	alias, isAlias := c.GetQuery("alias")
	result := map[string]any{
		"msg":    "",
		"status": Success,
	}

	if !isAlias {
		result["status"] = MissingParameters
		result["msg"] = "missing 'alias' parameters"
		c.JSON(http.StatusBadRequest, result)
		return
	}

	alias = strings.TrimSpace(alias)

	if alias == "" {
		result["status"] = EmptyString
		result["msg"] = "alias cannot be an empty string"
		c.JSON(http.StatusBadRequest, result)
		return
	}

	status, err := f.SDVXManager.DelAlias(alias)
	result["status"] = status
	if err != nil {
		result["msg"] = err.Error()
		c.JSON(http.StatusInternalServerError, result)
		return
	}

	c.JSON(http.StatusOK, result)
}
