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

	username := httpSession.Values["username"].(string)
	if "playground" == username { // reserved user for Playground
		http.Redirect(w, r, conf.IDE.Context+"/login", http.StatusFound)

		return
	}

	httpSession.Options.MaxAge = conf.IDE.HTTPSessionMaxAge
	if "" != conf.IDE.Context {
		httpSession.Options.Path = conf.IDE.Context
	}
	httpSession.Save(r, w)

	user := conf.GetUser(username)
	if user == nil {
		log.Warn("Not found user [%s]", username)
		http.Redirect(w, r, conf.IDE.Context+"/login", http.StatusFound)
		return
	}

	ideSessions := session.GetUserSessions(user)

	//model := NewHtmlParam([]string{"ide.css", "index.css"}, []string{"lib/jquery-3.3.1.min.js", "ide.js"})
	model := NewHtmlParam([]string{"lib/webix.css", "lib/skins/compact.css", "ide.css"},
		[]string{"lib/webix.js", "ide.js"})
	model["i18n"] = i18n.GetAll(user.Locale)
	model["locale"] = user.Locale
	model["sid"] = session.GenId()

	log.Debug("User [%s] has [%d] sessions", user, len(ideSessions))

	if indexTemplate == nil || log.LevelDebug == log.GetLevel() {
		indexTemplate = util.LoadTemplate("views/ide.tmpl", "views/index.html", "views/menu.html")
	}

	if err := indexTemplate.ExecuteTemplate(w, "layout", model); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
}
