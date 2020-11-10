**Inshorts-api-golang-mongodb**

api developed with golang and mongodb server

only with default golang packages

    "go.mongodb.org/mongo-driver/bson"

    "go.mongodb.org/mongo-driver/bson/primitive"

 

<span style="text-decoration:underline;">without gorilla or mux</span>

**only with mongo-driver**

you might have to install mongo-driver if its not already installed

to install the same please look into - https://github.com/mongodb/mongo-go-driver

#######

**WORKING**

**GET -  "/articles - show all articles**

**GET - "/articles/&lt;id> - show article with id**

**POST - "/articles - Add to record to articles**

While posting only give title, subtitle and content of the article. ID, Object_ID  and Time Stamp will be automatically generated.

**GET - "/articles/search?q=&lt;search_Word> - to search articles with search word in it(in content, subtitle)**

unit testing will be uploaded soon..

any issues - contact me on instagram @outsploiter
