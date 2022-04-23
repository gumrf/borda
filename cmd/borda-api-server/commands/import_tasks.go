package commands

import (
	"borda/internal/domain"
	"borda/internal/repository/postgres"
	"borda/pkg/pg"
	"encoding/json"
	"errors"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/cobra"
)

func ImportTasksCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "import-tasks",
		Short: "Import tasks from GitHub repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := cmd.Flags().GetString("token")
			if err != nil {
				return err
			}

			repoURL, err := cmd.Flags().GetString("repo")
			if err != nil {
				return err
			}

			dbURL, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}

			if err := importTasks(token, repoURL, dbURL); err != nil {
				return err
			}

			return nil
		},
	}

	command.Flags().StringP("token", "t", "nil", "GitHub personal access token")
	command.Flags().StringP("repo", "r", "nil", "Git repository URL")
	command.Flags().String("db", "db", "Database URL")

	command.MarkFlagRequired("token")
	command.MarkFlagRequired("repo")
	command.MarkFlagRequired("db")

	return command
}

type Task struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Category    string        `json:"category"`
	Complexity  string        `json:"complexity"`
	Points      int           `json:"points"`
	Hint        string        `json:"hint"`
	Flag        string        `json:"flag"`
	Author      domain.Author `json:"author"`
	Link        string        `json:"link"`
	Files       []string      `json:"files"`
}

func importTasks(token, repoURL, dbURL string) error {
	a := strings.Split(repoURL, "/")
	path := os.TempDir() + "/" + a[len(a)-1]

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		Auth: &http.BasicAuth{
			Username: "nil", // yes, this can be anything except an empty string
			Password: token,
		},
		URL:      repoURL,
		Progress: os.Stdout,
	})

	if err != nil {
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
			return err
		}
	}

	file, err := os.ReadFile(path + "/data.json")
	if err != nil {
		return err
	}

	var data []Task
	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	db, err := pg.Open(dbURL)
	if err != nil {
		return err
	}

	for _, task := range data {
		if len(task.Files) > 0 {
			taskPath := strings.ReplaceAll(strings.ToLower(task.Title), " ", "-")

			if err := os.MkdirAll("static/"+taskPath, 0777); err != nil {
				return err
			}

			for _, file := range task.Files {
				src, err := os.ReadFile(strings.Join([]string{path, task.Category, taskPath, file}, "/"))
				if err != nil {
					return err
				}

				if err := os.WriteFile("static/"+taskPath+"/"+file, src, 0777); err != nil {
					return err
				}
			}
		}

		if _, err := postgres.NewTaskRepository(db).SaveTask(domain.Task{
			Title:       task.Title,
			Description: task.Description,
			Category:    task.Category,
			Complexity:  task.Complexity,
			Points:      task.Points,
			Hint:        task.Hint,
			Flag:        task.Flag,
			IsActive:    true,
			IsDisabled:  false,
			Author:      task.Author,
			Link:        task.Link,
		}); err != nil {
			return err
		}
	}

	return nil
}
