## PANIZZA

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
Panizza -n "your project name like `my_first_pro`"
```
And then, the cli created a project named my_first_pro on the current path,like this:
```cmd
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
