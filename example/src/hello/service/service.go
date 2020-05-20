package service

import (
	. "github.com/CP-Panizza/Panizza"
	"fmt"
)

type HelloService struct{

}

//在结构体初始化后调用，用于做构造函数使用
func (this *HelloService)Init(){
	fmt.Println("init run!")
}

func (this *HelloService)Hello(ctx *HandleContext){
	ctx.JSON_write(200, JSONObject{
		"data":"hello",
	})
}

