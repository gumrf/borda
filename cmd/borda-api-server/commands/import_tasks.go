package main

import (
	"borda/pkg/pg"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/jmoiron/sqlx"
)

type Task struct {
	Id          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Complexity  string   `json:"complexity"`
	Points      int      `json:"points"`
	Hint        string   `json:"hint"`
	Flag        string   `json:"flag"`
	Author      Author   `json:"author"`
	IsActive    bool     `json:"isActive"`
	IsDisabled  bool     `json:"isDisabled"`
	Link        string   `json:"link"`
	Files       []string `json:"files"`
}

type Author struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

func main() {
	token := flag.String("token", "nil", "GitHub personal access token")
	url := flag.String("url", "nil", "Git repository URL")
	dsn := flag.String("db", "nil", "Database URL")

	flag.Parse()

	fmt.Println("Token:", *token)
	fmt.Println("URL:", *url)
	fmt.Println("DB_URL:", *dsn)

	a := strings.Split(*url, "/")
	path := os.TempDir() + "/" + a[len(a)-1]

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

	file, err := os.ReadFile(path + "/data.json")
	if err != nil {
		log.Fatal(err)
	}

	var data []Task
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	db, err := pg.Open(*dsn)
	if err != nil {
		fmt.Print(err)
	}

	for _, task := range data {
		if len(task.Files) > 0 {
			taskPath := strings.ReplaceAll(strings.ToLower(task.Title), " ", "-")

			if err := os.MkdirAll("static/"+taskPath, 0777); err != nil {
				log.Fatal(err)
			}

			for _, file := range task.Files {
				src, err := os.ReadFile(strings.Join([]string{path, task.Category, taskPath, file}, "/"))
				if err != nil {
					log.Fatal(err)
				}

				if err := os.WriteFile("static/"+taskPath+"/"+file, src, 0777); err != nil {
					log.Fatal(err)
				}
			}
		}

		tx, err := db.Beginx()
		if err != nil {
			fmt.Println(err)
		}
		defer tx.Rollback() //nolint

		if err := findOrCreateAuthor(tx, &task.Author); err != nil {
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
			task.Points, task.Hint, task.Flag, true, false, task.Author.Id,
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
