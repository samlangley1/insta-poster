package instagram

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"

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
		fmt.Printf("Login failed: %v\n", err)

		// Check if this is a 2FA challenge that requires manual code entry
		if insta.TwoFactorInfo != nil {
			fmt.Println("2FA challenge detected.")

			// Check if TOTP_NAME environment variable is set for automated 2FA
			totpName := os.Getenv("TOTP_NAME")
			totpPassword := os.Getenv("TOTP_PASSWORD")
			if totpName != "" {
				fmt.Printf("Generating 2FA code using totp-cli for account: %s\n", totpName)

				// Execute totp-cli to generate the code
				// Handle both "namespace:account" and "account" formats
				// Try multiple possible paths for totp-cli
				totpCliPaths := []string{"totp-cli", "/root/go/bin/totp-cli", "/usr/local/bin/totp-cli", "~/go/bin/totp-cli"}

				var cmd *exec.Cmd
				var totpCliPath string

				// Find working totp-cli path
				for _, path := range totpCliPaths {
					if _, err := exec.LookPath(path); err == nil {
						totpCliPath = path
						break
					}
				}

				if totpCliPath == "" {
					fmt.Println("totp-cli not found in PATH or common locations")
					fmt.Println("Falling back to manual code entry...")
				} else {
					// Create command with password input
					if strings.Contains(totpName, ":") {
						parts := strings.Split(totpName, ":")
						cmd = exec.Command(totpCliPath, "generate", parts[0], parts[1])
					} else {
						cmd = exec.Command(totpCliPath, "generate", totpName)
					}

					// If password is provided, pipe it to stdin
					if totpPassword != "" {
						cmd.Stdin = strings.NewReader(totpPassword + "\n")
					}
				}

				if totpCliPath != "" {

					output, err := cmd.Output()
					if err != nil {
						fmt.Printf("Failed to generate TOTP code automatically: %v\n", err)
						if totpPassword == "" {
							fmt.Println("Hint: Set TOTP_PASSWORD environment variable if your totp-cli requires a password")
						}
						fmt.Println("Falling back to manual code entry...")
					} else {
						code := strings.TrimSpace(string(output))
						fmt.Printf("Generated 2FA code: %s\n", code)

						// Validate code format
						if len(code) == 6 {
							// Attempt 2FA login with the generated code
							if err := insta.TwoFactorInfo.Login2FA(code); err != nil {
								fmt.Printf("Auto-generated 2FA code failed: %v\n", err)
								fmt.Println("Falling back to manual code entry...")
							} else {
								fmt.Println("Automated 2FA code verified! Login successful.")
								return insta, nil
							}
						} else {
							fmt.Printf("Invalid auto-generated code length: %d digits\n", len(code))
							fmt.Println("Falling back to manual code entry...")
						}
					}
				}
			}

			// Manual fallback
			fmt.Print("Please enter the 6-digit code from your authenticator app: ")
			reader := bufio.NewReader(os.Stdin)
			code, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("failed to read 2FA code: %w", err)
			}

			// Clean up the input
			code = strings.TrimSpace(code)

			// Validate code format
			if len(code) != 6 {
				return nil, fmt.Errorf("invalid 2FA code - must be exactly 6 digits, got %d digits", len(code))
			}

			// Attempt 2FA login with the provided code
			if err := insta.TwoFactorInfo.Login2FA(code); err != nil {
				return nil, fmt.Errorf("2FA code verification failed: %w", err)
			}

			fmt.Println("2FA code verified! Login successful.")
			return insta, nil
		}

		return nil, fmt.Errorf("login failed: %w", err)
	}

	fmt.Println("Login successful!")

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
