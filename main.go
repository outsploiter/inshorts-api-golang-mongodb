package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"strconv"
	
	"sync"
	"time"
	"context"
	"os"
	"reflect"
	"log"
	
	"helper" 					//package for DB Initialization self-created

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Article struct {
	ID         string `json:"id"`
	Title 	   string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Content    string `json:"content"`
	TimeStamp  string `json:"timestamp"`
}

type articleHandlers struct {
	sync.Mutex
	store map[string]Article
}

// Inserting a element into the DataBase

func insertData(article Article){

	result, insertErr := col.InsertOne(context.TODO(), article)
	if insertErr != nil {// Print if any Error
	fmt.Println("InsertOne ERROR:", insertErr)
	os.Exit(1)
	} else { // Print Success Message
		fmt.Println("InsertOne() result type: ", reflect.TypeOf(result))
		fmt.Println("InsertOne() API result:", result)
		}

}

//function to choose handle requests

func (h *articleHandlers) articles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {//Choose Function with request
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

//function to handle get requests

func (h *articleHandlers) get(w http.ResponseWriter, r *http.Request) {

	h.Lock()
	var allArticles []Article
	cur, err := col.Find(context.TODO(), bson.M{})                         //returns all documents

	if err != nil {
		initDB.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var article Article
		err := cur.Decode(&article)
		if err != nil {
			log.Fatal(err)
		}

		// add items our array
		allArticles = append(allArticles, article)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(allArticles)                          // encode similar to serialize process.
	 h.Unlock()
}


//function to handle search and select article by ID

func (h *articleHandlers) getArticle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
  temp := parts[2]
  id, err := strconv.Atoi(temp)
  if err != nil {
	query:=r.URL.Query()                                                 // SEARCH WITH q in url GET
    h.Lock()
  	var allArticles []Article
 
    filter := bson.D{{"$or", bson.A{
        bson.D{{"title", primitive.Regex{Pattern: query["q"][0], Options: "i"}}},
        bson.D{{"content", primitive.Regex{Pattern: query["q"][0], Options: "i"}}},}},}
    
    cur,err := col.Find(context.TODO(), filter)
  	if err != nil {
  		initDB.GetError(err, w)
  		return
  	}
  	defer cur.Close(context.TODO())

  	for cur.Next(context.TODO()) {
  		var article Article
  		err := cur.Decode(&article)
  		if err != nil {
  			log.Fatal(err)
  		}

  		allArticles = append(allArticles, article)                 // add item our array
  	}

  	if err := cur.Err(); err != nil {
  		log.Fatal(err)
  	}

  	json.NewEncoder(w).Encode(allArticles)
  	
  	 h.Unlock()                                                      // encode similar to serialize process.
	} else{
		h.Lock()
		var allArticles []Article
	cur, err := col.Find(context.TODO(), bson.D{{"id",strconv.Itoa(id) }})
	
																		// gET USING ARTICLE ID NOT OBJECT_ID
	if err != nil {
		initDB.GetError(err, w)
		return
	}
	
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var article Article
		err := cur.Decode(&article)
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		
		allArticles = append(allArticles, article)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(allArticles)
	h.Unlock()

	}

	}
	
//function to handle post requests

func (h *articleHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("require content-type 'application/json',got '%s' instead", ct)))
		return
	}

	var article Article
	err = json.Unmarshal(bodyBytes, &article)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	t := time.Now()
	article.TimeStamp = fmt.Sprintf(t.String())
  	page ,e := col.CountDocuments(context.TODO(),bson.M{})
  	fmt.Printf("%v %v",int64(page),e)
 	article.ID = strconv.Itoa(int(page + 1))                    //using our own id instead of using object id created by default
 	
 	
 	//Inserting Data into Database
 	
 	
	insertData(article)
	h.Lock()
	h.store[article.TimeStamp] = article
	defer h.Unlock()
}

func newArticleHandlers() *articleHandlers {
	return &articleHandlers{
		store: map[string]Article{},
	}

}

//Initializing DB
var col = initDB.ConnectDB()


func main() {
	articleHandlers := newArticleHandlers()
	http.HandleFunc("/articles", articleHandlers.articles)
	http.HandleFunc("/articles/", articleHandlers.getArticle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
