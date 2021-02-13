package manifest

// FilesConfig holdes configuration about the files used by the project.
type FilesConfig struct {
	// Sources lists the project's source files.
	Sources []string `yaml:"sources"`
	// Includes lists the project's header files.
	Includes []string `yaml:"includes"`
	// IncludeDirs lists the projects additional header directories.
	IncludeDirs []string `yaml:"include_dirs"`
	// ExportedDir is the path to the directory containing exported headers.
	ExportedDir         string   `yaml:"exported_dir"`
	ExportedIncludeDirs []string `yaml:"exported_include_dirs"`
}

//IsEmpty checks if a FilesConfig structure holds no information
func (f *FilesConfig) IsEmpty() bool {
	return len(f.Sources) == 0 && len(f.Includes) == 0 && len(f.IncludeDirs) == 0 && len(f.ExportedDir) == 0
}
