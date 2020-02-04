package Panizza

import (
	"errors"
	"reflect"
)

type FilterImpl struct {
	FilterHandles []Handle
}

func (mw *FilterImpl) Use(s interface{}) {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumMethod(); i++ {
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
