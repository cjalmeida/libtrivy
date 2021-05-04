package scan

import (
	"encoding/json"
	"os"

	ftypes "github.com/aquasecurity/fanal/types"
	"github.com/aquasecurity/trivy-db/pkg/db"
	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy/pkg/commands/config"
	"github.com/aquasecurity/trivy/pkg/commands/operation"
	"github.com/aquasecurity/trivy/pkg/detector/ospkg"
	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/aquasecurity/trivy/pkg/scanner/local"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/aquasecurity/trivy/pkg/utils"
	"github.com/aquasecurity/trivy/pkg/vulnerability"
	"golang.org/x/xerrors"
)

type FileApplier struct{}

type Config struct {
	CacheDir     string
	CacheBackend string
	DBConfig     config.DBConfig
}

var (
	version      = "v0.17.2"
	alreadySetup = false
)

func Scan(sourceFile, destFile string) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	// defer cancel()

	c := Config{
		CacheDir:     utils.DefaultCacheDir(),
		CacheBackend: "fs",
	}

	// setup
	if !alreadySetup {
		err := setup(c)
		if err != nil {
			return err
		}
		alreadySetup = true
	}

	applier := FileApplier{}
	detector := ospkg.Detector{}
	scanner := local.NewScanner(applier, detector)
	options := types.ScanOptions{
		VulnType: []string{"os", "library"},
	}

	// Do scan
	res, _, _, err := scanner.Scan(sourceFile, sourceFile, nil, options)
	if err != nil {
		return err
	}

	// Fill info
	vulnClient := vulnerability.NewClient(db.Config{})
	for i := range res {
		vulnClient.FillInfo(res[i].Vulnerabilities, res[i].Type)
	}

	outStream := os.Stdout
	if destFile != "" {
		outStream, err = os.Create(destFile)
		if err != nil {
			return err
		}
	}

	// write results
	if err = report.WriteResults("json", outStream, []dbTypes.Severity{},
		res, "", false); err != nil {
		return err
	}

	outStream.Close()

	return nil
}

func setup(c Config) error {

	// configure cache dir
	utils.SetCacheDir(c.CacheDir)
	cache, err := operation.NewCache(c.CacheBackend)
	if err != nil {
		return xerrors.Errorf("server cache error: %w", err)
	}
	defer cache.Close()

	// download the database file
	if err = operation.DownloadDB(version, c.CacheDir, true, false, false); err != nil {
		return err
	}

	if err = db.Init(c.CacheDir); err != nil {
		return xerrors.Errorf("error in vulnerability DB initialize: %w", err)
	}

	return nil
}

func (applier FileApplier) ApplyLayers(artifactID string, blobIDs []string) (detail ftypes.ArtifactDetail, err error) {

	file, err := os.Open(artifactID)
	if err != nil {
		return ftypes.ArtifactDetail{}, err
	}

	dec := json.NewDecoder(file)

	var artifact ftypes.ArtifactDetail
	err = dec.Decode(&artifact)
	if err != nil {
		return ftypes.ArtifactDetail{}, err
	}

	return artifact, nil
}
