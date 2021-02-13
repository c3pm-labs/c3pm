package dependencies

type DependencyHandler interface {
	FetchDeps(request PackageRequest) (Dependencies, error)
	Install(request PackageRequest) error
}