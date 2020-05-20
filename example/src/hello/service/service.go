package service

import ."github.com/CP-Panizza/Panizza"

type HelloService struct{

}

func (this *HelloService)Hello(ctx *HandleContext){
	ctx.JSON_write(200, JSONObject{
		"data":"hello",
	})
}

