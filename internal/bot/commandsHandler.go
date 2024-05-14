package bot

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tgholiday "git.foxminded.ua/foxstudent107051/tgholiday"
	tgbotapi "github.com/crocone/telegram-bot-api"
	"main.go/internal/dataProvider"
	"main.go/internal/logger"
	"main.go/internal/openWeather"
	"main.go/internal/openWeather/model"
)

const (
	CommandWeather     = "weather"
	CommandLinks       = "links"
	CommandAbout       = "about"
	CommandHolidays    = "holidays"
	CommandStart       = "start"
	CommandHelp        = "help"
	CommandUnSubscribe = "unsubscribe"

	ParseModeMarkdown   = "markdown"
	ParseModeMarkdownV2 = "MarkdownV2"
)

var countries = map[string]string{
	"🇺🇦": "UA",
	"🇺🇸": "US",
	"🇲🇨": "MC",
	"🇲🇪": "ME",
	"🇪🇸": "ES",
	"🇵🇱": "PL",
}

var YesNo = map[string]bool{
	"YES": true,
	"NO":  false,
}

type SLogger interface {
	GetSLog() *slog.Logger
}

type SLog struct{}

func (L SLog) GetSLog() *slog.Logger { return logger.GetLogger() }

type button struct {
	name string
	data string
}

var sL = &SLog{}
var sLogger = new(SLog).GetSLog()

func HandleWeatherSubscription() {
	subscriptions, _ := openWeather.GetAllSubscriptions()
	for _, subscription := range subscriptions {
		user, _ := openWeather.GetUserById(subscription.UserID)
		weatherSubs, _ := openWeather.GetWeatherSubscriptionById(subscription.WeatherSubscriptionID)

		now := time.Now()

		notificationTime, _ := time.Parse("2006-01-02 15:04:05", weatherSubs[0].NotificationTime)
		if notificationTime.Before(now) {
			notificationTime = notificationTime.Add(24 * time.Hour)
		}
		duration := notificationTime.Sub(now)
		ticker := time.NewTicker(duration)

		go func() {
			for {
				select {
				case <-ticker.C:
					msg := tgbotapi.NewMessage(user.ChatId, prepareWeatherResponse(openWeather.GetWeather(user.Lat, user.Lon)))
					msg.ParseMode = ParseModeMarkdownV2
					sendMessage(msg)
					notificationTime = notificationTime.Add(24 * time.Hour)
					ticker = time.NewTicker(notificationTime.Sub(now))
				}
			}
		}()
	}
}

func HandleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallbacks(update, updates)
		} else if update.Message.IsCommand() {
			handleCommands(update, updates)
		} else {
			handleMessage(update)
		}
	}
}

func handleCallbacks(update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	data := update.CallbackQuery.Data
	chatID := update.CallbackQuery.Message.Chat.ID

	switch data {
	case CommandAbout:
		executeAboutCommand(chatID)
	case CommandLinks:
		executeLinksCommand(chatID)
	case CommandWeather:
		executeWeatherCommand(chatID, updates)
	case CommandHolidays:
		executeHolidaysCommand(chatID, updates)
	}
}

func handleCommands(update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	command := update.Message.Command()
	chatID := update.Message.Chat.ID

	switch command {
	case CommandStart:
		executeHolidaysCommand(chatID, updates)
	case CommandHelp:
		executeHelpCommand(chatID)
	case CommandWeather:
		executeWeatherCommand(chatID, updates)
	case CommandLinks:
		executeLinksCommand(chatID)
	case CommandAbout:
		executeAboutCommand(chatID)
	case CommandUnSubscribe:
		executeUnSubscribeCommand(chatID, update.Message.Chat.FirstName)
	default:
		UnknownCommand(chatID)
		executeHelpCommand(chatID)
	}
}

func executeUnSubscribeCommand(id int64, name string) {
	deletedCount := openWeather.UnSubscribe(id, name)

	if deletedCount > 0 {
		msg := tgbotapi.NewMessage(id, "You subscription has been declined")
		sendMessage(msg)
	}
}

func executeWeatherCommand(chatID int64, updates tgbotapi.UpdatesChannel) {
	msg := tgbotapi.NewMessage(chatID, "Please provide Longitude")
	msg.ParseMode = ParseModeMarkdown
	sendMessage(msg)

	lon := waitForUserResponse(updates)
	msg = tgbotapi.NewMessage(chatID, "Please provide Latitude")
	sendMessage(msg)

	lat := waitForUserResponse(updates)
	fLat, _ := strconv.ParseFloat(lat, 64)
	fLon, _ := strconv.ParseFloat(lon, 64)
	msg = tgbotapi.NewMessage(chatID, prepareWeatherResponse(openWeather.GetWeather(fLat, fLon)))
	msg.ParseMode = ParseModeMarkdownV2
	sendMessage(msg)

	msg = tgbotapi.NewMessage(chatID, "Do you want to receive everyday notification about weather?")
	msg.ReplyMarkup = YesNoKeyMarkup()
	msg.ParseMode = ParseModeMarkdown
	sendMessage(msg)
	for update := range updates {
		user := model.Users{
			ChatId: update.Message.Chat.ID,
			Name:   update.Message.Chat.FirstName,
			Lat:    fLat,
			Lon:    fLon,
		}

		if YesNo[update.Message.Text] {
			openWeather.Subscribe(user)
			msg = tgbotapi.NewMessage(chatID, "Thanks for the subscription")
			sendMessage(msg)
		} else {
			msg = tgbotapi.NewMessage(chatID, "See you next time")
			sendMessage(msg)
		}

		break
	}

}

