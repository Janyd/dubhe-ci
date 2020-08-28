package repos

import (
	"context"
	"dubhe-ci/common"
	"dubhe-ci/core"
	"dubhe-ci/errors"
	"dubhe-ci/logger"
	"dubhe-ci/scm"
	"dubhe-ci/service/rpc/pb"
	"dubhe-ci/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
	"time"
)

func NewRepoService(
	repoStore core.RepositoryStore,
	branchStore core.BranchStore,
	credStore core.CredentialStore,
	gitService scm.GitService,
	triggerService core.TriggerService,
) *RepositoryService {
	return &RepositoryService{
		repoStore:      repoStore,
		branchStore:    branchStore,
		credStore:      credStore,
		gitService:     gitService,
		triggerService: triggerService,
	}
}

type RepositoryService struct {
	repoStore      core.RepositoryStore
	credStore      core.CredentialStore
	branchStore    core.BranchStore
	gitService     scm.GitService
	triggerService core.TriggerService
}

func (r *RepositoryService) Register(grpcServer *grpc.Server) {
	pb.RegisterRepositoryServer(grpcServer, r)
}

func (r *RepositoryService) List(ctx context.Context, page *pb.Page) (*pb.RepoRecords, error) {
	p := &common.Page{
		Current: page.Current,
		Size:    page.Size,
		Desc:    page.Desc,
		Asc:     page.Asc,
	}
	p, err := r.repoStore.List(ctx, p)
	if err != nil {
		return nil, err
	}
	page.Total = p.Total

	repos := make([]*pb.Repo, 0)
	err = copier.Copy(&repos, p.Records)
	if err != nil {
		return nil, err
	}

	return &pb.RepoRecords{
		Page:  page,
		Repos: repos,
	}, nil
}

