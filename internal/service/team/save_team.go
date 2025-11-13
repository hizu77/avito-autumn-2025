package team

import (
	"context"

	"github.com/hizu77/avito-autumn-2025/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (s *Service) SaveTeam(ctx context.Context, team model.Team) (model.Team, error) {
	var saved model.Team

	err := s.transactor.WithTx(ctx, func(ctx context.Context) error {
		savedTeam, err := s.teamStorage.SaveTeam(ctx, team)
		if err != nil {
			return errors.Wrap(err, "team storage saving team")
		}

		members := make([]model.User, 0, len(team.Members))
		for _, member := range team.Members {
			member.TeamName = team.Name
			savedMember, err := s.userStorage.SaveUser(ctx, member)
			if err != nil {
				return errors.Wrap(err, "user storage saving user")
			}

			members = append(members, savedMember)
		}

		savedTeam.Members = members
		saved = savedTeam

		return nil
	})
	if err != nil {
		s.logger.Error("saving team", zap.Error(err))
		return model.Team{}, errors.Wrap(err, "saving team")
	}

	return saved, nil
}
