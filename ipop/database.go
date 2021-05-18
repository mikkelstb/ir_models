package ipop

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	irmodels "github.com/mikkelstb/ir_models"
	"github.com/mikkelstb/simplelog"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {

	login_credentials string
	initialized bool
	estimated_articles int
	processed_articles int
	cached_articles []irmodels.Article
	offset int
}


func (this *Database) Init(preferences map[string]string) {

	this.login_credentials = preferences["username"] + ":" + preferences["password"] + "@/" + preferences["dbname"]
	simplelog.InitLoggers("./database.log")
	this.initialized = true
	this.offset = 0
}

func (this *Database) checkDB() {

	connection, err := sql.Open("mysql", this.login_credentials)
	if err != nil {
		simplelog.Error.Println(err.Error())
	}
	defer connection.Close()

	connection.SetConnMaxLifetime(time.Second*5)
	connection.SetMaxOpenConns(10)
	connection.SetConnMaxIdleTime(10)
	
	query := `select id, headline, story from article limit ?, ?`

	rows, err := connection.Query(query, this.offset, 100)
	if err != nil {
		simplelog.Error.Println(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		this.offset++
		var intro string
		var story string
		var doc_id int

		if err := rows.Scan(&doc_id, &intro, &story); err != nil {
			simplelog.Error.Println(err.Error())
		}
		article := irmodels.Article{Doc_id: doc_id, Text: strings.Join([]string{intro, story}, " ")}
		this.cached_articles = append(this.cached_articles, article)
	}
	fmt.Printf("Found %v articles so far\n", this.offset)
}


func (this *Database) HasNext() bool {

	if len(this.cached_articles) > 0 {
		return true
	}
	this.checkDB()
	if len(this.cached_articles) > 0 {
		return true
	}
	return false
}


func (this *Database) GetNext() *irmodels.Article {
	if len(this.cached_articles) == 0 {
		return nil
	}
	art := this.cached_articles[0]
	this.cached_articles = this.cached_articles[1:]
	return &art
}


