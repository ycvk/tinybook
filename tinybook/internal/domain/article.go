package domain

type Article struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  Author `json:"author"`
}

type Author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
