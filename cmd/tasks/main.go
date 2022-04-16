package main

import (
	"borda/internal/config"
	"borda/pkg/pg"
	"strings"

	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/jmoiron/sqlx"
)

type Task struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Complexity  string `json:"complexity"`
	Points      int    `json:"points"`
	Hint        string `json:"hint"`
	Flag        string `json:"flag"`
	Author      Author `json:"author"`
	IsActive    bool   `json:"isActive"`
	IsDisabled  bool   `json:"isDisabled"`
	Files       string `json:"files"`
}

type Author struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

func main() {
	token := flag.String("token", "GITHUB_ACCESS_TOKEN", "GitHub personal access token")
	url := flag.String("url", "GIT_REPO_URL", "Git repository url")

	flag.Parse()

	fmt.Println("Token:",*token)
	fmt.Println("URL:",*url)

	a := strings.Split(*url, "/")
	path := "./." + a[len(a)-1]

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		Auth: &http.BasicAuth{
			Username: "nil", // yes, this can be anything except an empty string
			Password: *token,
		},
		URL:      *url,
		Progress: os.Stdout,
	})

	if err != nil {
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
			log.Fatal(err)
		}
	}

	file, err := os.ReadFile(".ctf_tasks_2022/data.json")
	if err != nil {
		log.Fatal(err)
	}

	var input []Task
	err = json.Unmarshal(file, &input)
	if err != nil {
		log.Fatal(err)
	}

	db, err := pg.Open(config.DatabaseURL())
	if err != nil {
		fmt.Print(err)
	}

	for _, task := range input {
		tx, err := db.Beginx()
		if err != nil {
			fmt.Println(err)
		}
		defer tx.Rollback() //nolint

		err = findOrCreateAuthor(tx, &task.Author)
		if err != nil {
			fmt.Print(err)
		}

		query := `
		INSERT INTO public.task (
			title, 
			description, 
			category, 
			complexity, 
			points, 
			hint, 
			flag, 
			is_active, 
			is_disabled, 
			author_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		if _, err := tx.Exec(
			query, task.Title, task.Description, task.Category, task.Complexity,
			task.Points, task.Hint, task.Flag, true, task.IsDisabled, task.Author.Id,
		); err != nil {
			log.Fatalln(err)
		}

		if err := tx.Commit(); err != nil {
			log.Fatalln(err)
		}
	}
}
func findOrCreateAuthor(tx *sqlx.Tx, author *Author) error {
	query := `
		SELECT id
		FROM author
		WHERE name = $1
		LIMIT 1`

	if err := tx.Get(&author.Id, query, author.Name); err != nil {
		insert := `
			INSERT INTO author (
				name,
				contact
			)
			VALUES ($1, $2)
			RETURNING id`

		if err := tx.Get(&author.Id, insert, author.Name, author.Contact); err != nil {
			return err
		}
	}

	return nil
}
