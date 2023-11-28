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
	ID         int64  `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	Content    string `json:"content,omitempty"`
	Abstract   string `json:"abstract,omitempty"`
	Author     string `json:"author,omitempty"`
	AuthorName string `json:"authorName,omitempty"`
	Status     string `json:"status,omitempty"`
	Ctime      string `json:"ctime,omitempty"`
	Utime      string `json:"utime,omitempty"`
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
