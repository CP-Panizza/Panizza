package Panizza

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

//自定义数据库主键生成

/**
idLength:主键生成的长度
table_name:表名，数据库内表的真实名字
pk:表的主键字段
sign: user00000012 中的user
*/
type PKGenarater struct {
	DB        *sql.DB `inject:"db"`
	TableName string
	IdLength  int
	PK        string
	Sign      string
}

func (p *PKGenarater) GetPK() (string, error) {
	str := strings.Builder{}
	str.WriteString("select max(a.")
	str.WriteString(p.PK)
	str.WriteString(") from ")
	str.WriteString(p.TableName)
	str.WriteString(" as a where a.")
	str.WriteString(p.PK)
	str.WriteString(" like '")
	str.WriteString(p.Sign)
	str.WriteString("%'")
	sql := str.String()
	fmt.Println(sql)
	row := p.DB.QueryRow(sql)
	var result string = ""
	//第一次查询，表的主键为空时
	if err := row.Scan(&result); err != nil {
		gpk := p.Sign
		for i := 0; i < (p.IdLength - len(p.Sign) - 1); i++ {
			gpk += "0"
		}
		gpk += "1"
		fmt.Println(gpk)
		return gpk, nil
	}

	fmt.Println("last val:", result)
	number := result[strings.LastIndex(result, "0")+1:]
	num, err := strconv.Atoi(number)
	if err != nil {
		return "", err
	}
	nextNumber := num + 1
	nextNumberStr := strconv.Itoa(nextNumber)
	gpk := p.Sign
	for i := 0; i < (p.IdLength - len(p.Sign) - len(nextNumberStr)); i++ {
		gpk += "0"
	}
	gpk += nextNumberStr
	fmt.Println("next val:", gpk)
	return gpk, nil
}
