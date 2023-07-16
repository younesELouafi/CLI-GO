package dialictl

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	dialictl "github.com/younesELouafi/cli/functions"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "start a minikube cluster",
	Run: func(cmd *cobra.Command, args []string) {
		profile, _ := cmd.Flags().GetString("profile")
		fmt.Println("SetUp starting ...")
		memory, _ := cmd.Flags().GetString("memory")
		cpu, _ := cmd.Flags().GetString("cpus")
		_, err := exec.Command("minikube", "start", "--memory", memory, "--cpus", cpu, "-p", profile).Output()
		if err != nil {
			fmt.Println("Error starting Minikube:", err)
			return
		}
		fmt.Println("your management cluster has been created successfully")
		//fmt.Println(string(out))
		dialictl.CrossPlane()
		dialictl.ArgoCD()
		provider, _ := cmd.Flags().GetString("provider")
		dialictl.Setup_EKS(provider)
		time.Sleep(60 * time.Second)
		for {
			if dialictl.PodsReady("crossplane-system") {
				break
			}
			fmt.Println("Waiting for providers to be ready...")
			time.Sleep(60 * time.Second) // Sleep for seconds before checking again
		}
		fmt.Println("Providers are ready")
		fmt.Printf("Your %s is ready for provisioning an %s cluster", profile, provider)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
	//setupCmd.Flags().StringVarP(&clusterName, "profile", "p", "management-cluster", "Name of the management cluster (default: management-cluster)")
	setupCmd.PersistentFlags().StringP("profile", "p", "minikube", "name of your management cluster")
	setupCmd.PersistentFlags().StringP("memory", "m", "2200", "memory")
	setupCmd.PersistentFlags().StringP("cpus", "c", "2", "CPUs")
	setupCmd.PersistentFlags().StringP("provider", "v", "", "Possible values:eks..")
	setupCmd.MarkPersistentFlagRequired("provider")
}
