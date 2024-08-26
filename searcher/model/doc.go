package model

// IndexDoc 索引实体
type IndexDoc struct {
	Id       string                 `json:"id,omitempty"`
	Title    string                 `json:"title,omitempty"`
	Text     string                 `json:"text,omitempty"`
	Flag     string                 `json:"flag,omitempty"`
	Time     string                 `json:"time,omitempty"`
	Tags     string                 `json:"tags,omitempty"`
	Document map[string]interface{} `json:"document,omitempty"`
}

// StorageIndexDoc 文档对象
type StorageIndexDoc struct {
	*IndexDoc
	Keys []string `json:"keys,omitempty"`
}

type ResponseDoc struct {
	IndexDoc
	OriginalText  string   `json:"originalText,omitempty"`
	OriginalTitle string   `json:"originalTitle,omitempty"`
	Flag          string   `json:"flag,omitempty"`
	Time          string   `json:"time,omitempty"`
	Tags          string   `json:"tags,omitempty"`
	Score         int      `json:"score,omitempty"` //得分
	Keys          []string `json:"keys,omitempty"`
}

type RemoveIndexModel struct {
	Id string `json:"id,omitempty"`
}

type ResponseDocSort []ResponseDoc

func (r ResponseDocSort) Len() int {
	return len(r)
}

func (r ResponseDocSort) Less(i, j int) bool {
	return r[i].Score < r[j].Score
}

func (r ResponseDocSort) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}
