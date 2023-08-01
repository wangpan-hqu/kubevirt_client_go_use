package main
import(
	"fmt"
	"github.com/spf13/pflag"
	"kubevirt.io/client-go/kubecli"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"
	"text/tabwriter"
	virtv1 "kubevirt.io/api/core/v1"
	v1 "kubevirt.io/api/core/v1"
)
func main(){
	// kubecli.DefaultClientConfig() prepares config using kubeconfig.
	// typically, you need to set env variable, KUBECONFIG=<path-to-kubeconfig>/.kubeconfig
	clientConfig := kubecli.DefaultClientConfig(&pflag.FlagSet{})

	// retrive default namespace.
	namespace, _, err := clientConfig.Namespace()
	if err != nil {
		log.Fatalf("error in namespace : %v\n", err)
	}

	// get the kubevirt client, using which kubevirt resources can be managed.
	virtClient, err := kubecli.GetKubevirtClientFromClientConfig(clientConfig)
	if err != nil {
		log.Fatalf("cannot obtain KubeVirt client: %v\n", err)
	}

	// Fetch list of VMs & VMIs
	vmList, err := virtClient.VirtualMachine(namespace).List(&k8smetav1.ListOptions{})
	if err != nil {
		log.Fatalf("cannot obtain KubeVirt vm list: %v\n", err)
	}
	vmiList, err := virtClient.VirtualMachineInstance(namespace).List(&k8smetav1.ListOptions{})
	if err != nil {
		log.Fatalf("cannot obtain KubeVirt vmi list: %v\n", err)
	}

	vm :=&virtv1.VirtualMachine{
		TypeMeta: k8smetav1.TypeMeta{
			Kind: "VirtualMachine",
			APIVersion: "kubevirt.io/v1",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name: "test",
			Namespace: "default",
		},
		Spec: virtv1.VirtualMachineSpec{Running: nil,
			                            RunStrategy: nil,
			                            Flavor: nil,
			                            Preference: nil,
			                            Template: nil,
			                            DataVolumeTemplates: nil},
		Status: virtv1.VirtualMachineStatus{SnapshotInProgress: nil,
			                                RestoreInProgress: nil,
			                                Created: nil,
			                                Ready: nil,
			                                PrintableStatus: nil,
			                                Conditions: nil,
			                                StateChangeRequests: nil,
			                                VolumeRequests: nil,
			                                VolumeSnapshotStatuses: nil,
			                                StartFailure: nil,
			                                MemoryDumpRequest: nil},
	}
	virtClient.VirtualMachine(namespace).Create(vm)
	virtClient.VirtualMachine(namespace).Start("",&v1.StartOptions{})
	virtClient.VirtualMachine(namespace).Delete("vmName",&k8smetav1.DeleteOptions{})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 5, ' ', 0)
	fmt.Fprintln(w, "Type\tName\tNamespace\tStatus")

	for _, vm := range vmList.Items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", vm.Kind, vm.Name, vm.Namespace, vm.Status.Ready)
	}
	for _, vmi := range vmiList.Items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", vmi.Kind, vmi.Name, vmi.Namespace, vmi.Status.Phase)
	}
	w.Flush()

}
