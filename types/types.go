package types

type LocalBuildOpts struct {
	EnvVars    map[string]string
	BuildFlags []string
	Dir        string
	In         string
	Out        string
}

type Platform struct {
	OS   string
	Arch string
}

type CrossBuildOpts struct {
	Platforms     []Platform
	EnvVars       map[string]string
	BuildFlags    []string
	Dir           string
	In            string
	OutFileFormat string
}
