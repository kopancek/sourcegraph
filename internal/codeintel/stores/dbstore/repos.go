package dbstore

import (
	"context"
	"database/sql"

	"github.com/keegancsmith/sqlf"
	"github.com/opentracing/opentracing-go/log"

	"github.com/sourcegraph/sourcegraph/internal/database/basestore"
	"github.com/sourcegraph/sourcegraph/internal/observation"
)

// RepoName returns the name for the repo with the given identifier.
func (s *Store) RepoName(ctx context.Context, repositoryID int) (_ string, err error) {
	ctx, endObservation := s.operations.repoName.With(ctx, &err, observation.Args{LogFields: []log.Field{
		log.Int("repositoryID", repositoryID),
	}})
	defer endObservation(1, observation.Args{})

	name, exists, err := basestore.ScanFirstString(s.Store.Query(ctx, sqlf.Sprintf(repoNameQuery, repositoryID)))
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrUnknownRepository
	}
	return name, nil
}

const repoNameQuery = `
-- source: enterprise/internal/codeintel/stores/dbstore/repos.go:RepoName
SELECT name FROM repo WHERE id = %s
`

type GetJVMDependencyReposOpts struct {
	ArtifactName string
	After        int
	Limit        int
}

type JVMDependencyRepo struct {
	Module  string
	Version string
	ID      int
}

func (s *Store) GetJVMDependencyRepos(ctx context.Context, filter GetJVMDependencyReposOpts) ([]JVMDependencyRepo, error) {
	conds := []*sqlf.Query{sqlf.Sprintf("scheme = %s", "semanticdb")}

	var limit interface{} = "ALL"

	if filter.ArtifactName != "" {
		conds = append(conds, sqlf.Sprintf("name = %s", filter.ArtifactName))
	}

	if filter.After > 0 {
		conds = append(conds, sqlf.Sprintf("id > %d", filter.After))
	}

	if filter.Limit != 0 {
		limit = filter.Limit
	}

	return scanJVMDependencyRepo(s.Query(ctx, sqlf.Sprintf(getLSIFDependencyReposQuery, sqlf.Join(conds, "AND"), limit)))
}

func scanJVMDependencyRepo(rows *sql.Rows, queryErr error) (dependencies []JVMDependencyRepo, err error) {
	if queryErr != nil {
		return nil, queryErr
	}
	defer func() { err = basestore.CloseRows(rows, err) }()

	for rows.Next() {
		var dep JVMDependencyRepo
		if err = rows.Scan(
			&dep.Module,
			&dep.Version,
		); err != nil {
			return nil, err
		}

		dependencies = append(dependencies, dep)
	}

	return dependencies, nil
}

const getLSIFDependencyReposQuery = `
SELECT id, name, version FROM lsif_dependency_repos
WHERE %s ORDER BY id LIMIT %v
`
