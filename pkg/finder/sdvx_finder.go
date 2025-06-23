package finder

import (
	"bytes"
	"encoding/json"
	l "finder/pkg/util/log"
	"fmt"
	"github.com/clbanning/mxj/v2"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type RadarInfo struct {
	Notes    uint8 `json:"notes"`     // 物量
	Peak     uint8 `json:"peak"`      // 爆发
	Tsumami  uint8 `json:"tsumami"`   // 旋钮
	Tricky   uint8 `json:"tricky"`    // 棘手
	HandTrip uint8 `json:"hand_trip"` // 出张
	OneHand  uint8 `json:"one_hand"`  // 片手
}

func NewRadarInfo() RadarInfo {
	return RadarInfo{
		Notes:    0,
		Peak:     0,
		Tsumami:  0,
		Tricky:   0,
		HandTrip: 0,
		OneHand:  0,
	}
}

type DifficultyInfo struct {
	Level       uint8     `json:"level"`        // 难度
	Illustrator string    `json:"illustrator"`  // 曲绘画师
	EffectedBy  string    `json:"effected_by"`  // *混响男孩*
	Price       int32     `json:"price"`        // price
	Limited     uint8     `json:"limited"`      // limited
	JacketPrint int32     `json:"jacket_print"` // jacket_print
	JacketMask  int32     `json:"jacket_mask"`  // jacket_mask
	MaxExscore  int32     `json:"max_exscore"`  // 最大得分
	Radar       RadarInfo `json:"radar"`        // 六维
}

type SDVXMusicInfo struct {
	Id int32 `json:"id"` // 曲目id
	// Label            string                    `json:"label"`             // 曲目标签
	TitleName        string                    `json:"title_name"`        // 曲名
	TitleYomigana    string                    `json:"title_yomigana"`    // 曲目发音
	ArtistName       string                    `json:"artist_name"`       // 曲师
	ArtistYomigana   string                    `json:"artist_yomigana"`   // 曲师发音
	Ascii            string                    `json:"ascii"`             // 不知道是做什么的x
	BPMMax           float32                   `json:"bpm_max"`           // 最大bpm
	BPMMin           float32                   `json:"bpm_min"`           // 最小bpm
	DistributionDate uint32                    `json:"distribution_date"` // 曲目发布日期
	Volume           uint16                    `json:"volume"`            // volume
	BGNo             uint16                    `json:"bg_no"`             // 背景id
	Genre            uint32                    `json:"genre"`             // 类型
	IsFixed          bool                      `json:"is_fixed"`          // 希腊奶
	Version          string                    `json:"version"`           // sdvx曲目更新版本
	DemoPri          int8                      `json:"demo_pri"`          // 希腊奶
	DiffVer4         string                    `json:"diff_ver4"`         // 第4难度追加版本
	Difficulties     map[string]DifficultyInfo `json:"difficulties"`      // 难度信息
	DifficultyList   []string                  `json:"difficulty_list"`   // 可选难度列表
}

type SDVXManager struct {
	SDVXMusicInfos map[int32]SDVXMusicInfo
	SDVXAliases    map[string][]string
	logger         *l.Log
	AliasesPath    string
	m              sync.RWMutex
}

// logf 打印日志(如果没有启用则打到控制台)
func (manager *SDVXManager) logf(format string, v ...interface{}) {
	if manager.logger != nil {
		manager.logger.Printf(format, v...)
	}
	log.Printf(format, v...)
}

// logln 打印日志(如果没有启用则打到控制台)
func (manager *SDVXManager) logln(v ...any) {
	if manager.logger != nil {
		manager.logger.Println(v...)
	}
	log.Println(v...)
}

// panic 崩溃
func (manager *SDVXManager) panic(v any) {
	manager.logf("%v", v)
	panic(v)
}

var SDVXVersionName = []string{"", "Booth", "Infinite Infection", "Gravity Wars", "Heavenly Haven", "Vivid Wave", "Exceed Gear"}

// LoadData 加载数据库
func (manager *SDVXManager) LoadData(DBPath string) error {
	manager.m.RLock()
	defer manager.m.RUnlock()
	manager.SDVXMusicInfos = make(map[int32]SDVXMusicInfo)
	// 打开文件
	file, err := os.Open(DBPath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			manager.panic(err)
		}
	}(file)

	shiftJISReader := transform.NewReader(file, japanese.ShiftJIS.NewDecoder())

	// var reader io.Reader = file
	var buf bytes.Buffer
	_, err = io.Copy(&buf, shiftJISReader)
	if err != nil {
		return err
	}

	// 获取文件内容并删除 XML 声明
	content := buf.Bytes()
	re := regexp.MustCompile(`(?i)^\s*<\?xml[^?>]+\?>\s*`)
	contentWithoutXMLDecl := re.ReplaceAll(content, []byte{})

	mv, err := mxj.NewMapXmlReader(bytes.NewReader(contentWithoutXMLDecl))
	if err != nil {
		return err
	}

	LevelMapper := map[string]string{
		"nov": "novice",
		"adv": "advanced",
		"exh": "exhaust",
		"inf": "infinite",
		"grv": "infinite",
		"hvn": "infinite",
		"vvd": "infinite",
		"xcd": "infinite",
		"mxm": "maximum",
	}

	musicList := mv["mdb"].(map[string]any)["music"].([]any)
	for _, music := range musicList {
		infoAll := music.(map[string]any)
		idx, _ := strconv.Atoi(infoAll["-id"].(string))
		id := int32(idx)
		var Info SDVXMusicInfo
		Info.Id = id

		musicInfo := infoAll["info"].(map[string]any)

		// 第四难度版本
		infVerText, _ := strconv.Atoi(musicInfo["inf_ver"].(map[string]any)["#text"].(string))
		infVer := uint8(infVerText)

		Info.DiffVer4 = SDVXVersionName[infVer]
		// Info.Label = musicInfo["label"].(string)
		Info.TitleName = musicInfo["title_name"].(string)
		Info.TitleYomigana = musicInfo["title_yomigana"].(string)
		Info.Ascii = musicInfo["ascii"].(string)
		Info.ArtistName = musicInfo["artist_name"].(string)
		Info.ArtistYomigana = musicInfo["artist_yomigana"].(string)

		musicVersion, _ := strconv.Atoi(musicInfo["version"].(map[string]any)["#text"].(string))
		volume, _ := strconv.Atoi(musicInfo["volume"].(map[string]any)["#text"].(string))
		isFixed, _ := strconv.Atoi(musicInfo["is_fixed"].(map[string]any)["#text"].(string))
		genre, _ := strconv.Atoi(musicInfo["genre"].(map[string]any)["#text"].(string))
		distributionDate, _ := strconv.Atoi(musicInfo["distribution_date"].(map[string]any)["#text"].(string))
		demoPri, _ := strconv.Atoi(musicInfo["demo_pri"].(map[string]any)["#text"].(string))
		bpmMin, _ := strconv.Atoi(musicInfo["bpm_min"].(map[string]any)["#text"].(string))
		bpmMax, _ := strconv.Atoi(musicInfo["bpm_max"].(map[string]any)["#text"].(string))
		bgNo, _ := strconv.Atoi(musicInfo["bg_no"].(map[string]any)["#text"].(string))

		Info.Version = SDVXVersionName[uint8(musicVersion)]
		Info.Volume = uint16(volume)
		Info.IsFixed = isFixed != 0
		Info.Genre = uint32(genre)
		Info.DistributionDate = uint32(distributionDate)
		Info.DemoPri = int8(demoPri)
		Info.BPMMax = float32(bpmMax) / 100.0
		Info.BPMMin = float32(bpmMin) / 100.0
		Info.BGNo = uint16(bgNo)

		// 难度列表
		difficultyList := make([]string, 0)
		difficulties := infoAll["difficulty"].(map[string]any)

		_, noviceExist := difficulties["novice"]
		_, advancedExist := difficulties["advanced"]
		_, exhaustExist := difficulties["exhaust"]
		_, infiniteExist := difficulties["infinite"]
		_, maximumExist := difficulties["maximum"]

		if noviceExist {
			difficultyList = append(difficultyList, "nov")
		}
		if advancedExist {
			difficultyList = append(difficultyList, "adv")
		}
		if exhaustExist {
			difficultyList = append(difficultyList, "exh")
		}
		if infiniteExist {
			if infVer == 2 {
				difficultyList = append(difficultyList, "inf")
			} else if infVer == 3 {
				difficultyList = append(difficultyList, "grv")
			} else if infVer == 4 {
				difficultyList = append(difficultyList, "hvn")
			} else if infVer == 5 {
				difficultyList = append(difficultyList, "vvd")
			} else if infVer == 6 {
				difficultyList = append(difficultyList, "xcd")
			}
		}
		if maximumExist {
			difficultyList = append(difficultyList, "mxm")
		}

		Info.Difficulties = make(map[string]DifficultyInfo)

		for _, levelKey := range difficultyList {
			levelKeyFull := LevelMapper[levelKey]
			var DifficultyInfos DifficultyInfo
			difficultyInfo := difficulties[levelKeyFull].(map[string]any)

			diffNum, _ := strconv.Atoi(difficultyInfo["difnum"].(map[string]any)["#text"].(string))
			jacketMask, _ := strconv.Atoi(difficultyInfo["jacket_mask"].(map[string]any)["#text"].(string))
			jacketPrint, _ := strconv.Atoi(difficultyInfo["jacket_print"].(map[string]any)["#text"].(string))
			limited, _ := strconv.Atoi(difficultyInfo["limited"].(map[string]any)["#text"].(string))
			price, _ := strconv.Atoi(difficultyInfo["price"].(map[string]any)["#text"].(string))

			DifficultyInfos.Level = uint8(diffNum)
			DifficultyInfos.JacketMask = int32(jacketMask)
			DifficultyInfos.JacketPrint = int32(jacketPrint)
			DifficultyInfos.Limited = uint8(limited)
			DifficultyInfos.Price = int32(price)
			DifficultyInfos.EffectedBy = difficultyInfo["effected_by"].(string)
			DifficultyInfos.Illustrator = difficultyInfo["illustrator"].(string)

			DifficultyInfos.MaxExscore = 0
			DifficultyInfos.Radar = NewRadarInfo()

			maxExscoreMap, extendExist := difficultyInfo["max_exscore"]
			if extendExist {
				maxExscore, _ := strconv.Atoi(maxExscoreMap.(map[string]any)["#text"].(string))
				DifficultyInfos.MaxExscore = int32(maxExscore)
				radar := difficultyInfo["radar"].(map[string]any)
				handTrip, _ := strconv.Atoi(radar["hand-trip"].(map[string]any)["#text"].(string))
				oneHand, _ := strconv.Atoi(radar["one-hand"].(map[string]any)["#text"].(string))
				notes, _ := strconv.Atoi(radar["notes"].(map[string]any)["#text"].(string))
				peak, _ := strconv.Atoi(radar["peak"].(map[string]any)["#text"].(string))
				tricky, _ := strconv.Atoi(radar["tricky"].(map[string]any)["#text"].(string))
				tsumami, _ := strconv.Atoi(radar["tsumami"].(map[string]any)["#text"].(string))

				DifficultyInfos.Radar.HandTrip = uint8(handTrip)
				DifficultyInfos.Radar.OneHand = uint8(oneHand)
				DifficultyInfos.Radar.Notes = uint8(notes)
				DifficultyInfos.Radar.Peak = uint8(peak)
				DifficultyInfos.Radar.Tricky = uint8(tricky)
				DifficultyInfos.Radar.Tsumami = uint8(tsumami)
			}

			Info.Difficulties[levelKey] = DifficultyInfos
		}

		Info.DifficultyList = difficultyList

		manager.SDVXMusicInfos[id] = Info
	}

	manager.logln("sdvx db loaded")
	return nil
}