func (r *RepositoryService) Create(ctx context.Context, repo *pb.Repo) (*pb.Repo, error) {

	_, err := r.credStore.Find(ctx, repo.CredentialId)
	if err != nil {
		return nil, err
	}

	repository := &core.Repository{}
	err = copier.Copy(repository, repo)
	if err != nil {
		return nil, err
	}

	repository, err = r.repoStore.Create(ctx, repository)
	if err != nil {
		return nil, err
	}
	err = copier.Copy(repo, repository)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *RepositoryService) Find(ctx context.Context, id *pb.Id) (*pb.Repo, error) {
	repository, err := r.repoStore.Find(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	repo := &pb.Repo{}
	err = copier.Copy(repo, repository)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *RepositoryService) Update(ctx context.Context, repo *pb.Repo) (*pb.Empty, error) {
	repository, err := r.convert(repo)
	if err != nil {
		return nil, err
	}

	err = r.repoStore.Update(ctx, repository)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (r *RepositoryService) Delete(ctx context.Context, id *pb.Id) (*pb.Empty, error) {
	repository, err := r.repoStore.Find(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	err = r.repoStore.Delete(ctx, repository)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (r *RepositoryService) convert(repo *pb.Repo) (*core.Repository, error) {
	repository := &core.Repository{}
	err := copier.Copy(repository, repo)
	if err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *RepositoryService) SyncBranch(ctx context.Context, repo *core.Repository) ([]string, error) {
	log := logger.WithAction("init repository")

	branches, err := r.branchStore.List(ctx, repo.Id)
	branchMap := make(map[string]*core.Branch, len(branches))

	for _, branch := range branches {
		branchMap[branch.Name] = branch
	}

	auth, err := r.auth(ctx, repo)

	gitRepo, _, err := r.gitService.Open(repo.Url, repo.Name, "master", auth)
	if err != nil {
		return nil, err
	}

	remoteBranches, err := gitRepo.RemoteBranches()
	if err != nil {
		return nil, err
	}

	newBranches := make([]string, 0)
	for _, branch := range remoteBranches {
		if _, ok := branchMap[branch]; !ok {
			_, err := r.branchStore.Create(ctx, repo.Id, branch)
			if err != nil {
				log.WithError(err).Warnln("cannot create new branch")
				continue
			}
			newBranches = append(newBranches, branch)
		}
		_, _, err = r.gitService.Open(repo.Url, repo.Name, branch, auth)
		if err != nil {
			log.WithError(err).Warnf("cannot open %s repository %s branch", repo.Name, branch)
		}
	}

	for k := range branchMap {
		if !utils.Contain(remoteBranches, k) {
			err := r.branchStore.InActivate(ctx, repo.Id, k)
			if err != nil {
				log.WithError(err).Warnf("cannot inactivate  %s branch", k)
			}
		}
	}

	return newBranches, nil
}

func (r *RepositoryService) Scan(ctx context.Context, repoId string) error {
	log := logger.WithAction("scan repository")
	repo, err := r.repoStore.Find(ctx, repoId)
	if err != nil {
		return err
	}

	newBranches, err := r.SyncBranch(ctx, repo)
	if err != nil {
		return err
	}

	auth, err := r.auth(ctx, repo)
	if err != nil {
		log.WithError(err).Errorln("cannot get auth info")
		return err
	}

	for _, branch := range newBranches {
		gitRepo, _, err := r.gitService.Open(repo.Url, repo.Name, branch, auth)
		if err != nil {
			log.WithError(err).Errorln("cannot open repository " + repo.Name)
			continue
		}

		head, err := gitRepo.Head()
		if err != nil {
			log.WithError(err).
				WithField("repo", repo.Name).
				Errorln("cannot get head")
			continue
		}

		h := &core.Hook{
			Trigger:     core.TriggerHook,
			Branch:      branch,
			Event:       core.EventPromote,
			Timestamp:   time.Now().Unix(),
			Title:       head.Message,
			Message:     head.Message,
			After:       head.Hash.String(),
			Ref:         plumbing.NewBranchReferenceName(branch).String(),
			Author:      head.Author.Name,
			AuthorEmail: head.Author.Email,
		}

		go func() {
			_, err := r.triggerService.Trigger(ctx, repo, h)
			if err != nil {
				log.WithError(err).
					WithField("repo", repo.Name).
					WithField("branch", h.Branch).
					Errorln("repo trigger fail")
			}
		}()
	}

	branches, err := r.branchStore.List(ctx, repo.Id)
	if err != nil {
		log.WithError(err).Errorln("cannot get branches")
		return err
	}

	for _, branch := range branches {
		if !utils.Contain(newBranches, branch.Name) {
			_ = r.Build(ctx, repoId, branch.Id, false)
			continue
		}
	}

	return nil
}

func (r *RepositoryService) auth(ctx context.Context, repo *core.Repository) (transport.AuthMethod, error) {
	cred, err := r.credStore.Find(ctx, repo.CredentialId)
	if err != nil {
		return nil, err
	}

	var auth transport.AuthMethod
	if cred.CredentialType == 1 {
		auth, err = r.gitService.AuthAccount(cred.Username, cred.Password)
	} else if cred.CredentialType == 2 {
		auth, err = r.gitService.AuthPrivateKey(cred.PrivateKey)
	} else {
		return nil, errors.New(999999)
	}

	return auth, nil
}

func (r *RepositoryService) Build(ctx context.Context, repoId, branchId string, force bool) error {
	log := logger.WithAction("build repository")
	repo, err := r.repoStore.Find(ctx, repoId)
	if err != nil {
		return err
	}

	branch, err := r.branchStore.Find(ctx, branchId)
	if err != nil {
		log.WithError(err).Errorln("cannot find branch info")
		return err
	}

	auth, err := r.auth(ctx, repo)
	if err != nil {
		log.WithError(err).Errorln("cannot get auth info")
		return err
	}

	gitRepo, _, err := r.gitService.Open(repo.Url, repo.Name, branch.Name, auth)
	if err != nil {
		log.WithError(err).
			WithField("repo", repo.Name).
			WithField("branch", branch.Name).
			Errorln("cannot open repository")
		return err
	}

	commits, err := gitRepo.Pull()
	if err != nil {
		log.WithError(err).
			WithField("branch", branch.Name).
			WithField("repo", repo.Name).
			Errorln("cannot open repository")
	}

	if len(commits) > 0 || force {
		msg := ""
		var changesText []string
		var prev *object.Commit
		for _, commit := range commits {

			if prev == nil {
				prev = commit
				continue
			}
			msg += commit.Message
			tree, err := commit.Tree()
			if err != nil {
				log.WithError(err).
					WithField("commit", commit.Hash).
					WithField("branch", branch.Name).
					WithField("repo", repo.Name).
					Errorln("cannot get commit tree")
				continue
			}
			prevTree, err := prev.Tree()
			changes, err := object.DiffTree(tree, prevTree)
			for _, change := range changes {
				changesText = append(changesText, change.String())
			}

			prev = commit
		}

		head, err := gitRepo.Head()
		if err != nil {
			log.WithError(err).
				WithField("branch", branch.Name).
				WithField("repo", repo.Name).
				Errorln("cannot get head commit")
			return err
		}
		var firstCommit *object.Commit
		if len(commits) > 0 {
			firstCommit = commits[0]
		} else {
			firstCommit = head
		}

		h := &core.Hook{
			Trigger:     core.TriggerHook,
			Branch:      branch.Name,
			Event:       core.EventPromote,
			Timestamp:   time.Now().Unix(),
			Title:       msg,
			Message:     msg,
			Before:      firstCommit.Hash.String(),
			After:       head.Hash.String(),
			Ref:         plumbing.NewBranchReferenceName(branch.Name).String(),
			Author:      head.Author.Name,
			AuthorEmail: head.Author.Email,
			Changes:     changesText,
		}

		go func() {
			_, err := r.triggerService.Trigger(ctx, repo, h)
			if err != nil {
				log.WithError(err).
					WithField("repo", repo.Name).
					WithField("branch", h.Branch).
					Errorln("repo trigger fail")
			}
		}()
	}

	return nil
}
