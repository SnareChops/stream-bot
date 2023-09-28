package admin

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/SnareChops/twitchbot/db"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/exp/slices"
)

const (
	MysteryCollection = "mystery"
)

func SyncMysteries() {
	// Sync files to DB
	files, err := findMysteryFiles()
	if err != nil {
		panic(err)
	}
	dbItems, err := findMysteryItems()
	if err != nil {
		panic(err)
	}
	// Build full list of mystery items
	results := []MysteryItem{}
	for _, fileName := range files {
		exists := false
		for _, dbItem := range dbItems {
			if fileName == dbItem.Name {
				// If the item exists in the database, add it to the results
				results = append(results, dbItem)
				exists = true
			}
		}
		if !exists {
			// Add newly discovered mystery items to database
			item, err := addMysteryItem(fileName)
			if err != nil {
				fmt.Printf("Failed to add mystery item to database: %s\n", err)
			}
			results = append(results, item)
		}
	}

	// Remove mystery items no longer present
	for _, item := range dbItems {
		if !slices.Contains(files, item.Name) {
			err := deleteMysteryItem(item.Name)
			if err != nil {
				fmt.Printf("Failed to delete mystery item from database: %s\n", err)
			}
		}
	}
}

func GETMysteries(w http.ResponseWriter, req *http.Request) {
	mysteries, err := findMysteryItems()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	for _, mystery := range mysteries {
		AdminTemplate.ExecuteTemplate(w, "mystery-item", mystery)
	}
}

type MysteryListSignal struct {
	Kind  string        `json:"kind"`
	Items []MysteryItem `json:"items"`
}

type MysteryItem struct {
	Selected bool    `json:"selected" bson:"selected"`
	Name     string  `json:"name" bson:"name"`
	Type     string  `json:"type" bson:"type"`
	Start    float32 `json:"start" bson:"start"`
	End      float32 `json:"end" bson:"end"`
	Speed    float32 `json:"speed" bson:"speed"`
}

func findMysteryFiles() ([]string, error) {
	files, err := os.ReadDir("files/mystery")
	if err != nil {
		return nil, err
	}
	fileNames := []string{}
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames, nil
}

func findMysteryItems() ([]MysteryItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := db.Database.Collection(MysteryCollection).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	dbItems := []MysteryItem{}
	err = cursor.All(context.Background(), &dbItems)
	return dbItems, err
}

func addMysteryItem(name string) (MysteryItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	item := MysteryItem{
		Selected: false,
		Name:     name,
		Start:    0,
		End:      0,
		Speed:    1,
	}
	switch path.Ext(name) {
	case ".mp4":
		item.Type = "video"
	case ".mp3":
		item.Type = "audio"
	default:
		item.Type = "unknown"
	}
	_, err := db.Database.Collection(MysteryCollection).InsertOne(ctx, item)
	return item, err
}

func deleteMysteryItem(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.Database.Collection(MysteryCollection).DeleteOne(ctx, bson.M{"name": name})
	return err
}
