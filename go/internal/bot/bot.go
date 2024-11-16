package bot

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"util/config"
	"go.uber.org/zap"
	compute "google.golang.org/api/compute/v1"
)

func MessageHandler() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	ctx := context.Background()
	config, err := util.LoadConfig(".")
	if err != nil {
		logger.Fatal("cannot load config:", zap.Error(err))
	}

	bot, err := discordgo.New(config.BotToken)
	if err != nil {
		logger.Fatal("Error creating Discord session.", zap.Error(err))
		return
	}

	bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		service, err := compute.NewService(ctx)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Something went wrong! @ekc0_")
			logger.Fatal("Error creating client.", zap.Error(err))
		}

		switch {
		case strings.Contains(m.Content, "!valheim-start"):
			game := "valheim"
			instanceName := game + config.Instance
			service.Instances.Start(config.ProjectID, config.Zone, instanceName).Do()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error starting Server! @ekc0_")
				logger.Error("Error starting Server.", zap.Error(err))
			}

			instance, err := service.Instances.Get(config.ProjectID, config.Zone, instanceName).Context(ctx).Do()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error getting Server IP! @ekc0_")
				logger.Error("Error getting Server IP", zap.Error(err))
			}

			var externalIP string
			for _, networkInterface := range instance.NetworkInterfaces {
				for _, accessConfig := range networkInterface.AccessConfigs {
					if accessConfig.Name == "External NAT" {
						externalIP = accessConfig.NatIP
						break
					}
				}
			}

			s.ChannelMessageSend(m.ChannelID, externalIP)

		case strings.Contains(m.Content, "!valheim-down"):
			game := "valheim"
			instanceName := game + config.Instance
			service.Instances.Stop(config.ProjectID, config.Zone, instanceName).Do()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Error stopping Server! @ekc0_")
				logger.Error("Error stopping Server", zap.Error(err))
			}
		}
	})

	err = bot.Open()
	if err != nil {
		logger.Error("Error openning connection", zap.Error(err))
		return
	}

	logger.Info("Bot is running. Press CTRL+C to exit.")
	select {} // Keep the program running
}
