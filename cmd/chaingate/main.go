// Command chaingate attributes dependencies introduced by AI coding agents and
// applies an admission policy before those changes are merged.
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Wapiti08/chaingate/internal/audit"
	"github.com/Wapiti08/chaingate/internal/dependency"
	"github.com/Wapiti08/chaingate/internal/model"
	"github.com/Wapiti08/chaingate/internal/policy"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "chaingate:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "check" {
		fmt.Fprintln(stderr, "usage: chaingate check --before OLD_GO_MOD --after NEW_GO_MOD --events TRACE.jsonl [--policy POLICY.json]")
		return errors.New("expected the check command")
	}

	flags := flag.NewFlagSet("check", flag.ContinueOnError)
	flags.SetOutput(stderr)
	beforePath := flags.String("before", "", "path to the go.mod before the agent run")
	afterPath := flags.String("after", "", "path to the go.mod after the agent run")
	eventsPath := flags.String("events", "", "path to dependency-relevant JSONL events")
	policyPath := flags.String("policy", "", "optional JSON admission policy")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if *beforePath == "" || *afterPath == "" || *eventsPath == "" {
		return errors.New("--before, --after, and --events are required")
	}

	before, err := loadGoMod(*beforePath)
	if err != nil {
		return fmt.Errorf("load before manifest: %w", err)
	}
	after, err := loadGoMod(*afterPath)
	if err != nil {
		return fmt.Errorf("load after manifest: %w", err)
	}
	events, err := loadEvents(*eventsPath)
	if err != nil {
		return err
	}
	admission, err := loadPolicy(*policyPath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(audit.AddedDependencies(before, after, events, admission))
}

func loadGoMod(path string) (map[string]model.Dependency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return dependency.ParseGoMod(file)
}

func loadEvents(path string) ([]model.Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open events: %w", err)
	}
	defer file.Close()

	events := make([]model.Event, 0)
	scanner := bufio.NewScanner(file)
	for line := 1; scanner.Scan(); line++ {
		var event model.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, fmt.Errorf("decode events line %d: %w", line, err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read events: %w", err)
	}
	return events, nil
}

func loadPolicy(path string) (policy.Policy, error) {
	admission := policy.Default()
	if path == "" {
		return admission, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return admission, fmt.Errorf("open policy: %w", err)
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(&admission); err != nil {
		return admission, fmt.Errorf("decode policy: %w", err)
	}
	return admission, nil
}
