package manifest

type FilesConfig struct {
	Sources             []string `yaml:"sources"`
	Includes            []string `yaml:"includes"`
	IncludeDirs         []string `yaml:"include_dirs"`
	ExportedDir         string   `yaml:"exported_dir"`
	ExportedIncludeDirs []string `yaml:"exported_include_dirs"`
}

func (f *FilesConfig) IsEmpty() bool {
	return len(f.Sources) == 0 && len(f.Includes) == 0 && len(f.IncludeDirs) == 0 && len(f.ExportedDir) == 0
}
