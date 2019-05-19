package controller

type VideoUrls map[string]MetadataType

type FormatType struct {
	Format    string
	URL       string
	filesize  int
	formatID  string
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
