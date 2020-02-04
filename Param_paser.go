package Panizza

import (
	"errors"
	"strconv"
	"strings"
)

type ParamParser struct {
	ParamMap map[string]string
	ParamKey []string
	ParamVal []string
}

func (pp *ParamParser) Add(key string, value string) {
	pp.ParamMap[key] = value
}

func (pp *ParamParser) Parame(key string) string {
	val, ok := pp.ParamMap[key]
	if ok {
		return val
	}
	return ""
}

func (pp *ParamParser) ParameInt(key string) (int, error) {
	val, ok := pp.ParamMap[key]
	if ok {
		data, err := strconv.Atoi(val)
		if err != nil {
			return data, err
		}
		return data, nil
	}
	return -1, errors.New("have not value:" + key)
}

func (pp *ParamParser) ParameInt64(key string) (int64, error) {
	val, ok := pp.ParamMap[key]
	if ok {
		data, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return data, err
		}
		return data, nil
	}
	return int64(-1), errors.New("have not value:" + key)
}

func (pp *ParamParser) ParameFloat64(key string) (float64, error) {
	val, ok := pp.ParamMap[key]
	if ok {
		data, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return data, err
		}
		return data, nil
	}
	return float64(-0.0), errors.New("have not value:" + key)
}

func (pp *ParamParser) ParameFloat32(key string) (float64, error) {
	val, ok := pp.ParamMap[key]
	if ok {
		data, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return data, err
		}
		return data, nil
	}

	return float64(-0.0), errors.New("have not value:" + key)
}

func (pp *ParamParser) ParameBool(key string) bool {
	val, ok := pp.ParamMap[key]
	if ok {
		if val == "true" {
			return true
		} else if val == "false" {
			return false
		}
	}
	return false
}

//检查自定义路由与访问路由是否匹配，如果匹配则返回true，
//若不匹配则返回false
func (pp *ParamParser) Check(registRouter string, url string) bool {

	var registList []string = strings.Split(registRouter, "/")

	var urlList []string = strings.Split(url, "/")

	if len(registList) != len(urlList) {
		//通过'/'拆分registRouter和url，判断长度是否相等，相等就做下一步操作，
		//不相等就return
		return false
	} else {
		for i := 0; i < len(registList); i++ {

			//获取路由数据并保存
			//如果registList[i]含有":"就continue,并把该值去除":"加入key数组，
			//并把下标相等的urlList[i]加入value数组
			if strings.Index(registList[i], ":") != -1 {
				continue
			}

			//如果registList[i]和urlList[i]不相等，说明路由不匹配，就return空的参数对象和false
			if strings.Compare(registList[i], urlList[i]) != 0 {
				return false
			}
		}

		return true
	}
}

//在Check方法通过之后便可以调用此方法获得url参数，参数放在对象的ParamMap之内
func (pp *ParamParser) GetParams(registRouter string, url string) {
	var registList []string = strings.Split(registRouter, "/")

	var urlList []string = strings.Split(url, "/")

	var key []string
	var value []string

	if len(registList) != len(urlList) {
		//通过'/'拆分registRouter和url，判断长度是否相等，相等就做下一步操作，
		//不相等就return
		return
	} else {
		for i := 0; i < len(registList); i++ {

			//获取路由数据并保存
			//如果registList[i]含有":"就continue,并把该值去除":"加入key数组，
			//并把下标相等的urlList[i]加入value数组
			if strings.Index(registList[i], ":") != -1 {
				key = append(key, string([]byte(registList[i])[1:]))
				value = append(value, urlList[i])
				continue
			}

			//如果registList[i]和urlList[i]不相等，说明路由不匹配，就return空的参数对象和false
			if strings.Compare(registList[i], urlList[i]) != 0 {
				return
			}
		}

		pp.ParamKey = key
		pp.ParamVal = value

		m := make(map[string]string)
		for i := 0; i < len(key); i++ {
			m[key[i]] = value[i]
		}

		pp.ParamMap = m

		return
	}
}
