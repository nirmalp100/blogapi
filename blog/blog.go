package blog

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	_ "github.com/lib/pq"
)

type book struct {
	Id          string `json:"Id"`
	Title       string `json:"Title"`
	Content     string `json:"Content"`
	Authorname  string `json:"Authorname"`
	Publishdate string `json:"Publishdate"`
	Publishtime string `json:"Publishtime"`
}

type books []book

func CheckErr(err error) { // function to check err
	if err != nil {
		log.Fatal(err)
	}
}

//CONNECT TO DATABASE
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "nirmal123"
	dbname   = "blog"
)

var db *sql.DB
var dberror error

func SetupDB() {
	dbinfo := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, dberror = sql.Open("postgres", dbinfo)
	CheckErr(dberror)

	if err := db.Ping(); err != nil {
		CheckErr(err)
	}

	fmt.Println("database is accessed")

}

//CONNECT TO CACHE
var rdb *redis.Client

func ConnectRedisCache() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	if _, redis_err := rdb.Ping().Result(); redis_err != nil {
		fmt.Println(redis_err.Error())
		panic("Error: Unable to connect to Redis")
	} else {
		fmt.Println("Caching active")
	}
}

//FUCNTION TO CREATE BLOGPOST

func CreateBlogPost(g *gin.Context) {
	// in database

	header := g.GetHeader("Authorization")
	splittedheader := strings.Fields(header)
	signedToken := splittedheader[1]

	email := ExtractClaims(signedToken)

	rows, err := db.Query("Select user_id from users where email=$1", email)
	CheckErr(err)
	//From the signed token, email is extracted
	//userid of the user is retrieved from the users database

	var userid int
	for rows.Next() {

		err := rows.Scan(&userid)

		CheckErr(err)

	}

	v, err := ioutil.ReadAll(g.Request.Body)
	CheckErr(err)
	g.Request.Body.Close()
	var data map[string]string
	json.Unmarshal(v, &data)

	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}

	_, err = db.Exec("INSERT INTO blogs(id,title,content,author_name,publish_date,publish_time,author_id) VALUES(nextval('blogid'),$1,$2,$3,$4,$5,$6)", data["title"], data["content"], data["author_name"], data["publish_date"], data["publish_time"], userid)
	CheckErr(err)

	g.JSON(http.StatusOK, "blog successfully created")

	// in cache

	rdb.HSet("blogposts", data["title"], string(v))

}

func FetchFromDatabase(g *gin.Context) {
	rows, err := db.Query("Select * from blogs order by id asc")

	CheckErr(err)
	defer rows.Close()

	var data books

	for rows.Next() {
		a := book{}

		if err := rows.Scan(&a.Id, &a.Title, &a.Content, &a.Authorname, &a.Publishdate, &a.Publishtime); err != nil {
			CheckErr(err)
		}
		data = append(data, a)

	}
	g.JSON(http.StatusOK, data)

}

func FetchByTitle(g *gin.Context) {
	title := g.Param("title")

	//fetch from cache
	val := rdb.HGet("blogposts", title)
	data, err := val.Result()

	if err != nil {
		g.JSON(http.StatusOK, data)
	} else {

		//fetch from databse
		rows, err := db.Query("Select * from blogs order by id asc")

		CheckErr(err)
		defer rows.Close()

		var data books

		for rows.Next() {
			a := book{}

			if err := rows.Scan(&a.Id, &a.Title, &a.Content, &a.Authorname, &a.Publishdate, &a.Publishtime); err != nil {
				if err != nil {
					g.JSON(http.StatusBadRequest, nil)
				}
			}

			data = append(data, a)

		}

		for _, a := range data {
			if a.Title == title {
				g.IndentedJSON(http.StatusOK, a)
			}
		}

	}

}

func UpdateBlogPosts(g *gin.Context) {

	header := g.GetHeader("Authorization")
	splittedheader := strings.Fields(header)
	signedToken := splittedheader[1]

	email := ExtractClaims(signedToken)

	rows, err := db.Query("Select user_id from users where email=$1", email)
	CheckErr(err)
	//From the signed token, email is extracted
	//userid of the user is retrieved from the users database

	var userid int
	for rows.Next() {

		err := rows.Scan(&userid)

		CheckErr(err)

	}

	_, err = db.Query("Select author_id from blogs where author_id=$1", userid)

	if err != nil {
		g.JSON(http.StatusBadRequest, "Only the author can update post")

	} else {
		//updating in database
		v, err := ioutil.ReadAll(g.Request.Body)
		CheckErr(err)
		g.Request.Body.Close()
		var data map[string]string
		json.Unmarshal(v, &data)

		_, err = db.Exec("UPDATE blogs SET title=$2,content=$3,author_name=$4,publish_date=$5,publish_time=$6 WHERE title=$7", data["title"], data["content"], data["author_name"], data["publish_date"], data["publish_time"], data["title"])

		CheckErr(err)

		//update in cache
		rdb.HSet("blogposts", data["title"], string(v))

	}

}

func DeleteBlogPosts(g *gin.Context) {

	header := g.GetHeader("Authorization")
	splittedheader := strings.Fields(header)
	signedToken := splittedheader[1]

	email := ExtractClaims(signedToken)

	rows, err := db.Query("Select user_id from users where email=$1", email)
	CheckErr(err)
	//From the signed token, email is extracted
	//userid of the user is retrieved from the users database

	var userid int
	for rows.Next() {

		err := rows.Scan(&userid)

		CheckErr(err)

	}

	_, err = db.Query("Select author_id from blogs where author_id=$1", userid)

	if err != nil {
		g.JSON(http.StatusBadRequest, "Only the author can update post")

	} else {
		//Delete from database
		title := g.Param("title")
		_, err := db.Query("Delete from blogs where title=$1", title)
		CheckErr(err)

		//Delete from cache
		rdb.HDel("blogposts", title)
	}
}
