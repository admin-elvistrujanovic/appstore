package apps

import (
	"encoding/json"
	"os"
	"sort"

	"github.com/hashicorp/go-version"
)

type RepoIndex struct {
	Repo     map[string]interface{}   `json:"repo"`
	Requests map[string]interface{}   `json:"requests"`
	Apps     []map[string]interface{} `json:"apps"`

	Packages map[string][]PackageInfo `json:"packages"`
}

type PackageInfo struct {
	Added            int64    `json:"added"`
	ApkName          string   `json:"apkName"`
	Hash             string   `json:"hash"`
	HashType         string   `json:"hashType"`
	MinSdkVersion    int      `json:"minSdkVersion"`
	Nativecode       []string `json:"nativecode"`
	PackageName      string   `json:"packageName"`
	Sig              string   `json:"sig"`
	Signer           string   `json:"signer"`
	Size             int      `json:"size"`
	TargetSdkVersion int      `json:"targetSdkVersion"`
	VersionCode      int      `json:"versionCode,omitempty"`
	VersionName      string   `json:"versionName"`
}

func (r *RepoIndex) FindLatestPackage(pkgName string) (p PackageInfo, ok bool) {
	pkgs, ok := r.Packages[pkgName]
	if !ok {
		return p, false
	}

	sort.Slice(pkgs, func(i, j int) bool {
		v1, err := version.NewVersion(pkgs[i].VersionName)
		if err != nil {
			return true
		}

		v2, err := version.NewVersion(pkgs[i].VersionName)
		if err != nil {
			return false
		}

		return v1.LessThan(v2)
	})

	// Return the one with the latest version
	return pkgs[len(pkgs)-1], true
}

func (r *RepoIndex) RemoveVersionCode() {
	for i := range r.Apps {
		delete(r.Apps[i], "suggestedVersionCode")
	}

	for pkg := range r.Packages {
		for i := range r.Packages[pkg] {
			r.Packages[pkg][i].VersionCode = 0
		}
	}
}

func ReadIndex(path string) (index *RepoIndex, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&index)

	return
}

func WriteIndex(path string, index *RepoIndex) (err error) {
	tmpPath := path + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return
	}

	err = json.NewEncoder(f).Encode(index)
	if err != nil {
		_ = f.Close()
		return
	}

	err = f.Close()
	if err != nil {
		return
	}

	return os.Rename(tmpPath, path)
}