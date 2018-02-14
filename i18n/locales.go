package i18n

import (
	"encoding/json"
	"github.com/AlkBur/GoIDE/log"
	"io/ioutil"
	"os"
	"strings"
)

var Locales = map[string]locale{}

type locale struct {
	Name     string
	Langs    map[string]interface{}
	TimeZone string
}

func Load() {
	f, _ := os.Open("i18n")
	names, _ := f.Readdirnames(-1)
	f.Close()

	if len(Locales) == len(names)-1 {
		return
	}

	for _, name := range names {
		if !strings.HasSuffix(name, ".json") {
			continue
		}

		loc := name[:strings.LastIndex(name, ".")]
		load(loc)
	}
}

func Get(locale, key string) interface{} {
	l := GetAll(locale)
	return l[key]
}

func GetAll(locale string) map[string]interface{} {
	return Locales[locale].Langs
}

func load(localeStr string) {
	bytes, err := ioutil.ReadFile("i18n/" + localeStr + ".json")
	if nil != err {
		log.Error(err)
		os.Exit(-1)
	}

	l := locale{Name: localeStr}

	err = json.Unmarshal(bytes, &l.Langs)
	if nil != err {
		log.Error(err)
		os.Exit(-1)
	}
	Locales[localeStr] = l
}
