package handlers

import (
	"github.com/AlkBur/GoIDE/conf"
	"github.com/AlkBur/GoIDE/i18n"
	"github.com/AlkBur/GoIDE/log"
	"github.com/AlkBur/GoIDE/session"
	"github.com/AlkBur/GoIDE/util"
	"html/template"
	"net/http"
	"time"
)

var indexTemplate *template.Template

func NewHtmlParam(css, js []string) util.H {
	if css == nil {
		css = make([]string, 0)
	}
	if js == nil {
		js = make([]string, 0)
	}

	return util.H{"conf": conf.IDE, "i18n": i18n.GetAll(conf.IDE.Locale),
		"locale": conf.IDE.Locale, "ver": conf.IdeVersion, "year": time.Now().Year(), "title": conf.IdeName,
		"files": struct {
			CSS []string
			JS  []string
		}{CSS: css, JS: js}}
}

func Index(w http.ResponseWriter, r *http.Request) {
	if conf.IDE.Context+"/" != r.RequestURI {
		http.Redirect(w, r, conf.IDE.Context+"/", http.StatusFound)
		return
	}
	httpSession, _ := session.HTTPSession.Get(r, "ide-session")
	if httpSession.IsNew {
		http.Redirect(w, r, conf.IDE.Context+"/login", http.StatusFound)
		return
	}

	if indexTemplate == nil || log.LevelDebug == log.GetLevel() {
		indexTemplate = util.LoadTemplate("views/ide.tmpl", "views/index.html")
	}

	model := NewHtmlParam([]string{"ide.css", "login.css"}, []string{"lib/jquery-3.3.1.min.js", "ide.js"})

	if err := indexTemplate.ExecuteTemplate(w, "layout", model); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
