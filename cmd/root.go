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
)

var flagNamespace string
var flagInputFile string

var logger elog.View

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "direkcli",
	Short: "A cli for direktiv that talks directly to nats",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger = log.GetLogger()
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// namespaceCmd
var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage namespaces",
	Long:  ``,
}

// namespaceListCmd
var namespaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists namespaces",
	Long:  ``,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		namespaces, err := namespace.List()
		if err != nil {
			logger.Errorf("Error listing namespaces: %s", err.Error())
		}

		if len(namespaces) == 0 {
			logger.Printf("No namespaces are running on this direktiv")
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Created"})

		// Build string array rows
		for _, namespace := range namespaces {
			timeString := namespace.Created.Format("15:04 02-01-06")
			table.Append([]string{
				namespace.Name,
				timeString,
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
		success, err := namespace.Create(args[0])
		if err != nil {
			logger.Errorf("Error creating namespace: %s", err.Error())
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
		success, err := namespace.Delete(args[0])
		if err != nil {
			logger.Errorf("Error deleting namespace: %s", err.Error())
			os.Exit(1)
		}
		logger.Printf(success)
	},
}

// workflowCmd
var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage workflows",
	Long:  ``,
}

// workflowListCmd
var workflowListCmd = &cobra.Command{
	Use:   `list`,
	Short: "Lists all workflows under a namespace",
	Args:  cobra.ExactArgs(0),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logger.Errorf("unable to retrieve namespace flag")
			os.Exit(1)
		}

		if namespace == "" {
			logger.Errorf("--namespace flag is required to use workflow commands")
			os.Exit(1)
		}

		list, err := workflow.List(namespace)
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		if len(list) == 0 {
			logger.Printf("No workflows are under %s", namespace)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Description", "Created"})

		// Build string array rows
		for _, wf := range list {
			timeString := wf.Created.Format("15:04 02-01-06")
			table.Append([]string{
				wf.ID,
				wf.Name,
				wf.Description,
				timeString,
			})
		}

		table.Render()
	},
}

// workflowGetCmd
var workflowGetCmd = &cobra.Command{
	Use:   "get [ID]",
	Short: "Get yaml from a workflow",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logger.Errorf("unable to retrieve namespace flag")
			os.Exit(1)
		}

		if namespace == "" {
			logger.Errorf("--namespace flag is required to use workflow commands")
			os.Exit(1)
		}

		success, err := workflow.Get(args[0], namespace)
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		logger.Printf(success)
	},
}

// workflowExecuteCmd
var workflowExecuteCmd = &cobra.Command{
	Use:   "execute [ID]",
	Short: "Executes workflow with given ID",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			logger.Errorf("unable to retrieve input flag")
			os.Exit(1)
		}

		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logger.Errorf("unable to retrieve namespace flag")
			os.Exit(1)
		}
		if namespace == "" {
			logger.Errorf("--namespace flag is required to use workflow commands")
			os.Exit(1)
		}

		instanceID, err := workflow.Execute(input, args[0], namespace)
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		logger.Printf("Instance ID of executed workflow '%s'", instanceID)
	},
}

// workflowAddCmd
var workflowAddCmd = &cobra.Command{
	Use:   "add [WORKFLOW]",
	Short: "Adds a new workflow",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logger.Errorf("unable to retrieve namespace flag")
			os.Exit(1)
		}

		if namespace == "" {
			logger.Errorf("--namespace flag is required to use workflow commands")
			os.Exit(1)
		}

		workflowID, err := workflow.Add(args[0], namespace)
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		logger.Printf("Successfully created workflow '%s' under namespace '%s'", workflowID, namespace)
	},
}

// workflowUpdateCmd
var workflowUpdateCmd = &cobra.Command{
	Use:   "update [WORKFLOW] [ID]",
	Short: "Updates an existing workflow",
	Args:  cobra.ExactArgs(2),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logger.Errorf("unable to retriee namespace flag")
		}

		if namespace == "" {
			logger.Errorf("--namespace flag is required to use workflow commands")
			os.Exit(1)
		}

		success, err := workflow.Update(args[0], args[1], namespace)
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		logger.Printf(success)
	},
}

// workflowDeleteCmd
var workflowDeleteCmd = &cobra.Command{
	Use:   "delete [ID]",
	Short: "Deletes an existing workflow",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logger.Errorf("unable to retrieve namespace flag")
		}

		if namespace == "" {
			logger.Errorf("--namespace flag is required to use workflow commands")
			os.Exit(1)
		}

		id, err := workflow.Delete(args[0], namespace)
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		logger.Printf("Sucessfully deleted workflow '%s' under namespace '%s'", id, namespace)
	},
}

// instanceCmd
var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Manage instances",
	Long:  ``,
}

var instanceGetCmd = &cobra.Command{
	Use:   "get [ID]",
	Short: "Get details about a workflow instance",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		wfinstance, err := instance.Get(args[0])
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		logger.Printf("%s", wfinstance)
	},
}

var instanceLogsCmd = &cobra.Command{
	Use:   "logs [ID]",
	Short: "Grabs all logs for the instance ID provided",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logs, err := instance.Logs(args[0])
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		logger.Printf("%s", logs)
	},
}
var instanceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workflow instances in a namespace",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logger.Errorf("unable to retrieve namespace flag")
		}

		if namespace == "" {
			logger.Errorf("--namespace flag is required to use workflow commands")
			os.Exit(1)
		}

		list, err := instance.List(namespace)
		if err != nil {
			logger.Errorf(err.Error())
			os.Exit(1)
		}

		if len(list) == 0 {
			logger.Printf("No workflow instances are under %s", namespace)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Status", "Created"})

		// Build string array rows
		for _, wf := range list {
			timeString := wf.Created.Format("15:04 02-01-06")
			table.Append([]string{
				wf.Name,
				wf.Status,
				timeString,
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

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		logger = log.GetLogger()
		return nil
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// instanceCmd flags
	instanceListCmd.PersistentFlags().StringVarP(&flagNamespace, "namespace", "", "", "defines the namespace to use for queries")

	// workflowCmd add flag for the namespace
	workflowCmd.PersistentFlags().StringVarP(&flagNamespace, "namespace", "", "", "defines the namespace to use for queries")
	workflowExecuteCmd.PersistentFlags().StringVarP(&flagInputFile, "input", "", "", "filepath to json input")
}
