package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/thevickypedia/filebrowser/v2/settings"
	"github.com/thevickypedia/filebrowser/v2/users"
)

// MethodJSONAuth is used to identify json auth.
const MethodJSONAuth settings.AuthMethod = "json"

type jsonCred struct {
	Password  string `json:"password"`
	Username  string `json:"username"`
	ReCaptcha string `json:"recaptcha"`
}

// JSONAuth is a json implementation of an Auther.
type JSONAuth struct {
	ReCaptcha *ReCaptcha `json:"recaptcha" yaml:"recaptcha"`
}

// decodeUnicodeEscape decodes Unicode escape sequences in a string
func decodeUnicodeEscape(value string) (string, error) {
	// Wrap the string in double quotes to make it a valid JSON string
	// quotedValue := fmt.Sprintf(`"%s"`, value) //nolint:govet
	quotedValue := `"` + value + `"`
	// Use json.Unmarshal to decode the Unicode escape sequences
	var decodedValue string
	err := json.Unmarshal([]byte(quotedValue), &decodedValue)
	if err != nil {
		return "", err
	}
	return decodedValue, nil
}

func decodeBase64(value string) (string, error) {
	// Decode base64
	decodedAuth, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	return string(decodedAuth), nil
}

// getCredentialParts breaks down the credentials into parts and decodes them
func getCredentialParts(value string) ([]string, error) {
	// Split the input string by commas
	parts := make([]string, 0, 3)
	unicodeParts := strings.Split(value, ",")
	for i, part := range unicodeParts {
		// Decode each part using unicode escape
		decodedUnicode, err := decodeUnicodeEscape(part)
		if err != nil {
			if i == 2 {
				// Handle the special case for the third part (recaptcha)
				parts = append(parts, "")
			} else {
				log.Printf("error: decodeAuth: %s", err)
				return nil, err
			}
		}
		parts = append(parts, decodedUnicode)
	}
	return parts, nil
}

func extractCredentials(value string) (*jsonCred, error) {
	decodedAuth, err := decodeBase64(value)
	if err != nil {
		return nil, err
	}
	// Convert decoded byte array to string
	parts, err := getCredentialParts(decodedAuth)
	if err != nil {
		return nil, err
	}
	// Check if we have enough parts
	if len(parts) < 3 {
		return nil, fmt.Errorf("insufficient parts extracted from the decoded string")
	}
	// Create jsonCred struct with the extracted values
	authDetails := &jsonCred{
		Username:  parts[0],
		Password:  parts[1],
		ReCaptcha: parts[2],
	}
	return authDetails, nil
}

// mutex to handle concurrent access to the errors map
// var mutex = &sync.Mutex{}

// authCounter is a map to store the error counts for each host
var authCounter = make(map[string]int)

// forbidden is an array to store the hosts that are temporarily forbidden
var forbidden []string

func handleAuthError(r *http.Request) {
	if count, exists := authCounter[r.Host]; exists {
		authCounter[r.Host] = count + 1
		log.Printf("Failed auth, attempt #%d for %s", authCounter[r.Host], r.Host)
		attempt := authCounter[r.Host]
		if attempt >= 10 {
			epoch := time.Now().Unix()
			until := epoch + 2_592_000
			formattedTime := time.Unix(until, 0).Format("2006-01-02 15:04:05 MST")
			log.Printf("%s is blocked until %s", r.Host, formattedTime)
			removeRecord(r.Host)
			putRecord(r.Host, until)
		} else if attempt > 3 {
			var alreadyBlocked bool
			alreadyBlocked = false
			for _, blocked := range forbidden {
				if r.Host == blocked {
					alreadyBlocked = true
					break
				}
			}
			if !alreadyBlocked {
				forbidden = append(forbidden, r.Host)
			}

			mapped := map[int]int{4: 5, 5: 10, 6: 20, 7: 40, 8: 80, 9: 160, 10: 220}
			minutes, ok := mapped[attempt]
			if !ok {
				log.Printf("Something went horribly wrong for %dth attempt", attempt)
				minutes = 60 // Default to 1 hour
			}
			epoch := time.Now().Unix()
			until := epoch + int64(minutes*60)
			formattedTime := time.Unix(until, 0).Format("2006-01-02 15:04:05 MST")
			log.Printf("%s is blocked (for %d minutes) until %s", r.Host, minutes, formattedTime)
			removeRecord(r.Host)
			putRecord(r.Host, until)
		}
	} else {
		log.Printf("Failed auth, attempt #1 for %s", r.Host)
		authCounter[r.Host] = 1
	}
}

func removeItem(slice []string, item string) []string {
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			// Found the item, remove it by slicing the slice
			return append(slice[:i], slice[i+1:]...)
		}
	}
	// If item is not found, return the original slice
	return slice
}

// Auth authenticates the user via a json in authorization header.
func (a JSONAuth) Auth(r *http.Request, usr users.Store, _ *settings.Settings, srv *settings.Server) (*users.User, error) {
	var block bool
	block = false
	for _, blocked := range forbidden {
		if r.Host == blocked {
			block = true
			break
		}
	}
	if block {
		timestamp, err := getRecord(r.Host)
		if err != nil {
			log.Printf("Unable to check if %s was forbidden, allowing..", r.Host)
		} else {
			epoch := time.Now().Unix()
			if timestamp > epoch {
				formattedTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05 MST")
				log.Printf("%s is forbidden until %s due to repeated login failures", r.Host, formattedTime)
				return nil, os.ErrPermission
			}
		}
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, os.ErrPermission
	}

	cred, err := extractCredentials(authHeader)
	if err != nil {
		log.Fatal("error:", err)
		return nil, err
	}

	// If ReCaptcha is enabled, check the code.
	if a.ReCaptcha != nil && a.ReCaptcha.Secret != "" {
		ok, err := a.ReCaptcha.Ok(cred.ReCaptcha) //nolint:govet

		if err != nil {
			return nil, err
		}

		if !ok {
			return nil, os.ErrPermission
		}
	}

	u, err := usr.Get(srv.Root, cred.Username)
	if err != nil || !users.CheckPwd(cred.Password, u.Password) {
		handleAuthError(r)
		return nil, os.ErrPermission
	}

	forbidden = removeItem(forbidden, r.Host)
	return u, nil
}

// LoginPage tells that json auth doesn't require a login page.
func (a JSONAuth) LoginPage() bool {
	return true
}

const reCaptchaAPI = "/recaptcha/api/siteverify"

// ReCaptcha identifies a recaptcha connection.
type ReCaptcha struct {
	Host   string `json:"host"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

// Ok checks if a reCaptcha responde is correct.
func (r *ReCaptcha) Ok(response string) (bool, error) {
	body := url.Values{}
	body.Set("secret", r.Secret)
	body.Add("response", response)

	client := &http.Client{}

	resp, err := client.Post(
		r.Host+reCaptchaAPI,
		"application/x-www-form-urlencoded",
		strings.NewReader(body.Encode()),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	var data struct {
		Success bool `json:"success"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return false, err
	}

	return data.Success, nil
}
