/*
Copyright 2015 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dockertools

import (
	"fmt"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

// This file contains helper functions to convert docker API types to runtime
// (kubecontainer) types.

func mapStatus(status string) kubecontainer.ContainerStatus {
	// Parse the status string in docker.APIContainers. This could break when
	// we upgrade docker.
	switch {
	case strings.HasPrefix(status, "Up"):
		return kubecontainer.ContainerStatusRunning
	case strings.HasPrefix(status, "Exited"):
		return kubecontainer.ContainerStatusExited
	default:
		return kubecontainer.ContainerStatusUnknown
	}
}

// Converts docker.APIContainers to kubecontainer.Container.
func toRuntimeContainer(c *docker.APIContainers) (*kubecontainer.Container, error) {
	if c == nil {
		return nil, fmt.Errorf("unable to convert a nil pointer to a runtime container")
	}

	dockerName, hash, err := getDockerContainerNameInfo(c)
	if err != nil {
		return nil, err
	}

	return &kubecontainer.Container{
		ID:      kubetypes.DockerID(c.ID).ContainerID(),
		Name:    dockerName.ContainerName,
		Image:   c.Image,
		Hash:    hash,
		Created: c.Created,
		Status:  mapStatus(c.Status),
	}, nil
}

// Converts docker.APIImages to kubecontainer.Image.
func toRuntimeImage(image *docker.APIImages) (*kubecontainer.Image, error) {
	if image == nil {
		return nil, fmt.Errorf("unable to convert a nil pointer to a runtime image")
	}

	return &kubecontainer.Image{
		ID:   image.ID,
		Tags: image.RepoTags,
		Size: image.VirtualSize,
	}, nil
}

// convert RawContainerSTatus to api.ContainerStatus.
func rawToAPIContainerStatus(raw *kubecontainer.RawContainerStatus) *api.ContainerStatus {
	containerID := DockerPrefix + raw.ID.ID
	status := api.ContainerStatus{
		Name:         raw.Name,
		RestartCount: raw.RestartCount,
		Image:        raw.Image,
		ImageID:      raw.ImageID,
		ContainerID:  containerID,
	}
	switch raw.Status {
	case kubecontainer.ContainerStatusRunning:
		status.State.Running = &api.ContainerStateRunning{StartedAt: unversioned.NewTime(raw.StartedAt)}
	case kubecontainer.ContainerStatusExited:
		status.State.Terminated = &api.ContainerStateTerminated{
			ExitCode:    raw.ExitCode,
			Reason:      raw.Reason,
			Message:     raw.Message,
			StartedAt:   unversioned.NewTime(raw.StartedAt),
			FinishedAt:  unversioned.NewTime(raw.FinishedAt),
			ContainerID: containerID,
		}
	default:
		status.State.Waiting = &api.ContainerStateWaiting{}
	}
	return &status
}
