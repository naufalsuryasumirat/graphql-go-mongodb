package database

import (
	"context"
	"log"
	"time"

	"github.com/naufalsuryasumirat/graphql-go-mongodb/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client 	*mongo.Client
}

type AuthorDB struct {
	ID        	string `json:"_id" bson:"_id,omitempty"`
	Name      	string `json:"name"`
	Birthdate 	primitive.DateTime `json:"birthdate"`
}

type BookDB struct {
	ID     	string  `json:"_id" bson:"_id,omitempty"`
	IDAuth 	string	`json:"id_auth"`
	Title  	string  `json:"title"`
}

const TimeFormat = "02-01-2006"

const short = 5 * time.Second
const medium = 10 * time.Second
const long = 30 * time.Second

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), medium)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{client: client}
}

func (db *DB) AddAuthor(input *model.AuthorInput) *model.Author {
	collection := db.client.Database("library").Collection("authors")
	ctx, cancel := context.WithTimeout(context.Background(), short)
	defer cancel()

	dt, _ := time.Parse(TimeFormat, input.Birthdate)
	toInsert := AuthorDB{
		Name: 		input.Name,
		Birthdate: 	primitive.NewDateTimeFromTime(dt),
	}

	res, err := collection.InsertOne(ctx, toInsert)
	if err != nil {
		log.Fatal(err)
	}

	return &model.Author{
		ID:			res.InsertedID.(primitive.ObjectID).Hex(),
		Name: 		input.Name,
		Birthdate: 	input.Birthdate,
	}
}

func (db *DB) FindAuthorByID(id string) *model.Author {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	collection := db.client.Database("library").Collection("authors")
	ctx, cancel := context.WithTimeout(context.Background(), short)
	defer cancel()

	res := collection.FindOne(ctx, bson.M{"_id": objId})
	authorDB := AuthorDB{}
	res.Decode(&authorDB)

	return convertAuthorType(&authorDB)
}

func (db *DB) AllAuthors() []*model.Author {
	collection := db.client.Database("library").Collection("authors")
	ctx, cancel := context.WithTimeout(context.Background(), long)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var authors []*model.Author
	for cur.Next(ctx) {
		authorDB := AuthorDB{}
		err := cur.Decode(&authorDB)
		if err != nil {
			log.Fatal(err)
		}
		author := convertAuthorType(&authorDB)
		authors = append(authors, author)
	}

	return authors
}

func (db *DB) AddBook(input *model.BookInput, idAuthor *string) *model.Book {
	author := db.FindAuthorByID(*idAuthor)

	collectionBooks := db.client.Database("library").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), short)
	defer cancel()

	toInsert := BookDB{
		IDAuth: *idAuthor,
		Title: 	input.Title,
	}

	resBook, err := collectionBooks.InsertOne(ctx, toInsert)
	if err != nil {
		log.Fatal(err)
	}

	return &model.Book{
		ID:		resBook.InsertedID.(primitive.ObjectID).Hex(),
		Title: 	input.Title,
		Author: author,
	}
}

func (db *DB) FindBookByID(id string) *model.Book {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}

	collection := db.client.Database("library").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), short)
	defer cancel()

	res := collection.FindOne(ctx, bson.M{"_id": objId})
	bookDB := BookDB{}
	res.Decode(&bookDB)

	author := db.FindAuthorByID(bookDB.IDAuth)

	return convertBookType(&bookDB, author)
}

func (db *DB) AllBooks() []*model.Book {
	collection := db.client.Database("library").Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), long)
	defer cancel()

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var books []*model.Book
	for cur.Next(ctx) {
		bookDB := BookDB{}
		err := cur.Decode(&bookDB)
		if err != nil {
			log.Fatal(err)
		}
		author := db.FindAuthorByID(bookDB.IDAuth)
		books = append(books, convertBookType(&bookDB, author))
	}

	return books
}

func (db *DB) AllBooksByAuthor(id string) []*model.Book {
	books := db.AllBooks()
	var booksBy []*model.Book
	for _, element := range books {
		if element.Author.ID == id {
			booksBy = append(booksBy, element)
		}
	}
	return booksBy
}

func convertAuthorType(authorDB *AuthorDB) *model.Author {
	return &model.Author{
		ID: 		authorDB.ID,
		Name: 		authorDB.Name,
		Birthdate: 	authorDB.Birthdate.Time().Format(TimeFormat),
	}
}

func convertBookType(bookDB *BookDB, author *model.Author) *model.Book {
	return &model.Book{
		ID: 	bookDB.ID,
		Title: 	bookDB.Title,
		Author: author,
	}
}