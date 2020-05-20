package Panizza

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"container/list"
)

//处理函数类型预定义
type Handle func(ctx *HandleContext)

func NewHandleContext(w *http.ResponseWriter, r *http.Request) *HandleContext {
	return &HandleContext{
		w,
		r,
		ParamParser{
			make(map[string]string),
			[]string{},
			[]string{},
		},
		sync.Map{},
		false,
	}
}

type HandleContext struct {
	ResponseWriter *http.ResponseWriter
	Request        *http.Request
	ParamParser
	Data      sync.Map
	CloseHttp bool
}

//终止当前http请求
func (ctx *HandleContext) Abort(code int) {
	(*ctx.ResponseWriter).WriteHeader(code)
	ctx.CloseHttp = true
}

//传入结构体指针，把请求参数绑定到结构体内部
func (ctx *HandleContext) Bind(i interface{}) error {

	parameResult := ctx.ParamMap
	if len(parameResult) != 0 {
		return bindVal2Struct(i, ctx.ParamParser)
	} else {
		//非url参数绑定到结构体
		fmt.Println("非url参数绑定到结构体")
		err := ctx.Request.ParseForm()
		if err != nil {
			return err
		}
		content_type := ctx.Request.Header.Get("Content-Type")
		if content_type == "application/json" {
			body, err := ioutil.ReadAll(ctx.Request.Body)
			if err != nil {
				return err
			}

			if len(body) == 0 {
				return errors.New("body len is zero")
			}
			valMap := map[string]interface{}{}
			err = json.Unmarshal(body, &valMap)
			if err != nil {
				return err
			}
			pp := ParamParser{}
			val := map[string]string{}
			for k, v := range valMap {
				switch reflect.ValueOf(valMap[k]).Type().String() {
				case "float64":
					val[k] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
					break
				case "string":
					val[k] = v.(string)
					break
				case "bool":
					val[k] = strconv.FormatBool(v.(bool))
				}
			}
			pp.ParamMap = val
			return bindVal2Struct(i, pp)
		} else if strings.HasPrefix(content_type, "multipart/form-data") {
			//设置存储容量
			mr, err := ctx.Request.MultipartReader()
			if err != nil {
				return err
			}
			form, err := mr.ReadForm(512)
			if err != nil {
				return err
			}

			pp := ParamParser{}
			valMap := map[string]string{}

			for k, v := range form.Value {
				valMap[k] = v[0]
			}

			pp.ParamMap = valMap
			return bindVal2Struct(i, pp)
		} else if len(ctx.Request.Form) != 0 {
			//从url参数里面绑定到结构体
			valMap := make(map[string]string)
			pp := ParamParser{}
			for k, v := range ctx.Request.Form {
				valMap[k] = v[0]
			}
			pp.ParamMap = valMap
			fmt.Println("form:", valMap)
			return bindVal2Struct(i, pp)
		}
	}
	return nil
}

//传入结构体实例和paramParser对象，把paramParser对象内的数据绑定到结构体中
func bindVal2Struct(i interface{}, paramParser ParamParser) error {

	v := reflect.ValueOf(i).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("json")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}

		name = strings.Split(name, ",")[0]
		//fmt.Println(name)
		typeof := v.Field(i).Type().String()
		//fmt.Println(typeof)
		switch typeof {
		case "int":
			data, err := paramParser.ParameInt(name)
			if err != nil {
				panic(err)
			}
			v.Field(i).Set(reflect.ValueOf(data))
			return nil
		case "int64":
			data, err := paramParser.ParameInt64(name)
			if err != nil {
				panic(err)
			}
			v.Field(i).Set(reflect.ValueOf(data))
			return nil
		case "string":
			data := paramParser.Parame(name)
			v.Field(i).Set(reflect.ValueOf(data))
			return nil
		case "float32":
			data, err := paramParser.ParameFloat32(name)
			if err != nil {
				panic(err)
			}
			v.Field(i).Set(reflect.ValueOf(data))
			return nil
		case "float64":
			data, err := paramParser.ParameFloat64(name)
			if err != nil {
				panic(err)
			}
			v.Field(i).Set(reflect.ValueOf(data))
			return nil
		case "time.Time":

			t := paramParser.Parame(name)

			var formatTime time.Time
			var err error

			if strings.Contains(t, ":") {
				formatTime, err = time.Parse("2006-01-02 15:04:05", t)
				if err != nil {
					log.Println(err)
				}
			} else {
				formatTime, err = time.Parse("2006-01-02", t)
				if err != nil {
					log.Println(err)
				}
			}

			v.Field(i).Set(reflect.ValueOf(formatTime))
			return nil
		case "bool":
			data := paramParser.ParameBool(name)
			v.Field(i).Set(reflect.ValueOf(data))
			return nil
		default:
			return errors.New(fmt.Sprintf("not soupport :", typeof))
		}
	}
	return nil
}

//---------------------------------------------------------------------------------------------
//传入map、array、slice、struct
func (ctx *HandleContext) JSON_write(code int, e interface{}) {
	if e == nil {
		return
	}

	(*ctx.ResponseWriter).Header().Add("Content-Type", "application/json; charset=utf-8")
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	(*ctx.ResponseWriter).WriteHeader(code)
	(*ctx.ResponseWriter).Write(data)
}

//向客户端写入一个字符串
func (ctx *HandleContext) Str_write(code int, msg string) {
	(*ctx.ResponseWriter).Header().Add("Content-Type", "text/plain; charset=utf-8")
	(*ctx.ResponseWriter).WriteHeader(code)
	(*ctx.ResponseWriter).Write([]byte(msg))
}

//---------------------------------------------------------------------------------------------
//静态资源访问的过滤器函数
type FilterFun func(w http.ResponseWriter, r *http.Request, Data sync.Map) bool

//---------------------------------------------------------------------------------------------
type MDContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Data           interface{}
}

//---------------------------------------------------------------------------------------------
//Filter interface
type Filter interface {
	Use(i interface{})
	RunFilter(ctx *HandleContext)
}

//---------------------------------------------------------------------------------------------
type ModelMap map[string]interface{}

func (M *ModelMap) AddAttribute(key string, val interface{}) {
	(*M)[key] = val
}

type JSONArray []interface{}

func (ja *JSONArray) Add(obj interface{}) {
	(*ja) = append((*ja), obj)
}

type JSONObject map[string]interface{}

func (jo *JSONObject) Put(key string, val interface{}) {
	(*jo)[key] = val
}

//---------------------------------------------------------------------------------------------

type node struct {
	path         string
	HandlerFunc  Handle
	HandleName   string
	HasAspect    bool `Description:"此节点是否有切面的标记"`
}

//---------------------------------------------------------------------------------------------

type methodMap map[string]*list.List

//头插法插入路由节点
func (m *methodMap) addRouter(method string, path string, handle Handle, HandleName string) {
	var Mux = sync.Mutex{}
	Mux.Lock()
	defer Mux.Unlock()
	val, ok := (*m)[method]

	if !ok {
		root := new(node)
		root.path = path
		root.HandlerFunc = handle
		root.HandleName = HandleName
		root.HasAspect = false
		listHead := list.New()
		listHead.PushBack(root)
		(*m)[method] = listHead
	} else {
		for e := val.Front(); e != nil; e = e.Next() {
			n := e.Value.(*node)
			if n.path == path {
				panic(errors.New("Path:\t" + path + "\t is areadly regist!"))
			}
		}

		next := new(node)
		next.HandlerFunc = handle
		next.HandleName = HandleName
		next.HasAspect = false
		next.path = path
		val.PushBack(next)
	}
}
