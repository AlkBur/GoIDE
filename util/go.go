package util

import (
	"path/filepath"
	"os"
	"runtime"
)

func GetExecutableInGOBIN(executable string) string {
	if IsWindows() {
		executable += ".exe"
	}

	gopaths := filepath.SplitList(os.Getenv("GOPATH"))

	for _, gopath := range gopaths {
		// $GOPATH/bin/$GOOS_$GOARCH/executable
		ret := filepath.Join(gopath, "bin",
			os.Getenv("GOOS") + "_" + os.Getenv("GOARCH"), executable)
		if IsExist(ret) {
			return ret
		}

		// $GOPATH/bin/{runtime.GOOS}_{runtime.GOARCH}/executable
		ret = filepath.Join(gopath,"bin",
			runtime.GOOS + "_" + runtime.GOARCH, executable)
		if IsExist(ret) {
			return ret
		}

		// $GOPATH/bin/executable
		ret = filepath.Join(gopath, "bin", executable)
		if IsExist(ret) {
			return ret
		}
	}

	// $GOBIN/executable
	gobin := os.Getenv("GOBIN")
	if "" != gobin {
		ret := filepath.Join(gobin, executable)
		if IsExist(ret) {
			return ret
		}
	}

	return "./" + executable
}
