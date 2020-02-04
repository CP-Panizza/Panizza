package Panizza

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

var GlobalApi = APIManeger{
	APIMapper: []API{},
	APILength: 0,
}

type APIManeger struct {
	APIMapper []API
	APILength int
}

func (api *APIManeger) Add(Method string, URL string, Description string, Param []string, Location string) {
	api_temp := API{
		Method,
		URL,
		Description,
		Param,
		Location,
	}
	api.APIMapper = append(api.APIMapper, api_temp)
	api.APILength++
}

//生成API文档字符串
func (api *APIManeger) ToHtmlString() string {
	data, _ := json.Marshal(GlobalApi.APIMapper)
	fmt.Println(string(data))
	html := ApiTemplate1 + string(data) + ApiTemplate2
	return html
}

//传入保存的文件的绝对路径
func (api *APIManeger) ToHtmlFile(saveFilePath string) {
	data, _ := json.Marshal(GlobalApi.APIMapper)
	//fmt.Println(string(data))
	html := ApiTemplate1 + string(data) + ApiTemplate2
	if IsDir(saveFilePath) {
		panic(errors.New("Must give a file path, not a dir path!"))
	}

	//如果传入的路径是一个已经存在的文件就把他删除
	if Existe(saveFilePath) && !IsDir(saveFilePath) {
		if err := os.Remove(saveFilePath); err != nil {
			panic(err)
		}
	}

	file, err := os.Create(saveFilePath)

	if err != nil {
		panic(err)
	}

	file.WriteString(html)
	log.Println("Document is generate at:", saveFilePath)
	defer file.Close()
	return
}

type API struct {
	Method      string
	URL         string
	Description string
	Param       []string
	Location    string
}

func (a *API) String() string {
	str := "Method:" + a.Method + "\n\r" + "URL:" + a.URL + "\n\r" + "Description:" + a.Description + "\n\r" + "params:"
	for _, v := range a.Param {
		str += v[1:] + ","
	}
	str += "\n\rLocation:" + a.Location
	return str
}
