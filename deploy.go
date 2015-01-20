package main

import (
	"fmt"
)

// Build is our basic wrapper for Build operations
type Deploy struct {
	*BasePipeline
	options *GlobalOptions
}

// ToDeploy converts a RawPipeline into a Deploy
func (p *RawPipeline) ToDeploy(options *GlobalOptions) (*Deploy, error) {
	var steps []*Step
	var afterSteps []*Step

	// Start with the secret step, wercker-init that runs before everything
	initStep, err := NewWerckerInitStep(options)
	if err != nil {
		return nil, err
	}
	steps = append(steps, initStep)

	realSteps, err := ExtraRawStepsToSteps(p.RawSteps, options)
	if err != nil {
		return nil, err
	}
	steps = append(steps, realSteps...)

	// For after steps we again need werker-init
	afterSteps = append(afterSteps, initStep)
	realAfterSteps, err := ExtraRawStepsToSteps(p.RawAfterSteps, options)
	if err != nil {
		return nil, err
	}
	afterSteps = append(afterSteps, realAfterSteps...)

	deploy := &Deploy{NewBasePipeline(options, steps, afterSteps), options}
	deploy.InitEnv()
	return deploy, nil
}

// InitEnv sets up the internal state of the environment for the build
func (d *Deploy) InitEnv() {
	env := d.Env()

	a := [][]string{
		[]string{"DEPLOY", "true"},
		[]string{"WERCKER_DEPLOY_ID", d.options.DeployID},
		[]string{"WERCKER_DEPLOY_URL", fmt.Sprintf("%s#deploy/%s", d.options.BaseURL, d.options.DeployID)},
	}

	env.Update(d.CommonEnv())
	env.Update(a)
	env.Update(d.MirrorEnv())
	env.Update(d.PassthruEnv())
}

func (d *Deploy) DockerRepo() string {
	return fmt.Sprintf("%s/%s", d.options.ApplicationOwnerName, d.options.ApplicationName)
}

func (d *Deploy) DockerTag() string {
	tag := d.options.Tag
	if tag == "" {
		tag = fmt.Sprintf("deploy-%s", d.options.DeployID)
	}
	return tag
}

func (d *Deploy) DockerMessage() string {
	message := d.options.Message
	if message == "" {
		message = fmt.Sprintf("Build %s", d.options.DeployID)
	}
	return message
}

// CollectArtifact copies the artifacts associated with the Build.
func (d *Deploy) CollectArtifact(sess *Session) (*Artifact, error) {
	artificer := NewArtificer(d.options)

	// Ensure we have the host directory

	artifact := &Artifact{
		ContainerID:   sess.ContainerID,
		GuestPath:     d.options.GuestPath("output"),
		HostPath:      d.options.HostPath("build.tar"),
		ApplicationID: d.options.ApplicationID,
		DeployID:      d.options.DeployID,
	}

	sourceArtifact := &Artifact{
		ContainerID:   sess.ContainerID,
		GuestPath:     d.options.SourcePath(),
		HostPath:      d.options.HostPath("build.tar"),
		ApplicationID: d.options.ApplicationID,
		DeployID:      d.options.DeployID,
	}

	// Get the output dir, if it is empty grab the source dir.
	fullArtifact, err := artificer.Collect(artifact)
	if err != nil {
		if err == ErrEmptyTarball {
			fullArtifact, err = artificer.Collect(sourceArtifact)
			if err != nil {
				return nil, err
			}
			return fullArtifact, nil
		}
		return nil, err
	}

	return fullArtifact, nil
}