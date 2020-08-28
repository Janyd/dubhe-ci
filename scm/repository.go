package scm

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"io"
)

type Repository interface {
	RemoteBranches() ([]string, error)

	Pull() ([]*object.Commit, error)

	PullWithChanges() ([]string, error)

	FindCommit(ref string) (*object.Commit, error)

	Head() (*object.Commit, error)
}

func NewRepository(repo *git.Repository, auth transport.AuthMethod, branch string) Repository {
	return &repository{
		repo:   repo,
		auth:   auth,
		Branch: branch,
	}
}

type repository struct {
	repo   *git.Repository
	auth   transport.AuthMethod
	Branch string
}

func (r *repository) Head() (*object.Commit, error) {
	head, err := r.repo.Head()
	if err != nil {
		return nil, err
	}
	return r.repo.CommitObject(head.Hash())
}

func (r *repository) RemoteBranches() ([]string, error) {
	refs, err := r.repo.References()
	if err != nil {
		return nil, err
	}

	branches := make([]string, 0)

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsRemote() {
			refName := ref.Name().Short()
			branch := refName[7:]
			branches = append(branches, branch)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *repository) Pull() ([]*object.Commit, error) {
	preRef, err := r.repo.Head()
	if err != nil {
		return nil, err
	}
	fmt.Println(preRef)

	w, err := r.repo.Worktree()
	if err != nil {
		return nil, err
	}

	err = w.Pull(&git.PullOptions{
		ReferenceName: plumbing.NewBranchReferenceName(r.Branch),
		RemoteName:    "origin",
		SingleBranch:  true,
		Auth:          r.auth,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, err
	}

	afterRef, err := r.repo.Head()
	if err != nil {
		return nil, err
	}

	if preRef.Hash().String() != afterRef.Hash().String() {
		o, err := r.repo.Log(&git.LogOptions{
			From: afterRef.Hash(),
			All:  false,
		})

		if err != nil {
			return nil, err
		}

		commits := make([]*object.Commit, 0)

		for {
			commit, err := o.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			if commit.Hash.String() == preRef.Hash().String() {
				commits = append(commits, commit)
				break
			}
			commits = append(commits, commit)
		}

		return commits, nil
	}

	return nil, nil
}

func (r *repository) FindCommit(ref string) (*object.Commit, error) {
	commit, err := r.repo.CommitObject(plumbing.NewHash(ref))
	if err != nil {
		return nil, err
	}

	return commit, nil
}

func (r *repository) PullWithChanges() ([]string, error) {
	commits, err := r.Pull()
	if err != nil {
		return nil, err
	}

	changeStrings := make([]string, 0)
	var prev *object.Commit
	for _, commit := range commits {
		if prev == nil {
			prev = commit
			continue
		}
		tree, err := commit.Tree()
		if err != nil {
			return nil, err
		}
		prevTree, err := prev.Tree()
		changes, err := object.DiffTree(tree, prevTree)
		for _, c := range changes {
			changeStrings = append(changeStrings, c.String())
		}

		prev = commit
	}

	return changeStrings, nil
}
