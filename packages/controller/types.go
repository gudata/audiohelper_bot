package controller

type VideoUrls map[string]MetadataType

type FormatType struct {
	Format    string
	Url       string
	filesize  int
	format_id string
}

// https://github.com/syndtr/goleveldb
type MetadataType struct {
	Fulltitle   string
	Filename    string
	Description string
	Track       string
	Formats     []FormatType
	Raw         string
}
