package collector

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"syscall"
)

var IgnoredFSTypes = regexp.MustCompile("^(autofs|binfmt_misc|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|iso9660|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs)$")

type StorageInfo struct {
	Size  uint64 `json:"size"`
	Avail uint64 `json:"avail"`
	Used  uint64 `json:"used"`
}

func GetStorage() (map[string]string, error) {
	mp := make(map[string]string)

	mountpoints, err := GetMountPoints()
	if err != nil {
		return nil, err
	}

	for _, mountpoint := range mountpoints {
		var stat syscall.Statfs_t

		syscall.Statfs(mountpoint, &stat)
		size := stat.Blocks * uint64(stat.Bsize)
		avail := stat.Bavail * uint64(stat.Bsize)
		used := size - avail

		mpstat := StorageInfo{
			Size:  size,
			Avail: avail,
			Used:  used,
		}

		m, _ := json.Marshal(mpstat)

		mp[mountpoint] = string(m)

	}

	return mp, nil
}

func GetMountPoints() ([]string, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseMountPoints(file)
}

func ParseMountPoints(file *os.File) ([]string, error) {
	var mountpoints []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if !IgnoredFSTypes.Match([]byte(parts[2])) {
			mountpoints = append(mountpoints, parts[1])
		}
	}

	return mountpoints, nil
}
