## PANIZZA
## web  https://cp-panizza.github.io/panizza/
Quickly to build a restapi web application!!!

## QUICK START

Install panizza

``` cmd
go get github.com/CP-Panizza/panizza
```

Add panizza to environment variable path
```cmd
Path=${GOPATH}\src\github.com\CP-Panizza\panizza\bin
```
Open cmd
```cmd
Pan -n "your project name like `my_first_pro`"
```
And then, the cli created a project named my_first_pro on the current path,like this:
```text
/my_first_pro
        |--src
            |--/conf
            |     |__application.conf 
            |
            |--/initPro
            |        |__initPro.go
            |
            |--/resource
            |
            |__my_first_pro.go
```
Finally, open browser enter http://127.0.0.1:8080, you will see the welcome page of Panizza!

# Document
## Controller & Service
First,make a dir named controller and defind tow structs as controller and service:
```golang
package controller

import(
    ."github.com/CP-panizza/panizza"
)

//@Controller
type MyController struct{
    Hello Handle `methods:"GET" path:"/hello" inject:"Hello"`
}

//@Service
type MyService struct {
    DB *sql.DB `inject:"db"`
}

func (this *MyService)Hello(ctx *HandleContext){
    ctx.JSON_write(200, JSONObject{
        "msg": "hello!",
    })
}
```
&nbsp;&nbsp;We can use <font color="green">//@Controller</font> to mark struct MyController as a controller, MyController like a interface, method is assign the http method, has GET, POST, PUT, DELETE...<br/>
<font color="green">tip:&nbsp;betwen // and @ do not has space!!</font></br>
&nbsp;&nbsp;And then, we defind a struct who has a function named Hello to implement MyController's Helle inline function.</br>
Open cmd and goto my_first_pro direction enter order:
```cmd
Pan -c    # -c it means complie.
```
This order is auto complie the components for current project. Add the MyController and MyService to the Ioc as a component.

Finally, build <font color="green">my_first_pro.go</font> and run, visit the http://127.0.0.1:8080/mycontroller/hello</br>
you will get the result.
```json
 {"msg": "hello!"}
 ```

### tip:&nbsp;&nbsp;panizza support <font color="green">@Controller,@Service,@Filter,@Aspect</font>.All of them are called component.</br>
If you make a new component,must be run the order <font color="green">Panizza -c</font> to make it work!!!

# project config
A project has a config file, it name is application.conf and in this file has three already exsist varibles.</br>
PROJECT_PACKAGE:  your project's rootpath.</br>
PORT:  your app listening port.</br>
FILE_SERVER:  your file rootpath.</br>
you can set the 'inject' tag on your component field to get config.</br>
## Example:
application.conf
```text
APPID=wx264sd6sd844c
```
Service.go
```go
//@Service
type MyService struct{
    AppId string `inject:"APPID"`  //you can get the config by inject:"APPID" tag
}

func (this *MyService)Handle(ctx *HandleContext){
        fmt.Println(this.AppId) //wx264sd6sd844c
}
```
And you can defind a list in config like this:</br>
```text
#[]string
MYLIST=[aaaa,bbbb,cccc,dddd]
#map[string]string
MYMAP={name:aaaa,school:bbb}
```
# inject
It likes springboot's @AotoWired
