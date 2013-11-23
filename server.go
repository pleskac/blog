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
)

type Post struct {
	Guid         string
	Post_excerpt string
	Id           string
	Meta_value   string
}

const postId = "postId"

func main() {
	router := mux.NewRouter()
        r := router.Host("{domain:pleskac.org|api.pleskac.org|localhost}").Subrouter()

        r.HandleFunc("/blog/{"+postId+":[0-9]+}", PostHandler)
        http.ListenAndServe(":1337", r)
	fmt.Println("shouldn't be here")
	//go endpoint()
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
//	/blog/{ID}		specific post
func endpoint() {
	router := mux.NewRouter()
	r := router.Host("{domain:pleskac.org|api.pleskac.org|localhost}").Subrouter()

	r.HandleFunc("/blog/{"+postId+":[0-9]+}", PostHandler)
	http.ListenAndServe(":1337", r)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross domain AJAX requests
	w.Header().Set("Access-Control-Allow-Origin", "http://pleskac.org")
	vars := mux.Vars(r)
	post := vars[postId]
	fmt.Println(post)

	output := getPost(post)

	enc := json.NewEncoder(w)
	enc.Encode(output)
}

func getPost(p string) []Post {
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

	var myPost []Post

	for _, row := range rows {
		myPost = append(myPost, Post{row.Str(0), row.Str(1), row.Str(2), row.Str(3)})
	}

	return myPost
}