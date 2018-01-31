package conf

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"github.com/AlkBur/GoIDE/log"
	"encoding/json"
	"github.com/AlkBur/GoIDE/util"
	"strings"
	"os/exec"
	"github.com/AlkBur/GoIDE/event"
)


var (
	// IDE configurations.
	IDE *conf
	// configurations of users.
	Users []*User
)

// Configuration.
type conf struct {
	IP                    string // server ip, ${ip}
	Port                  int    // server port
	LogLevel              string // logging level: trace/debug/info/warn/error
	HTTPSessionMaxAge     int    // HTTP session max age (in seciond)
	WD                    string // current working direcitory, ${pwd}
	Locale                string // default locale
	Playground            string // playground directory
	UsersWorkspaces       string // users' workspaces directory (admin defaults to ${GOPATH}, others using this)
}

func Load(confPath, confIP string, confPort int, confLogLevel string) {
	initIDE(confPath, confIP, confPort, confLogLevel)
	initUsers()
}

func initIDE(confPath, confIP string, confPort int, confLogLevel string) {
	IDE = &conf{
		IP: confIP,
		Port: confPort,
		LogLevel: confLogLevel,
		HTTPSessionMaxAge: 86400,
		Locale: "en_US",
		WD: "${pwd}",
		Playground: "${WD}/workspaces/playground",
		UsersWorkspaces: "${WD}/workspaces",
	}
	if !util.IsExist(confPath) {
		err := createDefaultConfig(confPath)
		if nil != err {
			log.Error(err)

			os.Exit(-1)
		}
	}
	bytes, err := ioutil.ReadFile(confPath)
	if nil != err {
		log.Error(err)

		os.Exit(-1)
	}

	err = json.Unmarshal(bytes, IDE)
	if err != nil {
		log.Error("Parses [IDE.json] error: ", err)

		os.Exit(-1)
	}

	// Logging Level
	switch IDE.LogLevel {
	case "error":
		log.SetLevel(log.LevelError)
	case "warm":
		log.SetLevel(log.LevelWarn)
	case "debug":
		log.SetLevel(log.LevelDebug)
	case "off":
		log.SetLevel(log.LevelOff)
	default:
		log.SetLevel(log.LevelInfo)
	}

	log.Debug("Conf: \n" + string(bytes))

	// User Home
	home, err := util.HomeDir()
	if nil != err {
		log.Error("Can't get user's home, please report this issue to developer", err)

		os.Exit(-1)
	}
	log.Debug("${user.home} [%s]", home)

	// Working Directory
	pwd := util.Pwd()

	// Config working directory
	IDE.WD = strings.Replace(IDE.WD, "${pwd}",  pwd, 1)
	IDE.WD = strings.Replace(IDE.WD, "${home}", home, 1)
	log.Debug("${pwd} [%s]", IDE.WD)

	// Playground Directory
	IDE.Playground = strings.Replace(IDE.Playground, "${WD}", IDE.WD, 1)
	IDE.Playground = strings.Replace(IDE.Playground, "${home}", home, 1)
	IDE.Playground = strings.Replace(IDE.Playground, "${pwd}", pwd, 1)
	IDE.Playground = strings.Replace(IDE.Playground, "${UsersWorkspaces}", IDE.UsersWorkspaces, 1)
	IDE.Playground = filepath.Clean(IDE.Playground)
	log.Debug("${Playground} [%s]", IDE.Playground)
	if !util.IsExist(IDE.Playground) {
		createDir(IDE.Playground)
	}

	// Users' workspaces Directory
	IDE.UsersWorkspaces = strings.Replace(IDE.UsersWorkspaces, "${WD}", IDE.WD, 1)
	IDE.UsersWorkspaces = strings.Replace(IDE.UsersWorkspaces, "${home}", home, 1)
	IDE.UsersWorkspaces = strings.Replace(IDE.UsersWorkspaces, "${pwd}", pwd, 1)
	IDE.UsersWorkspaces = strings.Replace(IDE.UsersWorkspaces, "${Playground}", IDE.Playground, 1)
	IDE.UsersWorkspaces = filepath.Clean(IDE.UsersWorkspaces)
	log.Debug("${UsersWorkspaces} [%s]", IDE.UsersWorkspaces)
	if !util.IsExist(IDE.UsersWorkspaces) {
		createDir(IDE.UsersWorkspaces)
	}

	time := strconv.FormatInt(time.Now().UnixNano(), 10)
	log.Debug("${time} [%s]", time)
}

