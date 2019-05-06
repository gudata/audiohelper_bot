package youtube

type youtubeMetadataType struct {
	UploadDate         string      `json:"upload_date"`
	Protocol           string      `json:"protocol"`
	Extractor          string      `json:"extractor"`
	Series             interface{} `json:"series"`
	Format             string      `json:"format"`
	FormatNote         string      `json:"format_note"`
	Chapters           interface{} `json:"chapters"`
	Acodec             string      `json:"acodec"`
	LikeCount          int         `json:"like_count"`
	Duration           int         `json:"duration"`
	Fulltitle          string      `json:"fulltitle"`
	PlayerURL          string      `json:"player_url"`
	Quality            int         `json:"quality"`
	PlaylistIndex      interface{} `json:"playlist_index"`
	Album              interface{} `json:"album"`
	ViewCount          int         `json:"view_count"`
	Playlist           interface{} `json:"playlist"`
	Title              string      `json:"title"`
	Filename           string      `json:"_filename"`
	Creator            string      `json:"creator"`
	Ext                string      `json:"ext"`
	ID                 string      `json:"id"`
	DislikeCount       int         `json:"dislike_count"`
	AverageRating      interface{} `json:"average_rating"`
	Abr                int         `json:"abr"`
	UploaderURL        string      `json:"uploader_url"`
	Categories         []string    `json:"categories"`
	SeasonNumber       interface{} `json:"season_number"`
	Annotations        interface{} `json:"annotations"`
	WebpageURLBasename string      `json:"webpage_url_basename"`
	Filesize           int         `json:"filesize"`
	DisplayID          string      `json:"display_id"`
	AutomaticCaptions  struct {
	} `json:"automatic_captions"`
	Description        string        `json:"description"`
	Tags               []interface{} `json:"tags"`
	Track              string        `json:"track"`
	RequestedSubtitles interface{}   `json:"requested_subtitles"`
	StartTime          interface{}   `json:"start_time"`
	Tbr                float64       `json:"tbr"`
	DownloaderOptions  struct {
		HTTPChunkSize int `json:"http_chunk_size"`
	} `json:"downloader_options"`
	Uploader      string      `json:"uploader"`
	FormatID      string      `json:"format_id"`
	EpisodeNumber interface{} `json:"episode_number"`
	UploaderID    string      `json:"uploader_id"`
	Subtitles     struct {
	} `json:"subtitles"`
	ReleaseYear interface{} `json:"release_year"`
	HTTPHeaders struct {
		AcceptCharset  string `json:"Accept-Charset"`
		AcceptLanguage string `json:"Accept-Language"`
		AcceptEncoding string `json:"Accept-Encoding"`
		Accept         string `json:"Accept"`
		UserAgent      string `json:"User-Agent"`
	} `json:"http_headers"`
	Thumbnails []struct {
		URL string `json:"url"`
		ID  string `json:"id"`
	} `json:"thumbnails"`
	License      interface{} `json:"license"`
	Artist       string      `json:"artist"`
	URL          string      `json:"url"`
	ExtractorKey string      `json:"extractor_key"`
	ReleaseDate  interface{} `json:"release_date"`
	AltTitle     string      `json:"alt_title"`
	Thumbnail    string      `json:"thumbnail"`
	ChannelID    string      `json:"channel_id"`
	IsLive       interface{} `json:"is_live"`
	EndTime      interface{} `json:"end_time"`
	WebpageURL   string      `json:"webpage_url"`
	Formats      []struct {
		HTTPHeaders struct {
			AcceptCharset  string `json:"Accept-Charset"`
			AcceptLanguage string `json:"Accept-Language"`
			AcceptEncoding string `json:"Accept-Encoding"`
			Accept         string `json:"Accept"`
			UserAgent      string `json:"User-Agent"`
		} `json:"http_headers"`
		FormatNote        string  `json:"format_note"`
		Protocol          string  `json:"protocol"`
		Format            string  `json:"format"`
		URL               string  `json:"url"`
		Vcodec            string  `json:"vcodec"`
		Tbr               float64 `json:"tbr,omitempty"`
		Abr               int     `json:"abr,omitempty"`
		PlayerURL         string  `json:"player_url"`
		DownloaderOptions struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options,omitempty"`
		Ext        string `json:"ext"`
		Filesize   int    `json:"filesize"`
		FormatID   string `json:"format_id"`
		Quality    int    `json:"quality"`
		Acodec     string `json:"acodec"`
		Container  string `json:"container,omitempty"`
		Height     int    `json:"height,omitempty"`
		Width      int    `json:"width,omitempty"`
		Fps        int    `json:"fps,omitempty"`
		Resolution string `json:"resolution,omitempty"`
	} `json:"formats"`
	ChannelURL string `json:"channel_url"`
	Vcodec     string `json:"vcodec"`
	AgeLimit   int    `json:"age_limit"`
}
