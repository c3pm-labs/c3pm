package ctpm

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/c3pm-labs/c3pm/dependencies"
	"strings"
)

type ListOptions struct {
	Tree bool
}

var ListDefaultOptions = ListOptions{
	Tree: false,
}

type DependencyNode struct {
	Children []*DependencyNode
	Name     string
}

func (d DependencyNode) String() string {
	sb := strings.Builder{}
	sb.WriteString(d.Name)
	sb.WriteString(": ")
	sb.WriteString("[")
	for i := 0; i < len(d.Children); i++ {
		sb.WriteString(d.Children[i].Name)
		if i < len(d.Children)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

type DependencyBranch []*DependencyNode

func (d DependencyBranch) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for i := 0; i < len(d); i++ {
		sb.WriteString(d[i].Name)
		if i < len(d)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

type DependencyFetcher struct {
	Done          manifest.Dependencies
	CurrentNode   *DependencyNode
	CurrentBranch DependencyBranch
}

func (d *DependencyFetcher) FetchDeps(request dependencies.PackageRequest) (dependencies.Dependencies, error) {
	if _, ok := d.Done[request.Name]; ok {
		return dependencies.Dependencies{}, nil
	}
	// Find new current node
	var currentNode *DependencyNode = nil
	found := false
	for !found {
		for _, node := range d.CurrentNode.Children {
			if node.Name == fmt.Sprintf("%s@%s", request.Name, request.Version) {
				found = true
				d.CurrentBranch = append(d.CurrentBranch, node)
				currentNode = node
			}
		}
		if !found {
			// Move up one level
			d.CurrentBranch = d.CurrentBranch[:len(d.CurrentBranch)-1]
			d.CurrentNode = d.CurrentBranch[len(d.CurrentBranch)-1]
		}
	}
	d.Done[request.Name] = request.Version
	libPath := config.LibCachePath(request.Name, request.Version)
	pc, err := config.Load(libPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	ret := make(dependencies.Dependencies)
	for k, v := range pc.Manifest.Dependencies {
		ret[k] = v
		currentNode.Children = append(currentNode.Children, &DependencyNode{Name: fmt.Sprintf("%s@%s", k, v), Children: []*DependencyNode{}})
	}
	d.CurrentNode = currentNode
	return ret, nil
}

func (d *DependencyFetcher) PreAct(_ dependencies.PackageRequest) error {
	return nil
}

func (d *DependencyFetcher) PostAct(_ dependencies.PackageRequest) error {
	return nil
}

func displayTree(root *DependencyNode, under string) {
	fmt.Printf("───%s\n", root.Name)
	for i := 0; i < len(root.Children); i++ {
		if i < len(root.Children)-1 {
			fmt.Printf("%s├", under)
			displayTree(root.Children[i], fmt.Sprintf("%s%s", under, "│   "))
		} else {
			fmt.Printf("%s└", under)
			displayTree(root.Children[i], fmt.Sprintf("%s%s", under, "    "))
		}
	}
}

func List(pc *config.ProjectConfig, opt ListOptions) error {
	var topLevelDeps []*DependencyNode
	for name, version := range pc.Manifest.Dependencies {
		topLevelDeps = append(topLevelDeps, &DependencyNode{Name: fmt.Sprintf("%s@%s", name, version), Children: []*DependencyNode{}})
	}
	treeRoot := &DependencyNode{Name: fmt.Sprintf("%s@%s", pc.Manifest.Name, pc.Manifest.Version), Children: topLevelDeps}

	allDeps := make(manifest.Dependencies)
	for name, version := range pc.Manifest.Dependencies {
		_, err := dependencies.Install(dependencies.PackageRequest{Name: name, Version: version}, &DependencyFetcher{Done: allDeps, CurrentNode: treeRoot, CurrentBranch: []*DependencyNode{treeRoot}})
		if err != nil {
			return err
		}
	}

	if !opt.Tree {
		for key, val := range allDeps {
			fmt.Printf("%s@%s\n", key, val)
		}
	} else {
		displayTree(treeRoot, "   ")
	}

	return nil
}
