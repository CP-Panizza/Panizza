package Panizza

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"log"
)

var (
	IocInstance        = NewIoc() //导入ioc包的时候实例化的全局单例
	ErrFactoryNotFound = errors.New("factory not found")
)

type factory func() (interface{}, error)

// 容器
type Ioc struct {
	sync.Mutex
	singletons map[string]interface{}
	factories  map[string]factory
}

// 容器实例化
func NewIoc() *Ioc {
	return &Ioc{
		singletons: make(map[string]interface{}),
		factories:  make(map[string]factory),
	}
}

// 注册单例对象
func (ioc *Ioc) AddBeen(name string, singleton interface{}) {
	ioc.Lock()
	ioc.singletons[name] = singleton
	ioc.Unlock()
}

// 获取单例对象
func (ioc *Ioc) GetBeen(name string) interface{} {
	return ioc.singletons[name]
}

// 获取实例对象
func (ioc *Ioc) GetPrototype(name string) (interface{}, error) {
	factory, ok := ioc.factories[name]
	if !ok {
		return nil, ErrFactoryNotFound
	}
	return factory()
}

// 设置实例对象工厂
func (ioc *Ioc) SetPrototypeFactory(name string, factory factory) {
	ioc.Lock()
	ioc.factories[name] = factory
	ioc.Unlock()
}

//传入结构体指针，向ioc容器注册组件
func RegisterComponents(been interface{}) {
	t := reflect.TypeOf(been).Elem()
	v := reflect.ValueOf(been).Elem()
	var FilterIndex = []int{}
	var ControllersIndex = []int{}
	var AspectIndex = []int{}
	for i := 0; i < v.NumField(); i++ {
		//fieldName := t.Field(i).Name
		fieldType := t.Field(i).Type
		isFilter := strings.Contains(string(t.Field(i).Tag), "@Filter")
		isController := strings.Contains(string(t.Field(i).Tag), "@Controller")
		isService := strings.Contains(string(t.Field(i).Tag), "@Service")
		isAspecter := strings.Contains(string(t.Field(i).Tag), "@Aspect")
		if isService {
			//1.首先注册service
			RegisterService(IocInstance.Inject(reflect.New(fieldType).Interface()))
		} else if isAspecter {
			AspectIndex = append(AspectIndex, i)
		} else if isController {
			ControllersIndex = append(ControllersIndex, i)
		} else if isFilter {
			FilterIndex = append(FilterIndex, i)
		}
	}

	//2.首先注册拦截器
	for _, v := range FilterIndex {
		fieldType := t.Field(v).Type
		instance := reflect.New(fieldType).Interface()
		panizzaInstance.Use(IocInstance.Inject(instance))
	}

	//3.注册controller
	for _, v := range ControllersIndex {
		fieldType := t.Field(v).Type
		instance := reflect.New(fieldType).Interface()
		RegisterController(IocInstance.Inject(instance))
	}

	//4.注册aop切面
	for _, v := range AspectIndex {
		fieldType := t.Field(v).Type
		instance := reflect.New(fieldType).Interface()
		RegistAspecter(IocInstance.Inject(instance))
	}

	defer func() {
		port := panizzaInstance.GetBeen("PORT")
		portString := ""
		if port == nil {
			portString = "8080"
		} else {
			portString = port.(string)
		}
		log.Println("Panizza started at:", portString)
	}()
}

//提供外部调用依赖注入,传入实例指针
func Inject(instance interface{}) {
	elemType := reflect.TypeOf(instance).Elem()
	ele := reflect.ValueOf(instance).Elem()
	for i := 0; i < elemType.NumField(); i++ { // 遍历字段
		fieldType := elemType.Field(i)
		tag := fieldType.Tag.Get("inject") // 获取tag
		diName := IocInstance.injectName(tag)
		if diName == "" {
			continue
		}
		var (
			diInstance interface{}
			err        error
		)
		if IocInstance.isSingleton(tag) {
			diInstance = IocInstance.GetBeen(diName)
		}
		if IocInstance.isPrototype(tag) {
			diInstance, err = IocInstance.GetPrototype(diName)
		}
		if err != nil {
			panic(err)
		}
		if diInstance == nil {
			panic(errors.New(diName + " dependency not found"))
		}
		ele.Field(i).Set(reflect.ValueOf(diInstance))
	}
}

// 注入依赖,此函数不应该recover():此函数运行成功就能保证运用运行起来
func (ioc *Ioc) Inject(instance interface{}) interface{} {
	elemType := reflect.TypeOf(instance).Elem()
	ele := reflect.ValueOf(instance).Elem()
	for i := 0; i < elemType.NumField(); i++ { // 遍历字段
		fieldType := elemType.Field(i)
		tag := fieldType.Tag.Get("inject") // 获取tag
		diName := ioc.injectName(tag)
		if diName == "" {
			continue
		}
		var (
			diInstance interface{}
			err        error
		)
		if ioc.isSingleton(tag) {
			diInstance = ioc.GetBeen(diName)
		}
		if ioc.isPrototype(tag) {
			diInstance, err = ioc.GetPrototype(diName)
		}
		if err != nil {
			panic(err)
		}
		if diInstance == nil {
			panic(errors.New(diName + " dependency not found"))
		}
		ele.Field(i).Set(reflect.ValueOf(diInstance))
	}

	return instance
}

// 获取需要注入的依赖名称
func (ioc *Ioc) injectName(tag string) string {
	tags := strings.Split(tag, ",")
	if len(tags) == 0 {
		return ""
	}
	return tags[0]
}

// 检测是否单例依赖
func (ioc *Ioc) isSingleton(tag string) bool {
	tags := strings.Split(tag, ",")
	for _, name := range tags {
		if name == "prototype" {
			return false
		}
	}
	return true
}

// 检测是否实例依赖
func (ioc *Ioc) isPrototype(tag string) bool {
	tags := strings.Split(tag, ",")
	for _, name := range tags {
		if name == "prototype" {
			return true
		}
	}
	return false
}

// 打印容器内部实例
func (ioc *Ioc) String() string {
	lines := make([]string, 0, len(ioc.singletons)+len(ioc.factories)+2)
	lines = append(lines, "singletons:")
	for name, item := range ioc.singletons {
		line := fmt.Sprintf("  %s: %x %s", name, &item, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	lines = append(lines, "factories:")
	for name, item := range ioc.factories {
		line := fmt.Sprintf("  %s: %x %s", name, &item, reflect.TypeOf(item).String())
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

//注册service，传入service结构体指针
func RegisterService(s interface{}) interface{} {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	for i := 0; i < v.NumMethod(); i++ {
		name := t.Method(i).Name
		if name == "OnCreate" {
			v.Method(i).Call([]reflect.Value{})
			continue
		}
		handle, ok := v.Method(i).Interface().(func(*HandleContext))
		if !ok {
			fmt.Println(name + " is not a Panizza.Handle!")
			continue
		}
		IocInstance.AddBeen(name, handle)
	}
	return s
}
