// go build ./01-core-process/main.go && sudo ./main

package main

import (
	"context"
	"log"
	"syscall"

	// "time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func main() {
	if err := containerExample(); err != nil {
		log.Fatal(err)
	}
}

func containerExample() error {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "default")
	// image, err := client.Pull(ctx, "docker.io/library/golang:1.20", containerd.WithPullUnpack)
	image, err := client.Pull(ctx, "docker.io/library/busybox:1.36", containerd.WithPullUnpack)
	if err != nil {
		return err
	}
	log.Printf("Successfully pulled %s image\n", image.Name())

	container, err := client.NewContainer(
		ctx,
		"busybox",
		containerd.WithNewSnapshot("busybox", image),
		containerd.WithNewSpec(
			oci.WithImageConfig(image),
			oci.WithProcessArgs("sleep", "infinity"),
		),
	)
	if err != nil {
		return err
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	defer task.Delete(ctx)
	err = task.Start(ctx)
	if err != nil {
		return err
	}
	defer task.Kill(ctx, syscall.SIGKILL)

	log.Printf("Successfully created container with ID %s and snapshot with ID busybox", container.ID())
	// time.Sleep(60 * time.Second)
	return nil
}
