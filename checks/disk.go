package checks

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"os"
)

func spaceAvailablePercent(disk *linuxproc.Disk) float64 {
	return (float64(disk.Free) / float64(disk.All) * 100)
}

func diskSpaceCheck(path string) error {
	d, err := linuxproc.ReadDisk(path)
	if err != nil {
		return fmt.Errorf("Cannot read disk info of %s file system.", path)
	}
	if spaceAvailablePercent(d) < 20 {
		return fmt.Errorf("Low free space on %s. Free disk space: %2.1f %%", path, spaceAvailablePercent(d))
	}
	return nil
}

func rootDiskSpaceCheck() error {
	return diskSpaceCheck(baseDir + "/")
}

func mountedDiskSpaceCheck() error {
	path := baseDir + "/vol"
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return nil
	}
	return diskSpaceCheck(path)
}

func DiskChecks(checks *[]fthealth.Check) {
	rootDiskSpaceCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be effected",
		Name:             "Root disk space check.",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "Please clear some disk space on the 'root' mount",
		Checker:          rootDiskSpaceCheck,
	}

	mountedDiskSpaceCheck := fthealth.Check{
		BusinessImpact:   "A part of the publishing workflow might be effected",
		Name:             "Persistent disk space check mounted on '/vol' (always true for stateless nodes)",
		PanicGuide:       "Please refer to technical summary",
		Severity:         2,
		TechnicalSummary: "Please clear some disk space on the 'vol' mount",
		Checker:          mountedDiskSpaceCheck,
	}

	*checks = append(*checks, rootDiskSpaceCheck)
	*checks = append(*checks, mountedDiskSpaceCheck)
}
