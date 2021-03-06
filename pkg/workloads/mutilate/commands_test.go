// Copyright (c) 2017 Intel Corporation
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

package mutilate

import (
	"fmt"
	"testing"
	"time"

	"github.com/intelsdi-x/swan/pkg/executor"
	. "github.com/smartystreets/goconvey/convey"
)

func (s *MutilateTestSuite) TestGetPopulateCommand() {
	command := getPopulateCommand(s.mutilate.config)

	Convey("Mutilate populate command should contain mutilate binary path", s.T(), func() {
		expected := fmt.Sprintf("%s", s.mutilate.config.PathToBinary)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate populate command should contain target server host:port", s.T(), func() {
		expected := fmt.Sprintf(
			"-s %s:%d", s.mutilate.config.MemcachedHost, s.mutilate.config.MemcachedPort)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate populate command should contain --loadonly switch", s.T(), func() {
		expected := fmt.Sprintf("--loadonly")
		So(command, ShouldContainSubstring, expected)
	})
}

func (s *MutilateTestSuite) soExpectBaseCommandOptions(command string) {
	Convey("Mutilate base command should contain mutilate binary path", s.T(), func() {
		expected := fmt.Sprintf("%s", s.mutilate.config.PathToBinary)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate base command should contain target server host:port", s.T(), func() {
		expected := fmt.Sprintf(
			"-s %s:%d", s.mutilate.config.MemcachedHost, s.mutilate.config.MemcachedPort)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate base command should contain warmup", s.T(), func() {
		expected := fmt.Sprintf("--warmup %d", int(s.config.WarmupTime.Seconds()))
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate base command should contain noload", s.T(), func() {
		expected := fmt.Sprintf("--noload")
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate base command should contain keySize, valuSize and interArrivalDist option", s.T(), func() {
		expected := fmt.Sprintf("-K %s", s.mutilate.config.KeySize)
		So(command, ShouldContainSubstring, expected)
		expected = fmt.Sprintf("-V %s", s.mutilate.config.ValueSize)
		So(command, ShouldContainSubstring, expected)
		expected = fmt.Sprintf("-i %s", s.mutilate.config.InterArrivalDist)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate base command should contain master threads option", s.T(), func() {
		expected := fmt.Sprintf("-T %d", s.mutilate.config.MasterThreads)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate base command should contain agents connection options", s.T(), func() {
		expected := fmt.Sprintf("-d %d", s.mutilate.config.AgentConnectionsDepth)
		So(command, ShouldContainSubstring, expected)
		expected = fmt.Sprintf("-c %d", s.mutilate.config.AgentConnections)
		So(command, ShouldContainSubstring, expected)
	})

	if s.mutilate.config.MasterQPS != 0 {
		Convey("Mutilate base command should contain masterQPS option", s.T(), func() {
			expected := fmt.Sprintf("-Q %d", s.mutilate.config.MasterQPS)
			So(command, ShouldContainSubstring, expected)
		})
	}
}

func (s *MutilateTestSuite) TestGetLoadCommand() {
	const load = 300
	const duration = 10 * time.Second

	s.mutilate.config.MasterQPS = 0
	s.mutilate.config.Records = 12345
	s.mutilate.config.Update = "0.5"
	command := getLoadCommand(s.mutilate.config, load, duration, []executor.TaskHandle{})

	s.soExpectBaseCommandOptions(command)

	Convey("Mutilate load command should contain load duration", s.T(), func() {
		expected := fmt.Sprintf("-t %d", int(duration.Seconds()))
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate load command should contain qps option", s.T(), func() {
		expected := fmt.Sprintf("-q %d", load)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate load command should contain qps option", s.T(), func() {
		expected := fmt.Sprintf("-q %d", load)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate load command should contain records option", s.T(), func() {
		expected := fmt.Sprintf("-r %d", s.mutilate.config.Records)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate load command should contain update option", s.T(), func() {
		expected := fmt.Sprintf("-u %s", s.mutilate.config.Update)
		So(command, ShouldContainSubstring, expected)
	})
}

