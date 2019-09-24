// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	// mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	Clean()
	must(Test())
	must(Amd64())
	must(Arm6())
	must(Arm7())

	return nil
}

func Amd64() error {
	fmt.Println("Building AMD64...")

	return buildForTarget("amd64", "")
}

func Arm6() error {
	fmt.Println("Building ARM6...")

	return buildForTarget("arm", "6")
}

func Arm7() error {
	fmt.Println("Building ARM7...")

	return buildForTarget("arm", "7")
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("dist")
}

func Test() error {
	return exec.Command("go", "test", "./...").Run()
}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildForTarget(arch string, arm string) error {
	fullDist := fmt.Sprintf("dist/%s/", arch+arm)

	os.MkdirAll(fullDist+"cmd", 0700)

	os.Mkdir(fullDist+"storage", 0700)

	if err := exec.Command("cp", "-r", "resources", fullDist+"resources").Run(); err != nil {
		return err
	}

	if err := exec.Command("cp", "-r", ".env.prod", fullDist+".env").Run(); err != nil {
		return err
	}

	if err := exec.Command("cp", "-r", "scripts/install.sh", fullDist).Run(); err != nil {
		return err
	}

	if err := exec.Command("cp", "-r", "scripts/uninstall.sh", fullDist).Run(); err != nil {
		return err
	}

	if err := exec.Command("cp", "-r", "scripts/goairmon.service", fullDist).Run(); err != nil {
		return err
	}

	if err := buildCommand(fullDist+"goairmon", ".", arch, arm); err != nil {
		return err
	}

	if err := buildCommand(fullDist+"cmd/adduser", "./cmd/adduser", arch, arm); err != nil {
		return err
	}
	if err := buildCommand(fullDist+"cmd/rmuser", "./cmd/rmuser", arch, arm); err != nil {
		return err
	}

	tarCmd := exec.Command("tar", "-czf", "dist/goairmon-"+arch+arm+".tar.gz", "-C", fullDist, ".")
	if out, err := tarCmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return err
	}

	return nil
}

func buildCommand(output string, input string, arch string, arm string) error {
	cmd := exec.Command("go", "build", "-o", output, input)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GOARCH="+arch)
	if arm != "" {
		cmd.Env = append(cmd.Env, "GOARM="+arm)
	}

	return cmd.Run()
}
