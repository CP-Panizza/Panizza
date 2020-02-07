package Panizza

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

//全局静态文件路径
var globleStaticPath string = ""

//文件服务
type fileServe struct {
	root              string //文件服务根路径
	FilterFunList     []FilterFun
	IsOpenFileService bool //文件服务是否开启
	StaticPath        string
}

//设置文件服务根路径
func (fs *fileServe) Static() {

	staticPath, ok := AppConfig.GetConfiger("FILE_SERVER")
	if ok {

		if IsDir(staticPath.(string)) {
			fmt.Println("指定了静态文件路径", staticPath)
			fs.StaticPath = staticPath.(string)
			globleStaticPath = staticPath.(string)
			fs.IsOpenFileService = true
			return
		}

		fs.root = staticPath.(string)
		projectPkg, ok1 := AppConfig.GetConfiger("PROJECT_PACKAGE")

		if !ok1 {
			return
		}

		dir := FindProjectPathByPKGName(projectPkg.(string))
		FindStaticPath(dir, fs.root, &fs.StaticPath)

		if !IsDir(fs.StaticPath) {
			panic(errors.New("current path have not :\t" + staticPath.(string)))
		}

		globleStaticPath = fs.StaticPath

		fs.FilterFunList = append(fs.FilterFunList, func(w http.ResponseWriter, r *http.Request, Data sync.Map) bool {
			log.Println("default file filter is run:\t", r.URL.Path, ", \tData:", Data)
			return true
		})
		fs.IsOpenFileService = true
	}
}

func FindProjectPathByPKGName(projectPKG string) string {
	found := false
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dir = strings.Replace(dir, "\\", "/", -1)

	var projectPath string
	for _, v := range strings.Split(dir, "/") {
		projectPath = projectPath + v + "/"
		if v == projectPKG {
			found = true
			break
		}
	}

	if !found {
		panic(errors.New("can not find project!"))
	}

	projectPath = projectPath[:len(projectPath)-1]

	return projectPath
}

//判断访问的文件路径是否存在
func Existe(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

//从给定的目录下寻找某个文件夹,结果为该文件夹的绝对路径
func FindStaticPath(ProjectPath string, rightPath string, outStr *string) {
	fs, _ := ioutil.ReadDir(ProjectPath)
	for _, file := range fs {
		if file.IsDir() && "/"+file.Name() == rightPath {
			right := ProjectPath + "/" + file.Name()
			fmt.Println("findPath:" + right)
			*outStr = right
			return
		} else {
			next := ProjectPath + "/" + file.Name()
			FindStaticPath(next, rightPath, outStr)
		}
	}

}

//从给定的目录下寻找某个文件，返回该文件绝对路径
func FindFilesFromStartPath(startPath string, rightFileName string, outStr *[]string) {
	fs, err := ioutil.ReadDir(startPath)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range fs {
		if file.IsDir() {

			next := startPath + "/" + file.Name()
			FindFilesFromStartPath(next, rightFileName, outStr)

		} else if file.Name() == rightFileName {
			right := startPath + "/" + file.Name()
			*outStr = append(*outStr, strings.Replace(right, "\\", "/", -1))
			return
		}
	}
}
