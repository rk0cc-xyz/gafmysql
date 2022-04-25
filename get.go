package gafmysql

import (
	"github.com/rk0cc-xyz/gaf/storage"
	"github.com/rk0cc-xyz/gaf/structure"
	"github.com/rk0cc-xyz/gafmysql/access"
)

func GetArchivedRepositoryAPI() ([]structure.GitHubRepositoryStructure, *string, error) {
	msqlinst, msqlinsterr := access.GetMySQLHandlerInstance()
	if msqlinsterr != nil {
		return nil, nil, msqlinsterr
	}

	mp, _ := msqlinst.GetMaxPage()

	dfca := make([]storage.DatabaseFieldContainer, 0)

	for cfp := int64(1); cfp <= *mp; cfp++ {
		dfc, dfcerr := storage.GetFieldContainerFromDatabase(cfp, msqlinst)
		if dfcerr != nil {
			return nil, nil, dfcerr
		}
		dfca = append(dfca, *dfc)
	}

	msqlinst.CloseCurrentSQL()

	return joinAllContent(dfca)
}

func joinAllContent(containers []storage.DatabaseFieldContainer) ([]structure.GitHubRepositoryStructure, *string, error) {
	ghrc := make([]structure.GitHubRepositoryStructure, 0)

	for _, c := range containers {
		ctx, ctxerr := c.GetContent()
		if ctxerr != nil {
			return nil, nil, ctxerr
		}

		ghrc = append(ghrc, ctx...)
	}

	ru := containers[len(containers)-1].GetUpdatedAt()

	return ghrc, &ru, nil
}
