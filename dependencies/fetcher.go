package dependencies

type DependencyHandler interface {
	FetchDeps(request PackageRequest) (Dependencies, error)
	Act(request PackageRequest) error
}
