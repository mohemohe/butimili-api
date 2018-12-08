package util

var (
	version string
	hash    string
	branch  string
)

func GetVersion() string {
	if version == "" {
		return "butimili API bleeding edge"
	} else {
		return "butimili API v" + version
	}
}

func GetHash() string {
	if version == "" {
		return "-"
	} else {
		return hash
	}
}

func GetBranch() string {
	if version == "" {
		return "-"
	} else {
		return branch
	}
}
