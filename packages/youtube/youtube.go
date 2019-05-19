package youtube

import (
	"encoding/json"
	"github.com/google/logger"
	"errors"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/alessio/shellescape.v1"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

// https://youtu.be/2O_r-d-HrjA
// /watch
var validMovie = regexp.MustCompile(`^https?://(www.)?(youtube.com|youtu.be)/`)
var audioOnly = regexp.MustCompile(`audio only`)

type YoutubeType struct {
	urlString string
	db        *leveldb.DB
}

func NewYoutube(url string) *YoutubeType {
	youtube := YoutubeType{
		urlString: "",
	}
	youtube.urlString = url

	return &youtube
}

func (youtube *YoutubeType) SetStorage(db *leveldb.DB) {
	youtube.db = db
}


func (youtube *YoutubeType) Detect() bool {
	matched := validMovie.MatchString(youtube.urlString)
	return matched
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func (youtube *YoutubeType) Download(filepath, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func (youtube *YoutubeType) getAudioMeta() (youtubeMetadataType, error) {
	var youtubeMetadata youtubeMetadataType

	var output []byte
	var err error

	output, err = youtube.db.Get([]byte(youtube.urlString), nil)
	if output == nil {
		command := exec.Command("/usr/local/bin/youtube-dl", "--print-json", "--skip-download", youtube.urlString)
		output, err = command.Output()
		if err != nil {
			return youtubeMetadataType{}, err
		}
	}

	err = json.Unmarshal([]byte(output), &youtubeMetadata)
	if err != nil {
		fmt.Println("error:", err)
	}

	youtube.db.Put([]byte(youtube.urlString), []byte(output), nil)

	return youtubeMetadata, nil
}

func (youtube *YoutubeType) ConvertToAudio(filePath, convertedFilePath string)  {
	command := exec.Command("/usr/bin/ffmpeg", `-i`, shellescape.Quote(filePath), `-c:a`, `mp3`, shellescape.Quote(convertedFilePath))
	println("/usr/bin/ffmpeg", `-i`, shellescape.Quote(filePath), `-c:a`, `mp3`, shellescape.Quote(convertedFilePath))
	output, err := command.Output()
	println(output)

	if err != nil {
		logger.Error(err)
	}
}


func (youtube *YoutubeType) Formats() map[string]string {
	formats := make(map[string]string)
	youtubeMetadata, err := youtube.getAudioMeta()

	if err != nil {
		return formats
	}

	for _, entry := range youtubeMetadata.Formats {
		if !audioOnly.MatchString(entry.Format) {
			continue
		}

		formats[entry.FormatID] = fmt.Sprintf("%s - %s", ByteCountSI(int64(entry.Filesize)), entry.FormatNote)
	}

	return formats
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func (youtube *YoutubeType) GetAudioURL(formatID string) (string, error) {
	youtubeMetadata, err := youtube.getAudioMeta()
	if err != nil {
		return "", err
	}

	for _, entry := range youtubeMetadata.Formats {
		if entry.FormatID == formatID {
			return entry.URL, nil
		}
	}

	return "", errors.New("Can't find the url")
}

// The meta needed for the app
func (youtube *YoutubeType) GetMeta() (map[string]string, error) {
	meta := make(map[string]string)
	youtubeMetadata, err := youtube.getAudioMeta()

	if err != nil {
		return meta, err
	}

	meta["title"] = youtubeMetadata.Title
	meta["filename"] = youtubeMetadata.Filename
	meta["duration"] = fmt.Sprintf("%d", youtubeMetadata.Duration)
	meta["id"] = youtubeMetadata.Filename
	meta["artist"] = youtubeMetadata.Artist
	return meta, nil
}
