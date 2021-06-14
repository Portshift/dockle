package pkg

import (
	"context"
	"fmt"
	l "log"
	"os"
	"strings"

	"github.com/containers/image/v5/transports/alltransports"
	deckodertypes "github.com/goodwithtech/deckoder/types"
	"github.com/urfave/cli"

	"github.com/Portshift/dockle/config"
	"github.com/Portshift/dockle/pkg/log"
	"github.com/Portshift/dockle/pkg/report"
	"github.com/Portshift/dockle/pkg/scanner"
	"github.com/Portshift/dockle/pkg/types"
	"github.com/Portshift/dockle/pkg/utils"
)

func RunFromCli(c *cli.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.Duration("timeout"))
	defer cancel()

	cliVersion := "v" + c.App.Version
	latestVersion, err := utils.FetchLatestVersion(ctx)

	// check latest version
	if err != nil {
		log.Logger.Infof("Failed to check latest version. %s", err)
	} else if cliVersion != latestVersion && c.App.Version != "dev" {
		log.Logger.Warnf("A new version %s is now available! You have %s.", latestVersion, cliVersion)
	}

	config.CreateFromCli(c)
	_, err = run(ctx)

	return err
}

func RunFromConfig(conf *config.Config) (types.AssessmentMap, error) {
	config.Conf = *conf
	ctx, cancel := context.WithTimeout(context.Background(), config.Conf.Timeout)
	defer cancel()

	return run(ctx)
}

func run(ctx context.Context) (ret types.AssessmentMap, err error) {
	debug := config.Conf.Debug
	if err = log.InitLogger(debug); err != nil {
		l.Fatal(err)
	}

	// set docker option
	dockerOption := deckodertypes.DockerOption{
		Timeout:               config.Conf.Timeout,
		UserName:              config.Conf.Username,
		Password:              config.Conf.Password,
		InsecureSkipTLSVerify: config.Conf.Insecure,
		SkipPing:              true,
	}

	var useLatestTag bool
	// Check whether 'latest' tag is used
	if config.Conf.ImageName != "" {
		if useLatestTag, err = useLatest(config.Conf.ImageName); err != nil {
			return nil, fmt.Errorf("invalid image: %w", err)
		}
	}
	log.Logger.Debug("Start assessments...")

	assessments, err := scanner.ScanImage(ctx, config.Conf.ImageName, config.Conf.FilePath, dockerOption)
	if err != nil {
		return nil, err
	}
	if useLatestTag {
		assessments = append(assessments, &types.Assessment{
			Code:     types.AvoidLatestTag,
			Filename: "image tag",
			Desc:     "Avoid 'latest' tag",
		})
	}

	log.Logger.Debug("End assessments...")

	assessmentMap := types.CreateAssessmentMap(assessments, config.Conf.IgnoreMap)
	// Store ignore checkpoint code
	o := config.Conf.Output
	output := os.Stdout
	if o != "" {
		if output, err = os.Create(o); err != nil {
			return nil, fmt.Errorf("failed to create an output file: %w", err)
		}
	}

	var writer report.Writer
	switch format := config.Conf.Format; format {
	case "json":
		writer = &report.JsonWriter{Output: output, ImageName: config.Conf.ImageName}
	case "sarif":
		writer = &report.SarifWriter{Output: output}
	default:
		writer = &report.ListWriter{Output: output}
	}

	abend, err := writer.Write(assessmentMap)
	if err != nil {
		return nil, fmt.Errorf("failed to write results: %w", err)
	}

	if config.Conf.ExitCode != 0 && abend {
		os.Exit(config.Conf.ExitCode)
	}

	return assessmentMap, nil
}

func useLatest(imageName string) (bool, error) {
	ref, err := alltransports.ParseImageName("docker://" + imageName)
	if err != nil {
		return false, err

	}
	return strings.HasSuffix(ref.DockerReference().String(), ":latest"), nil
}
