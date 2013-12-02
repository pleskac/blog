package main

import (
	//change to github.com/pleskac/SpotLocator/dblayer
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
		if strings.Contains(row.Str(1), "\x96") {
			strings.Replace(title, "\x96", "|", -1)
			fmt.Println("New title:", title)
		}
		allPosts = append(allPosts, Post{row.Str(0), title})
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
