/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/sisatech/tablewriter"
	"github.com/spf13/cobra"
	"github.com/vorteil/direkcli/pkg/instance"
	log "github.com/vorteil/direkcli/pkg/log"
	"github.com/vorteil/direkcli/pkg/namespace"
	"github.com/vorteil/direkcli/pkg/workflow"
	"github.com/vorteil/vorteil/pkg/elog"
	"google.golang.org/grpc"
)

var flagNamespace string
var flagInputFile string
var flagGRPC string

var conn *grpc.ClientConn
var logger elog.View

const grpcConnection = "127.0.0.1:6666"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "direkcli",
	Short: "A cli for direktiv that talks directly a grpc server direktiv hosts",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger = log.GetLogger()
		var err error
		connF, err := cmd.Flags().GetString("grpc")
		if err != nil {
			return err
		}

		if connF == "" {
			connF = grpcConnection
		}

		conn, err = grpc.Dial(connF, grpc.WithInsecure())
		if err != nil {
			return err
		}

		return nil
	},
}

// namespaceCmd
var namespaceCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "List, Create and Delete namespaces",
	Long:  ``,
}

// namespaceListCmd
var namespaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of namespaces",
	Long:  ``,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		list, err := namespace.List(conn)
		if err != nil {
			logger.Errorf("%s", err.Error())
			os.Exit(1)
		}
		if len(list) == 0 {
			logger.Printf("No namespaces exist")
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name"})

		for _, namespace := range list {
			table.Append([]string{
				namespace.GetName(),
			})
		}

		table.Render()
	},
}

// namespaceCreateCmd
var namespaceCreateCmd = &cobra.Command{
	Use:   "create [NAME]",
	Short: "Creates a new namespace",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		success, err := namespace.Create(args[0], conn)
		if err != nil {
			logger.Errorf("%s", err.Error())
			os.Exit(1)
		}
		logger.Printf(success)
	},
}

// namespaceDeleteCmd
var namespaceDeleteCmd = &cobra.Command{
	Use:   "delete [NAME]",
	Short: "Deletes a namespace",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		success, err := namespace.Delete(args[0], conn)
		if err != nil {
			logger.Errorf("%s", err.Error())
			os.Exit(1)
		}
		logger.Printf(success)
	},
}

// workflowCmd
var workflowCmd = &cobra.Command{
	Use:   "workflows",
	Short: "List, Create, Get and Execute workflows",
	Long:  ``,
}

// workflowListCmd
var workflowListCmd = &cobra.Command{
	Use:   `list [NAMESPACE]`,
	Short: "Lists all workflows under a namespace",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		list, err := workflow.List(conn, args[0])
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		if len(list) == 0 {
			logger.Printf("No workflows exist under '%s'", args[0])
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID"})

		// Build string array rows
		for _, wf := range list {
			table.Append([]string{
				wf.GetId(),
			})
		}
		table.Render()
	},
}

// workflowGetCmd
var workflowGetCmd = &cobra.Command{
	Use:   "get [NAMESPACE] [ID]",
	Short: "Get yaml from a workflow",
	Args:  cobra.ExactArgs(2),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		success, err := workflow.Get(conn, args[0], args[1])
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}
		logger.Printf(success)
	},
}

// workflowExecuteCmd
var workflowExecuteCmd = &cobra.Command{
	Use:   "execute [NAMESPACE] [ID]",
	Short: "Executes workflow with given ID",
	Args:  cobra.ExactArgs(2),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			logger.Errorf("unable to retrieve input flag")
			os.Exit(1)
		}

		success, err := workflow.Execute(conn, args[0], args[1], input)
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		logger.Printf(success)
	},
}

// workflowAddCmd
var workflowAddCmd = &cobra.Command{
	Use:   "add [NAMESPACE] [WORKFLOW]",
	Short: "Adds a new workflow",
	Args:  cobra.ExactArgs(2),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// args[0] should be namespace, args[1] should be path to the workflow file
		success, err := workflow.Add(conn, args[0], args[1])
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}
		logger.Printf(success)
	},
}

// workflowUpdateCmd
var workflowUpdateCmd = &cobra.Command{
	Use:   "update [NAMESPACE] [ID] [WORKFLOW]",
	Short: "Updates an existing workflow",
	Args:  cobra.ExactArgs(3),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		success, err := workflow.Update(conn, args[0], args[1], args[2])
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}
		logger.Printf(success)
	},
}

// workflowDeleteCmd
var workflowDeleteCmd = &cobra.Command{
	Use:   "delete [NAMESPACE] [ID]",
	Short: "Deletes an existing workflow",
	Args:  cobra.ExactArgs(2),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		success, err := workflow.Delete(conn, args[0], args[1])
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}
		logger.Printf(success)
	},
}

// instanceCmd
var instanceCmd = &cobra.Command{
	Use:   "instances",
	Short: "List, Get and Retrieve Logs for instances",
	Long:  ``,
}

var instanceGetCmd = &cobra.Command{
	Use:   "get [ID]",
	Short: "Get details about a workflow instance",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := instance.Get(conn, args[0])
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}
		logger.Printf("ID: %s", resp.GetId())
		logger.Printf("Input: %s", string(resp.GetInput()))
		logger.Printf("Output: %s", string(resp.GetOutput()))
	},
}

var instanceLogsCmd = &cobra.Command{
	Use:   "logs [ID]",
	Short: "Grabs all logs for the instance ID provided",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logs, err := instance.Logs(conn, args[0])
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}
		for _, log := range logs {
			fmt.Printf("%s", log.GetMessage())
		}
	},
}

var instanceListCmd = &cobra.Command{
	Use:   "list [NAMESPACE]",
	Short: "List all workflow instances in a namespace",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		list, err := instance.List(conn, args[0])
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		if len(list) == 0 {
			logger.Printf("No instances exist under '%s'", args[0])
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Status"})

		// Build string array rows
		for _, instance := range list {
			table.Append([]string{
				instance.GetId(),
				instance.GetStatus(),
			})
		}
		table.Render()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	// Namespace command
	namespaceCmd.AddCommand(namespaceListCmd)
	namespaceCmd.AddCommand(namespaceCreateCmd)
	namespaceCmd.AddCommand(namespaceDeleteCmd)

	// Workflow commands
	workflowCmd.AddCommand(workflowAddCmd)
	workflowCmd.AddCommand(workflowDeleteCmd)
	workflowCmd.AddCommand(workflowListCmd)
	workflowCmd.AddCommand(workflowUpdateCmd)
	workflowCmd.AddCommand(workflowGetCmd)
	workflowCmd.AddCommand(workflowExecuteCmd)

	// Workflow instance commands
	instanceCmd.AddCommand(instanceGetCmd)
	instanceCmd.AddCommand(instanceListCmd)
	instanceCmd.AddCommand(instanceLogsCmd)

	// Root Commands
	rootCmd.AddCommand(namespaceCmd)
	rootCmd.AddCommand(workflowCmd)
	rootCmd.AddCommand(instanceCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&flagGRPC, "grpc", "", "", "ip and port for connection GRPC default is 127.0.0.1:6666")

	// workflowCmd add flag for the namespace
	workflowExecuteCmd.PersistentFlags().StringVarP(&flagInputFile, "input", "", "", "filepath to json input")
}
