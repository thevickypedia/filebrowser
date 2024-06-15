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

func decodeUnicodeEscape(value string) (string, error) {
	// First, wrap the string in double quotes to make it a valid JSON string
	quotedValue := fmt.Sprintf("%q", value)
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

func getCredentialParts(value string) ([]string, error) {
	// Break down credentials into username, signature and recaptcha
	// Decode each of them separately using unicode escape
	parts := make([]string, 0, 3)
	unicodeParts := strings.Split(value, ",")
	for i, part := range unicodeParts {
		decodedUnicode, err := decodeUnicodeEscape(part)
		if err != nil {
			if i == 2 {
				formatError := fmt.Sprintf("ReCaptcha is null: %s", err)
				log.Print(formatError)
				parts = append(parts, "")
			} else {
				log.Fatal("error: decodeAuth:", err)
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

// Auth authenticates the user via a json in authorization header.
func (a JSONAuth) Auth(r *http.Request, usr users.Store, _ *settings.Settings, srv *settings.Server) (*users.User, error) {
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
		return nil, os.ErrPermission
	}

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
