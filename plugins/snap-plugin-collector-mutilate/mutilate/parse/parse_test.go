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

package parse

import (
	"bytes"
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStdoutParser(t *testing.T) {
	Convey("Opening non-existing file should fail", t, func() {
		data, err := File("/non/existing/file")

		So(data.Raw, ShouldBeEmpty)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "open /non/existing/file: no such file or directory")
	})

	Convey("Opening readable and correct file should provide meaningful results", t, func() {
		path, err := getCurrentDirFilePath("mutilate.stdout")
		So(err, ShouldBeNil)

		data, err := File(path)

		So(err, ShouldBeNil)
		So(data.Raw, ShouldHaveLength, 10)
		So(data.Raw[MutilateAvg], ShouldResemble, 20.8)
		So(data.Raw[MutilateStd], ShouldResemble, 23.1)
		So(data.Raw[MutilateMin], ShouldResemble, 11.9)
		So(data.Raw[MutilatePercentile5th], ShouldResemble, 13.3)
		So(data.Raw[MutilatePercentile10th], ShouldResemble, 13.4)
		So(data.Raw[MutilatePercentile90th], ShouldResemble, 33.4)
		So(data.Raw[MutilatePercentile95th], ShouldResemble, 43.1)
		So(data.Raw[MutilatePercentile99th], ShouldResemble, 59.5)
		So(data.Raw[MutilateQPS], ShouldResemble, 4993.1)
		So(data.Raw[MutilateMisses], ShouldEqual, 1234)
	})

	Convey("Attempting to read file with wrong number of read columns should return an error and no metrics", t, func() {
		path, err := getCurrentDirFilePath("mutilate_incorrect_count_of_columns.stdout")
		So(err, ShouldBeNil)

		data, err := File(path)

		So(data.Raw, ShouldHaveLength, 0)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "Incorrect number of fields: expected 8 but got 2")
	})

	Convey("Attempting to read a file with no read row at all should return no metrics", t, func() {
		path, err := getCurrentDirFilePath("mutilate_missing_read_row.stdout")
		So(err, ShouldBeNil)

		data, err := File(path)
		So(err, ShouldBeNil)

		// QPS and Misses are stil available, thus 2.
		So(data.Raw, ShouldHaveLength, 2)

		So(data.Raw, ShouldNotContainKey, MutilateAvg)
		So(data.Raw, ShouldNotContainKey, MutilateStd)
		So(data.Raw, ShouldNotContainKey, MutilateMin)
		So(data.Raw, ShouldNotContainKey, MutilatePercentile5th)
		So(data.Raw, ShouldNotContainKey, MutilatePercentile10th)
		So(data.Raw, ShouldNotContainKey, MutilatePercentile90th)
		So(data.Raw, ShouldNotContainKey, MutilatePercentile95th)
		So(data.Raw, ShouldNotContainKey, MutilatePercentile99th)
	})

	Convey("Attempting to read a file with read row containing incorrect values should return an error and no results", t, func() {
		path, err := getCurrentDirFilePath("/mutilate_non_numeric_default_metric_value.stdout")
		So(err, ShouldBeNil)

		data, err := File(path)

		So(data.Raw, ShouldHaveLength, 0)
		So(err.Error(), ShouldStartWith, "'thisIsNotANumber' latency value must be a float")
	})

	Convey("Trying to reorder columns and parsing should pass", t, func() {
		in := bytes.NewReader([]byte("#type min 5th\nread 5 10\n"))
		data, err := Parse(in)
		So(err, ShouldBeNil)
		So(data.Raw, ShouldHaveLength, 2)
		So(data.Raw[MutilateMin], ShouldResemble, 5.0)

		in = bytes.NewReader([]byte("#type 5th min\nread 10 5"))
		data, err = Parse(in)
		So(err, ShouldBeNil)
		So(data.Raw, ShouldHaveLength, 2)
		So(data.Raw[MutilateMin], ShouldResemble, 5.0)

		in = bytes.NewReader([]byte("#type 99th 5th min\nread 90 50 4.5"))
		data, err = Parse(in)
		So(err, ShouldBeNil)
		So(data.Raw, ShouldHaveLength, 3)
		So(data.Raw[MutilatePercentile99th], ShouldResemble, 90.0)
		So(data.Raw[MutilatePercentile5th], ShouldResemble, 50.0)
	})
}

func getCurrentDirFilePath(name string) (string, error) {
	gwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(gwd, name), nil
}
