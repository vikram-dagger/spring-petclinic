package main

import (
	"context"
	"dagger/spring-petclinic/internal/dagger"
	"fmt"
	"math"
	"math/rand"
)

type SpringPetclinic struct{}

func (m *SpringPetclinic) Build(ctx context.Context, source *dagger.Directory) *dagger.File {
	return dag.Java().
		WithJdk("17").
		WithMaven("3.9.5").
		WithProject(source.WithoutDirectory("dagger")).
		Maven([]string{"package"}).
		File("target/spring-petclinic-3.4.0-SNAPSHOT.jar")
}

func (m *SpringPetclinic) Publish(ctx context.Context, source *dagger.Directory) (string, error) {
	return dag.Container(dagger.ContainerOpts{Platform: "linux/amd64"}).
		From("eclipse-temurin:17-alpine").
		WithFile("/app/spring-petclinic-3.4.0-SNAPSHOT.jar", m.Build(ctx, source)).
		WithEntrypoint([]string{"java", "-jar", "/app/spring-petclinic-3.4.0-SNAPSHOT.jar"}).
		Publish(ctx, fmt.Sprintf("ttl.sh/spring-petclinic-%.0f", math.Floor(rand.Float64()*10000000))) //#nosec
}
