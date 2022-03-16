package usecase

import (
	"borda/internal/domain"
	"borda/internal/repository"
)

type UserUsecase struct {
	taskRepo repository.TaskRepository
}

func NewUserUsecase(tsr repository.TaskRepository) *UserUsecase {
	return &UserUsecase{taskRepo: tsr}
}

func (u *UserUsecase) ShowAllTasks(filter domain.TaskFilter) ([]*domain.Task, error) {
	var tasks []*domain.Task

	tasks, err := u.taskRepo.GetTasks(filter)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (a *UserUsecase) TryToSolveTask(submission domain.SubmitTaskRequest) (string, error) {
	var task *domain.Task
	var err error
	//Надо как-то заполнить таблицу в бд task_submissions
	// нет нужных методов в репозитории

	task, err = a.taskRepo.GetTaskById(submission.TaskId)
	if err != nil {
		return "", err
	}

	if submission.Flag == task.Flag {
		err = a.taskRepo.SolveTask(task.Id, submission.TeamId)
		// relation \"public.solved_task\" does not exist"
		// Где-то тут надо заполнить task_submissions с помощью domain.TaskSubmission
		if err != nil {
			return "", err
		}
		return "Submission is correct", nil
	} else {
		return "Submission is incorrect", nil
	}

}
