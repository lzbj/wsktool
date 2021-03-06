/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"bytes"
	"fmt"
	"github.com/apache/incubator-openwhisk-wskdeploy/utils"
	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"os"
	"os/exec"
	"strings"
)

const cmd = "wskdeploy"

type Wskdeploy struct {
	Path string
	Dir  string
}

func NewWskdeploy() *Wskdeploy {
	return NewWskWithPath(os.Getenv("GOPATH") + "/src/github.com/apache/incubator-openwhisk-wskdeploy/")
}

func NewWskWithPath(path string) *Wskdeploy {
	var dep Wskdeploy
	dep.Path = cmd
	dep.Dir = path
	return &dep
}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err string) {
	outputStream := colorable.NewColorableStderr()
	if len(err) > 0 {
		fmt.Fprintf(outputStream, "==> Error: %s.\n", color.RedString(err))
	} else {
		fmt.Fprintf(outputStream, "==> Error: %s.\n", color.RedString("No error message"))
	}
}

func printOutput(outs string) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s.\n", outs)
	}
}

func (wskdeploy *Wskdeploy) RunCommand(s ...string) (string, error) {
	command := exec.Command(wskdeploy.Path, s...)
	command.Dir = wskdeploy.Dir

	fmt.Println("wskdeploy.Path is " + wskdeploy.Path)
	//fmt.Println("s is " + string(s))
	printCommand(command)

	var outb, errb bytes.Buffer
	command.Stdout = &outb
	command.Stderr = &errb
	err := command.Run()

	var returnError error = nil
	if err != nil {
		returnError = err
	} else {
		if len(errb.String()) > 0 {
			returnError = utils.NewTestCaseError(errb.String())
		}
	}
	printOutput(outb.String())
	if returnError != nil {
		printError(returnError.Error())
	}
	return outb.String(), returnError
}

func (wskdeploy *Wskdeploy) Deploy(manifestPath string, deploymentPath string) (string, error) {
	return wskdeploy.RunCommand("-m", manifestPath, "-d", deploymentPath)
}

func (wskdeploy *Wskdeploy) Undeploy(manifestPath string, deploymentPath string) (string, error) {
	return wskdeploy.RunCommand("undeploy", "-m", manifestPath, "-d", deploymentPath)
}

func (wskdeploy *Wskdeploy) DeployProjectPathOnly(projectPath string) (string, error) {
	return wskdeploy.RunCommand("-p", projectPath)
}

func (wskdeploy *Wskdeploy) UndeployProjectPathOnly(projectPath string) (string, error) {
	return wskdeploy.RunCommand("undeploy", "-p", projectPath)

}

func (wskdeploy *Wskdeploy) DeployManifestPathOnly(manifestpath string) (string, error) {
	return wskdeploy.RunCommand("-m", manifestpath)
}

func (wskdeploy *Wskdeploy) UndeployManifestPathOnly(manifestpath string) (string, error) {
	return wskdeploy.RunCommand("undeploy", "-m", manifestpath)

}
