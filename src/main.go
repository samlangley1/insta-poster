package main

import (
	"go-insta/config"
	"go-insta/filesystem"
	"go-insta/instagram"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Info("failed to fetch env file: %s", err)
	}
}

func main() {
	// Generate new config object from config package
	conf := config.New()

	// Get random image file from image directory
	image, fileName, err := filesystem.GetRandomContent(conf.Filesystem.ImageDirectory)
	if err != nil {
		slog.Info("failed to fetch image content: %s", err)
	}

	// Log into instagram account and create session
	slog.Info("logging into: " + conf.Instagram.Username)
	sess, err := instagram.CreateSession(conf.Instagram.Username, conf.Instagram.Password)
	if err != nil {
		slog.Info("failed to log into account: %s", err)
	}

	// Set proxy settings if provided
	if len(conf.Network.ProxyAddress) > 0 {
		err = sess.SetProxy(conf.Network.ProxyAddress, true, true)
		if err != nil {
			slog.Info("failed to set proxy: %s", err)
		}
	}

	// Post random image file to Instagram
	err = instagram.PostContent(sess, image)
	if err != nil {
		slog.Info("failed to post content: %s", err)
	}
	slog.Info("uploaded content to instagram")

	err = filesystem.MoveFileToPostedDirectory(conf.Filesystem.ImageDirectory, fileName)
	if err != nil {
		slog.Info("failed to move file to posted directory: %s", err)
	}
	slog.Info("image file moved to posted directory")
}