func initUsers() {
	usersPath := filepath.Join(IDE.UsersWorkspaces, "users")
	if !util.IsExist(usersPath) {
		createDir(usersPath)
	}
	f, err := os.Open(usersPath)
	if nil != err {
		log.Error(err)

		os.Exit(-1)
	}

	names, err := f.Readdirnames(-1)
	if nil != err {
		log.Error(err)

		os.Exit(-1)
	}
	f.Close()

	if len(names) == 0 {
		err = createUser("admin", usersPath)
		if nil != err {
			log.Error(err)
			os.Exit(-1)
		}
	}

	for _, name := range names {
		if strings.HasPrefix(name, ".") { // hiden files that not be created by Wide
			continue
		}

		if ".json" != filepath.Ext(name) { // such as backup (*.json~) not be created by Wide
			continue
		}

		user := &User{}

		bytes, _ := ioutil.ReadFile(filepath.Join(usersPath, name))

		err := json.Unmarshal(bytes, user)
		if err != nil {
			log.Error("Parses [%s] error: %v, skip loading this user", name, err)

			continue
		}
		//Add user
		Users = append(Users, user)
	}

	initWorkspaceDirs()
}

// Creates directories if not found on path of workspace.
func initWorkspaceDirs() {
	paths := []string{}

	for _, user := range Users {
		paths = append(paths, filepath.SplitList(user.WorkspacePath())...)
	}

	for _, path := range paths {
		CreateWorkspaceDir(path)
	}
}

func CreateWorkspaceDir(path string) {
	createDir(path)
	createDir(filepath.Join(path, "src"))
	createDir(filepath.Join(path, "pkg"))
	createDir(filepath.Join(path, "bin"))
}

func createDir(path string) {
	if !util.IsExist(path) {
		if err := os.MkdirAll(path, 0775); nil != err {
			log.Error(err)
			os.Exit(-1)
		}
		log.Info("Created a dir [%s]", path)
	}
}

func createDefaultConfig(filename string) error  {
	jsonFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	jsonData, err := json.Marshal(IDE)
	if err != nil {
		return err
	}
	jsonFile.Write(jsonData)
	log.Info("Created a config file [%s]", filename)
	return nil
}

func CheckEnv() {
	cmd := exec.Command("go", "version")
	buf, err := cmd.CombinedOutput()
	if nil != err {
		log.Error("Not found 'go' command, please make sure Go has been installed correctly")
		os.Exit(-1)
	}
	log.Info(string(buf))

	if "" == os.Getenv("GOPATH") {
		log.Error("Not found $GOPATH, please configure it before running Wide")
		os.Exit(-1)
	}

	gocode := util.GetExecutableInGOBIN("gocode")
	cmd = exec.Command(gocode)
	_, err = cmd.Output()
	if nil != err {
		event.EventQueue <- &event.Event{Code: event.EvtCodeGocodeNotFound}
		log.Warn("Not found gocode [%s], please install it with this command: go get github.com/nsf/gocode", gocode)
	}

	//LiteIDE Golang Tools
	ideStub := util.GetExecutableInGOBIN("gotools")
	cmd = exec.Command(ideStub, "version")
	_, err = cmd.Output()
	if nil != err {
		event.EventQueue <- &event.Event{Code: event.EvtCodeIDEStubNotFound}
		log.Warn("Not found gotools [%s], please install it with this command: go get github.com/visualfc/gotools", ideStub)
	}
}