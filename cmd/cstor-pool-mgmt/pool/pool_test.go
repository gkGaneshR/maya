package pool

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/golang/glog"
	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type TestRunner struct{}

func (r TestRunner) RunCombinedOutput(command string, args ...string) ([]byte, error) {
	var cs []string
	var cmd *exec.Cmd
	cs = append(cs, args...)
	cmd = exec.Command(os.Args[0], cs...)
	switch args[0] {
	case "create":
		cs = []string{"-test.run=TestCreaterProcess", "--"}
		cmd.Env = []string{"createErr=nil"}
		break
	case "import":
		cs = []string{"-test.run=TestImporterProcess", "--"}
		cmd.Env = []string{"importErr=nil"}
		break
	case "destroy":
		cs = []string{"-test.run=TestDestroyerProcess", "--"}
		cmd.Env = []string{"destroyErr=nil"}
		break
	case "labelclear":
		cs = []string{"-test.run=TestLabelClearerProcess", "--"}
		cmd.Env = []string{"labelClearErr=nil"}
		break
	case "status":
		cs = []string{"-test.run=TestStatusProcess", "--"}
		cmd.Env = []string{"StatusErr=nil"}
		break
	}
	stdout, err := cmd.CombinedOutput()
	return stdout, err
}

func (r TestRunner) RunStdoutPipe(command string, args ...string) ([]byte, error) {
	var cs []string
	var cmd *exec.Cmd
	cs = append(cs, args...)
	cmd = exec.Command(os.Args[0], cs...)
	switch args[0] {
	case "get":
		cs = []string{"-test.run=TestGetterProcess", "--"}
		cmd.Env = []string{"poolName=cstor-123abc"}
		break
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		glog.Errorf(err.Error())
		return []byte{}, err
	}
	if err := cmd.Start(); err != nil {
		glog.Errorf(err.Error())
		return []byte{}, err
	}
	data, _ := ioutil.ReadAll(stdout)
	if err := cmd.Wait(); err != nil {
		glog.Errorf(err.Error())
		return []byte{}, err
	}
	return data, nil
}
func TestCreaterProcess(*testing.T) {
	if os.Getenv("createErr") != "nil" {
		return
	}
	fmt.Println(nil)
	defer os.Exit(0)

}
func TestImporterProcess(*testing.T) {
	if os.Getenv("importErr") != "nil" {
		return
	}
	defer os.Exit(0)
}
func TestGetterProcess(*testing.T) {
	if os.Getenv("poolName") != "cstor-123abc" {
		return
	}
	defer os.Exit(0)
	fmt.Println("cstor-123abc")
}
func TestDestroyerProcess(*testing.T) {
	if os.Getenv("destroyErr") != "nil" {
		return
	}
	defer os.Exit(0)
	fmt.Println(nil)
}
func TestLabelClearerProcess(*testing.T) {
	if os.Getenv("labelClearErr") != "nil" {
		return
	}
	defer os.Exit(0)
	fmt.Println(nil)
}
func TestStatusProcess(*testing.T) {
	if os.Getenv("StatusErr") != "nil" {
		return
	}
	defer os.Exit(0)
	fmt.Println(nil)
}

func TestCreatePool(t *testing.T) {
	testPoolResource := map[string]struct {
		expectedError error
		test          *apis.CStorPool
	}{
		"img1PoolResource": {
			expectedError: nil,
			test: &apis.CStorPool{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					UID: types.UID("abc"),
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img1.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool1.cache",
						PoolType:         "mirror",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{},
			},
		},
	}
	RunnerVar = TestRunner{}
	obtainedErr := CreatePool(testPoolResource["img1PoolResource"].test)
	if testPoolResource["img1PoolResource"].expectedError != obtainedErr {
		t.Fatalf("Expected: %v, Got: %v", testPoolResource["img1PoolResource"].expectedError, obtainedErr)
	}
}

