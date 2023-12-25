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
	// Initialize structured logging configuration
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load environment variables", "error", err)
		os.Exit(1)
	}
}

func main() {
	// Generate new config object from config package
	conf := config.New()

	// Get random image file from image directory
	image, fileName, err := filesystem.GetRandomContent(conf.Filesystem.ImageDirectory)
	if err != nil {
		slog.Error("failed to retrieve image file", "error", err)
		os.Exit(1)
	}

	// Log into instagram account and create session
	slog.Info("logging into: " + conf.Instagram.Username)
	sess, err := instagram.CreateSession(conf.Instagram.Username, conf.Instagram.Password)
	if err != nil {
		slog.Error("failed to log into account", "error", err)
		os.Exit(1)
	}

	// Set proxy settings if provided
	if len(conf.Network.ProxyAddress) > 0 {
		err = sess.SetProxy(conf.Network.ProxyAddress, true, true)
		if err != nil {
			slog.Error("failed to set proxy", "error", err)
			os.Exit(1)
		}
	}

	// Post image to Instagram
	err = instagram.PostContent(sess, image)
	if err != nil {
		slog.Error("failed to post content to Instagram", "error", err)
		os.Exit(1)
	}
	slog.Info("uploaded content to Instagram")

	// Move posted image file to /posted subsdirectory
	err = filesystem.MoveFileToPostedDirectory(conf.Filesystem.ImageDirectory, fileName)
	if err != nil {
		slog.Warn("failed to move file to /posted sub directory", "error", err)
	}
	slog.Info("file moved to /posted sub directory")
}
