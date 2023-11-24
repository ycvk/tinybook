package domain

type Article struct {
	ID       int64         `json:"id"`
	Title    string        `json:"title"`
	Content  string        `json:"content"`
	Abstract string        `json:"abstract"`
	Author   Author        `json:"author"`
	Status   ArticleStatus `json:"status"`
	Ctime    int64         `json:"ctime"`
	Utime    int64         `json:"utime"`
}

type ArticleVo struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
	Status  string `json:"status"`
	Ctime   string `json:"ctime"`
	Utime   string `json:"utime"`
}

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ArticleStatus uint8

const (
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)
