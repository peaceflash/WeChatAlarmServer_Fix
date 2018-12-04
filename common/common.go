package common

import (
	"os"
	"os/exec"
	"path/filepath"
)

func GetPath() (dir, path string) {

	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		path = "\\"
	} else {
		path = "/"
	}

	osPath, _ := exec.LookPath(os.Args[0])
	//文件相对路径
	dir = filepath.Dir(osPath)

	return dir, path
}
