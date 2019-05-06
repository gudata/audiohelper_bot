package controller

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gudata/audiohelper_bot/packages/youtube"
	"github.com/syndtr/goleveldb/leveldb"
)

type ControllerType struct {
	bot *tgbotapi.BotAPI
	db  *leveldb.DB
}

func NewController(db *leveldb.DB) *ControllerType {
	controller := ControllerType{db: db}

	return &controller
}

func (controller *ControllerType) Formats(videoUrl string) map[string]string {
	youtube := youtube.NewYoutube(videoUrl, controller.db)
	return youtube.Formats()
}

func (controller *ControllerType) GetAudioUrl(videoUrl string, FormatID string) (string, error) {
	youtube := youtube.NewYoutube(videoUrl, controller.db)
	return youtube.GetAudioUrl(FormatID)
}

func (controller *ControllerType) GetMeta(videoUrl string) (map[string]string, error) {
	youtube := youtube.NewYoutube(videoUrl, controller.db)
	return youtube.GetMeta()
}