func (s *MutilateTestSuite) TestGetMultinodeLoadCommand() {
	const load = 300
	const duration = 10 * time.Second
	const agentAddress1 = "255.255.255.001"
	const agentAddress2 = "255.255.255.002"

	s.mutilate.config.MasterQPS = 0

	s.mAgentHandle1.On("Address").Return(agentAddress1).Once()
	s.mAgentHandle2.On("Address").Return(agentAddress2).Once()
	command := getLoadCommand(s.mutilate.config, load, duration, []executor.TaskHandle{
		s.mAgentHandle1,
		s.mAgentHandle2,
	})

	s.soExpectBaseCommandOptions(command)

	Convey("Mutilate load command should contain load duration", s.T(), func() {
		expected := fmt.Sprintf("-t %d", int(duration.Seconds()))
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate load command should contain qps option", s.T(), func() {
		expected := fmt.Sprintf("-q %d", load)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate load command should contain agents", s.T(), func() {
		expected := fmt.Sprintf("-a %s -a %s", agentAddress1, agentAddress2)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate base command should contain master connection options", s.T(), func() {
		expected := fmt.Sprintf("-D %d", s.mutilate.config.MasterConnectionsDepth)
		So(command, ShouldContainSubstring, expected)
		expected = fmt.Sprintf("-C %d", s.mutilate.config.MasterConnections)
		So(command, ShouldContainSubstring, expected)
	})

	// Check with MasterQPS different to 0.
	s.mutilate.config.MasterQPS = 24234

	s.mAgentHandle1.On("Address").Return(agentAddress1).Once()
	s.mAgentHandle2.On("Address").Return(agentAddress2).Once()
	command = getLoadCommand(s.mutilate.config, load, duration, []executor.TaskHandle{
		s.mAgentHandle1,
		s.mAgentHandle2,
	})

	s.soExpectBaseCommandOptions(command)

	Convey("Assert expectation should be met", s.T(), func() {
		So(s.mAgentHandle1.AssertExpectations(s.T()), ShouldBeTrue)
		So(s.mAgentHandle2.AssertExpectations(s.T()), ShouldBeTrue)
	})
}

func (s *MutilateTestSuite) TestGetTuneCommand() {
	const slo = 300

	s.mutilate.config.MasterQPS = 0
	command := getTuneCommand(s.mutilate.config, slo, []executor.TaskHandle{})

	s.soExpectBaseCommandOptions(command)

	Convey("Mutilate tuning command should contain tuning phase duration", s.T(), func() {
		expected := fmt.Sprintf("-t %d", int(s.mutilate.config.TuningTime.Seconds()))
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate tuning command should contain search option", s.T(), func() {
		expected := fmt.Sprintf("--search %s:%d",
			s.mutilate.config.LatencyPercentile, slo)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Assert expectation should be met", s.T(), func() {
		So(s.mAgentHandle1.AssertExpectations(s.T()), ShouldBeTrue)
		So(s.mAgentHandle2.AssertExpectations(s.T()), ShouldBeTrue)
	})
}
func (s *MutilateTestSuite) TestGetMultinodeTuneCommand() {
	const slo = 300
	const agentAddress1 = "255.255.255.001"
	const agentAddress2 = "255.255.255.002"

	s.mutilate.config.MasterQPS = 0

	s.mAgentHandle1.On("Address").Return(agentAddress1).Once()
	s.mAgentHandle2.On("Address").Return(agentAddress2).Once()
	command := getTuneCommand(s.mutilate.config, slo, []executor.TaskHandle{
		s.mAgentHandle1,
		s.mAgentHandle2,
	})

	s.soExpectBaseCommandOptions(command)

	Convey("Mutilate tuning command should contain tuning phase duration", s.T(), func() {
		expected := fmt.Sprintf("-t %d", int(s.mutilate.config.TuningTime.Seconds()))
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate tuning command should contain search option", s.T(), func() {
		expected := fmt.Sprintf("--search %s:%d",
			s.mutilate.config.LatencyPercentile, slo)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate tuning command should contain agents", s.T(), func() {
		expected := fmt.Sprintf("-a %s -a %s", agentAddress1, agentAddress2)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate tuning command should contain port options", s.T(), func() {
		expected := fmt.Sprintf("-p %d", s.mutilate.config.AgentPort)
		So(command, ShouldContainSubstring, expected)
	})

	// Check with MasterQPS different to 0.
	s.mutilate.config.MasterQPS = 24234

	s.mAgentHandle1.On("Address").Return(agentAddress1).Once()
	s.mAgentHandle2.On("Address").Return(agentAddress2).Once()
	command = getTuneCommand(s.mutilate.config, slo, []executor.TaskHandle{
		s.mAgentHandle1,
		s.mAgentHandle2,
	})

	s.soExpectBaseCommandOptions(command)

	Convey("Mutilate tuning command should contain tuning phase duration", s.T(), func() {
		expected := fmt.Sprintf("-t %d", int(s.mutilate.config.TuningTime.Seconds()))
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate tuning command should contain search option", s.T(), func() {
		expected := fmt.Sprintf("--search %s:%d",
			s.mutilate.config.LatencyPercentile, slo)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Mutilate tuning command should contain agents", s.T(), func() {
		expected := fmt.Sprintf("-a %s -a %s", agentAddress1, agentAddress2)
		So(command, ShouldContainSubstring, expected)
	})

	Convey("Assert expectation should be met", s.T(), func() {
		So(s.mAgentHandle1.AssertExpectations(s.T()), ShouldBeTrue)
		So(s.mAgentHandle2.AssertExpectations(s.T()), ShouldBeTrue)
	})
}