func executeHolidaysCommand(chatID int64, updates tgbotapi.UpdatesChannel) {
	msg := tgbotapi.NewMessage(chatID, "Choose an country")
	msg.ReplyMarkup = flagsIcon()
	msg.ParseMode = ParseModeMarkdown
	sendMessage(msg)

	ExecuteTodayHolidayCommand(chatID, waitForUserResponse(updates))
}

func executeHelpCommand(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Choose an action")
	msg.ReplyMarkup = startMenu()
	msg.ParseMode = ParseModeMarkdown
	sendMessage(msg)
}

func executeLinksCommand(chatID int64) {
	for name, value := range dataProvider.GetUserSocialLinks() {
		text := fmt.Sprintf("%s : %s", name, value)
		msg := tgbotapi.NewMessage(chatID, text)

		sendMessage(msg)
	}
}

func UnknownCommand(chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Unknown Command")
	sendMessage(msg)
}

func sendMessage(msg tgbotapi.Chattable) {
	if _, err := Bot.Send(msg); err != nil {
		sLogger.Error("Send message error", "err", err)
	}
}

func waitForUserResponse(updates tgbotapi.UpdatesChannel) string {
	for update := range updates {
		if update.Message.IsCommand() {
			handleCommands(update, updates)
			break
		}

		return update.Message.Text
	}

	return ""
}

func startMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{
			name: "About Me",
			data: CommandAbout,
		},
		{
			name: "Social Links",
			data: CommandLinks,
		},
		{
			name: "Weather",
			data: CommandWeather,
		},
		{
			name: "Holiday",
			data: CommandHolidays,
		},
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, len(states))
	for index, state := range states {
		buttons[index] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(state.name, state.data))
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func flagsIcon() tgbotapi.ReplyKeyboardMarkup {
	flags := tgbotapi.NewReplyKeyboard()
	row := make([]tgbotapi.KeyboardButton, len(countries))

	counter := 0
	for flag, _ := range countries {
		row[counter] = tgbotapi.NewKeyboardButton(flag)
		counter++
	}

	flags.Keyboard = append(flags.Keyboard, row)

	return flags
}

func YesNoKeyMarkup() tgbotapi.ReplyKeyboardMarkup {
	buttons := tgbotapi.NewReplyKeyboard()
	row := make([]tgbotapi.KeyboardButton, len(YesNo))

	counter := 0
	for button, _ := range YesNo {
		row[counter] = tgbotapi.NewKeyboardButton(button)
		counter++
	}

	buttons.Keyboard = append(buttons.Keyboard, row)

	return buttons
}

func executeAboutCommand(chatID int64) {
	for name, value := range dataProvider.GetAboutMe() {
		text := fmt.Sprintf("%s : %s", name, value)
		msg := tgbotapi.NewMessage(chatID, text)

		sendMessage(msg)
	}
}

func handleMessage(update tgbotapi.Update) {
	//TODO: For future implementation
}

func ExecuteTodayHolidayCommand(chatID int64, countryFlag string) {
	holidays := tgholiday.GetTodayHolidays(countries[countryFlag])

	if len(holidays) == 0 {
		text := fmt.Sprintf("Today in %s is no holidays:\n", countryFlag)
		msg := tgbotapi.NewMessage(chatID, text)
		sendMessage(msg)
	} else {
		text := fmt.Sprintf("Today in %s %v holidays:\n", countryFlag, len(holidays))
		for _, value := range holidays {
			text += fmt.Sprintf(" %s: %s\n", value.Location, value.Name)
		}

		msg := tgbotapi.NewMessage(chatID, text)
		sendMessage(msg)
	}
}

func prepareWeatherResponse(weatherResponse openWeather.WeatherResponse) string {
	return fmt.Sprintf(
		"*Your location: City: _%v_* \n"+
			"*Longitude:* _%v_*, Latitude:* _%v_\n"+
			"*Precipitation:* _%v_\n"+
			"*Description:* _%v_\n"+
			"*Current Temperature:* _%v_\n"+
			"*Max Temperature:* _%v_, *Min Temperature:* _%v_\n"+
			"*Wind speed:* _%v_",
		weatherResponse.Name,
		weatherResponse.Coord.Lon,
		weatherResponse.Coord.Lat,
		weatherResponse.Weather[0].Main,
		weatherResponse.Weather[0].Description,
		escapeFloatValue(weatherResponse.Main.Temp),
		escapeFloatValue(weatherResponse.Main.TempMax),
		escapeFloatValue(weatherResponse.Main.TempMin),
		escapeFloatValue(weatherResponse.Wind.Speed),
	)
}

func escapeFloatValue(value float64) string {
	newValue := fmt.Sprintf("%.2f", value)

	return strings.ReplaceAll(newValue, ".", "\\.")
}
