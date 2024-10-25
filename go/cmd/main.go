package main

import bot "github.com/tayl0rm/bluefield/go/internal/bot"

func main() {
	bot.BotToken = ""
	bot.Run() // call the run function of bot/bot.go
}