func TestMutilateAffinityCommand(t *testing.T) {
	Convey("Mutilate commands should not contain affinity by default", t, func() {
		config := Config{}
		So(getLoadCommand(config, 0, 0, nil), ShouldNotContainSubstring, "--affinity")
		So(getPopulateCommand(config), ShouldNotContainSubstring, "--affinity")
		So(getTuneCommand(config, 0, nil), ShouldNotContainSubstring, "--affinity")
		So(getAgentCommand(config), ShouldNotContainSubstring, "--affinity")
	})

	Convey("Mutilate master commands should contain affinity when requested with MasterAffinity", t, func() {
		config := Config{MasterAffinity: true}
		So(getLoadCommand(config, 0, 0, nil), ShouldContainSubstring, "--affinity")
		So(getTuneCommand(config, 0, nil), ShouldContainSubstring, "--affinity")
		Convey("but agent and populate commands not", func() {
			So(getPopulateCommand(config), ShouldNotContainSubstring, "--affinity")
			So(getAgentCommand(config), ShouldNotContainSubstring, "--affinity")
		})
	})

	Convey("Mutilate master and populate commands should not contain affinity when requested with AgentAffinity", t, func() {
		config := Config{AgentAffinity: true}
		So(getLoadCommand(config, 0, 0, nil), ShouldNotContainSubstring, "--affinity")
		So(getTuneCommand(config, 0, nil), ShouldNotContainSubstring, "--affinity")
		So(getPopulateCommand(config), ShouldNotContainSubstring, "--affinity")
		Convey("but agent commands should", func() {
			So(getAgentCommand(config), ShouldContainSubstring, "--affinity")
		})
	})
}

func TestMutilateBlockingFlag(t *testing.T) {
	Convey("Mutilate commands should contain -B by default", t, func() {
		config := DefaultMutilateConfig()
		So(getLoadCommand(config, 0, 0, nil), ShouldContainSubstring, "-B")
		So(getTuneCommand(config, 0, nil), ShouldContainSubstring, "-B")
		So(getAgentCommand(config), ShouldContainSubstring, "-B")
	})

	Convey("Mutilate master commands should not contain blocking when requested with MasterBlocking set to false", t, func() {
		config := DefaultMutilateConfig()
		config.MasterBlocking = false
		So(getLoadCommand(config, 0, 0, nil), ShouldNotContainSubstring, "-B")
		So(getTuneCommand(config, 0, nil), ShouldNotContainSubstring, "-B")
		Convey("but agent command still should", func() {
			So(getAgentCommand(config), ShouldContainSubstring, "-B")
		})
	})

	Convey("Mutilate agent commands should not contain blocking when requested with AgentBlocking set to false", t, func() {
		config := DefaultMutilateConfig()
		config.AgentBlocking = false
		So(getAgentCommand(config), ShouldNotContainSubstring, "-B")
		Convey("but master commands still should", func() {
			So(getLoadCommand(config, 0, 0, nil), ShouldContainSubstring, "-B")
			So(getTuneCommand(config, 0, nil), ShouldContainSubstring, "-B")
		})
	})

	Convey("Mutilate populate commands should not contain affinity regardless to MasterBlocking and AgentBlocing", t, func() {
		config := DefaultMutilateConfig()
		So(getPopulateCommand(config), ShouldNotContainSubstring, "-B")

		Convey("but master commands still should", func() {
			config.MasterBlocking = false
			config.AgentBlocking = false
			So(getPopulateCommand(config), ShouldNotContainSubstring, "-B")
		})
	})

}
