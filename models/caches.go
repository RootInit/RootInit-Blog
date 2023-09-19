package models

type TagCache struct {
	List  []Tag
	IdMap map[int]*Tag
}

type CategoryCache struct {
	List  []Category
	IdMap map[int]*Category
}

type UserCache struct {
	List    []User
	IdMap   map[int]*User
	NameMap map[string]*User
}

type ArticleFileCache struct {
	IdList []int          // List of Ids
	IdMap  map[int]string // map file_names
}
