package main

import (
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var adminID int64 = 123456789 // Замените на настоящий Telegram ID админа

type Expense struct {
	Item     string
	Price    float64
	Quantity int
}

type Service struct {
	Name  string
	Amount float64
}

var expenses []Expense
var services []Service

func main() {
	bot, err := tgbotapi.NewBotAPI("YOUR_BOT_API_KEY")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать! Используйте кнопки для ввода данных.")
			bot.Send(msg)
			showMainMenu(bot, update.Message.Chat.ID)
		default:
			handleInput(bot, update.Message)
		}
	}
}

func showMainMenu(bot *tgbotapi.BotAPI, chatID int64) {
	var keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Добавить товар"),
			tgbotapi.NewKeyboardButton("Добавить услугу"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отправить отчет"),
		),
	)
	msg := tgbotapi.NewMessage(chatID, "Выберите действие:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
}

func handleInput(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	text := message.Text
	chatID := message.Chat.ID

	switch text {
	case "Добавить товар":
		msg := tgbotapi.NewMessage(chatID, "Введите название товара, цену и количество через запятую (например, 'Товар, 100.50, 2'):")
		bot.Send(msg)
	case "Добавить услугу":
		msg := tgbotapi.NewMessage(chatID, "Введите название услуги и полученную сумму через запятую (например, 'Услуга, 150.75'):")
		bot.Send(msg)
	case "Отправить отчет":
		sendReport(bot, chatID)
	default:
		if strings.Contains(text, ",") {
			handleDataEntry(bot, message)
		} else {
			msg := tgbotapi.NewMessage(chatID, "Неизвестная команда. Пожалуйста, используйте кнопки.")
			bot.Send(msg)
		}
	}
}

func handleDataEntry(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text
	parts := strings.Split(text, ",")

	if len(parts) == 3 {
		// Обработка товара
		item := strings.TrimSpace(parts[0])
		price, err1 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		quantity, err2 := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err1 != nil || err2 != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка ввода. Попробуйте еще раз.")
			bot.Send(msg)
			return
		}
		expenses = append(expenses, Expense{Item: item, Price: price, Quantity: quantity})
		msg := tgbotapi.NewMessage(chatID, "Товар добавлен.")
		bot.Send(msg)
	} else if len(parts) == 2 {
		// Обработка услуги
		name := strings.TrimSpace(parts[0])
		amount, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка ввода. Попробуйте еще раз.")
			bot.Send(msg)
			return
		}
		services = append(services, Service{Name: name, Amount: amount})
		msg := tgbotapi.NewMessage(chatID, "Услуга добавлена.")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Ошибка ввода. Попробуйте еще раз.")
		bot.Send(msg)
	}
}

func sendReport(bot *tgbotapi.BotAPI, chatID int64) {
	expenseReport := "Отчет по расходам:\n"
	totalExpenses := 0.0
	for _, expense := range expenses {
		expenseReport += expense.Item + ": " + strconv.Itoa(expense.Quantity) + " x " + strconv.FormatFloat(expense.Price, 'f', 2, 64) + " = " + strconv.FormatFloat(expense.Price*float64(expense.Quantity), 'f', 2, 64) + "\n"
		totalExpenses += expense.Price * float64(expense.Quantity)
	}

	serviceReport := "Отчет по доходам:\n"
	totalIncome := 0.0
	for _, service := range services {
		serviceReport += service.Name + ": " + strconv.FormatFloat(service.Amount, 'f', 2, 64) + "\n"
		totalIncome += service.Amount
	}

	report := expenseReport + "Итого расходов: " + strconv.FormatFloat(totalExpenses, 'f', 2, 64) + "\n\n" + serviceReport + "Итого доходов: " + strconv.FormatFloat(totalIncome, 'f', 2, 64) + "\n"

	msg := tgbotapi.NewMessage(adminID, report)
	bot.Send(msg)

	msg = tgbotapi.NewMessage(chatID, "Отчет отправлен админу.")
	bot.Send(msg)

	// Очистка данных после отправки отчета
	expenses = []Expense{}
	services = []Service{}
}
