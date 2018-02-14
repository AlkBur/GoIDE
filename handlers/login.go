package handlers

import (
	"github.com/AlkBur/GoIDE/conf"
	"github.com/AlkBur/GoIDE/log"
	"github.com/AlkBur/GoIDE/session"
	"github.com/AlkBur/GoIDE/util"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
)

var (
	loginTemplate *template.Template
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		getLoginHandler(w, r)
	} else if r.Method == "POST" {
		postLoginHandler(w, r)
	}
}

func getLoginHandler(w http.ResponseWriter, r *http.Request) {
	if loginTemplate == nil || log.LevelDebug == log.GetLevel() {
		loginTemplate = util.LoadTemplate("views/ide.tmpl", "views/login.html")
	}

	model := NewHtmlParam([]string{"ide.css", "login.css"}, []string{"lib/jquery-3.3.1.min.js", "ide.js"})

	if err := loginTemplate.ExecuteTemplate(w, "layout", model); err != nil {
		log.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}
}

func postLoginHandler(w http.ResponseWriter, r *http.Request) {
	result := util.NewResult()
	defer result.Send(w, r)

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	user := conf.GetUser(username)
	if user == nil {
		result.Result = false
		return
	}
	if !user.CheckPassword(password) {
		result.Result = false
		return
	}
	// create a HTTP session
	httpSession, _ := session.HTTPSession.Get(r, "ide-session")
	httpSession.Values["username"] = username
	httpSession.Values["id"] = strconv.Itoa(rand.Int())
	httpSession.Options.MaxAge = conf.IDE.HTTPSessionMaxAge
	if "" != conf.IDE.Context {
		httpSession.Options.Path = conf.IDE.Context
	}
	httpSession.Save(r, w)

	log.Debug("Created a HTTP session [%s] for user [%s]", httpSession.Values["id"].(string), username)
}
