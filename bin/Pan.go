package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	Controller = "`@Controller`"
	Service    = "`@Service`"
	Aspect     = "`@Aspect`"
	Filter     = "`@Filter`"
	Been       = "`@Been`"

	AnnController = `//@Controller`
	AnnService    = `//@Service`
	AnnAspect     = `//@Aspect`
	AnnFilter     = `//@Filter`
	AnnBeen       = `//@Been`
)

//注解结构体
type Annotation struct {
	filePkg    string //当前文件所在的包
	RowDataStr string //扫描到的最初始数据
	StructName string
	Type       string
	Pkg        string
}

func (a *Annotation) GetPkgString() string {
	return a.Pkg
}

func (a *Annotation) String() string {
	return "\n\t" + a.StructName + " \t" + a.Type + "\r"
}

const (
	all_reg  = `//@\w+\stype\s\w+\sstruct`
	been_reg = `//@Been\([\w\W]+?\)\sfunc?\s[\w\W]+?\(?\)`
)

var Anns = []Annotation{}
var goFiles = []string{}

var rootPath string
var src string
var conf string
var resourse string
var initPro string
var confFile string
var initProFile string
var mainFile string

var initProFileContent = `package initPro

import (
	."github.com/18788567655/Panizza"
)

var	App = New()

type Components struct {
	
}

func init() {
	RegisterComponents(new(Components))
}`

var mainFileContent = `package main

import(
	."./initPro"
	)

func main() {
	App.StartServer()
}`

var confFileContent string

const (
	FILE = 0
	DIR  = 1
)

type Obj struct {
	path    string
	kind    int
	content string
}

func newObj(path string, kind int, content string) Obj {
	return Obj{path, kind, content}
}

var ObjList = []Obj{}

func CreateProj(pro []Obj) {
	for _, v := range pro {
		if v.kind == FILE {
			file, err := os.Create(v.path)
			mustPanic(err)
			file.WriteString(v.content)
			file.Close()
		} else {
			mustPanic(os.Mkdir(v.path, os.ModeDir))
		}
		log.Println("build\t" + v.path + "\tseccess!")
	}
}

