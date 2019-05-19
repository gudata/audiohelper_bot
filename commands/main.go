package main

import (
  "fmt"
  "log"
  "net/url"
  "os"
  "sort"
  "strconv"

  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
  "github.com/google/logger"
  "github.com/gudata/audiohelper_bot/packages/config"
  "github.com/gudata/audiohelper_bot/packages/controller"
  s "github.com/gudata/audiohelper_bot/packages/storage"
  "github.com/gudata/audiohelper_bot/packages/youtube"
  "github.com/syndtr/goleveldb/leveldb"
  "github.com/syndtr/goleveldb/leveldb/errors"
)

// Pair is a data structure to hold a key/value pair.
type Pair struct {
  Key   string
  Value string
}

// PairList is a slice of Pairs that implements sort.Interface to sort by Value.
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

func sendDownloadOptions(bot *tgbotapi.BotAPI, update tgbotapi.Update, videoURL string) {
  formats := controller.NewController(db).Formats(videoURL)

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

func sendDownloadMessage(bot *tgbotapi.BotAPI, chatID int64, videoURL, formatID string) {
  controller := controller.NewController(db)
  meta, _ := controller.GetMeta(videoURL)

  audioURL, _ := controller.GetAudioURL(videoURL, formatID)

  youtube := youtube.NewYoutube(videoURL)
  youtube.SetStorage(db)
  filePath := storage.DownloadPath(meta, formatID)

  storage.EnsureFolder(filePath)

  convertedFilePath := storage.ConvertedDownloadPath(filePath)

  logger.Info(convertedFilePath)

  if _, err := os.Stat(filePath); os.IsNotExist(err) {
    youtube.Download(filePath, audioURL)

    msg := tgbotapi.NewMessage(chatID, "File downloaded - now converting it music format, please wait...")
    msg.DisableWebPagePreview = true
    bot.Send(msg)

    youtube.ConvertToAudio(filePath, convertedFilePath)
    msg = tgbotapi.NewMessage(chatID, "File donverted - now telegraming it, please wait...")
    msg.DisableWebPagePreview = true
    bot.Send(msg)
  }

  audioMessage := tgbotapi.NewAudioUpload(chatID, filePath) // or NewAudioShare(chatID int64, fileID string)

  if duration, err := strconv.Atoi(meta["duration"]); err == nil {
    audioMessage.Duration = duration
  }

  audioMessage.Title = meta["title"]
  audioMessage.Performer = meta["artist"]
  // msg.MimeType = "audio/mpeg"
  // msg.FileSize = 688
  _, err := bot.Send(audioMessage)
  if err != nil {
    logger.Error(err)
  }
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

  bot.Debug = false

  log.Printf("Authorized on account %s", bot.Self.UserName)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 120

  updates, _ := bot.GetUpdatesChan(u)

  for update := range updates {

    if update.CallbackQuery != nil {
      bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Download started..."))

      oldMessage := update.CallbackQuery.Message.ReplyToMessage
      videoURL := oldMessage.Text
      formatID := update.CallbackQuery.Data

      controller := controller.NewController(db)
      audioURL, _ := controller.GetAudioURL(videoURL, formatID)
      meta, err := controller.GetMeta(videoURL)

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
      values.Add("url", audioURL)
      streamURL.RawQuery = values.Encode()

      downloadURL, _ := url.Parse("vlc-x-callback://x-callback-url/download")
      downloadURL.Path = "/download"
      values = url.Values{}
      values.Add("url", audioURL)
      values.Add("filename", meta["filename"])
      downloadURL.RawQuery = values.Encode()

      messageWithLinks := fmt.Sprintf("Psst - [The URL](%s) if you want to [Download](%s) or [Stream](%s) in VLC", audioURL, streamURL.String(), downloadURL.String())

      msg := tgbotapi.NewMessage(oldMessage.Chat.ID, messageWithLinks)

      msg.ReplyToMessageID = oldMessage.MessageID
      msg.ParseMode = "markdown"
      msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
      bot.Send(msg)

      go sendDownloadMessage(bot, oldMessage.Chat.ID, videoURL, formatID)

      continue
    }

    if update.Message == nil { // ignore any non-Message Updates
      continue
    }

    videoURL := update.Message.Text
    if youtube := youtube.NewYoutube(videoURL); youtube.Detect() {
      go sendDownloadOptions(bot, update, videoURL)
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
