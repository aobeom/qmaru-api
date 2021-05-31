package models

const (
	CrondTimeTable   = "crond_time"
	DramaInfoTable   = "drama_info"
	MediaInfoTable   = "media_info"
	RadikoInfoTable  = "radiko_info"
	StTokenTable     = "st_token"
	StInfoTable      = "st_info"
	FanclubInfoTable = "fanclub_info"
	CinemaInfoTable  = "cinema_info"
)

// CommonModel 公用结构
type CommonModel struct {
	ID        int `json:"id" db:"serial;PRIMARY KEY"`
	CreatedAt int `json:"created_at" db:"integer;DEFAULT 0"`
	UpdatedAt int `json:"updated_at" db:"integer;DEFAULT 0"`
	Status    int `json:"status" db:"boolean;DEFAULT true"`
}

// CrondTime 定时任务更新时间
type CrondTime struct {
	ID   int    `json:"id" db:"serial;PRIMARY KEY"`
	Type string `json:"type" db:"varchar(30);DEFAULT ''"`
	Time string `json:"time" db:"varchar(19);DEFAULT ''"`
}

// DramaInfo 剧集结构
type DramaInfo struct {
	CommonModel
	Type   string      `json:"type" db:"varchar(30);DEFAULT ''"`
	URL    string      `json:"url" db:"varchar(255);DEFAULT ''"`
	Title  string      `json:"title" db:"varchar(100);DEFAULT ''"`
	Date   string      `json:"date" db:"varchar(19);DEFAULT ''"`
	DLURLs QDramaArray `json:"dlurls" db:"jsonb;DEFAULT '[]'"`
}

// MediaInfo 媒体结构
type MediaInfo struct {
	CommonModel
	Type    string      `json:"type" db:"varchar(30);DEFAULT ''"`
	Website string      `json:"website" db:"varchar(100);DEFAULT ''"`
	URL     string      `json:"url" db:"varchar(100);DEFAULT ''"`
	Source  QMediaArray `json:"source" db:"jsonb;DEFAULT '[]'"`
}

// RadikoInfo 电台结构
type RadikoInfo struct {
	CommonModel
	Name string `json:"name" db:"varchar(100);DEFAULT ''"`
	URL  string `json:"url" db:"varchar(100);DEFAULT ''"`
}

// StToken STChannel Token 结构
type StToken struct {
	ID    int    `json:"id" db:"serial;PRIMARY KEY"`
	Token string `json:"token" db:"varchar(255);DEFAULT ''"`
}

// StInfo STChannel 结构
type StInfo struct {
	CommonModel
	Title      string `json:"title" db:"varchar(1024);DEFAULT ''"`
	PictureURL string `json:"picture_url" db:"varchar(255);DEFAULT ''"`
	MediaURL   string `json:"media_url" db:"varchar(255);DEFAULT ''"`
	Date       string `json:"date" db:"varchar(19);DEFAULT ''"`
	Path       string `json:"path" db:"varchar(255);DEFAULT ''"`
}

// FanclubInfo 373 fan club
type FanclubInfo struct {
	CommonModel
	Title string `json:"title" db:"varchar(200);DEFAULT ''"`
	Type  string `json:"type" db:"varchar(50);DEFAULT ''"`
	Time  string `json:"time" db:"varchar(50);DEFAULT ''"`
	URL   string `json:"url" db:"varchar(200);DEFAULT ''"`
}

// CinemaInfo cinema sales
type CinemaInfo struct {
	CommonModel
	Title  string `json:"title" db:"varchar(100);DEFAULT ''"`
	Time   string `json:"time" db:"varchar(8);DEFAULT ''"`
	Type   int    `json:"type" db:"integer;DEFAULT 0"`
	Sales  int    `json:"sales" db:"integer;DEFAULT 0"`
	Income int    `json:"income" db:"integer;DEFAULT 0"`
}