func mustPanic(err error) {
	if err != nil {
		panic(err)
	}
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

var projectName *string

func main() {
	complie := flag.Bool("c", false, "complie project!")
	projectName = flag.String("n", "awesomePanizza", "Your application name!")
	flag.Parse()
	if *complie == true {
		Complie()
	} else {
		NewProject()
	}
}

//新建一个新项目
func NewProject() {
	confFileContent = `PROJECT_PACKAGE=` + *projectName + `

PORT=8080

#FILE_SERVER=/resourses`
	fmt.Println(*projectName)
	dir, err := os.Getwd()
	mustPanic(err)
	dir = strings.Replace(dir, `\`, `/`, -1) + `/` + *projectName
	rootPath = dir
	ObjList = append(ObjList, newObj(rootPath, DIR, ""))
	src = dir + "/src"
	ObjList = append(ObjList, newObj(src, DIR, ""))
	conf = src + "/conf"
	ObjList = append(ObjList, newObj(conf, DIR, ""))
	resourse = src + "/resourses"
	ObjList = append(ObjList, newObj(resourse, DIR, ""))
	initPro = src + "/initPro"
	ObjList = append(ObjList, newObj(initPro, DIR, ""))
	confFile = conf + "/application.conf"
	ObjList = append(ObjList, newObj(confFile, FILE, confFileContent))
	initProFile = initPro + "/initPro.go"
	ObjList = append(ObjList, newObj(initProFile, FILE, initProFileContent))
	mainFile = src + "/" + *projectName + ".go"
	ObjList = append(ObjList, newObj(mainFile, FILE, mainFileContent))
	CreateProj(ObjList)
	log.Println("Project is build seccess!")
}

type BeenAnn struct {
	FilePkg string
	RowDataStr string
	FuncName string
	BeenName string
	Pkg      string
}

func (b *BeenAnn)String() string {
	return "\n\t" + "App.AddBeen(" + b.BeenName + "," + b.FuncName + ")" + "\r"
}

var beens = []BeenAnn{}

var targetFile string

var targetPkg string
//进行项目编译
func Complie() {

	filePath, err := FindFileAbsPath("initPro.go")
	if err != nil {
		panic(err)
	}
	targetFile = filePath

	targetPkg = targetFile[:strings.LastIndex(targetFile, "/")]
	fmt.Println(targetFile)
	//fmt.Println("targetPkg:",targetPkg)
	reg, err := regexp.Compile(all_reg)
	if err != nil {
		panic(err)
	}

	beenReg, err := regexp.Compile(been_reg)
	if err != nil {
		panic(err)
	}

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//fmt.Println("start path:", dir)

	SerachSrcFile(dir, &goFiles)

	for _, f := range goFiles {
		absPkgPath := f[:strings.LastIndex(f, "/")]
		data, err := ioutil.ReadFile(f)
		if err != nil {
			panic(err)
		}

		been := beenReg.FindAllString(string(data), len(data))
		if len(been) != 0 {
			for _, v := range been {
				b := BeenAnn{
					FilePkg: absPkgPath,
					RowDataStr: v,
				}
				beens = append(beens, b)
			}
		}

		strs := reg.FindAllString(string(data), len(data))
		if len(strs) != 0 {
			for _, val := range strs {
				ann := Annotation{
					filePkg:    absPkgPath,
					RowDataStr: val,
				}
				Anns = append(Anns, ann)
			}
		}
	}


	for k,v := range Anns {
		a , err := ParserAnns(v)
		if err != nil {
			panic(err)
		}
		Anns[k] = a
	}


	for k,v :=range beens {
		b, err := ParserBeens(v)
		if err != nil {
			panic(err)
		}
		beens[k] = b
	}

	pkgsString := ""

	componentsString := ""
	for _, v := range Anns {
		componentsString += v.String()
		if !strings.Contains(pkgsString, v.Pkg) && v.Pkg != `."."` {
			pkgsString += "\n\t" + v.Pkg + "\r"
		}
	}


	beensString := ""
	for _, v := range beens {
		beensString += v.String()
		if !strings.Contains(pkgsString, v.Pkg) && v.Pkg != `."."` {
			pkgsString += "\n\t" + v.Pkg + "\r"
		}
	}


	panizzaPath, err := FindFileAbsPath("Panizza.go")
	if err != nil {
		fmt.Println(err)
		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			panic(errors.New("Please set GOPATH!!!"))
		} else {
			file := []string{}
			FindFilesFromStartPath(goPath, "Panizza.go", &file)
			if len(file) == 0 {
				panic("can not find Panizza from GOPATH!!!")
			}
			panizzaPath = "." + `"` + file[0][strings.Index(file[0], "github"):strings.LastIndex(file[0], "/")] + `"`
		}
	} else {
		panizzaPath, err = GetRelPath(targetPkg, panizzaPath)
		if err != nil {
			panic(err)
		}
		panizzaPath = "." + `"` + panizzaPath[:strings.LastIndex(panizzaPath, "/")] + `"`
	}

	pkgsString = "\n\t" + panizzaPath + "\r" + pkgsString

	content := CreateInitProContent(pkgsString, componentsString, beensString)
	//create函数如果存在此文件就删除并从新创建
	f, err := os.Create(targetFile)
	if err != nil {
		panic(err)
	}
	f.WriteString(content)
	f.Close()

	log.Println("COMPILE SECCESS!")
}

//生成initPro文件内容
func CreateInitProContent(pkgs, comps, beens string) string {
	temp := `package initPro

import (
	`+ pkgs +`
)

var	App = New()

type Components struct {
	` + comps + `
}

func init() {
`+ beens +`
	RegisterComponents(new(Components))
}`
	return temp
}

//从字符串中获取been的名字
func GetBeensNameFromStr(str string)(string, error){

	regName, err := regexp.Compile(`\([\w\W]+?\)`)
	if err != nil {
		panic(err)
	}
	name := string(regName.Find([]byte(str)))
	if len(name) == 0 {
		return "", errors.New("can not getBeensName at :" + str)
	}
	name = name[1 : len(name) - 1]

	name = strings.TrimSpace(strings.Split(name, "=")[1])

	return name , nil
}

//获取return been的方法名
func GetBeenFuncNameFromStr(str string)(string, error){
	funcName := strings.TrimSpace(str[strings.Index(str, "func") + 4:])
	if len(funcName) == 0 {
		return "", errors.New("can not getFuncName at:" + str)
	}
	return funcName, nil
}

//解析Been注解
func ParserBeens(b BeenAnn) (BeenAnn, error){
	beenName, err := GetBeensNameFromStr(b.RowDataStr)
	if err != nil{
		return BeenAnn{}, err
	}
	b.BeenName = beenName

	beenFuncName, err := GetBeenFuncNameFromStr(b.RowDataStr)
	if err != nil{
		return BeenAnn{}, err
	}
	b.FuncName = beenFuncName

	pkg, err := GetRelPath(targetPkg, b.FilePkg)
	if err != nil{
		return BeenAnn{}, err
	}
	pkg = "." + `"` + pkg + `"`
	b.Pkg = pkg

	return b, nil
}

//字符串数组转化为字符串
func StringArrToString(s []string) string {
	str := ""
	for _, v := range s {
		str += v
	}
	return str
}

//遍历此目录下的所有.go文件
func SerachSrcFile(startPath string, files *[]string) {
	fs, err := ioutil.ReadDir(startPath)

	if err != nil {
		panic(err)
	}

	for _, file := range fs {
		if file.IsDir() {
			next := startPath + "/" + file.Name()
			SerachSrcFile(next, files)
		} else if strings.Contains(file.Name(), ".go") {
			right := startPath + "/" + file.Name()
			*files = append(*files, strings.Replace(right, "\\", "/", -1))
		}
	}
}

//通过字符串和路径生成Annotation
func ParserAnns(a Annotation) (Annotation, error) {
	structName, err := GetStructNameFromString(a.RowDataStr)
	if err != nil {
		return Annotation{}, err
	}
	a.StructName = structName
	pkg, err := GetRelPath(targetPkg, a.filePkg)
	if err != nil {
		return Annotation{}, err
	}
	pkg = "." + `"` + pkg + `"`
	a.Pkg = pkg

	switch {
	case strings.Contains(a.RowDataStr, AnnController):
		a.Type = Controller
		break
	case strings.Contains(a.RowDataStr, AnnService):
		a.Type = Service
		break
	case strings.Contains(a.RowDataStr, AnnAspect):
		a.Type = Aspect
		break
	case strings.Contains(a.RowDataStr, AnnFilter):
		a.Type = Filter
		break
	}

	return a, nil
}

//通过解析字符串获取结构体名称
func GetStructNameFromString(str string) (string, error) {
	reg, err := regexp.Compile(`type\s\w+\s`)
	if err != nil {
		return "", err
	}

	byte := reg.Find([]byte(str))

	if len(byte) == 0 {
		return "", errors.New("not find string!")
	}

	s := strings.Replace(string(byte), "type", "", 1)
	s = strings.TrimSpace(s)
	return s, nil
}

//解析两个绝对路径之间的相对路径
func GetRelPath(basePath, targetPath string) (string, error) {
	relPath, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		return "", err
	}
	return strings.Replace(relPath, `\`, `/`, len(relPath)), nil
}

//从当前程序运行的路径的子目录寻找某个文件的绝对路径
func FindFileAbsPath(FileName string) (string, error) {
	dir, _ := os.Getwd()
	var files []string
	FindFilesFromStartPath(dir, FileName, &files)
	if len(files) == 0 {
		return "", errors.New("can not find" + FileName + "file from project!")
	}
	return files[0], nil
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
