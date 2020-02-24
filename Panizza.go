package Panizza

import (
	"bufio"
	"container/list"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var panizzaInstance *Panizza

//全局配置实例
var AppConfig Configer

type Configer map[string]interface{}

//通过key获取配置文件中的配置
func (c *Configer) GetConfiger(key string) (interface{}, bool) {
	val, ok := (*c)[key]
	if !ok {
		return nil, false
	}
	return val, true
}

type Panizza struct {
	Configer
	controller
	fileServe
	Filter
	*Ioc
}

//实例化panizza，首先寻找项目中的配置文件application.conf并加载配置
func New() *Panizza {

	fmt.Println(miniLogo)

	propertiesPath := FindPropertiesFile("application.conf")

	AppConfig = ReadConfigFromProperties(propertiesPath)

	log.Println("AppConfig:", AppConfig)

	//把配置文件的配置全部放入ioc容器
	for key, val := range AppConfig {
		IocInstance.AddBeen(key, val)
	}

	//初始化静态资源服务
	fs := fileServe{}
	fs.Static()

	panizza := &Panizza{
		AppConfig,
		controllerInstence,
		fs,
		&FilterImpl{
			[]Handle{},
		},
		IocInstance,
	}

	panizzaInstance = panizza
	return panizza
}

//寻找application.conf配置文件
func FindPropertiesFile(confFileName string) string {
	dir, _ := os.Getwd()
	var files []string
	FindFilesFromStartPath(dir, confFileName, &files)
	if len(files) == 0 {
		panic(errors.New("can not find" + confFileName + "file from project!"))
	}
	return files[0]
}

func ReadConfigFromProperties(propertiesPath string) map[string]interface{} {
	file, err := os.Open(propertiesPath)
	if err != nil {
		panic(errors.New("Opne" + propertiesPath + "error!"))
	}
	defer file.Close()

	buf := bufio.NewReader(file)

	config := make(map[string]interface{})
	lineIndex := 0 //解析到的行号

	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		lineIndex++

		//判断是否为map
		if len(line) != 0 && strings.Contains(line, "=") && !strings.Contains(line, "#") && strings.Contains(line, "{") {
			list := strings.Split(line, "=")
			key := list[0]
			val := strings.Split(list[1][1:len(list[1])-1], ",")
			valsMap := make(map[string]string)
			for _, v := range val {
				result := strings.Split(v, ":")
				if len(result) != 2 {
					panic(errors.New("Parse map config err at line:" + string(lineIndex)))
				}
				valsMap[result[0]] = result[1]
			}
			config[key] = valsMap
			//fmt.Println("valsMap:",valsMap)
		} else if len(line) != 0 && strings.Contains(line, "=") && !strings.Contains(line, "#") && strings.Contains(line, "[") {
			//按照数组解析
			list := strings.Split(line, "=")
			key := list[0]
			val := strings.Split(list[1][1:len(list[1])-1], ",")
			valsSlice := []string{}
			for _, v := range val {
				valsSlice = append(valsSlice, v)
			}
			config[key] = valsSlice
			//fmt.Println("valsSlice:",valsSlice)
		} else if len(line) != 0 && strings.Contains(line, "=") && !strings.Contains(line, "#") {
			//按照key-val等式解析
			list := strings.Split(line, "=")
			config[list[0]] = list[1]
		}

		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(errors.New("read file err at line:" + string(lineIndex)))
			}
		}
	}

	return config
}

//开启panizza服务，如果没有配置PORT则默认监听8080端口
func (pz *Panizza) StartServer() {
	port, ok := AppConfig.GetConfiger("PORT")
	if !ok {
		port = "8080"
		pz.AddBeen("PORT", 8080)
	}

	if err := http.ListenAndServe(":"+(port).(string), pz); err != nil {
		panic(err)
	}
}

//全局注册过得url都会放在这里
var registUrl = []string{}