func TestImportPool(t *testing.T) {
	testPoolResource := map[string]struct {
		expectedError error
		test          *apis.CStorPool
	}{
		"img1PoolResource": {
			expectedError: nil,
			test: &apis.CStorPool{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					UID: types.UID("abc"),
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img1.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool1.cache",
						PoolType:         "mirror",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{},
			},
		},
	}
	RunnerVar = TestRunner{}
	obtainedErr := ImportPool(testPoolResource["img1PoolResource"].test)
	if testPoolResource["img1PoolResource"].expectedError != obtainedErr {
		t.Fatalf("Expected: %v, Got: %v", testPoolResource["img1PoolResource"].expectedError, obtainedErr)
	}
}

func TestDeletePool(t *testing.T) {
	testPoolResource := map[string]struct {
		expectedError error
		poolName      string
	}{
		"img1PoolResource": {
			expectedError: nil,
			poolName:      "pool1-a2b",
		},
	}
	RunnerVar = TestRunner{}
	obtainedErr := DeletePool(testPoolResource["img1PoolResource"].poolName)
	if testPoolResource["img1PoolResource"].expectedError != obtainedErr {
		t.Fatalf("Expected: %v, Got: %v", testPoolResource["img1PoolResource"].expectedError, obtainedErr)
	}
}

func TestLabelClear(t *testing.T) {
	testResource := map[string]struct {
		expectedError error
		disks         []string
	}{
		"Resource1": {
			expectedError: nil,
			disks:         []string{"/dev/sdb1"},
		},
	}
	RunnerVar = TestRunner{}
	obtainedErr := LabelClear(testResource["Resource1"].disks)
	if testResource["Resource1"].expectedError != obtainedErr {
		t.Fatalf("Expected: %v, Got: %v", testResource["Resource1"].expectedError, obtainedErr)
	}
}

func TestCheckForZrepl(t *testing.T) {
	done := make(chan bool)
	RunnerVar = TestRunner{}
	go func(done chan bool) {
		CheckForZrepl()
		done <- true
	}(done)

	select {
	case <-time.After(3 * time.Second):
		t.Fatalf("Zrepl test failure - Timed out")
	case <-done:

	}
}
func TestGetPool(t *testing.T) {
	testPoolResource := map[string]struct {
		expectedPoolName string
		expectedError    error
	}{
		"img1PoolResource": {
			expectedPoolName: "cstor-123abc\n",
			expectedError:    nil,
		},
	}
	RunnerVar = TestRunner{}
	obtainedPoolName, obtainedErr := GetPoolName()
	fmt.Println(obtainedPoolName, obtainedErr)
	if testPoolResource["img1PoolResource"].expectedPoolName != obtainedPoolName {
		t.Fatalf("Expected: %v, Got: %v", testPoolResource["img1PoolResource"].expectedPoolName, obtainedPoolName)
	}
	if testPoolResource["img1PoolResource"].expectedError != obtainedErr {
		t.Fatalf("Expected: %v, Got: %v", testPoolResource["img1PoolResource"].expectedError, obtainedErr)
	}
}

// TestCheckValidPool tests pool related operations
func TestCheckValidPool(t *testing.T) {
	testPoolResource := map[string]struct {
		expectedPoolName string
		expectedError    error
		test             *apis.CStorPool
	}{
		"Invalid-poolNameEmpty": {
			expectedPoolName: "",
			expectedError:    fmt.Errorf("Poolname cannot be empty"),
			test: &apis.CStorPool{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					UID: types.UID(""),
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img1.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool1.cache",
						PoolType:         "mirror",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{},
			},
		},
		"Invalid-DiskListEmpty": {
			expectedPoolName: "",
			expectedError:    fmt.Errorf("Disk name(s) cannot be empty"),
			test: &apis.CStorPool{
				TypeMeta: v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{
					UID: types.UID("abc"),
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool1.cache",
						PoolType:         "mirror",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{},
			},
		},
	}

	for desc, ut := range testPoolResource {
		Obtainederr := CheckValidPool(ut.test)
		if Obtainederr != nil {
			if Obtainederr.Error() != ut.expectedError.Error() {
				t.Fatalf("Desc : %v, Expected error: %v, Got : %v",
					desc, ut.expectedError, Obtainederr)
			}
		}

	}
}
