package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const dbName = "beerraters"

type Beer struct {
	BreweryID primitive.ObjectID `json:"breweryId" bson:"brewery_id"`
	BreweryName string `json:"breweryName" bson:"brewery_name"`
	CategoryID primitive.ObjectID `json:"categoryId" bson:"category_id"`
	Consumes []string `json:"consumes" bson:"consumes"`
	CountryID primitive.ObjectID `json:"countryId" bson:"country_id"`
	CountryName string `json:"countryName" bson:"country_name"`
	ID primitive.ObjectID `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Ratings []string `json:"ratings" bson:"ratings"`
	Reviews []string `json:"reviews" bson:"reviews"`
	StyleID primitive.ObjectID `json:"styleId" bson:"style_id"`
	StyleName string `json:"styleName" bson:"style_name"`

}

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})
	app.Get("/beer/:id?", getBeer)
	app.Get("/beers", getBeers)
	// app.Post("/beer", createBeer)
	// app.Put("/beer/:id", updateBeer)
	// app.Delete("/beer/:id", deleteBeer)

}

//GetMongoDbConnection get connection of mongodbgo
func GetMongoDbConnection() (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	// client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://adminuser:password123@mongo-nodeport-svc:27017"))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func getMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := GetMongoDbConnection()

	if err != nil {
		return nil, err
	}

	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}

func getBeers(c *fiber.Ctx) error {
	const collectionName = "beers"
	collection, err := getMongoDbCollection(dbName, collectionName)
	if err != nil {
		return c.Status(500).Send([]byte(err.Error()))
	}

	var filter bson.M = bson.M{}

	var results []bson.M
	cur, err := collection.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err != nil {
		return c.Status(500).Send([]byte(err.Error()))
	}

	cur.All(context.Background(), &results)

	if results == nil {
		return c.SendStatus(404)
	}

	json, _ := json.Marshal(results)
	return c.Send(json)
}

func getBeer(c *fiber.Ctx) error {
	const collectionName = "beers"
	collection, err := getMongoDbCollection(dbName, collectionName)
	if err != nil {
		return c.Status(500).Send([]byte(err.Error()))
	}

	var filter bson.M = bson.M{}

	if c.Params("id") != "" {
		id := c.Params("id")
		objID, _ := primitive.ObjectIDFromHex(id)
		filter = bson.M{"_id": objID}
	}

	var beer Beer
	err = collection.FindOne(context.Background(), filter).Decode(&beer)


	if err != nil {
		return c.Status(500).Send([]byte(err.Error()))
	}

	json, _ := json.Marshal(beer)
	return c.Send(json)
}

func main() {
	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	// Pass engine to Fiber's Views Engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:30084, http://localhost:5000",
    AllowHeaders:  "Origin, Content-Type, Accept",
}))

	app.Static("/", "./public")

	setupRoutes(app)

	app.Listen(":3000")
}
