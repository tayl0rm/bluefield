package bot

import (
	"fmt"
	"os"
	"os/signal"

	discord "github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var BotToken string

func Run() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	discord, err := discord.New("Bot " + BotToken)
	if err != nil {
		logger.Error("Error creating client.", zap.Error(err))
	}

	discord.AddHandler(newMessage)

	discord.Open()
	defer discord.Close()

	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
