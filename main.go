package main

import (
	"blog/database"
	"blog/models"
	"blog/static_gen"
	"log"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	var err error
	// Load Config
	config, err := GetConfig("config.yaml")
	if err != nil {
		log.Fatal(`Unable to load "config.yaml" Error: `, err)
	}

	// Initialize database
	db, err := database.InitializeDB(config.SourcePaths.DbPath)
	if err != nil {
		log.Fatal(`Failed to initialize database. Error: `, err)
	}

	// Initialize static generator
	resources := static_gen.GetResources(
		db.TagCache,
		db.CategoryCache,
		db.ArticleFileCache,
	)
	sG, err := static_gen.Initialize(&config, &resources)
	if err != nil {
		log.Fatal(`Failed to initialize static site generator. Error: `, err)
	}

	// Compile site assets
	if err := sG.BuildSiteAssets(); err != nil {
		log.Fatal(`Failed to build site assets. Error: `, err)
	}

	// Compile Article
	aI := db.GetArticleIttr()
	for aI.Next() {
		if err := sG.AddArticle(&aI.Article); err != nil {
			log.Fatal(`Failed to build Article. Error: `, err)
		}
	}
	// Compile Indices
	err = sG.MakeIndexes()
	if err != nil {
		log.Fatal(err)
	}

}

func GetConfig(configPath string) (models.Config, error) {
	var config models.Config
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return config, err
	}
	// Remove leading and trailing `/` from RootUrl
	config.RootUrl = strings.Trim(config.RootUrl, `/`)

	// Function to clean relative paths
	cleanRel := func(s string) string {
		s = strings.Trim(s, `/\`)
		s = path.Clean(s)
		return s
	}
	// Clean paths
	config.SourcePaths.Root = path.Clean(config.SourcePaths.Root)
	config.SourcePaths.DbPath = path.Clean(config.SourcePaths.DbPath)
	config.OutputPaths.Root = path.Clean(config.OutputPaths.Root)
	config.OutputPaths.IndexPageDir = cleanRel(config.OutputPaths.IndexPageDir)
	config.OutputPaths.ArticleDir = cleanRel(config.OutputPaths.ArticleDir)
	return config, err
}
