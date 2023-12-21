package main

import (
	"log"
  "go-insta/instagram"
  "go-insta/filesystem"
  "go-insta/config"
  "github.com/joho/godotenv"
)

func init() {
  // loads values from .env into the system
  if err := godotenv.Load(); err != nil {
      log.Printf("No .env file found: %s", err)
  }
}

func main() { 
  // Generate new config object from config package
  conf := config.New()

  // Get random image file from image directory
  image, fileName, err := filesystem.GetRandomContent(conf.Filesystem.ImageDirectory)
  if err != nil {
    log.Fatal(err)
  }

  // Log into instagram account and create session
  log.Println("logging into: " + conf.Instagram.Username)
  sess, err := instagram.CreateSession(conf.Instagram.Username, conf.Instagram.Password)
  if err != nil {
    log.Fatal(err)
  }

  // Post random image file to Instagram
  err = instagram.PostContent(sess, image)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("uploaded content to instagram")

  err = filesystem.MoveFileToPostedDirectory(conf.Filesystem.ImageDirectory, fileName)
  if err != nil {
    log.Fatal(err)
  }
  log.Println("image file moved to posted directory")
}