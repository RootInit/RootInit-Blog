package static_gen

import (
	"blog/models"
	"html/template"
	"io/fs"
	"log"
	"path"
	"path/filepath"
	"strings"
)

/* Holds resources to be passed to the StaticGen. */
type resources struct {
	tagCache         *models.TagCache
	categoryCache    *models.CategoryCache
	articleFileCache *models.ArticleFileCache
	// More in future
}

type StaticGen struct {
	config      *models.Config
	resources   *resources
	tmplIndex   *template.Template
	tmplArticle *template.Template
	cards       *cardCache
}

/* Returns a `resources` object to be passed to StaticGen Initializer. */
func GetResources(tagCache *models.TagCache, categoryCache *models.CategoryCache, aFileCache *models.ArticleFileCache) resources {
	r := resources{
		tagCache:         tagCache,
		categoryCache:    categoryCache,
		articleFileCache: aFileCache,
	}
	return r
}

/*
Initializes and returns an instance of StaticGenerator.
Loads templates required for site building.
*/
func Initialize(config *models.Config, data *resources) (StaticGen, error) {
	sG := StaticGen{
		config:    config,
		resources: data,
		cards:     newCardCache(),
	}
	var err error
	loadTemplates := func(files []string) (*template.Template, error) {
		var filePaths []string
		for _, file := range files {
			fullPath := path.Join(config.SourcePaths.Root, "templates", file)
			filePaths = append(filePaths, fullPath)
		}
		template, err := template.ParseFiles(filePaths...)
		return template, err
	}
	sG.tmplIndex, err = loadTemplates(config.Templates.Index)
	if err != nil {
		log.Fatal(`Failed to load Index templates. Error: `, err)
	}
	sG.tmplArticle, err = loadTemplates(config.Templates.Article)
	if err != nil {
		log.Fatal(`Failed to load Article templates. Error: `, err)
	}
	return sG, err
}

/*
Itterates through asset source path and processes files by extention
before writing to output directory defined in config
*/
func (sG StaticGen) BuildSiteAssets() error {
	assetSrcDir := path.Join(
		sG.config.SourcePaths.Root,
		`assets`,
	)
	assetOutDir := path.Join(
		sG.config.OutputPaths.Root,
		sG.config.OutputPaths.AssetsDir,
	)

	minify := sG.getMinifier()
	processAsset := func(srcPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileName := d.Name()
		// Skip folders and _name files (used for imports)
		if d.IsDir() || strings.HasPrefix(fileName, `_`) {
			return nil
		}
		relativePath := srcPath[len(assetSrcDir):]
		destPath := path.Join(assetOutDir, relativePath)
		switch filepath.Ext(fileName) {
		case `.css`:
			err = sG.processCss(minify, srcPath, destPath)
			return err
		case `.html`:
			err = sG.processHtml(minify, srcPath, destPath)
			return err
		case `.js`:
			err = sG.processJs(minify, srcPath, destPath)
			return err
		case `.svg`:
			err = sG.processSvg(minify, srcPath, destPath)
			return err
		default:
			err = copyFile(srcPath, destPath)
			return err
		}
	}
	err := filepath.WalkDir(assetSrcDir, processAsset)
	return err
}

/* Write a models.Article article page and store an article card to be used for indexes  */
func (sG StaticGen) AddArticle(article *models.Article) error {
	err := sG.writeArticlePage(article)
	if err != nil {
		log.Printf(`Error writing article page for articleId %d`, article.Id)
		return err
	}
	sG.addCardToCache(article)
	return nil
}

/*
Called when all articles have been added and StaticGenerator
cardCache is ready to be processed into the various indexes
*/
func (sG StaticGen) MakeIndexes() error {
	if err := sG.buildMainIndex(); err != nil {
		return err
	}
	if err := sG.buildTagIndexes(); err != nil {
		return err
	}

	err := sG.buildCategoryIndexes()
	return err
}
