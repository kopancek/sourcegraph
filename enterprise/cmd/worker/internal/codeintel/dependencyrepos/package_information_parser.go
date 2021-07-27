package dependencyrepos

import (
	"strings"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/stores/lsifstore"
	"github.com/sourcegraph/sourcegraph/internal/extsvc"
)

var transformerForScheme map[string]func(lsifstore.Package) (DependencyRepoInfo, string, error) = map[string]func(lsifstore.Package) (DependencyRepoInfo, string, error){
	"semanticdb": transformJVMPackageReference,
}

func transformJVMPackageReference(packageReference lsifstore.Package) (DependencyRepoInfo, string, error) {
	replaced := strings.ReplaceAll(strings.TrimPrefix(packageReference.Name, "maven/"), "/", ":")
	return DependencyRepoInfo{
		Scheme:     packageReference.Scheme,
		Identifier: replaced,
		Version:    packageReference.Version,
	}, extsvc.KindJVMPackages, nil
}
