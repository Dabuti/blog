package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Dabuti/blog/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Env has all the environment data.
type Env struct {
	db *gorm.DB
}

func main() {
	// Create a database connection
	db, err := models.NewDB("blog:blog@tcp(localhost:3306)/go_blog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Panic(err)
	}
	db.SingularTable(true)
	defer db.Close()
	env := &Env{db: db}

	go watch("./posts/current")

	r := mux.NewRouter()
	r.HandleFunc("/secret", secret)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/new_post", newPost)
	r.HandleFunc("/posts/new", env.createPost).Methods("POST")
	r.HandleFunc("/posts", env.postsIndex).Methods("GET")
	r.HandleFunc("/posts/{id}", env.getPost).Methods("GET")
	r.HandleFunc("/posts2/{id}", env.getPost2).Methods("GET")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func (env *Env) createPost(w http.ResponseWriter, r *http.Request) {
	maxUploadSize := int64(4096)

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		fmt.Println(w, "FILE_TOO_BIG", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("data")
	if err != nil {
		fmt.Println(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println(w, "INVALID_FILE", http.StatusBadRequest)
		return
	}

	// Insert the post into the DB
	id, err := models.CreatePost(env.db, fileBytes)
	if err != nil {
		log.Panic(err)
	}
	fmt.Fprintf(w, "New post created with id: %d", id)
}

func (env *Env) getPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var post *models.Post
	id64, _ := strconv.ParseUint(vars["id"], 10, 64)
	post, err := models.GetPost(env.db, id64)
	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, "%d, %s, \n%s", post.ID, post.Title, post.Body)
}

func (env *Env) getPost2(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id64, _ := strconv.ParseUint(vars["id"], 10, 64)
	html, err := models.GetPost2(env.db, id64)
	if err != nil {
		log.Fatal(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(w, "%s", html)
}

func (env *Env) postsIndex(w http.ResponseWriter, r *http.Request) {
	posts, err := models.AllPosts(env.db)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	for _, post := range posts {
		fmt.Fprintf(w, "%d, %s, %s, %s\n", post.ID, post.Title, post.CreatedAt, post.Body)
	}
}
