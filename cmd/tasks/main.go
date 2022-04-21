package main

import (
	"borda/pkg/pg"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	token := flag.String("token", "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", "GitHub personal access token")
	url := flag.String("url", "https://github.com/gumrf/ctf_tasks_2022", "Git repository url")
	dburl := flag.String("db", "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable", "Database url")

	flag.Parse()

	fmt.Println("Token:", *token)
	fmt.Println("URL:", *url)
	fmt.Println("DB_URL:", *dburl)

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

	file, err := os.ReadFile(path + "/data.json")
	if err != nil {
		log.Fatal(err)
	}

	var input []Task
	err = json.Unmarshal(file, &input)
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range input {
		replacer := strings.NewReplacer(" ", "_")
		fileName := replacer.Replace(task.Title)

		err := os.MkdirAll("static/"+fileName, 0777)
		if err != nil {
			log.Fatal(err)
		}

		if err := CopyDir("static/"+fileName, ".ctf_tasks_2022/"+task.Category+"/"+task.Title); err != nil {
			fmt.Println(err)
		}

	}

	db, err := pg.Open(*dburl)
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

func CopyDir(dst, src string) error {
	src, err := filepath.EvalSymlinks(src)
	if err != nil {
		return err
	}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == src {
			return nil
		}

		if strings.HasPrefix(filepath.Base(path), ".") {
			// Skip any dot files
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		// The "path" has the src prefixed to it. We need to join our
		// destination with the path without the src on it.
		dstPath := filepath.Join(dst, path[len(src):])

		// we don't want to try and copy the same file over itself.
		if eq, err := SameFile(path, dstPath); eq {
			return nil
		} else if err != nil {
			return err
		}

		// If we have a directory, make that subdirectory, then continue
		// the walk.
		if info.IsDir() {
			if path == filepath.Join(src, dst) {
				// dst is in src; don't walk it.
				return nil
			}

			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}

			return nil
		}

		// If the current path is a symlink, recreate the symlink relative to
		// the dst directory
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			target, err := os.Readlink(path)
			if err != nil {
				return err
			}

			return os.Symlink(target, dstPath)
		}

		// If we have a file, copy the contents.
		srcF, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcF.Close()

		dstF, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstF.Close()

		if _, err := io.Copy(dstF, srcF); err != nil {
			return err
		}

		// Chmod it
		return os.Chmod(dstPath, info.Mode())
	}

	return filepath.Walk(src, walkFn)
}

// SameFile returns true if the two given paths refer to the same physical
// file on disk, using the unique file identifiers from the underlying
// operating system. For example, on Unix systems this checks whether the
// two files are on the same device and have the same inode.
func SameFile(a, b string) (bool, error) {
	if a == b {
		return true, nil
	}

	aInfo, err := os.Lstat(a)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	bInfo, err := os.Lstat(b)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return os.SameFile(aInfo, bInfo), nil
}
