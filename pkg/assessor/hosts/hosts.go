package hosts

import (
	"os"

	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/Portshift/dockle/pkg/log"
	"github.com/Portshift/dockle/pkg/types"
)

type HostsAssessor struct{}

func (a HostsAssessor) Assess(fileMap deckodertypes.FileMap) ([]*types.Assessment, error) {
	log.Logger.Debug("Start scan : /etc/hosts")

	assesses := []*types.Assessment{}
	// TODO : check hosts setting
	return assesses, nil
}

func (a HostsAssessor) RequiredFiles() []string {
	return []string{"etc/hosts"}
}

func (a HostsAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{}
}
