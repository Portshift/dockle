package privilege

import (
	"fmt"
	"os"

	"github.com/Portshift/dockle/pkg/types"
)

type PrivilegeAssessor struct{}

func (a PrivilegeAssessor) Assess(imageData *types.ImageData) ([]*types.Assessment, error) {
	var assesses []*types.Assessment

	for filename, filedata := range imageData.FileMap {
		if filedata.FileMode&os.ModeSetuid != 0 {
			assesses = append(
				assesses,
				&types.Assessment{
					Code:     types.CheckSuidGuid,
					Filename: filename,
					Desc:     fmt.Sprintf("setuid file: %s %s", filedata.FileMode, filename),
				})
		}
		if filedata.FileMode&os.ModeSetgid != 0 {
			assesses = append(
				assesses,
				&types.Assessment{
					Code:     types.CheckSuidGuid,
					Filename: filename,
					Desc:     fmt.Sprintf("setgid file: %s %s", filedata.FileMode, filename),
				})
		}

	}
	return assesses, nil
}

func (a PrivilegeAssessor) RequiredFiles() []string {
	return []string{}
}

func (a PrivilegeAssessor) RequiredExtensions() []string {
	return []string{}
}

//const GidMode os.FileMode = 4000
func (a PrivilegeAssessor) RequiredPermissions() []os.FileMode {
	return []os.FileMode{os.ModeSetgid, os.ModeSetuid}
}
