// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

func main() {
	cmd := &cobra.Command{
		Use: "release",
	}
	cmd.AddCommand(operator(), lemming())
	if err := cmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}

func operator() *cobra.Command {
	return &cobra.Command{
		Use:  "operator",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Validating working directory")
			sha, err := validateWorkDir()
			if err != nil {
				return err
			}
			fmt.Println("Running prerelease tests")
			if err := triggerBuild(cmd.Context(), "operator-test", sha, false); err != nil {
				return err
			}

			tag := fmt.Sprintf("operator/%s", args[0])
			fmt.Println("Creating and Pushing Tag")
			if err := createAndPushTag(tag); err != nil {
				return err
			}
			fmt.Println("Building and Pushing container")
			if err := triggerBuild(cmd.Context(), "operator-release", tag, true); err != nil {
				return err
			}
			return nil
		},
	}
}

func lemming() *cobra.Command {
	return &cobra.Command{
		Use:  "lemming",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Validating working directory")
			sha, err := validateWorkDir()
			if err != nil {
				return err
			}
			fmt.Println("Running prerelease tests")
			if err := triggerBuild(cmd.Context(), "lemming-test", sha, false); err != nil {
				return err
			}

			tag := args[0]
			fmt.Println("Creating and Pushing Tag")
			if err := createAndPushTag(tag); err != nil {
				return err
			}
			fmt.Println("Building and Pushing container")
			if err := triggerBuild(cmd.Context(), "lemming-release", tag, true); err != nil {
				return err
			}
			return nil
		},
	}
}

func createAndPushTag(tag string) error {
	if out, err := exec.Command("git", "tag", tag).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create tag: out %s, error %v", string(out), err)
	}
	if out, err := exec.Command("git", "push", "origin", tag).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to push tag: out %s, error %v", string(out), err)
	}
	return nil
}

const (
	// cloudBuildEndpoint is the regional endpoint for the cloud build API.
	// Related to issue https://github.com/googleapis/google-cloud-go/issues/5095.
	cloudBuildEndpoint = "us-central1-cloudbuild.googleapis.com:443"
	triggerNamePrefix  = "projects/openconfig-lemming/locations/us-central1/triggers"
)

// triggerBuild runs a cloud build trigger at the given tag if set, or the main branch if unset.
func triggerBuild(ctx context.Context, trigger, tagOrSHA string, tag bool) error {
	c, err := cloudbuild.NewClient(ctx, option.WithEndpoint(cloudBuildEndpoint))
	if err != nil {
		return err
	}
	defer c.Close()

	src := &cloudbuildpb.RepoSource{
		Revision: &cloudbuildpb.RepoSource_CommitSha{
			CommitSha: tagOrSHA,
		},
	}
	if tag {
		src = &cloudbuildpb.RepoSource{
			Revision: &cloudbuildpb.RepoSource_TagName{
				TagName: tagOrSHA,
			},
		}
	}

	op, err := c.RunBuildTrigger(ctx, &cloudbuildpb.RunBuildTriggerRequest{
		Name:   fmt.Sprintf("%s/%s", triggerNamePrefix, trigger),
		Source: src,
	})
	if err != nil {
		return err
	}
	if _, err := op.Poll(ctx); err != nil {
		return err
	}
	md, err := op.Metadata()
	if err != nil {
		return err
	}
	fmt.Printf("Build ID: %s\n Logs: %s", md.GetBuild().GetId(), md.GetBuild().GetLogUrl())
	fmt.Println("Waiting for Op")
	b, err := op.Wait(ctx)
	if err != nil {
		return err
	}
	fmt.Println(b.Id, b.Status)
	return nil
}

// validateWorkDir checks the status of the working dir to make sure it is clean state.
func validateWorkDir() (string, error) {
	stOut, err := exec.Command("git", "status", "--porcelain").CombinedOutput()
	if err != nil {
		return "", err
	}
	status := strings.TrimSpace(string(stOut))
	brOut, err := exec.Command("git", "branch", "--show-current").CombinedOutput()
	if err != nil {
		return "", err
	}
	branch := strings.TrimSpace(string(brOut))
	revOut, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		return "", err
	}
	sha := strings.TrimSpace(string(revOut))
	ready := true
	if branch != "main" {
		fmt.Println("Not on main branch")
		ready = false
	}
	if status != "" {
		fmt.Println("Working directory dirty")
		ready = false
	}
	if !ready {
		ok, err := promptBool("Are you sure you want to continue")
		if err != nil {
			return "", err
		}
		if !ok {
			return "", fmt.Errorf("repository in invalid state")
		}
	}
	return sha, nil
}

// promptBool is a yes/no command line prompt.
func promptBool(prompt string) (bool, error) {
	fmt.Print(prompt + " (y/n): ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case "y":
			return true, nil
		case "n":
			return false, nil
		default:
			fmt.Println("invalid input")
		}
	}
	return false, scanner.Err()
}
