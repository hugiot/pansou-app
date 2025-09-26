package pansou

type MergedLinks map[string][]MergedLink

type HealthRequest struct {
}

type HealthResponse struct {
	Status         string   `json:"status"`
	PluginsEnabled bool     `json:"plugins_enabled"`
	Channels       []string `json:"channels"`
	ChannelsCount  int      `json:"channels_count"`
	PluginCount    int      `json:"plugin_count"`
	Plugins        []string `json:"plugins"`
}

type SearchResponse struct {
	Total        int            `json:"total" sonic:"total"`
	Results      []SearchResult `json:"results" sonic:"results"`
	MergedByType MergedLinks    `json:"merged_by_type" sonic:"merged_by_type"`
}

type SearchResult struct {
	MessageID string   `json:"message_id" sonic:"message_id"`
	UniqueID  string   `json:"unique_id" sonic:"unique_id"` // 全局唯一ID
	Channel   string   `json:"channel" sonic:"channel"`
	Datetime  string   `json:"datetime" sonic:"datetime"`
	Title     string   `json:"title" sonic:"title"`
	Content   string   `json:"content" sonic:"content"`
	Links     []Link   `json:"links" sonic:"links"`
	Tags      []string `json:"tags,omitempty" sonic:"tags,omitempty"`
	Images    []string `json:"images,omitempty" sonic:"images,omitempty"` // TG消息中的图片链接
}

type Link struct {
	Type     string `json:"type" sonic:"type"`
	URL      string `json:"url" sonic:"url"`
	Password string `json:"password" sonic:"password"`
}

// MergedLink 合并后的网盘链接
type MergedLink struct {
	URL      string   `json:"url" sonic:"url"`
	Password string   `json:"password" sonic:"password"`
	Note     string   `json:"note" sonic:"note"`
	Datetime string   `json:"datetime" sonic:"datetime"`
	Source   string   `json:"source,omitempty" sonic:"source,omitempty"` // 数据来源：tg:频道名 或 plugin:插件名
	Images   []string `json:"images,omitempty" sonic:"images,omitempty"` // TG消息中的图片链接
}
