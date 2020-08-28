package core

import "context"

type (
	File struct {
		Data []byte
		Hash []byte
	}

	FileService interface {
		Find(ctx context.Context, repo, branch, commit, ref, path string) (*File, error)
	}
)
