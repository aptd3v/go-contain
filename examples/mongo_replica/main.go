// This example demonstrates how to create a MongoDB replica set with x members
// and a Mongo Express instance to manage the replica set.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/aptd3v/go-contain/pkg/client"
	"github.com/aptd3v/go-contain/pkg/client/options/container/execopt"
	"github.com/aptd3v/go-contain/pkg/compose"
	"github.com/aptd3v/go-contain/pkg/compose/options/up"
	"github.com/aptd3v/go-contain/pkg/create"
	"github.com/aptd3v/go-contain/pkg/create/config/cc"
	"github.com/aptd3v/go-contain/pkg/create/config/cc/health"
	"github.com/aptd3v/go-contain/pkg/create/config/hc"
	"github.com/aptd3v/go-contain/pkg/create/config/nc"
	"github.com/aptd3v/go-contain/pkg/create/config/sc"
	"github.com/aptd3v/go-contain/pkg/tools"
)

const NumReplicas = 3 // Number of MongoDB replicas in the replica set

type RSet struct {
	ID      string     `json:"_id"`
	Members []RSMember `json:"members"`
}

type RSMember struct {
	ID   int    `json:"_id"`
	Host string `json:"host"`
}

func main() {

	project := create.NewProject("mongo-db-cluster")

	members := []RSMember{}
	urlParts := make([]string, 0, NumReplicas)
	for i := range NumReplicas {
		serviceName := fmt.Sprintf("db-%d", i)

		project.WithService(serviceName,
			WithMongoReplica(i),
			// depends on the previous  db-0 <- db-1 <- db-2
			tools.WhenTrue(i > 0,
				sc.WithDependsOn(fmt.Sprintf("db-%d", i-1)),
			),
		)
		members = append(members, RSMember{
			Host: serviceName,
			ID:   i,
		})
		urlParts = append(urlParts, fmt.Sprintf("%s:27017", serviceName))
	}
	url := fmt.Sprintf("mongodb://%s/?replicaSet=rs0", strings.Join(urlParts, ","))

	project.WithService("mongo-express", WithMongoExpress(url))

	project.WithNetwork("mongo-cluster").WithVolume("mongo-data")

	database := compose.NewCompose(project)

	if err := database.Up(
		context.Background(),
		up.WithForceRecreate(),
		up.WithRemoveOrphans(),
		up.WithWait(), // Wait for the containers to be healthy before initializing the replica set
	); err != nil {
		log.Fatal(err)
	}

	err := Initialize(context.Background(), "mongodb-0", members)
	if err != nil {
		log.Fatal(err)
	}
	signalsChan := make(chan os.Signal, 1)
	signal.Notify(signalsChan, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-signalsChan
		defer cancel()
		fmt.Println("Received interrupt signal, shutting down...")
		if err := database.Down(context.Background()); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	err = project.Export("./examples/mongo_replica/docker-compose.yaml", 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB replica set", "rs0", "url", url)
	fmt.Println("MongoDB replica set initialized and running.")
	fmt.Println("You can access Mongo Express at http://localhost:8081")
	fmt.Println("Press Ctrl+C to stop the containers.")
	<-ctx.Done()
}

func WithMongoReplica(index int) *create.Container {
	containerName := fmt.Sprintf("mongodb-%d", index)
	return create.NewContainer(containerName).
		WithContainerConfig(
			cc.WithImage("mongo:latest"),
			cc.WithCommand("mongod", "--replSet", "rs0", "--bind_ip_all"),
			cc.WithHealthCheck(
				health.WithTest("CMD", "mongosh", "--eval", `db.adminCommand("ping")`),
				health.WithInterval("1s"),
				health.WithTimeout("10s"),
				health.WithStartPeriod("0s"),
				health.WithRetries(5),
			),
			cc.WithExposedPort("tcp", "27017"),
		).
		WithHostConfig(
			hc.WithRestartPolicyUnlessStopped(),
		).
		WithNetworkConfig(
			nc.WithEndpoint("mongo-cluster"),
		)
}

// Initialize initializes the MongoDB replica set with the provided members.
// It runs the `rs.initiate` command in the specified container.
func Initialize(ctx context.Context, initContainer string, members []RSMember) error {
	cli, err := client.NewClient(
		tools.WhenTrue(os.Getenv("DOCKER_HOST") != "",
			client.FromEnv(),
		),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	if len(members) == 0 {
		return fmt.Errorf("no members provided for replica set initialization")
	}
	// Prepare the members for the rs.initiate command
	initiate := RSet{
		ID:      "rs0",
		Members: members,
	}

	init, err := json.Marshal(initiate)
	if err != nil {
		return fmt.Errorf("failed to marshal members: %w", err)
	}
	command := []string{"mongosh", "--eval", fmt.Sprintf("rs.initiate(%s)", string(init))}

	fmt.Println(strings.Join(command, " "))
	res, err := cli.ContainerExecCreate(
		ctx,
		initContainer,
		execopt.WithCommand(command...),
		execopt.WithAttachStdout(),
		execopt.WithAttachStderr(),
	)
	if err != nil {
		return fmt.Errorf("failed to create exec command: %w", err)
	}
	attached, err := cli.ContainerExecAttach(ctx, res.ID)
	if err != nil {
		return fmt.Errorf("failed to start exec command: %w", err)
	}
	defer attached.Close()
	_, _ = io.Copy(os.Stdout, attached.Reader)
	return nil
}

func WithMongoExpress(url string) *create.Container {
	return create.NewContainer("mongo-express").
		WithContainerConfig(
			cc.WithImage("mongo-express:latest"),
			cc.WithEnv("ME_CONFIG_MONGODB_URL", url),
			cc.WithEnv("ME_CONFIG_MONGODB_AUTH_USERNAME", "admin"),
			cc.WithEnv("ME_CONFIG_MONGODB_AUTH_PASSWORD", "password"),
		).
		WithHostConfig(
			hc.WithPortBindings("tcp", "0.0.0.0", "8081", "8081"),
			hc.WithRestartPolicyAlways(),
		).
		WithNetworkConfig(
			nc.WithEndpoint("mongo-cluster"),
		)
}