// GetAll 获取全部曲目信息
func (manager *SDVXManager) GetAll() *map[int32]SDVXMusicInfo {
	return &manager.SDVXMusicInfos
}

// Get 通过ID获取曲目信息
func (manager *SDVXManager) Get(id any) (*SDVXMusicInfo, error) {
	var sid int32

	switch v := id.(type) {
	case int:
		sid = int32(v)
	case string:
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("string is not a number: %v", err)
		}
		sid = int32(val)
	case int32:
		sid = v
	default:
		return nil, fmt.Errorf("id type error")
	}

	musicInfo, exists := manager.SDVXMusicInfos[sid]
	if !exists {
		return nil, fmt.Errorf("music info with id %d not found", sid)
	}

	return &musicInfo, nil
}

// Exist 判断曲子是否存在
func (manager *SDVXManager) Exist(id any) (bool, error) {
	var sid int32

	switch v := id.(type) {
	case int:
		sid = int32(v)
	case string:
		val, err := strconv.Atoi(v)
		if err != nil {
			return false, fmt.Errorf("string is not a number: %v", err)
		}
		sid = int32(val)
	case int32:
		sid = v
	default:
		return false, fmt.Errorf("id type error")
	}

	_, exists := manager.SDVXMusicInfos[sid]
	return exists, nil
}

