package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Author struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type Blog struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Author *Author            `json:"author,omitempty" bson:"author,omitempty"`
	Title  string             `json:"title,omitempty" bson:"title,omitempty"`
	Body   string             `json:"body,omitempty" bson:"body,omitempty"`
}

func GetAuthors() ([]Author, error) {
	var authors []Author
	collection := client.Database("app_db").Collection("authors")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var author Author
		cursor.Decode(&author)
		authors = append(authors, author)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

func GetAuthor(id string) (*Author, error) {
	var author Author
	collection := client.Database("app_db").Collection("authors")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	oid, _ := primitive.ObjectIDFromHex(id)
	err := collection.FindOne(ctx, Author{ID: oid}).Decode(&author)
	if err != nil {
		return nil, err
	}

	return &author, nil
}

func CreateAuthor(author Author) (*Author, error) {
	collection := client.Database("app_db").Collection("authors")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := collection.InsertOne(ctx, author)
	if err != nil {
		return nil, err
	}

	author.ID = result.InsertedID.(primitive.ObjectID)
	return &author, nil
}

func GetBlogs() ([]Blog, error) {
	var blogs []Blog
	collection := client.Database("app_db").Collection("blogs")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var blog Blog
		cursor.Decode(&blog)
		blogs = append(blogs, blog)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func GetBlog(id string) (*Blog, error) {
	var blog Blog
	collection := client.Database("app_db").Collection("blogs")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	oid, _ := primitive.ObjectIDFromHex(id)
	err := collection.FindOne(ctx, Blog{ID: oid}).Decode(&blog)
	if err != nil {
		return nil, err
	}

	return &blog, nil
}

func CreateBlog(blog Blog) (*Blog, error) {
	collection := client.Database("app_db").Collection("blogs")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := collection.InsertOne(ctx, blog)
	if err != nil {
		return nil, err
	}

	blog.ID = result.InsertedID.(primitive.ObjectID)
	return &blog, nil
}

func main() {
	var err error
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	authorType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Author",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"firstname": &graphql.Field{
				Type: graphql.String,
			},
			"lastname": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	blogType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Blog",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"body": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"authors": &graphql.Field{
				Type: graphql.NewList(authorType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetAuthors()
				},
			},
			"author": &graphql.Field{
				Type: authorType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return GetAuthor(params.Args["id"].(string))
				},
			},
			"blogs": &graphql.Field{
				Type: graphql.NewList(blogType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetBlogs()
				},
			},
			"blog": &graphql.Field{
				Type: blogType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return GetBlog(params.Args["id"].(string))
				},
			},
		},
	})

	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"createAuthor": &graphql.Field{
				Type: authorType,
				Args: graphql.FieldConfigArgument{
					"firstname": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"lastname": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					author := Author{
						Firstname: params.Args["firstname"].(string),
						Lastname:  params.Args["lastname"].(string),
					}

					return CreateAuthor(author)
				},
			},
			"createBlog": &graphql.Field{
				Type: blogType,
				Args: graphql.FieldConfigArgument{
					"authorID": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"body": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					var blog = Blog{}
					author, err := GetAuthor(params.Args["authorID"].(string))
					if err != nil {
						return blog, err
					}

					blog.Author = author
					blog.Title = params.Args["title"].(string)
					blog.Body = params.Args["body"].(string)

					return CreateBlog(blog)
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootMutation,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	})

	if err := http.ListenAndServe(":8088", nil); err != nil {
		log.Fatal(err.Error())
	}
}
