package conf

import (
	"os"
	"path/filepath"
	"strings"
	"github.com/AlkBur/GoIDE/log"
	"encoding/json"
	"time"
)

type User struct {
	Name                  string
	Password              string
	Email                 string
	Workspace             string // the GOPATH of this user
	Locale                string
	Created               int64  // user create time in unix nano
	Updated               int64  // preference update time in unix nano
	Lived                 int64  // the latest session activity in unix nano
	Editor                *editor
	LatestSessionContent  *LatestSessionContent
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

func createUser(name, dir string) error {
	filename := filepath.Join(dir, name+".json")
	usr := &User{
		Name: name,
		Workspace: "${GOPATH}",
		Locale: IDE.Locale,
		Editor: &editor{
			FontFamily: "Consolas, 'Courier New', monospace",
			FontSize: "13px",
			LineHeight: "17px",
			Theme: "default",
			TabSize: 4,
		},
		LatestSessionContent: &LatestSessionContent{
			FileTree: make([]string, 0),
			Files: make([]string, 0),
		},
		Created: time.Now().UnixNano(),
		Updated: time.Now().UnixNano(),
	}

	jsonFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := json.Marshal(usr)
	if err != nil {
		return err
	}
	jsonFile.Write(jsonData)
	log.Info("Created a user file [%s]", filename)
	return nil
}