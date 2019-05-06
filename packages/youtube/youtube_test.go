package youtube

import (
	"fmt"
	"testing"
)

func TestGetAudioMeta(t *testing.T) {
	youtube := NewYoutube("https://www.youtube.com/watch?v=kkJC8p48g6g")
	youtubeMetadata, _ := youtube.getAudioMeta()
	if youtubeMetadata.Title != "A. VIVALDI: «Filiae maestae Jerusalem» RV 638 [II.Sileant Zephyri], Ph.Jaroussky/Ensemble Artaserse" {
		t.Errorf("Title should be correct, got: %+v", youtubeMetadata.Title)
	}
}

func TestFormats(t *testing.T) {
	youtube := NewYoutube("https://www.youtube.com/watch?v=kkJC8p48g6g")

	formats := youtube.Formats()

	fmt.Printf("add %+v", formats)

	if len(formats) != 5 {
		t.Errorf("Formats should be 5, got %d: ", len(formats))
	}

	if formats["249"] != "2.5 MB - DASH audio" {
		t.Errorf("249 format should be '2.5 MB - DASH audio', got %s", formats["249"])
	}
}
