package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"syscall"

	authz "github.com/OreCast/common/authz"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cobra"
	term "golang.org/x/term"
)

// User represents structure used by users DB in Authz service to handle incoming requests
type User struct {
	Login    string
	Password string
}

// helper function to get orecast token
func getToken(login, pass string) (string, error) {
	var token string
	// make a call to Authz service to check for a user
	rurl := fmt.Sprintf(
		"%s/oauth/authorize?client_id=%s&response_type=code",
		_oreConfig.Services.AuthzURL,
		_oreConfig.Authz.ClientId)
	if verbose > 0 {
		fmt.Println("HTTP GET", rurl)
	}
	user := User{Login: login, Password: pass}
	data, err := json.Marshal(user)
	if err != nil {
		return token, err
	}
	resp, err := http.Post(rurl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return token, err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if verbose > 1 {
		fmt.Println("## authorize response", string(data))
	}
	var response authz.Response
	err = json.Unmarshal(data, &response)
	if err != nil {
		return token, err
	}
	if response.Status != "ok" {
		msg := fmt.Sprintf("No user %s found in Authz service", user)
		return token, errors.New(msg)
	}

	// make request to get authz token

	var aToken authz.Token
	rurl = fmt.Sprintf(
		"%s/oauth/token?client_id=%s&client_secret=%s&grant_type=client_credentials&scope=read",
		_oreConfig.Services.AuthzURL,
		_oreConfig.Authz.ClientId,
		_oreConfig.Authz.ClientSecret)

	if verbose > 0 {
		fmt.Println("HTTP GET", rurl)
	}

	resp, err = http.Get(rurl)
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if verbose > 1 {
		fmt.Println("## token data", string(data))
	}
	if err != nil {
		return token, err
	}
	err = json.Unmarshal(data, &aToken)
	if err != nil {
		return token, err
	}
	reqToken := aToken.AccessToken
	if verbose > 1 {
		fmt.Println("## request token", reqToken)
	}

	// validate our token
	var jwtKey = []byte(_oreConfig.Authz.ClientId)
	claims := &authz.Claims{}
	tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return token, errors.New("invalid token signature")
		}
		return token, err
	}
	if !tkn.Valid {
		return token, errors.New("invalid token validity")
	}
	return reqToken, nil
}

// helper function to get user input
func inputPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, label+" ")
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

// helper function to get user password
func passwordPrompt(label string) string {
	var s string
	for {
		fmt.Fprint(os.Stderr, label+" ")
		pw, _ := term.ReadPassword(int(syscall.Stdin))
		s = string(pw)
		if s != "" {
			break
		}
	}
	fmt.Println()
	return s
}

// helper function to get access token
func accessToken() (string, error) {
	user := inputPrompt("OreCast username:")
	pass := passwordPrompt("OreCast password:")
	return getToken(user, pass)
}

func authCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "OreCast token command",
		Long: `OreCast token command
                Complete documentation is available at https://orecast.com/documentation/`,
		Run: func(cmd *cobra.Command, args []string) {
			if token, err := accessToken(); err == nil {
				fmt.Println(token)
			} else {
				fmt.Println("ERROR", err)
			}
		},
	}
	return cmd
}
