package Panizza

import (
	"errors"
	"log"
	"reflect"
	"strings"
)

//全局路由注册实例
var controllerInstence = controller{
	make(methodMap),
}

//注册controller，传入结构体指针，此函数不应该recover():此函数运行成功就能保证运用运行起来
func RegisterController(i interface{}) interface{} {
	controllerInstence.controller(i)
	return i
}

//路由绑定
type controller struct {
	methodTree methodMap
}

//controller开启
var IsOpenController bool = false

//注册controller：接收结构体指针，
func (mt *controller) controller(s interface{}) {
	v := reflect.ValueOf(s).Elem()
	t := reflect.TypeOf(s).Elem()
	pkgPath := v.Type().PkgPath()
	//用结构体名称作为父路由前缀
	structName := strings.ToLower(t.Name())
	groupNameFunc := reflect.ValueOf(s).MethodByName("GroupName")
	var groupName string = ""
	if groupNameFunc.IsValid() {
		result := groupNameFunc.Call([]reflect.Value{})
		groupName = result[0].Interface().(string)
	}


	if t.Kind().String() != "struct" {
		panic(errors.New("only accept struct"))
	}

loop:
	for i := 0; i < t.NumField(); i++ {
		handle, Convok := v.Field(i).Interface().(Handle)

		if handle == nil {
			panic(errors.New("controller   '" + structName + "'   method   '" + t.Field(i).Name + "'   need to impl! "))
		}

		if !Convok {
			panic(errors.New("convert handle error!"))
		}

		if t.Field(i).Tag == "" {
			continue loop
		}

		path, ok1 := t.Field(i).Tag.Lookup("path")

		method, ok2 := t.Field(i).Tag.Lookup("method")

		if !ok2 {
			panic(errors.New("must give Method"))
		}

		description, ok3 := t.Field(i).Tag.Lookup("description")

		if !ok3 {
			description = "no description!"
		}

		if groupName != "" {
			if !ok1 {
				path = "/" + groupName
			} else {
				path = "/" + groupName + path
			}
		} else if structName != "" {
			if !ok1 {
				path = "/" + structName
			} else {
				path = "/" + structName + path
			}
		}

		method = strings.ToUpper(method)

		mt.methodTree.addRouter(method, path, handle, t.Field(i).Name)

		list := strings.Split(path, "/")
		params := []string{}
		for _, v := range list {
			if strings.Contains(v, ":") {
				params = append(params, v)
			}
		}
		IsOpenController = true
		//注册api
		GlobalApi.Add(method, path, description, params, pkgPath+"."+structName+"."+t.Field(i).Name)
		log.Println("URLMapping:" + method + ": " + path + " at:" + pkgPath + "." + structName + "." + t.Field(i).Name)
	}
}
