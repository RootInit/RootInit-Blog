package models

import (
	"time"
)

type Config struct {
	SiteName  string `yaml:"Site_Name"`
	RootUrl   string `yaml:"Root_URL"`
	IndexDesc string `yaml:"Index_Description"`
	IndexSize int    `yaml:"Index_Size"`

	SourcePaths struct {
		DbPath string `yaml:"Database_Path"`
		Root   string `yaml:"Root_Path"`
	} `yaml:"Source_Paths"`

	OutputPaths struct {
		Root         string `yaml:"Root_Path"`
		IndexPageDir string `yaml:"Index_Page_Directory"`
		ArticleDir   string `yaml:"Article_Directory"`
		AssetsDir    string `yaml:"Assets_Directory"`
		TagIndexDir  string `yaml:"Tag_Index_Directory"`
		CatIndexDir  string `yaml:"Category_Index_Directory"`
	} `yaml:"Output_Paths"`

	OutputOpts struct {
		PageFileExt string `yaml:"Page_File_Ext"`
		DateFormat  string `yaml:"Date_Format"`
	} `yaml:"Output_Settings"`

	Templates struct {
		Index   []string `yaml:"Index"`
		Article []string `yaml:"Article"`
	} `yaml:"Templates"`

	SubStrings map[string]string `yaml:"Strings"`
}

type Tag struct {
	Id   int
	Name string
	Icon string
}

type Category struct {
	Id   int
	Name string
}

type User struct {
	Id       int
	Username string
	Password string
	Email    string
}

type Article struct {
	Id          int
	Title       string
	UrlTitle    string
	Description string
	Date        time.Time
	Body        string
	Thumbnail   string
	Category    Category
	Tags        []Tag
	Comments    []Comment
}

type Comment struct {
	Id        int
	ArticleId int
	Author    User
	Date      time.Time
	Body      string
	Replies   []Comment
}
