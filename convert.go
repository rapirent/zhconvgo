package zhconvgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const dictPath = "zhcdict.json"
const (
	ZHCN   = "zh-cn"
	ZHHK   = "zh-hk"
	ZHTW   = "zh-tw"
	ZHSG   = "zh-sg"
	ZHMY   = "zh-my"
	ZHMO   = "zh-mo"
	ZHHANT = "zh-hant"
	ZHHANS = "zh-hans"
	ZH     = "zh"
)

var Locales = map[string][]string{
	ZHCN:   {ZHCN, ZHHANS, ZHSG, ZH},
	ZHHK:   {ZHHK, ZHHANT, ZHTW, ZH},
	ZHTW:   {ZHTW, ZHHANT, ZHHK, ZH},
	ZHSG:   {ZHSG, ZHHANS, ZHCN, ZH},
	ZHMY:   {ZHMY, ZHSG, ZHHANS, ZHCN, ZH},
	ZHMO:   {ZHMO, ZHHK, ZHHANT, ZHTW, ZH},
	ZHHANT: {ZHHANT, ZHTW, ZHHK, ZH},
	ZHHANS: {ZHHANS, ZHCN, ZHSG, ZH},
	ZH:     {ZH},
}

var zhData map[string]interface{}
var zhCNDICT map[string]string
var zhHKDICT map[string]string
var zhTWDICT map[string]string
var zhSGDICT map[string]string
var zhHANTDICT map[string]string
var zhHANSDICT map[string]string
var pfsDICT map[string][]string

func getCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func loadDict(zhdata *map[string]interface{}) {
	pwd := getCurrentAbPath()
	fileName := fmt.Sprintf("%s/%s", pwd, dictPath)
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(errors.New(fmt.Sprintf("can`t load file from %s", dictPath)))
	}
	err = json.Unmarshal(content, &zhdata)
	if err != nil {
		panic(err)
	}
}

func isSupportLocale(locale string) bool {
	switch locale {
	case ZH, ZHCN, ZHHK, ZHTW, ZHMO, ZHHANT, ZHHANS, ZHMY, ZHSG:
		return true
	default:
		return false
	}
}

func getDict(locale string) map[string]string {
	result := make(map[string]string)
	switch locale {
	case ZHCN:
		if zhCNDICT == nil {
			zhCNDICT = genDict("zh2Hans", "zh2CN")
		}
		result = zhCNDICT
	case ZHTW:
		if zhTWDICT == nil {
			zhTWDICT = genDict("zh2Hant", "zh2TW")
		}
		result = zhTWDICT
	case ZHHK, ZHMO:
		if zhHKDICT == nil {
			zhHKDICT = genDict("zh2Hant", "zh2HK")
		}
		result = zhHKDICT
	case ZHSG, ZHMY:
		if zhSGDICT == nil {
			zhSGDICT = genDict("zh2Hans", "zh2SG")
		}
		result = zhSGDICT
	case ZHHANS:
		if zhHANSDICT == nil {
			zhHANSDICT = genDict("zh-hans")
		}
		result = zhHANSDICT
	case ZHHANT:
		if zhHANTDICT == nil {
			zhHANTDICT = genDict("zh2Hant")
		}
		result = zhHANTDICT
	}
	if pfsDICT == nil {
		pfsDICT = make(map[string][]string)
	}
	if _, ok := pfsDICT[locale]; !ok {
		pfsDICT[locale] = getPFSet(result)
	}
	return result
}

func getPFSet(data map[string]string) []string {
	pfset := make([]string, 0)
	for k := range data {
		runeK := []rune(k)
		for idx := range runeK {
			pfset = append(pfset, string(runeK[:idx+1]))
		}
	}
	return pfset
}

func genDict(locales ...string) map[string]string {
	if zhData == nil {
		loadDict(&zhData)
	}
	dict := make(map[string]string)
	if len(locales) == 0 {
		return dict
	}
	for _, locale := range locales {
		data, ok := zhData[locale]
		if !ok {
			continue
		}
		dataDict, ok := data.(map[string]interface{})
		if ok {
			for k, v := range dataDict {
				dict[k] = v.(string)
			}
		}
	}
	return dict
}

func Convert(s string, locale string) string {
	if locale == ZH || !isSupportLocale(locale) {
		return s
	}
	zhDict := getDict(locale)
	if zhDict == nil {
		return s
	}
	pfSet := pfsDICT[locale]
	newSet := make([]string, 0)
	ch := make([]string, 0)
	pos := 0
	sRune := []rune(s)
	N := len(sRune)
	for pos < N {
		i := pos
		frag := string(sRune[pos])
		maxPos := 0
		maxWord := ""
		for i < N && (containWord(frag, pfSet) || containWord(frag, newSet)) {
			matchWord, ok := zhDict[frag]
			if ok {
				maxPos = i
				maxWord = matchWord
			}
			i += 1
			frag = string(sRune[pos : i+1])
		}
		if maxWord == "" {
			maxWord = string(sRune[pos])
			pos += 1
		} else {
			pos = maxPos + 1
		}
		ch = append(ch, maxWord)
	}
	return strings.Join(ch, "")
}

func containWord(word string, data []string) bool {
	for _, v := range data {
		if v == word {
			return true
		}
	}
	return false
}
