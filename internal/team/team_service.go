package team

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"borda/internal/pkg/core"
	"borda/pkg/transaction"
)

type TeamService struct {
	userRepository core.UserRepository
	teamRepository core.TeamRepository
	txManager      transaction.Manager
}

func NewTeamService(userRepository core.UserRepository,
	teamRepository core.TeamRepository, txManager transaction.Manager) *TeamService {
	return &TeamService{
		userRepository: userRepository,
		teamRepository: teamRepository,
		txManager:      txManager,
	}
}

func (ts *TeamService) JoinTeam(userId int, request JoinTeamRequest) error {
	if err := request.Validate(); err != nil {
		return fmt.Errorf("can't validate input: %w", err)
	}

	return ts.txManager.Run(
		context.Background(),
		func(ctx context.Context) error {
			team, err := ts.teamRepository.FindByToken(ctx, request.Token)
			if err != nil {
				return fmt.Errorf("can't find team: %w", err)
			}

			if err := ts.teamRepository.AddMember(ctx, team.Id, userId); err != nil {
				return fmt.Errorf("can't add member: %w", err)
			}

			return nil
		},
		true, // commit transaction
	)
}

func (ts *TeamService) CreateTeam(userId int, request CreatTeamRequest) error {
	if err := request.Validate(); err != nil {
		return fmt.Errorf("can't validate input: %w", err)
	}

	// TODO: token generator
	uuid, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("can't create token: %w", err)
	}

	token := uuid.String()

	return ts.txManager.Run(
		context.Background(),
		func(ctx context.Context) error {
			team, err := ts.teamRepository.Save(ctx,
				Team{
					Name:  request.Name,
					Token: token,
				},
			)
			if err != nil {
				return fmt.Errorf("can't save team: %w", err)
			}

			if err := ts.teamRepository.AddMember(ctx, team.Id, userId); err != nil {
				return fmt.Errorf("can't add member: %w", err)
			}

			return nil
		},
		true, // commit transaction
	)
}
