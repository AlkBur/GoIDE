package handlers

import (
	"github.com/AlkBur/GoIDE/conf"
	"github.com/AlkBur/GoIDE/log"
	"html/template"
	"net/http"
)

var tamplateIndex *template.Template

func init() {
	var err error
	tamplateIndex, err = template.ParseFiles("views/index.html")
	if err != nil {
		log.Error(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	if conf.IDE.Context+"/" != r.RequestURI {
		http.Redirect(w, r, conf.IDE.Context+"/", http.StatusFound)
		return
	}
	var err error
	if log.GetLevel() == log.LevelDebug {
		tamplateIndex, err = template.ParseFiles("views/index.html")
		if nil != err {
			log.Error(err)
			http.Error(w, err.Error(), 500)
			return
		}
	}
	model := make(map[string]interface{})

	err = tamplateIndex.Execute(w, model)
	if nil != err {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