func (pz *Panizza) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//初始化上下文实例，对应到每一个请求都有一个上下文
	var ctx = NewHandleContext(&w, r)

	defer UseTimePrint(ctx)()

	if !IsOpenController && !pz.IsOpenFileService {
		w.Header().Add("Content-Type", "text/html;charset=utf-8")
		fmt.Fprint(w, PanizzaWelcomePage())
		return
	}

	pz.Filter.RunFilter(ctx)

	if ctx.CloseHttp {
		log.Println("中间件提前断开了请求！")
		return
	}

	//通过url取得get对应方法并执行
	if r.Method == "GET" {
		//静态资源访问支持
		if pz.IsOpenFileService {
			//fmt.Println("file server isOpen")
			file := pz.StaticPath + r.URL.Path
			file = strings.Replace(file, "\\", "/", -1)
			indexHtml := file + "index.html"

			if IsDir(file) && Existe(indexHtml) {
				http.ServeFile(w, r, indexHtml)
				log.Println("end file index.html!")
				return
			} else if Existe(file) && !IsDir(file) {
				//执行静态资源过滤器
				var flag bool = true
				for _, filter := range pz.FilterFunList {
					flag = filter(w, r, ctx.Data)
					if !flag {
						//log.Fatalln("文件过滤器过滤请求:", r.URL.Path)
						return
					}
				}
				log.Println("file request:", file)
				start := time.Now().UnixNano()
				http.ServeFile(w, r, file)
				end := time.Now().UnixNano()
				Milliseconds := float64((end - start) / 100000)
				log.Print("Panizza [\tend file server:\t", file+"\t", Milliseconds, "ms\t]\n")
				return
			}
		}
	}

	nodeList, has := pz.controller.methodTree[r.Method]

	if !has {
		handle404((*ctx.ResponseWriter), ctx.Request)
		return
	}

	routerHandle(ctx, pz, *nodeList)
	return
}

//未匹配404
func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprint(w, "page is not found!")
}

//遍历list寻path找对应的node
func getNodeByPath(handlerList list.List, path string) *node {
	for e := handlerList.Front(); e != nil; e = e.Next() {
		n := e.Value.(*node)
		if n.path == path {
			return n
		}
	}
	return nil
}

func routerHandle(ctx *HandleContext, pz *Panizza, handlerList list.List) {

	n := getNodeByPath(handlerList, ctx.Request.URL.Path)
	var aspect interface{}
	if n != nil {
		if n.HasAspect {
			aspect, _ = Aspecters.Load(n.HandleName)
			RecoverAspectHandle(n.HandlerFunc, aspect, n.HandleName)(ctx)
			return
		} else {
			RecoverHandle(n.HandlerFunc)(ctx)
			return
		}
	}

	for e := handlerList.Front(); e != nil; e = e.Next() {
		n1 := e.Value.(*node)
		if ctx.Check(n1.path, ctx.Request.URL.Path) {
			ctx.GetParams(n1.path, ctx.Request.URL.Path)
			if n1.HasAspect {
				aspect, _ = Aspecters.Load(n1.HandleName)
				RecoverAspectHandle(n1.HandlerFunc, aspect, n1.HandleName)(ctx)
				return
			} else {
				RecoverHandle(n1.HandlerFunc)(ctx)
				return
			}
		}
	}

	handle404(*ctx.ResponseWriter, ctx.Request)
}

//异常处理和切面执行
func RecoverAspectHandle(handle Handle, aspect interface{}, HandleName string) Handle {
	return func(ctx *HandleContext) {
		defer func() {
			if pr := recover(); pr != nil {
				fmt.Printf("panic recover: %v\r\n", pr)
				debug.PrintStack()
				((aspect).(Aspect)).AfterPanic(pr, ctx, HandleName)
			}
		}()
		((aspect).(Aspect)).Before(ctx, HandleName)
		handle(ctx)
		((aspect).(Aspect)).After(ctx, HandleName)
	}
}

//异常处理
func RecoverHandle(handle Handle) Handle {
	return func(ctx *HandleContext) {
		defer func() {
			if pr := recover(); pr != nil {
				fmt.Printf("panic recover: %v\r\n", pr)
				debug.PrintStack()
			}
		}()
		handle(ctx)
	}
}

//用于打印程序执行时间
func UseTimePrint(ctx *HandleContext) func() {
	start := time.Now().UnixNano()
	return func() {
		end := time.Now().UnixNano()
		Milliseconds := float64((end - start) / 100000)
		status := reflect.ValueOf(*ctx.ResponseWriter).Elem().FieldByName("status")
		code := status.Int()
		string := strconv.FormatInt(code, 10)
		log.Print("Panizza\t[\t", ctx.Request.Method, "\t", ctx.Request.URL.Path+"\t", string, "\t", Milliseconds, "ms\t", "]\n")
	}
}
