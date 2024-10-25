package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	computeMetadata "cloud.google.com/go/compute/metadata"
	discord "github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	compute "google.golang.org/api/compute/v1"
)

const (
	zone      = "europe-west1-b"
	instance  = "-server"
	projectID = "ga-test-project-503ca"
)

func newMessage(discord *discord.Session, message *discord.MessageCreate) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	ctx := context.Background()
	game := "valheim"
	instanceName := game + instance

	if message.Author.ID == discord.State.User.ID { // Tell Bot ignore its own messages

	}
	// Create a new Compute Service client
	service, err := compute.NewService(ctx)
	if err != nil {
		logger.Error("Error creating client.", zap.Error(err))
	}

	switch {
	case strings.Contains(message.Content, "!valheim-start"):
		opr, err := service.Instances.Start(projectID, zone, instanceName).Do()
		if err != nil {
			logger.Error("Error creating client.", zap.Error(err))
		}
		serverIP, err := computeMetadata.ExternalIPWithContext(ctx)
		if err != nil {
			logger.Error("Error creating client.", zap.Error(err))
		}
		discord.ChannelMessageSend(message.ChannelID, serverIP)
		return waitForZoneOperation(service, projectID, zone, opr.Name)

	case strings.Contains(message.Content, "!valheim-down"):
		opr, err := service.Instances.Stop(projectID, zone, instanceName).Do()
		if err != nil {
			logger.Error("Error creating client.", zap.Error(err))
		}
		return waitForZoneOperation(service, projectID, zone, opr.Name)
	}

	return nil
}

func waitForZoneOperation(service *compute.Service, projectID, zone, operationName string) error {
	for {
		op, err := service.ZoneOperations.Get(projectID, zone, operationName).Do()
		if err != nil {
			return err
		}
		if op.Status == "DONE" {
			break
		}
		fmt.Printf("Waiting more 10 secs for the operation\n")
		time.Sleep(10 * time.Second)
	}
	return nil
}
