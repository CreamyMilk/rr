package courses

import (
	"aytp/engine"
	"aytp/engine/githubstuff"
	"aytp/graph/model"
	"errors"
)

func CreateRepoForCourse(input model.RepoCreateionForCourseInput) (bool, error) {
	students, err := engine.FetchStudentsDoingCourseById(input.CourseID)
	if err != nil {
		return false, errors.New("course doesn't exists")
	}

	for _, student := range students {
		if student.GithubToken != nil {
			go func(githubOAuthToken string) {
				githubstuff.CreateRepo(githubOAuthToken, &githubstuff.CreateRepoRequest{
					Name:        input.RepoName,
					Description: input.RepoDescription,
				})
			}(*student.GithubToken)
		}
	}

	return true, nil
}
