package metaparticle

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/go-openapi/swag"
	"github.com/metaparticle-io/metaparticle-ast/models"
)

type MetaparticleExecutor struct {
	SpecPath string
}

func (m *MetaparticleExecutor) Cancel(name string) error {
	cmd := exec.Command("mp-compiler", "-f", m.SpecPath, "--delete")
	return cmd.Run()
}

func (m *MetaparticleExecutor) Logs(name string, stdout io.Writer, stderr io.Writer) error {
	cmd := exec.Command("mp-compiler", "-f", m.SpecPath, "--deploy=false", "--attach=true")
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	return cmd.Run()
}

func makeReplicatedService(image string, name string, config *Runtime) models.Service {
	serveSpec := &models.ServeSpecification{Name: &name, Public: true}

	envProps := []*models.EnvVar{&models.EnvVar{Name: swag.String("METAPARTICLE_IN_CONTAINER"), Value: swag.String("true")}}

	containers := []*models.Container{&models.Container{Image: &image, Env: envProps}}

	ports := []*models.ServicePort{}
	for _, port := range config.Ports {
		ports = append(ports, &models.ServicePort{Number: swag.Int32(port), Protocol: "TCP"})
	}

	serviceSpecs := []*models.ServiceSpecification{&models.ServiceSpecification{Name: &name, Replicas: config.Replicas, Containers: containers, Ports: ports}}

	return models.Service{Name: &name, GUID: swag.Int64(1234567), Services: serviceSpecs, Serve: serveSpec}
}

func makeShardedService(image string, name string, config *Runtime) models.Service {
	serveSpec := &models.ServeSpecification{Name: &name, Public: true}

	envProps := []*models.EnvVar{&models.EnvVar{Name: swag.String("METAPARTICLE_IN_CONTAINER"), Value: swag.String("true")}}

	containers := []*models.Container{&models.Container{Image: &image, Env: envProps}}

	ports := []*models.ServicePort{}
	for _, port := range config.Ports {
		ports = append(ports, &models.ServicePort{Number: swag.Int32(port), Protocol: "TCP"})
	}

	shardSpec := &models.ShardSpecification{Shards: config.Shards, URLPattern: config.URLShardPattern}

	serviceSpecs := []*models.ServiceSpecification{&models.ServiceSpecification{Name: &name, ShardSpec: shardSpec, Containers: containers, Ports: ports}}

	return models.Service{Name: &name, GUID: swag.Int64(1234567), Services: serviceSpecs, Serve: serveSpec}
}

func (m *MetaparticleExecutor) Run(image string, name string, config *Runtime, stdout io.Writer, stderr io.Writer) error {
	if len(image) == 0 {
		return fmt.Errorf("An image must be specified")
	}

	if len(name) == 0 {
		return fmt.Errorf("The container's name must be specified")
	}

	if config == nil {
		return fmt.Errorf("config must not be nil")
	}

	var s models.Service
	if config.Shards > 0 {
		s = makeShardedService(image, name, config)
	} else {
		s = makeReplicatedService(image, name, config)
	}
	serviceJSON, err := json.Marshal(s)

	if err != nil {
		return fmt.Errorf("Could not create service json: %v", err)
	}

	specJSONFile, err := ioutil.TempFile("", "spec.json")
	if err != nil {
		return fmt.Errorf("Could not create temporary service json file: %v", err)
	}

	defer os.Remove(specJSONFile.Name())

	if _, err := specJSONFile.Write(serviceJSON); err != nil {
		return fmt.Errorf("Coult not write to temporary service json file")
	}
	if err := specJSONFile.Close(); err != nil {
		return fmt.Errorf("Could not close temporary service json file")
	}

	m.SpecPath = specJSONFile.Name()

	cmd := exec.Command("mp-compiler", "-f", specJSONFile.Name())
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	return cmd.Run()
}
