package qt

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type QtAdapter struct{}

func New() *QtAdapter {
	return &QtAdapter{}
}

// Working
// TODO can't see the log with this configuration / maybe we should use executeCli for every command?
func execConfigure(execPath string, path string, args ...string) error {
	cmd := exec.Command("/bin/sh", filepath.Join(execPath, "configure"))
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	return err
}


// FIXME we should not make the build fail when the repository is already initialize
func initRepository(pc *config.ProjectConfig, path string) error {
	fmt.Println("PATH: " + filepath.Join(path, "init-repository"))
	cmd := exec.Command("perl", filepath.Join(path, "init-repository"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start %s: %w", cmd, err)
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("%s process failed: %w", cmd, err)
	}
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	return err
}

func executeMakeCli(path string, args ...string) error {
	cmd := exec.Command("make", args...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start make: %w", err)
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("make process failed: %w", err)
	}
	return nil
}

func buildOnLinux(pc *config.ProjectConfig) error {
	// var path = filepath.Join(pc.ProjectRoot, "qt5")
	var pathBuild = filepath.Join(pc.ProjectRoot, "qt6-build")
	// err := initRepository(pc, path)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	// err := os.Mkdir(pathBuild, 0755)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	// err1 := execConfigure(pc, path, pathBuild)
	// if err1 != nil {
	// 	return err1
	// }
	err2 := execCmake(pathBuild)
	if err2 != nil {
		return err2
	}
	return nil
	// return execCmakeInstall(pathBuild)
}

func execCmakeInstall(path string, args ...string) error {
	fmt.Println("PATH: " + filepath.Join(path, "init-repository"))

	cmd := exec.Command("cmake", "--install", filepath.Join(path, "."))
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	return err
	// cmd := exec.Command("cmake", args...)
	// cmd.Dir = path
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Start()
	// if err != nil {
	// 	return fmt.Errorf("failed to start cmake: %w", err)
	// }
	// if err = cmd.Wait(); err != nil {
	// 	return fmt.Errorf("cmake process failed: %w", err)
	// }
	// return nil
}

func execCmake(path string) error {
	cmd := exec.Command("cmake", "--build", filepath.Join(path, "."))
	cmd.Dir = path
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	return err
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Start()
	// if err != nil {
	// 	return fmt.Errorf("failed to start cmake: %w", err)
	// }
	// if err = cmd.Wait(); err != nil {
	// 	return fmt.Errorf("cmake process failed: %w", err)
	// }
	// return nil
}

func executeCli(command string, dirPath string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dirPath
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start %s: %w", command, err)
	}
	cmd.Wait()
	//if err = cmd.Wait(); err != nil {
	//	return fmt.Errorf("%s process failed: %w", command, err)
	//}
	return nil
}

func getCliOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	return string(out), nil
}

func buildOnMacOS(pc *config.ProjectConfig) error {
	err := executeCli("git", pc.ProjectRoot, "clone", "https://github.com/qt/qt5.git")
	if err != nil {
		log.Fatal(err)
		return err
	}
	path := filepath.Join(pc.ProjectRoot, "qt5")
	err = executeCli("git", path, "checkout", "6.0")
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(path)
	pathBuild := filepath.Join(path, "qt6-build")

	err = executeCli("perl", path, filepath.Join(path, "init-repository"))
	if err != nil {
		log.Fatal(err)
		return err
	}

	if _, err = os.Stat(pathBuild); os.IsNotExist(err) {
		err = os.Mkdir(pathBuild, 0755)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	if err != nil {
		log.Fatal(err)
		return err
	}

	macOSSdk, err := getCliOutput("xcrun", "-sdk", "macosx", "--show-sdk-path")
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println(macOSSdk)

	args := []string{
		//configure args
		"-release",
		"-prefix",  pc.ProjectRoot,
		"-extprefix", pc.ProjectRoot,
		"-sysroot", macOSSdk,
		"-cmake-generator", "Ninja",
		"-archdatadir", "share/qt",
		"-examplesdir", "share/qt/example",
		"-testsdir", "share/qt/tests",
		"-libproxy",
		"-no-feature-relocatable",
		"-system-sqlite",
		"-no-sql-mysql",
		"-no-sql-odbc",
		"-no-sql-psql",

		"--",

		// cmake args
		"-DCMAKE_FIND_FRAMEWORK=FIRST",
		"-DINSTALL_MKSPECSDIR=share/qt/mkspecs",
		"-DFEATURE_pkg_config=ON",
	}

	err = executeCli(filepath.Join(path, "configure"), pathBuild, args...)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = executeCli("cmake", pathBuild, "--build", ".", "--parallel", "4")
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = executeCli("cmake", pathBuild, "--install", ".")
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

func (a *QtAdapter) Build(pc *config.ProjectConfig) error {
	fmt.Println(runtime.GOOS)
	switch runtime.GOOS {
	case "darwin":
		err := buildOnMacOS(pc)
		if err != nil {
			return err
		}
		return err
	case "linux":
		err := buildOnLinux(pc)
		if err != nil {
			return err
		}
		return err
	case "windows":
		return nil

	}
	return nil
}

func (a *QtAdapter) CmakeConfig(_ *config.ProjectConfig) (string, error) {
	return CmakeConfig, nil
}

func (a *QtAdapter) Targets(pc *config.ProjectConfig) (targets []string, err error) {
	return nil, nil
}