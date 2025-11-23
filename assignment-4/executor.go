package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	DESC string   `json:"desc"`
	CMD  string   `json:"cmd"`
	CWD  string   `json:"cwd"`
	DEPS []string `json:"deps"`
}

type ConfigFile map[string]Task

func main()  {
	// cmd := exec.Command("ls")

	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Printf("command failed: %v\noutput: \n%s\n", err, string(output))
	// 	return
	// }

	// fmt.Printf("output:\n%s\n", string(output))

	f, err := os.Open("cli-tasks.json")
	if err != nil {
		fmt.Printf("Reading file error %s", err)
		return
	}
	defer f.Close()

	var cfg ConfigFile
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields() // helps catch typos in JSON keys
	if err := dec.Decode(&cfg); err != nil {
		fmt.Println("json decode error:", err)
		return
	}

	// for name, task := range cfg {
	// 	fmt.Printf("Task %s: %+v\n", name, task)
	// }

	// Split CMD into command + args
	parts := strings.Fields(cfg["lint"].CMD) 
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Env = os.Environ()
	cmd.Dir = cfg["lint"].CWD

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("command failed: %v\noutput: \n%s\n", err, string(out))
		return
	}
	fmt.Printf("output:\n%s\n", string(out))

}