package conf

import (
	"encoding/json"
	"github.com/AlkBur/GoIDE/log"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type User struct {
	Name                 string
	Password             string
	Email                string
	Workspace            string // the GOPATH of this user
	Locale               string
	Created              int64 // user create time in unix nano
	Updated              int64 // preference update time in unix nano
	Lived                int64 // the latest session activity in unix nano
	Editor               *editor
	LatestSessionContent *LatestSessionContent
}

type LatestSessionContent struct {
	FileTree    []string `json:"fileTree"`    // paths of expanding nodes of file tree
	Files       []string `json:"files"`       // paths of files of opening editor tabs
	CurrentFile string   `json:"currentFile"` // path of file of the current focused editor tab
}

type editor struct {
	FontFamily string
	FontSize   string
	LineHeight string
	Theme      string
	TabSize    int
}

func (u *User) WorkspacePath() string {
	w := strings.Replace(u.Workspace, "{WD}", IDE.WD, 1)
	w = strings.Replace(w, "${GOPATH}", os.Getenv("GOPATH"), 1)

	return filepath.FromSlash(w)
}

func (u *User) Save() bool {
	usersPath := filepath.Join(IDE.UsersWorkspaces, "users")
	filename := filepath.Join(usersPath, u.Name+".json")

	bytes, err := json.MarshalIndent(u, "", "    ")

	if nil != err {
		log.Error(err)

		return false
	}

	if "" == string(bytes) {
		log.Error("Truncated user [" + u.Name + "]")

		return false
	}

	if err = ioutil.WriteFile(filename, bytes, 0644); nil != err {
		log.Error(err)

		return false
	}

	return true
}

func NewUser(username, password, email, workspace string) *User {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	now := time.Now().UnixNano()

	return &User{Name: username, Password: string(hash), Email: email, Workspace: workspace,
		Locale:  IDE.Locale,
		Created: now, Updated: now, Lived: now,
		Editor: &editor{
			FontFamily: "Consolas, 'Courier New', monospace",
			FontSize:   "13px", LineHeight: "17px",
			Theme: "default", TabSize: 4,
		},
		LatestSessionContent: &LatestSessionContent{
			FileTree: make([]string, 0),
			Files:    make([]string, 0),
		},
	}
}

func (u *User) CheckPassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}
	return true
}

func GetUser(username string) *User {
	if "playground" == username { // reserved user for Playground
		// mock it
		return NewUser("playground", "", "", "")
	}

	for _, user := range Users {
		if user.Name == username {
			return user
		}
	}
	return nil
}
