package setup

import (
	"borda/internal/core/entities"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create persons table
		{
			ID: "000000000001",
			Migrate: func(tx *gorm.DB) error {
				// it's a good pratice to copy the struct inside the function,
				// so side effects are prevented if the original struct changes during the time
				return tx.AutoMigrate(
					&entities.User{},
					&entities.Author{},
					&entities.Task{}, &entities.SolvedTask{}, &entities.TaskSubmission{},
					&entities.Role{}, &entities.UserRole{},
					&entities.Team{},
				)
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(
					"teams", "team_members",
					"roles", "user_roles",
					"task", "solved_tasks", "task_submissions",
					"authors", "author_tasks",
					"users",
				)
			},
		},
	})

	return m.Migrate()
}
