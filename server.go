package main

import (
	//change to github.com/pleskac/SpotLocator/dblayer
	"bytes"
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"fmt"
	z_mysql "github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	_ "log"
	"net/http"
	"strings"
	_ "unicode/utf8"
)

type Picture struct {
	Guid         string
	Post_excerpt string
	Id           string
	Meta_value   string
}

type Post struct {
	Id       string
	Title    string
	Content  string
	Pictures []Picture
}

type PostStub struct {
	Id         string
	Post_title string
}

const postId = "postId"

func main() {
	fmt.Println("Starting endpoint.go. Will block")
	endpoint()
}

func Connect() z_mysql.Conn {
	//Set up database connection
	db := z_mysql.New("tcp", "", "127.0.0.1:3306", "root", "rootroot", "wordpress")
	err := db.Connect()
	if err != nil {
		fmt.Println("ERROR CONNECTING:", err)
		panic(err)
	}

	return db
}

//JSON endpoints:
//	/{ID}		specific post
//	/blog		list of all posts
func endpoint() {
	router := mux.NewRouter()
	r := router.Host("{domain:pleskac.org|www.pleskac.org|localhost}").Subrouter()
	r.HandleFunc("/blog", HomeHandler)
	r.HandleFunc("/blog/{"+postId+":[0-9]+}", PostHandler)
	r.HandleFunc("/{"+postId+":[0-9]+}", PostDataHandler)

	http.ListenAndServe(":1337", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	output := getAllPosts()

	enc := json.NewEncoder(w)
	err := enc.Encode(output)

	if err != nil {
		fmt.Println(err)
	}
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post := vars[postId]

	top := `<head>
	<meta http-equiv="content-type" content="text/html; charset=utf-8">
	<title>dog food</title>
	<link rel="stylesheet" href="http://pleskac.org/blog/main.css">
	<script type="text/javascript" src="https://ajax.googleapis.com/ajax/libs/jquery/1.7.2/jquery.min.js"></script>
	<script type="text/javascript" src="http://pleskac.org/blog/page.js"></script>
	<style type="text/css"></style></head>`

	postTitle := GetPostTitle(post)

	fmt.Fprintf(w, "%s<body><h1>%s</h1><div id=\"post_content\"></div></body>", top, postTitle)
}

func PostDataHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")
	vars := mux.Vars(r)
	post := vars[postId]

	pictures := getPictures(post)

	title, content := GetPostData(post)

	//Id, title, content, pictures
	thePost := Post{
		Id:       post,
		Title:    title,
		Content:  content,
		Pictures: pictures,
	}

	enc := json.NewEncoder(w)
	enc.Encode(thePost)
}

func getAllPosts() []PostStub {
	db := Connect()
	defer db.Close()
	query := "SELECT id,post_title FROM wp_posts WHERE post_type = 'post'"

	rows, _, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	var allPosts []PostStub

	for _, row := range rows {
		title := row.Str(1)

		allPosts = append(allPosts, PostStub{row.Str(0), fromWindows1252(title)})
	}

	return allPosts
}

func getPictures(p string) []Picture {
	db := Connect()
	defer db.Close()

	//Get the current trip, if it exists
	query := "SELECT wp_posts.guid,wp_posts.post_excerpt,wp_posts.id,wp_postmeta.meta_value "
	query += "FROM wp_posts "
	query += "LEFT JOIN wp_postmeta "
	query += "ON wp_posts.ID = wp_postmeta.post_id "
	query += "WHERE wp_posts.post_parent = \"" + p + "\" "
	query += "AND wp_posts.post_type = \"attachment\" "
	query += "AND wp_postmeta.meta_key = \"_wp_attachment_metadata\""

	rows, _, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	var pictures []Picture

	for _, row := range rows {
		pictures = append(pictures, Picture{row.Str(0), row.Str(1), row.Str(2), row.Str(3)})
	}

	return pictures
}

//used from Stack Overflow
//http://stackoverflow.com/questions/6927611/go-language-how-to-convert-ansi-text-to-utf8
func fromWindows1252(str string) string {
	var arr = []byte(str)
	var buf = bytes.NewBuffer(make([]byte, 512))
	var r rune

	for _, b := range arr {
		switch b {
		case 0x96:
			r = 0x2013
		default:
			r = rune(b)
		}

		buf.WriteRune(r)
	}

	return strings.Trim(string(buf.Bytes()), "\u0000")
}

func GetPostData(id string) (string, string) {
	db := Connect()
	defer db.Close()

	query := "select post_title,post_content from wp_posts where id = " + id

	rows, _, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	if len(rows) == 1 {
		return rows[0].Str(0), rows[0].Str(1)
	}

	return "Post not found. How did you do this? Email me.", "This is where the content should go."
}

func GetPostTitle(id string) string {
	db := Connect()
	defer db.Close()

	query := "select post_title from wp_posts where id = " + id

	rows, _, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	if len(rows) == 1 {
		return rows[0].Str(0)
	}

	return "Post not found. How did you do this? Email me."
}
