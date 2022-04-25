package gafmysql

import (
	"github.com/rk0cc-xyz/gaf/fetch"
	"github.com/rk0cc-xyz/gaf/storage"
	"github.com/rk0cc-xyz/gaf/structure"
	"github.com/rk0cc-xyz/gafmysql/access"
)

func dlFromAPI() ([][]structure.GitHubRepositoryStructure, error) {
	graf, graferr := fetch.NewGitHubRepositoryAPIFetch()
	if graferr != nil {
		return nil, graferr
	}

	has_next := true
	current_page := 0
	ctx := make([][]structure.GitHubRepositoryStructure, 0)

	for has_next {
		cp := uint64(current_page + 1)

		sctxp, hnx, sctxerr := graf.FetchPage(cp)
		if sctxerr != nil {
			return nil, sctxerr
		}

		ctx = append(ctx, sctxp)
		has_next = *hnx
	}

	return ctx, nil
}

func asmContainer(api [][]structure.GitHubRepositoryStructure) ([]storage.DatabaseFieldContainer, *int64, error) {
	psdfc := make([]storage.DatabaseFieldContainer, len(api))

	var last_page int64

	for page, apiresult := range api {
		adjp := int64(page + 1)
		dfc, dfcerr := storage.CreateDatabaseFieldContainer(adjp, apiresult)

		if dfcerr != nil {
			return nil, nil, dfcerr
		}

		psdfc[page] = *dfc

		last_page = adjp
	}

	return psdfc, &last_page, nil
}

// Save current content of GitHub API to MySQL server.
func ArchiveCurrentAPIToDB() error {
	api, apierr := dlFromAPI()
	if apierr != nil {
		return apierr
	}

	msqlhdl, msqlhdlerr := access.GetMySQLHandlerInstance()

	if msqlhdlerr != nil {
		return msqlhdlerr
	}

	psdfc, lp, psdfcerr := asmContainer(api)
	if psdfcerr != nil {
		return psdfcerr
	}

	for _, convdfc := range psdfc {
		convdfc.SaveToDatabase(msqlhdl)
	}

	msqlhdl.ClearExtraPages(*lp)

	msqlhdl.CloseCurrentSQL()

	return nil
}
