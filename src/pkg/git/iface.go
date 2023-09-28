package git

type Git interface {
	FetchCommitSHA(owner, repo, ref string) (string, error)
	GetContents(owner, repo, ref string) (string, error)
}
