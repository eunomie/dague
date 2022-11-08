package types

type BuildOpts struct {
	EnvVars    map[string]string
	BuildFlags []string
	Dir        string
	In         string
}

type LocalBuildOpts struct {
	BuildOpts
	Out string
}

type Platform struct {
	OS   string
	Arch string
}

type CrossBuildOpts struct {
	BuildOpts
	Platforms     []Platform
	OutFileFormat string
}
