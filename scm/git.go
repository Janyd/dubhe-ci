package scm

import (
	"dubhe-ci/utils"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"path/filepath"
)

var NotEmptyRepositoryPath = errors.New("repository is not empty!")

type (
	GitService interface {

		//构建私钥认证
		AuthPrivateKey(privateKey string) (transport.AuthMethod, error)

		//构建用户名密码认证
		AuthAccount(username, password string) (transport.AuthMethod, error)

		//克隆代码库指定分支
		CloneWithBranch(url, name, branch string, auth transport.AuthMethod) (Repository, error)

		//克隆代码库默认分支master
		Clone(url, name string, auth transport.AuthMethod) (Repository, error)

		//打开本地代码库，如果不存在则将创建
		Open(url, name, branch string, auth transport.AuthMethod) (Repository, bool, error)
	}

	gitService struct {
		Workspace string
	}
)

func New(workspace string) GitService {
	return &gitService{Workspace: workspace}
}

func (g *gitService) AuthPrivateKey(privateKey string) (transport.AuthMethod, error) {
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, err
	}

	auth := &gitssh.PublicKeys{User: "git", Signer: signer, HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}}
	return auth, nil
}

func (g *gitService) AuthAccount(username, password string) (transport.AuthMethod, error) {
	auth := &gitssh.Password{
		User:     username,
		Password: password,
		HostKeyCallbackHelper: gitssh.HostKeyCallbackHelper{
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		},
	}
	return auth, nil
}

func (g *gitService) CloneWithBranch(url, name, branch string, auth transport.AuthMethod) (Repository, error) {
	branchPath := filepath.Join(g.Workspace, name, branch)
	err := utils.PathExistsOrCreate(branchPath)
	if err != nil {
		return nil, err
	}

	if !utils.EmptyDir(branchPath) {
		return nil, NotEmptyRepositoryPath
	}

	repo, err := git.PlainClone(branchPath, false, &git.CloneOptions{
		URL:           url,
		Auth:          auth,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		Progress:      os.Stdout,
	})
	if err != nil {
		return nil, err
	}

	repository := NewRepository(repo, auth, branch)

	return repository, nil
}

func (g *gitService) Clone(url, name string, auth transport.AuthMethod) (Repository, error) {
	return g.CloneWithBranch(url, name, "master", auth)
}

func (g *gitService) Open(url, name, branch string, auth transport.AuthMethod) (Repository, bool, error) {
	branchPath := filepath.Join(g.Workspace, name, branch)
	exist, err := utils.PathExists(branchPath)
	if err != nil {
		return nil, false, err
	}
	empty := utils.EmptyDir(branchPath)

	var repo Repository
	first := false
	if !exist || empty {
		first = true
		repo, err = g.CloneWithBranch(url, name, branch, auth)
		if err != nil {
			return nil, first, err
		}
	} else {
		r, err := git.PlainOpen(branchPath)
		if err != nil {
			return nil, first, err
		}
		repo = NewRepository(r, auth, branch)
	}

	return repo, first, nil
}