// Match 曲目匹配
// query 匹配的名称
// isNoCase 禁用大小写
// isFuzzy 模糊匹配
func (manager *SDVXManager) Match(query string, isNoCase bool, isFuzzy bool) []int32 {
	if isNoCase {
		query = strings.ToLower(query)
	}
	var matches []int32

	for id, value := range manager.SDVXMusicInfos {
		title := value.TitleName
		if isNoCase {
			title = strings.ToLower(title)
		}
		if isFuzzy {
			if strings.Contains(title, query) {
				matches = append(matches, id)
			}
		} else {
			if title == query {
				matches = append(matches, id)
			}
		}
	}

	return matches
}

// LoadAliases 加载别名
func (manager *SDVXManager) LoadAliases(aliasesPath string) error {
	manager.m.RLock()
	defer manager.m.RUnlock()
	file, err := os.Open(aliasesPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	manager.AliasesPath = aliasesPath

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			manager.panic(err)
		}
	}(file)

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	err = json.Unmarshal(fileContent, &manager.SDVXAliases)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	manager.logln("sdvx aliases loaded")
	return nil
}

const ( // 3 无法找到这个别名
	UnknownError       = iota - 1 // 未知错误
	Success                       // 0 成功
	AliasAlreadyExists            // 1 别名已经存在
	MusicIDNotExist               // 2 曲目id不存在
	NotFoundAlias
	MissingParameters
	EmptyString
)

