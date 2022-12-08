package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dietzy1/discordbot/src/bot/emotes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Repository interface {
	IncrementEmote(ctx context.Context, emote *emotes.Emote) error
	GetUserEmotes(ctx context.Context, emote *emotes.Emote) ([]emotes.Emote, error)
	GetServerEmote(ctx context.Context, emote *emotes.Emote) ([]emotes.Emote, error)
}

type Db struct {
	client *mongo.Client
}

func New() (*Db, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv(("DB"))))
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	a := &Db{client: client}

	return a, nil
}

func (db *Db) IncrementEmote(ctx context.Context, emote *emotes.Emote) error {
	collection := db.client.Database(emote.Guild).Collection(emote.User)
	_, err := collection.UpdateOne(ctx, bson.M{"emote": emote.Emote}, bson.M{"$inc": bson.D{{Key: "count", Value: 1}}, "$set": bson.D{{Key: "user", Value: emote.User}}}, options.Update().SetUpsert(true))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (db *Db) GetUserEmotes(ctx context.Context, emote *emotes.Emote) ([]emotes.Emote, error) {
	collection := db.client.Database(emote.Guild).Collection(emote.User)
	opts := options.Find().SetSort(bson.D{{Key: "count", Value: -1}})
	cursor, err := collection.Find(ctx, bson.D{{Key: "count", Value: bson.D{{Key: "$gt", Value: 0}}}}, opts)
	if err != nil {
		log.Println(err)
	}
	emotes := []emotes.Emote{}
	if err = cursor.All(ctx, &emotes); err != nil {
		log.Println(err)
	}
	return emotes, nil
}

// Takes in an emote and returns a leaderboard of the top users who have used that emote
func (db *Db) GetServerEmote(ctx context.Context, emote *emotes.Emote) ([]emotes.Emote, error) {
	database := db.client.Database(emote.Guild)
	results := []emotes.Emote{}
	cursor, err := database.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, collection := range cursor {
		aggregate, err := database.Collection(collection).Find(ctx, bson.M{"emote": emote.Emote})
		if err != nil {
			return nil, err //Could potentially be a bug here depends if the find returns an error if the emote is not found
		}
		for aggregate.Next(ctx) {
			result := emotes.Emote{}
			err := aggregate.Decode(&result)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			results = append(results, result)
		}
	}
	for i := 0; i < len(results); i++ {
		for j := 0; j < len(results); j++ {
			if results[i].Count > results[j].Count {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	return results, nil
}

func (a *Db) NewIndex(database string, collectionName string, field string, unique bool) {
	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1},
		Options: options.Index().SetUnique(unique),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := a.client.Database(database).Collection(collectionName)

	index, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Created new index:", index)
}
