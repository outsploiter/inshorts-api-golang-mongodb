# inshorts-api-golang-mongodb

api developed with golang and mongodb server

only with default golang packages 

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
  
without gorilla or mux

only with mongo-driver

##########################################

Pre req


you might have to install mongo-driver if its not already installed
to install the same please look into - https://github.com/mongodb/mongo-go-driver

#######
**working**

GET -  "/articles - show all articles
GET - "/articles/<id> - show article with id
POST - "/articles - Add to record to articles
GET - "/articles/search?q=<search_Word> - to search articles with search word in it(in content, subtitle)

unit testing will be uploaded soon.. 

any issues - contact me on instagram @outsploiter 
