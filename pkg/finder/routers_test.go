package finder

import (
	"encoding/json"
	"testing"
	"time"
)

func initTest() *Finder {
	return New(
		WithLog("log.test", 24, 1024*1024*1024),
		WithServer("0.0.0.0", 9999),
	)
}

func TestFunctions(t *testing.T) {
	finder := initTest()

	go func() {
		_ = finder.Start()
	}()

	time.Sleep(time.Hour)
}

func TestSDVXFinder(t *testing.T) {
	manager := SDVXManager{}

	err := manager.LoadData("E:\\datas\\programmingLanguage\\lang\\golang\\finder\\bin\\music_db.xml")
	err = manager.LoadAliases("E:\\datas\\programmingLanguage\\lang\\golang\\finder\\bin\\aliases.json")

	// manager.DelAlias("海神王")

	manager.AddAlias(1044, "海神王")

	// info, err := manager.Get(24)

	// dataJson, _ := json.MarshalIndent(info, "", "  ")
	// println(string(dataJson))

	// dataJson, _ = json.MarshalIndent(manager.Match("I", true, true), "", "  ")
	// println(string(dataJson))

	// println(manager.Exist(11361))

	// contents, _, _ := manager.GetAlias(1111)
	// dataJson, _ := json.MarshalIndent(contents, "", "  ")
	// println(string(dataJson))

	// dataJson, _ := json.MarshalIndent(manager.MatchAlias("船", true, true), "", "  ")
	// println(string(dataJson))

	dataJson, _ := json.MarshalIndent(manager.SimpleMatch("晕船"), "", "  ")
	println(string(dataJson))

	if err != nil {
		println("boom")
	}
}
