package main

import (
	"fmt"

	"github.com/PancyStudios/PancyBotGo/pkg/config"
	"github.com/PancyStudios/PancyBotGo/pkg/database"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	cfg, _ := config.Load()
	db, _ := database.Init(cfg.MongoDBURL, cfg.DBName)
	defer db.Disconnect()
	database.InitGlobalDataManagers(db)

	doc, err := database.GlobalGuildDM.Get(bson.M{"id": "1029032049853608026"})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Prefix: '%s'\n", doc.Configuration.Prefix)
}
