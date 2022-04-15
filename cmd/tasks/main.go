package main

import (
	"borda/internal/config"
	"borda/pkg/pg"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/jmoiron/sqlx"
)

type Task struct {
	Id          int    `json: "id"`
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
	url := "https://github.com/gumrf/ctf_tasks_2022"
	token := "token" //delete token !!!
	_, err := git.PlainClone("./.ctf_tasks_2022", false, &git.CloneOptions{
		// The intended use of a GitHub personal access token is in replace of your password
		// because access tokens can easily be revoked.
		// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
		Auth: &http.BasicAuth{
			Username: "abc123", // yes, this can be anything except an empty string
			Password: token,
		},
		URL:      url,
		Progress: os.Stdout,
	})

	if err != nil {
		if !errors.Is(err, git.ErrRepositoryAlreadyExists) {
			log.Fatal(err)
		}
	}

	var input []Task

	file, err := os.ReadFile(".ctf_tasks_2022/data.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &input)
	if err != nil {
		log.Fatal(err)
	}

	db, err := pg.Open(config.DatabaseURL())
	if err != nil {
		fmt.Print(err)
	}

	if err := pg.Migrate(db, config.MigrationsPath()); err != nil {
		fmt.Print(err)
	}

	for _, task := range input {

		tx, err := db.Beginx()
		if err != nil {
			fmt.Println(err)
		}
		defer tx.Rollback() //nolint

		_, err = findOrCreateAuthor(tx, &task.Author)
		if err != nil {
			fmt.Print(err)
		}

		query := fmt.Sprintf(`
		INSERT INTO task (
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`,
		)

		result := tx.QueryRow(
			query,
			task.Title, task.Description, task.Category, task.Complexity,
			task.Points, task.Hint, task.Flag, task.IsActive, task.IsDisabled,
			task.Author.Id)

		if err := result.Scan(&task.Id); err != nil {
			fmt.Print(err)
		}

		if err := tx.Commit(); err != nil {
			fmt.Print(err)
		}

		fmt.Println(task.Id)
	}
}
func findOrCreateAuthor(tx *sqlx.Tx, author *Author) (int, error) {
	query := fmt.Sprintf(`
		SELECT id
		FROM author
		WHERE name = $1
		LIMIT 1`,
	)

	result := tx.QueryRowx(query, author.Name)
	err := result.Scan(&author.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_query := fmt.Sprintf(`
				INSERT INTO author (
					name,
					contact
				)
				VALUES ($1, $2)
				RETURNING id`,
			)

			_result := tx.QueryRowx(_query, author.Name, author.Contact)
			if err := _result.Scan(&author.Id); err != nil {
				return -1, err
			}
		} else {
			return -1, err
		}
	}

	return author.Id, nil
}