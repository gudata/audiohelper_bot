package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/logger"
	"github.com/gudata/audiohelper_bot/packages/config"
	"github.com/gudata/audiohelper_bot/packages/controller"
	s "github.com/gudata/audiohelper_bot/packages/storage"
	"github.com/gudata/audiohelper_bot/packages/youtube"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"log"
	"net/url"
	"os"
	"sort"
)

// A data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value string
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]string) PairList {
	println(len(m))

	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	return p
}

func sendDownloadOptions(bot *tgbotapi.BotAPI, update tgbotapi.Update, videoUrl string) {
	formats := controller.NewController(db).Formats(videoUrl)

	sortedFormats := sortMapByValue(formats)

	chatID := update.Message.Chat.ID

	var row []tgbotapi.InlineKeyboardButton
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, pair := range sortedFormats {
		label := pair.Value
		key := pair.Key

		row = make([]tgbotapi.InlineKeyboardButton, 1)
		row[0] = tgbotapi.NewInlineKeyboardButtonData(label, key)
		rows = append(rows, row)
	}

	msg := tgbotapi.NewMessage(chatID, "Choose format")
	msg.DisableWebPagePreview = true
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
	bot.Send(msg)
}

func sendDownloadMessage(bot *tgbotapi.BotAPI, chatID int64, videoUrl, formatID string) {
	controller := controller.NewController(db)
	meta, _ := controller.GetMeta(videoUrl)

	audioUrl, _ := controller.GetAudioUrl(videoUrl, formatID)

	youtube := youtube.NewYoutube(videoUrl, db)
	file_path := storage.DownloadPath(meta, formatID)

	storage.EnsureFolder(file_path)

	logger.Info(file_path)

	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		youtube.Download(file_path, audioUrl)
		msg := tgbotapi.NewMessage(chatID, "File downloaded - now telegraming it, please wait...")
		msg.DisableWebPagePreview = true
		bot.Send(msg)
	}

	audioMessage := tgbotapi.NewAudioUpload(chatID, file_path) // or NewAudioShare(chatID int64, fileID string)
	bot.Send(audioMessage)
}

var db *leveldb.DB
var err error
var storage *s.StorageType

func main() {
	// https://www.progville.com/go/bolt-embedded-db-golang/
	dbFile := "audio-helper-telegram-bot-database"
	db, err = leveldb.OpenFile(dbFile, nil)
	if errors.IsCorrupted(err) {
		fmt.Println("ErrCorrupted", err)
		leveldb.RecoverFile(dbFile, nil)
		panic("Database recovered. Restart the application.")
	}

	if err != nil {
		panic(err)
	}

	defer db.Close()

	config := config.Config()
	defer config.InitLogging().Close()
	logger.Info("audio-helper Start")

	storage = s.NewStorage(config.OutputFolder, config.Debug)
	storage.CreateOutputFolder()

	bot, err := tgbotapi.NewBotAPI(config.Secret)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 120

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.CallbackQuery != nil {
			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Download started..."))

			oldMessage := update.CallbackQuery.Message.ReplyToMessage
			videoUrl := oldMessage.Text
			formatID := update.CallbackQuery.Data

			controller := controller.NewController(db)
			audioUrl, err := controller.GetAudioUrl(videoUrl, formatID)
			meta, err := controller.GetMeta(videoUrl)

			if err != nil {
				msg := tgbotapi.NewMessage(oldMessage.Chat.ID, "Can't find the audio url for the this video.")
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				bot.Send(msg)
				continue
			}

			// https://wiki.videolan.org/Documentation:IOS/#x-callback-url
			streamURL, _ := url.Parse("vlc-x-callback://x-callback-url/stream")
			streamURL.Path = "/stream"
			values := url.Values{}
			values.Add("url", audioUrl)

			downloadURL, _ := url.Parse("vlc-x-callback://x-callback-url/download")
			downloadURL.Path = "/download"
			values = url.Values{}
			values.Add("url", audioUrl)
			values.Add("filename", meta["filename"])

			msg := tgbotapi.NewMessage(oldMessage.Chat.ID, fmt.Sprintf("Psst - [The URL](%s) if you want to [Download](%s) or [Stream](%s) in VLC", audioUrl, streamURL.String(), downloadURL.String()))
			msg.ReplyToMessageID = oldMessage.MessageID
			msg.ParseMode = "markdown"
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)

			go sendDownloadMessage(bot, oldMessage.Chat.ID, videoUrl, formatID)

			continue
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		videoUrl := update.Message.Text
		if youtube := youtube.NewYoutube(videoUrl, db); youtube.Detect() {
			go sendDownloadOptions(bot, update, videoUrl)
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "Paste an URL for a youtube video or type /sayhi or /status."
			case "sayhi":
				msg.Text = "Hi :)"
			case "status":
				msg.Text = "I'm ok."
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Give me youtube URL!")
			msg.DisableWebPagePreview = true
			msg.ReplyToMessageID = update.Message.MessageID

			switch update.Message.Text {
			case "close":
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			}

			bot.Send(msg)
		}
	}
}
