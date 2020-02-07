package Panizza

import (
	"errors"
	"reflect"
)

type FilterImpl struct {
	FilterHandles []Handle
}

func (mw *FilterImpl) Use(s interface{}) {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumMethod(); i++ {
		name := t.Method(i).Name
		if name == "OnCreate" {
			v.Method(i).Call([]reflect.Value{})
			continue
		}
		handle, ok := v.Method(i).Interface().(func(*HandleContext))
		if !ok {
			panic(errors.New("Method must is func(*HandleContext)"))
		}
		mw.FilterHandles = append(mw.FilterHandles, RecoverHandle(handle))
	}
}

func (mw *FilterImpl) RunFilter(ctx *HandleContext) {
	for _, handle := range mw.FilterHandles {
		handle(ctx)
		if ctx.CloseHttp {
			return
		}
	}
}
