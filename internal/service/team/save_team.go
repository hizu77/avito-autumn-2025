package team

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/pkg/errors"
)

func (s *Service) SaveTeam(ctx context.Context, team model.Team) (model.Team, error) {
	usersWithTeam := collection.Map(
		team.Members,
		func(user model.User) model.User {
			return model.User{
				ID:       user.ID,
				Name:     user.Name,
				TeamName: team.Name,
				IsActive: user.IsActive,
			}
		},
	)

	var savedTeam model.Team
	var savedUsers []model.User
	err := s.trManager.Do(ctx, func(ctx context.Context) error {
		txTeam, err := s.teamStorage.SaveTeam(ctx, team)
		if err != nil {
			return errors.Wrap(err, "team storage saving team")
		}

		txUsers, err := s.userStorage.SaveUsers(ctx, usersWithTeam)
		if err != nil {
			return errors.Wrap(err, "user storage saving users")
		}

		savedTeam = txTeam
		savedUsers = txUsers

		return nil
	})
	if err != nil {
		return model.Team{}, errors.Wrap(err, "saving team with users")
	}

	savedTeam.Members = savedUsers

	return savedTeam, nil
}
