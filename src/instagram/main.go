package instagram

import (
	"fmt"
	"io"
	"os"
	"strings"
	"math/rand"
	"github.com/Davincible/goinsta/v3"
)

type caption struct {
	summary string
	cta string
	hashtag string
}

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

func randomHashtagAssortment(list []string, amountOfHashtags int) string {
	hashtags := ""
	for i := 0; i < amountOfHashtags; i++ {
		randomIndex := rand.Intn(len(list))
		randomHashtag := list[randomIndex]
		hashtags = fmt.Sprintf("%s %s", hashtags, randomHashtag)
	}
	return hashtags
}

func CreateSession(accountName string, accountPassword string) (*goinsta.Instagram, error) {
	insta := goinsta.New(accountName, accountPassword)
	if err := insta.Login(); err != nil {
		return nil, err
	}
	return insta, nil
}

func PostContent(insta *goinsta.Instagram, uploadContent io.Reader) (error) {
	var postCaption string
	if (os.Getenv("POST_CAPTION") == "") {
		summaryList := strings.Split(os.Getenv("CAPTION_SUMMARY"), ",")
		summary := randomItemFromList(summaryList)

		ctaList := strings.Split(os.Getenv("CAPTION_CTA"), ",")
		cta := randomItemFromList(ctaList)

		hashtagList := strings.Split(os.Getenv("CAPTION_HASHTAG"), ",")
		hashtag := randomHashtagAssortment(hashtagList, 15)

		c := caption{
			summary: summary,
			cta: cta,
			hashtag: hashtag,
		}
		postCaption = c.buildCaption()
	} else {
		postCaption = os.Getenv("POST_CAPTION")
	}

	_, err := insta.Upload(
		&goinsta.UploadOptions{
			File:    uploadContent,
			Caption:  postCaption},
	)
	if err != nil {
	  return err
	}
	return nil
}