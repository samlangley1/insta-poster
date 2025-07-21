package instagram

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Davincible/goinsta/v3"
)

type SessionOptions struct {
	ProxyAddress string
}

// Full caption type for creating randomly generated captions
type caption struct {
	summary string
	cta     string
	hashtag string
}

// Build full Instagram caption from parts
func (caption *caption) buildCaption() string {
	fullCaption := fmt.Sprintf(`%s

	-----------

	%s

	-----------

	%s`, caption.summary, caption.cta, caption.hashtag)
	return fullCaption
}

func randomItemFromList(list []string) string {
	randomIndex := rand.Intn(len(list))
	randomItem := list[randomIndex]
	return randomItem
}

// Return a single string of concatenated randomly selected hashtags from a list of hashtags
func randomHashtagAssortment(list []string, amountOfHashtags int) string {
	hashtags := ""
	for i := 0; i < amountOfHashtags; i++ {
		randomIndex := rand.Intn(len(list))
		randomHashtag := list[randomIndex]
		hashtags = fmt.Sprintf("%s %s", hashtags, randomHashtag)
	}
	return hashtags
}

// Log into Instagram
func CreateSession(accountName string, accountPassword string, o *SessionOptions) (*goinsta.Instagram, error) {
	insta := goinsta.New(accountName, accountPassword)

	// Set proxy settings if provided
	if len(o.ProxyAddress) > 0 {
		if err := insta.SetProxy(o.ProxyAddress, true, true); err != nil {
			return nil, fmt.Errorf("failed to set proxy %s: %w", o.ProxyAddress, err)
		}
	}

	// Attempt initial login
	if err := insta.Login(); err != nil {
		fmt.Printf("Initial login failed, this might be a 2FA request: %v\n", err)
		fmt.Println("Please approve the 2FA request on your phone if prompted. Waiting 60 seconds...")

		// Wait 60 seconds for user to approve 2FA
		time.Sleep(60 * time.Second)

		// Retry login after waiting
		if err := insta.Login(); err != nil {
			return nil, fmt.Errorf("login failed after 2FA wait: %w", err)
		}
		fmt.Println("Login successful after 2FA wait!")
	} else {
		fmt.Println("Login successful!")
	}

	return insta, nil
}

// Post content to logged in Instagram account
func PostContent(insta *goinsta.Instagram, uploadContent io.Reader) error {
	var postCaption string
	// Check for POST_CAPTION env var, which indicates the full caption is provided, otherwise, generate random caption from summary, cta, and hashtags
	if os.Getenv("POST_CAPTION") == "" {
		summaryList := strings.Split(os.Getenv("CAPTION_SUMMARY"), ",")
		summary := randomItemFromList(summaryList)

		ctaList := strings.Split(os.Getenv("CAPTION_CTA"), ",")
		cta := randomItemFromList(ctaList)

		hashtagList := strings.Split(os.Getenv("CAPTION_HASHTAG"), ",")
		hashtag := randomHashtagAssortment(hashtagList, 15)

		c := caption{
			summary: summary,
			cta:     cta,
			hashtag: hashtag,
		}
		postCaption = c.buildCaption()
	} else {
		postCaption = os.Getenv("POST_CAPTION")
	}
	// Upload image to Instagram with caption
	_, err := insta.Upload(
		&goinsta.UploadOptions{
			File:    uploadContent,
			Caption: postCaption},
	)
	if err != nil {
		return err
	}
	return nil
}
