// go build ./01-core-process/main.go && sudo ./main

package main

import (
	"context"
	"log"

	// "time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func main() {
	if err := busyboxExample(); err != nil {
		log.Fatal(err)
	}
}

func busyboxExample() error {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "default")
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
	log.Printf("Successfully created container with ID %s and snapshot with ID busybox", container.ID())
	// time.Sleep(60 * time.Second)
	return nil
}
