package dependencies

type DependencyHandler interface {
	FetchDeps(request PackageRequest) (Dependencies, error)
	PreAct(request PackageRequest) error
	PostAct(request PackageRequest) error
}
