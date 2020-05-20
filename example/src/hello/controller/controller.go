package controller

import ."github.com/CP-Panizza/Panizza"

type HelloController struct{
	Hello Handle `method:"GET" path:"" inject:"Hello"`
	Say Handle `method:"GET" path:"/say" inject:"Hello"`
}

func (this *HelloController)GroupName()string {
	return "hello"
}
