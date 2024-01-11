package routes

import (
	"aytp/engine/githubstuff"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type GHAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func GithubCallbackHandler(c *fiber.Ctx) error {
	// Variable is only valid within this handler
	tempCode := c.Query("code")
	theUserIdThisTokenBelongsTo := c.Query("state")

	permanentTokenStuff, err := turnGHCodeToAuthToken(tempCode)
	if err != nil {
		return err
	}

	err = githubstuff.PersitGHTokenForUser(theUserIdThisTokenBelongsTo, permanentTokenStuff.AccessToken)
	if err != nil {
		return err
	}

	c.Redirect("https://aytp-learning.vercel.app/profile/student")

	return nil
}

func turnGHCodeToAuthToken(tempToken string) (*GHAuthTokenResponse, error) {

	url := "https://github.com/login/oauth/access_token"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    "client_id":"4c626bc93a0859a895f1" ,
    "client_secret" :"803128e0606b7b0729ea64ceec0ef769a50cb650",
    "code":"%s"
}`, tempToken))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var githubTokenStuff GHAuthTokenResponse
	err = json.Unmarshal([]byte(responseBody), &githubTokenStuff)
	if err != nil {
		return nil, err
	}

	return &githubTokenStuff, nil
}
