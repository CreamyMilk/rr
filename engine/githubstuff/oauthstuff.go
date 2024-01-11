package githubstuff

import (
	"aytp/engine"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func PersitGHTokenForUser(userToken, ghToken string) error {
	userInfo, _, err := engine.FetchUserAndModelByToken(userToken)
	if err != nil {
		return err
	}

	_, err = engine.UpdateStudentGHToken(userInfo.ID.String(), ghToken)
	if err != nil {
		return err
	}

	log.Printf("✅ GH Stuff ripped out (%s) has an new token (%s) \n", userInfo.ID, ghToken)

	return nil
}

type RepoStuff struct {
	RepoName          string
	RepoDescription   string
	TemplateRepoMaybe string
}

// incase you wanna play around with the creation flow
// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#create-a-repository-for-the-authenticated-user
type CreateRepoRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Private     bool   `json:"private"`
}

type RepoCreationResp struct {
	ID       int    `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Private  bool   `json:"private"`
}

func CreateRepo(ghToken string, creq *CreateRepoRequest) (*RepoCreationResp, error) {
	url := "https://api.github.com/user/repos"

	requestBody, err := json.Marshal(creq)
	if err != nil {
		return nil, err
	}

	byteBody := []byte(requestBody)
	body := bytes.NewBuffer(byteBody)

	req, err := http.NewRequest(http.MethodPost, url, body)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ghToken))
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var repoCreatedStuff RepoCreationResp
	err = json.Unmarshal([]byte(responseBody), &repoCreatedStuff)
	if err != nil {
		return nil, err
	}

	log.Println("✅ Repo created", repoCreatedStuff.ID)

	return &repoCreatedStuff, nil
}