// saveAliases 将 SDVXAliases 数据写入 JSON 文件
func (manager *SDVXManager) saveAliases() error {
	// 将 SDVXAliases 数据编码为 JSON 格式
	data, err := json.MarshalIndent(manager.SDVXAliases, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal JSON: %v", err)
	}

	// 写入文件
	err = os.WriteFile(manager.AliasesPath, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to write file: %v", err)
	}

	return nil
}

// AddAlias 添加别名
func (manager *SDVXManager) AddAlias(id any, newAlias string) (int, error) {
	manager.m.Lock()         // 获取写锁
	defer manager.m.Unlock() // 释放写锁

	newAlias = strings.TrimSpace(newAlias) // 去除首尾空格

	var sid string

	switch v := id.(type) {
	case int:
		sid = strconv.Itoa(v)
	case int32:
		sid = strconv.Itoa(int(v))
	case string:
		sid = v
	default:
		return UnknownError, fmt.Errorf("id type error")
	}

	exist, err := manager.Exist(sid)

	if !exist {
		return MusicIDNotExist, err
	}

	_, isNotEmpty := manager.SDVXAliases[sid]

	if !isNotEmpty {
		manager.SDVXAliases[sid] = make([]string, 0)
		err = manager.saveAliases()
		if err != nil {
			return UnknownError, err
		}
	}

	for _, aliasList := range manager.SDVXAliases {
		for _, alias := range aliasList {
			if newAlias == alias {
				return AliasAlreadyExists, fmt.Errorf("alias already exists")
			}
		}
	}

	// 修改数据
	manager.SDVXAliases[sid] = append(manager.SDVXAliases[sid], newAlias)

	err = manager.saveAliases()
	if err != nil {
		return UnknownError, err
	}
	return Success, nil
}

