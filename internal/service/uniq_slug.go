package service

import (
	"context"
	"errors"

	"github.com/Thoustick/SlugKiller/internal/repository"
)

func (s *urlService) generateUniqueSlug(ctx context.Context) (string, error) {
	for i := 0; i < s.cfg.MaxAttempts; i++ {
		slug, err := s.slugGen.Generate(ctx)
		if err != nil {
			return "", err
		}

		_, err = s.repo.GetBySlug(ctx, slug)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return slug, nil
			}
			return "", err
		}
		// если ошибки не было, slug занят, идём на следующую итерацию
	}
	return "", errors.New("could not generate unique slug after several attempts")
}
