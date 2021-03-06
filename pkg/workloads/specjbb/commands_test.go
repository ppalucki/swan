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

package specjbb

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	transactionInjectorIndex = 1
	injectionRate            = 6000
	customerNumber           = 100
	duration                 = 100 * time.Second
	productsNumber           = 100
	jarPath                  = "/opt/swan/share/specjbb/specjbb2015.jar"
	propertiesFilePath       = "/opt/swan/share/specjbb/config/specjbb2015.props"
	outputDir                = "/opt/swan/share/specjbb"
	rawFileName              = "abc"
)

func TestCommandsWithDefaultConfig(t *testing.T) {
	Convey("While having default config", t, func() {
		defaultConfig := DefaultLoadGeneratorConfig()

		Convey("and SPECjbb transaction injector command", func() {
			command := getTxICommand(defaultConfig, transactionInjectorIndex)
			Convey("Should contain txinjector mode", func() {
				So(command, ShouldContainSubstring, "-m txinjector")
			})
			Convey("Should contain controller address host property", func() {
				So(command, ShouldContainSubstring, "-Dspecjbb.controller.host=127.0.0.1")
			})
			Convey("Should contain proper group", func() {
				So(command, ShouldContainSubstring, "-G GRP1")
			})
			Convey("Should contain proper JVM id", func() {
				So(command, ShouldContainSubstring, fmt.Sprintf("-J JVM%d", transactionInjectorIndex))
			})
			Convey("Should contain path to binary", func() {
				So(command, ShouldContainSubstring, jarPath)
			})
			Convey("Should contain path to properties file", func() {
				So(command, ShouldContainSubstring, propertiesFilePath)
			})
		})

		Convey("and SPECjbb load command", func() {
			command := getControllerLoadCommand(defaultConfig, injectionRate, duration)
			Convey("Should contain controller mode", func() {
				So(command, ShouldContainSubstring, "-m distcontroller")
			})
			Convey("Should injection controller type property", func() {
				So(command, ShouldContainSubstring, "-Dspecjbb.controller.type=PRESET")
			})
			Convey("Should contain injection rate property", func() {
				So(command, ShouldContainSubstring, fmt.Sprintf("-Dspecjbb.controller.preset.duration=%d", int(duration.Seconds())*1000))
			})
			Convey("Should contain preset duration property", func() {
				So(command, ShouldContainSubstring, fmt.Sprintf("-Dspecjbb.controller.preset.ir=%d", injectionRate))
			})
			Convey("Should contain customer number property", func() {
				So(command, ShouldContainSubstring, fmt.Sprintf("-Dspecjbb.input.number_customers=%d", customerNumber))
			})
			Convey("Should contain product number property", func() {
				So(command, ShouldContainSubstring, fmt.Sprintf("-Dspecjbb.input.number_products=%d", productsNumber))
			})
			Convey("Should contain path to binary", func() {
				So(command, ShouldContainSubstring, jarPath)
			})
			Convey("Should contain path to properties file", func() {
				So(command, ShouldContainSubstring, propertiesFilePath)
			})
			Convey("Should contain path to output dir", func() {
				So(command, ShouldContainSubstring, outputDir)
			})
		})

		Convey("and SPECjbb HBIR RT command", func() {
			command := getControllerTuneCommand(defaultConfig)
			Convey("Should contain controller mode", func() {
				So(command, ShouldContainSubstring, "-m distcontroller")
			})
			Convey("Should injection controller type property", func() {
				So(command, ShouldContainSubstring, "-Dspecjbb.controller.type=HBIR")
			})
			Convey("Should contain customer number property", func() {
				So(command, ShouldContainSubstring, fmt.Sprintf("-Dspecjbb.input.number_customers=%d", customerNumber))
			})
			Convey("Should contain product number property", func() {
				So(command, ShouldContainSubstring, fmt.Sprintf("-Dspecjbb.input.number_products=%d", productsNumber))
			})
			Convey("Should contain path to binary", func() {
				So(command, ShouldContainSubstring, jarPath)
			})
			Convey("Should contain path to properties file", func() {
				So(command, ShouldContainSubstring, propertiesFilePath)
			})
			Convey("Should contain path to output dir", func() {
				So(command, ShouldContainSubstring, outputDir)
			})
		})

		Convey("and SPECjbb reporter command", func() {
			command := getReporterCommand(defaultConfig, rawFileName, 50000)
			Convey("Should contain reporter mode", func() {
				So(command, ShouldContainSubstring, "-m reporter")
			})
			Convey("Should contain path to binary", func() {
				So(command, ShouldContainSubstring, jarPath)
			})
			Convey("Should contain path to output dir", func() {
				So(command, ShouldContainSubstring, rawFileName)
			})
		})
	})
}