// DelAlias 删除别名
func (manager *SDVXManager) DelAlias(delAlias string) (int, error) {
	manager.m.Lock()         // 获取写锁
	defer manager.m.Unlock() // 释放写锁

	for sid, aliasList := range manager.SDVXAliases {
		for index, alias := range aliasList {
			if delAlias == alias {
				manager.SDVXAliases[sid] = append(aliasList[:index], aliasList[index+1:]...)
				goto final
			}
		}
	}

	return NotFoundAlias, fmt.Errorf("alias not found")

final:
	err := manager.saveAliases()
	if err != nil {
		return UnknownError, err
	}

	return Success, nil
}

// GetAlias 通过曲目id获取别名
func (manager *SDVXManager) GetAlias(id any) ([]string, int, error) {
	var sid string

	switch v := id.(type) {
	case int:
		sid = strconv.Itoa(v)
	case int32:
		sid = strconv.Itoa(int(v))
	case string:
		sid = v
	default:
		return nil, UnknownError, fmt.Errorf("id type error")
	}

	exist, err := manager.Exist(sid)

	if !exist {
		return nil, MusicIDNotExist, err
	}

	_, isNotEmpty := manager.SDVXAliases[sid]

	if !isNotEmpty {
		manager.SDVXAliases[sid] = make([]string, 0)
		manager.m.Lock()
		err = manager.saveAliases()
		manager.m.Unlock()
		if err != nil {
			return nil, UnknownError, err
		}
	}

	return manager.SDVXAliases[sid], Success, err
}

// GetAliases 获取全部别名信息
func (manager *SDVXManager) GetAliases() *map[string][]string {
	return &manager.SDVXAliases
}

// MatchAlias 曲目匹配
// query 匹配的别名
// isNoCase 禁用大小写
// isFuzzy 模糊匹配
func (manager *SDVXManager) MatchAlias(query string, isNoCase, isFuzzy bool) []struct {
	Id    int32
	Alias string
} {
	// 如果需要忽略大小写，先将 query 转换为小写
	if isNoCase {
		query = strings.ToLower(query)
	}

	// 存储匹配结果
	matches := make([]struct {
		Id    int32
		Alias string
	}, 0)

	// 遍历 Aliases map
	for id, aliases := range manager.SDVXAliases {
		sid, _ := strconv.Atoi(id)

		for _, alias := range aliases {
			originalAlias := alias // 保留原始 alias 以便返回
			// 如果需要忽略大小写，转换 alias 为小写
			if isNoCase {
				alias = strings.ToLower(alias)
			}

			// 模糊匹配或完全匹配
			if isFuzzy {
				if strings.Contains(alias, query) {
					matches = append(matches, struct {
						Id    int32
						Alias string
					}{Id: int32(sid), Alias: originalAlias})
				}
			} else {
				if alias == query {
					matches = append(matches, struct {
						Id    int32
						Alias string
					}{Id: int32(sid), Alias: originalAlias})
				}
			}
		}
	}

	return matches
}

// SimpleMatch 简易匹配曲目(整合别名匹配+曲名匹配)
func (manager *SDVXManager) SimpleMatch(query string) []int32 {
	emptyList := make([]int32, 0)
	ids := emptyList

	// 精确曲名获取
	ids = manager.Match(query, false, false)
	if len(ids) != 0 {
		return ids
	}
	ids = emptyList

	// 精确别名获取
	idAliases := manager.MatchAlias(query, false, false)
	for _, lists := range idAliases {
		ids = append(ids, lists.Id)
	}
	if len(ids) != 0 {
		return ids
	}
	ids = emptyList

	// 不区分大小写曲名获取
	ids = manager.Match(query, true, false)
	if len(ids) != 0 {
		return ids
	}
	ids = emptyList

	// 不区分大小写别名获取
	idAliases = manager.MatchAlias(query, true, false)
	for _, lists := range idAliases {
		ids = append(ids, lists.Id)
	}
	if len(ids) != 0 {
		return ids
	}
	ids = emptyList

	// 模糊曲名匹配
	ids = manager.Match(query, true, true)
	if len(ids) != 0 {
		return ids
	}
	ids = emptyList

	// 模糊别名匹配
	idAliases = manager.MatchAlias(query, true, true)
	for _, lists := range idAliases {
		ids = append(ids, lists.Id)
	}
	if len(ids) != 0 {
		return ids
	}
	ids = emptyList

	return ids
}
