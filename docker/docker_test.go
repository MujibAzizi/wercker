//   Copyright 2016 Wercker Holding BV
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package dockerlocal

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/wercker/wercker/util"
)

type DockerSuite struct {
	*util.TestSuite
}

func TestDockerSuite(t *testing.T) {
	suiteTester := &DockerSuite{&util.TestSuite{}}
	suite.Run(t, suiteTester)
}

func (s *DockerSuite) TestPing() {
	client := DockerOrSkip(s.T())
	err := client.Ping()
	s.Nil(err)
}

// Another test to see if password interpolation is blowing up
func (s *DockerSuite) TestPasswordInterpolation() {
	env := util.NewEnvironment("X_PRIVATE=somethingwitha$sign", "XXX_OTHER=somethingwitha#sign")
	testStep := DockerPushStep{
		data: map[string]string{
			"username": "$PRIVATE",
			"password": "$OTHER",
		},
	}
	env.Update(env.GetPassthru().Ordered())
	env.Hidden.Update(env.GetHiddenPassthru().Ordered())
	testStep.InitEnv(env)
	s.Equal("somethingwitha$sign", testStep.username)
	s.Equal("somethingwitha#sign", testStep.password)
}

func (s *DockerSuite) TestGenerateDockerID() {
	id, err := GenerateDockerID()
	s.Require().NoError(err, "Unable to generate Docker ID")

	// The ID needs to be a valid hex value
	b, err := hex.DecodeString(id)
	s.Require().NoError(err, "Generated Docker ID was not a hex value")

	// The ID needs to be 256 bits
	s.Equal(256, len(b)*8)
}
