package pullrequest

import (
	"crypto/rand"
	"math/big"

	"github.com/hizu77/avito-autumn-2025/pkg/utils/collection"
	"github.com/pkg/errors"
)

func (s *Service) getRandomReviewers(candidates []string, maxCount int) ([]string, error) {
	if len(candidates) <= maxCount {
		return candidates, nil
	}

	selected := make(map[string]struct{}, maxCount)
	for len(selected) < maxCount {
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates))))
		if err != nil {
			return nil, errors.Wrap(err, "generating random number")
		}

		idx := int(nBig.Int64())
		id := candidates[idx]
		selected[id] = struct{}{}
	}

	reviewers := collection.Keys(selected)

	return reviewers, nil
}
