## PANIZZA
## web  https://cp-panizza.github.io/panizza/
Quickly to build a restapi web application!!!

## QUICK START
#web  https://cp-panizza.github.io/Panizza/
Install panizza

``` cmd
go get github.com/CP-Panizza/Panizza
```

Add panizza to environment variable path
```cmd
Path=${GOPATH}\src\github.com\CP-Panizza\Panizza\bin
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
And you can defind another data type in config like this:</br>
```text
#[]string
MYLIST=[aaaa,bbbb,cccc,dddd]
#map[string]string
MYMAP={name:aaaa,school:bbb}
```
# inject
It likes springboot's @AotoWired

Set it on your component field like this:</br>
```go
//@Service
type MyService struct{
    GormDB *gorm.DB `inject:"db"`
}
```
ioc will inject value into that field.

# Add been to Ioc.
You can defind a function to add a been to ioc by use @Been:</br>
```go
var db *gorm.DB

//db will be add to ioc. and you can use `inject:"gorm_db"` tag get this been.
//@Been(name="gorm_db")
func DB_Been()interface{}{
	return db
}

func init() {
	var err error
	db, err = gorm.Open("mysql", "localhost:3306@xxxxxx")
	if err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(10)

	db.DB().SetMaxOpenConns(100)

	if err := db.DB().Ping(); err != nil {
		panic("connect err!")
	}

	fmt.Println("connect seccess!")
}
```

#router group
```go
type MyController struct{
    Hello Handle `method:"GET" path:"/Hello" inject:"Hello"`
    World Handle `method:"GET" path:"/World" inject:"World"`
    Say Handle `method:"GET" path:"/Say" inject:"Say"`
}
//accordding to this controller struct,the router default perfix is struct name,
//the router is build to /MyController/Hello, /MyController/World,/MyController/Say
//if you want to change router perfix, you can make MyController implement method "GroupName()string" like this:

func (this *MyController)GroupName()string {
	return "xxx"
}
//the router is build to /xxx/Hello, /xxx/World,/xxx/Say
```




