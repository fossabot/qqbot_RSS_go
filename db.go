package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type groupUrl struct {
	GroupCode int
	Url       string
	UrlName   string
}

type AutoGenerated struct {
	Mysql Mysql `yaml:"mysql"`
	QqBot QqBot `yaml:"qq_bot"`
}

type Mysql struct {
	Db       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

type QqBot struct {
	Url         string `yaml:"url"`
	AccessToken string `yaml:"access_token"`
	Port        string `yaml:"port"`
}

var DB *sql.DB
var setting AutoGenerated

func InitDB() *sql.DB {
	dsn := setting.Mysql.User + ":" + setting.Mysql.Password + "@tcp(" + setting.Mysql.Host + ":" + setting.Mysql.Port + ")/" + setting.Mysql.Db
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Println("Init MySQL Success")
	return db
}

func queryUrl(rssType int) []string {
	var result []string
	row, err := DB.Query("select a.GroupCode, b.Url from group_info as a INNER JOIN url_info as b ON a.GroupId = b.GroupId inner join rss_type as c on b.RssTypeId = c.RssTypeId where a.Status = 1 and b.status = 1 and c.RssTypeId = ?", rssType)
	if err != nil {
		log.Fatal(err)
	}
	for row.Next() {
		var gUrl groupUrl
		err := row.Scan(&gUrl.GroupCode, &gUrl.Url)
		if err != nil {
			log.Fatal(err)
		}
		dict, _ := json.Marshal(gUrl)
		data := string(dict)
		result = append(result, data)
	}
	return result
}

func queryGroup(groupId int) []string {
	var result []string
	row, err := DB.Query("select b.UrlName from group_info as a INNER JOIN url_info as b ON a.GroupId = b.GroupId where a.Status = 1 and b.status = 1 and a.GroupCode = ?", groupId)
	if err != nil {
		log.Fatal(err)
	}
	for row.Next() {
		var gUrl groupUrl
		err := row.Scan(&gUrl.UrlName)
		if err != nil {
			log.Fatal(err)
		}
		dict, _ := json.Marshal(gUrl)
		data := string(dict)
		result = append(result, data)
	}
	return result
}
