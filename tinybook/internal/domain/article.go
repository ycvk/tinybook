package domain

type Article struct {
	ID      int64         `json:"id"`
	Title   string        `json:"title"`
	Content string        `json:"content"`
	Author  Author        `json:"author"`
	Status  ArticleStatus `json:"status"`
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
