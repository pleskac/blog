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
	_ "unicode/utf8"
)

type Picture struct {
	Guid         string
	Post_excerpt string
	Id           string
	Meta_value   string
}

type Post struct {
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
	r.HandleFunc("/{"+postId+":[0-9]+}", PostHandler)
	fmt.Println("Router:", r)
	http.ListenAndServe(":1337", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")

	output := getAllPosts()
	fmt.Println(output) //THIS WORKS

	enc := json.NewEncoder(w)
	err := enc.Encode(output) //WHY WONT YOU RETURN???
	fmt.Println(err)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")
	vars := mux.Vars(r)
	post := vars[postId]

	output := getPictures(post)

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func getAllPosts() []Post {
	db := Connect()
	defer db.Close()
	query := "SELECT id,post_title FROM wp_posts WHERE post_type = 'post'"

	rows, _, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	var allPosts []Post

	for _, row := range rows {
		title := row.Str(1)

		allPosts = append(allPosts, Post{row.Str(0), fromWindows1252(title)})
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
		case 0x80:
			r = 0x20AC
		case 0x82:
			r = 0x201A
		case 0x83:
			r = 0x0192
		case 0x84:
			r = 0x201E
		case 0x85:
			r = 0x2026
		case 0x86:
			r = 0x2020
		case 0x87:
			r = 0x2021
		case 0x88:
			r = 0x02C6
		case 0x89:
			r = 0x2030
		case 0x8A:
			r = 0x0160
		case 0x8B:
			r = 0x2039
		case 0x8C:
			r = 0x0152
		case 0x8E:
			r = 0x017D
		case 0x91:
			r = 0x2018
		case 0x92:
			r = 0x2019
		case 0x93:
			r = 0x201C
		case 0x94:
			r = 0x201D
		case 0x95:
			r = 0x2022
		case 0x96:
			r = 0x2013
		case 0x97:
			r = 0x2014
		case 0x98:
			r = 0x02DC
		case 0x99:
			r = 0x2122
		case 0x9A:
			r = 0x0161
		case 0x9B:
			r = 0x203A
		case 0x9C:
			r = 0x0153
		case 0x9E:
			r = 0x017E
		case 0x9F:
			r = 0x0178
		default:
			r = rune(b)
		}

		buf.WriteRune(r)
	}

	return string(buf.Bytes())
}
