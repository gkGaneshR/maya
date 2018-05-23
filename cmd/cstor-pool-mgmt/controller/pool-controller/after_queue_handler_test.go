package poolcontroller

import (
	"os"
	"testing"
	"time"

	apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
	openebsFakeClientset "github.com/openebs/maya/pkg/client/clientset/versioned/fake"
	informers "github.com/openebs/maya/pkg/client/informers/externalversions"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
)

// TestGetPoolResource checks if pool resource created is successfully got.
func TestGetPoolResource(t *testing.T) {
	fakeKubeClient := fake.NewSimpleClientset()
	fakeOpenebsClient := openebsFakeClientset.NewSimpleClientset()

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(fakeKubeClient, time.Second*30)
	openebsInformerFactory := informers.NewSharedInformerFactory(fakeOpenebsClient, time.Second*30)

	// Instantiate the cStor Pool controllers.
	poolController := NewCStorPoolController(fakeKubeClient, fakeOpenebsClient, kubeInformerFactory,
		openebsInformerFactory)

	testPoolResource := map[string]struct {
		expectedPoolName string
		test             *apis.CStorPool
	}{
		"img1PoolResource": {
			expectedPoolName: "abc",
			test: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: "pool1",
					UID:  types.UID("abc"),
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
		"img2PoolResource": {
			expectedPoolName: "abcd",
			test: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name: "pool2",
					UID:  types.UID("abcd"),
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{},
			},
		},
	}
	for desc, ut := range testPoolResource {
		// Create Pool resource
		_, err := poolController.clientset.OpenebsV1alpha1().CStorPools().Create(ut.test)
		if err != nil {
			t.Fatalf("Desc:%v, Unable to create resource : %v", desc, ut.test.ObjectMeta.Name)
		}
		// Get the created pool resource using name
		cStorPoolObtained, err := poolController.getPoolResource(ut.test.ObjectMeta.Name)
		if string(cStorPoolObtained.ObjectMeta.UID) != ut.expectedPoolName {
			t.Fatalf("Desc:%v, PoolName mismatch, Expected:%v, Got:%v", desc, ut.expectedPoolName,
				string(cStorPoolObtained.ObjectMeta.UID))
		}
	}
}

// TestRemoveFinalizer is to remove pool resource.
func TestRemoveFinalizer(t *testing.T) {
	fakeKubeClient := fake.NewSimpleClientset()
	fakeOpenebsClient := openebsFakeClientset.NewSimpleClientset()

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(fakeKubeClient, time.Second*30)
	openebsInformerFactory := informers.NewSharedInformerFactory(fakeOpenebsClient, time.Second*30)

	// Instantiate the cStor Pool controllers.
	poolController := NewCStorPoolController(fakeKubeClient, fakeOpenebsClient, kubeInformerFactory,
		openebsInformerFactory)

	testPoolResource := map[string]struct {
		expectedError error
		test          *apis.CStorPool
	}{
		"img2PoolResource": {
			expectedError: nil,
			test: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:       "pool2",
					UID:        types.UID("abcd"),
					Finalizers: []string{"openebs"},
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{},
			},
		},
	}
	for desc, ut := range testPoolResource {
		// Create Pool resource
		_, err := poolController.clientset.OpenebsV1alpha1().CStorPools().Create(ut.test)
		if err != nil {
			t.Fatalf("Desc:%v, Unable to create resource : %v", desc, ut.test.ObjectMeta.Name)
		}
		obtainedErr := poolController.removeFinalizer(ut.test)
		if obtainedErr != ut.expectedError {
			t.Fatalf("Desc:%v, Expected:%v, Got:%v", desc, ut.expectedError,
				obtainedErr)
		}
	}
}

func TestIsRightCStorPoolMgmt(t *testing.T) {
	testPoolResource := map[string]struct {
		expectedOutput bool
		test           *apis.CStorPool
	}{
		"img2PoolResource": {
			expectedOutput: true,
			test: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:       "pool2",
					UID:        types.UID("abcd"),
					Finalizers: []string{"openebs"},
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{},
			},
		},
	}
	for desc, ut := range testPoolResource {
		os.Setenv("cstorid", string(ut.test.UID))
		obtainedOutput := IsRightCStorPoolMgmt(ut.test)
		if obtainedOutput != ut.expectedOutput {
			t.Fatalf("Desc:%v, Expected:%v, Got:%v", desc, ut.expectedOutput,
				obtainedOutput)
		}
	}
}

func TestIsRightCStorPoolMgmtNegative(t *testing.T) {
	testPoolResource := map[string]struct {
		expectedOutput bool
		test           *apis.CStorPool
	}{
		"img2PoolResource": {
			expectedOutput: false,
			test: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:       "pool2",
					UID:        types.UID("abcd"),
					Finalizers: []string{"openebs"},
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{},
			},
		},
	}
	for desc, ut := range testPoolResource {
		os.Setenv("cstorid", string("awer"))
		obtainedOutput := IsRightCStorPoolMgmt(ut.test)
		if obtainedOutput != ut.expectedOutput {
			t.Fatalf("Desc:%v, Expected:%v, Got:%v", desc, ut.expectedOutput,
				obtainedOutput)
		}
	}
}

func TestIsDestroyEvent(t *testing.T) {
	deletionTimeStamp := metav1.Now()
	testPoolResource := map[string]struct {
		expectedOutput bool
		test           *apis.CStorPool
	}{
		"img2PoolResource": {
			expectedOutput: true,
			test: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pool2",
					UID:               types.UID("abcd"),
					Finalizers:        []string{"openebs"},
					DeletionTimestamp: &deletionTimeStamp,
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{},
			},
		},
		"img1PoolResource": {
			expectedOutput: false,
			test: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pool1",
					UID:               types.UID("abcd"),
					Finalizers:        []string{"openebs"},
					DeletionTimestamp: nil,
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{},
			},
		},
	}
	for desc, ut := range testPoolResource {
		obtainedOutput := IsDestroyEvent(ut.test)
		if obtainedOutput != ut.expectedOutput {
			t.Fatalf("Desc:%v, Expected:%v, Got:%v", desc, ut.expectedOutput,
				obtainedOutput)
		}
	}
}

func TestIsOnlyStatusChange(t *testing.T) {
	deletionTimeStamp := metav1.Now()
	testPoolResource := map[string]struct {
		expectedOutput bool
		testOld        *apis.CStorPool
		testNew        *apis.CStorPool
	}{
		"img2PoolResource": {
			expectedOutput: true,
			testOld: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pool2",
					UID:               types.UID("abcd"),
					Finalizers:        []string{"openebs"},
					DeletionTimestamp: &deletionTimeStamp,
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{Phase: "init"},
			},
			testNew: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pool2",
					UID:               types.UID("abcd"),
					Finalizers:        []string{"openebs"},
					DeletionTimestamp: &deletionTimeStamp,
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{Phase: "offline"},
			},
		},
		"img1PoolResource": {
			expectedOutput: false,
			testOld: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pool1",
					UID:               types.UID("abc"),
					Finalizers:        []string{"openebs"},
					DeletionTimestamp: &deletionTimeStamp,
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{Phase: "init"},
			},
			testNew: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pool2",
					UID:               types.UID("abcd"),
					Finalizers:        []string{"openebs"},
					DeletionTimestamp: &deletionTimeStamp,
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{Phase: "init"},
			},
		},
	}
	for desc, ut := range testPoolResource {
		obtainedOutput := IsOnlyStatusChange(ut.testOld, ut.testNew)
		if obtainedOutput != ut.expectedOutput {
			t.Fatalf("Desc:%v, Expected:%v, Got:%v", desc, ut.expectedOutput,
				obtainedOutput)
		}
	}
}

func TestIsInitStatus(t *testing.T) {
	deletionTimeStamp := metav1.Now()
	testPoolResource := map[string]struct {
		expectedOutput bool
		test           *apis.CStorPool
	}{
		"img2PoolResource": {
			expectedOutput: true,
			test: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pool2",
					UID:               types.UID("abcd"),
					Finalizers:        []string{"openebs"},
					DeletionTimestamp: &deletionTimeStamp,
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{Phase: "init"},
			},
		},
		"img1PoolResource": {
			expectedOutput: false,
			test: &apis.CStorPool{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:              "pool1",
					UID:               types.UID("abcde"),
					Finalizers:        []string{"openebs"},
					DeletionTimestamp: &deletionTimeStamp,
				},
				Spec: apis.CStorPoolSpec{
					Disks: apis.DiskAttr{
						DiskList: []string{"/tmp/img2.img"},
					},
					PoolSpec: apis.CStorPoolAttr{
						CacheFile:        "/tmp/pool2.cache",
						PoolType:         "striped",
						OverProvisioning: false,
					},
				},
				Status: apis.CStorPoolStatus{Phase: "online"},
			},
		},
	}
	for desc, ut := range testPoolResource {
		obtainedOutput := IsInitStatus(ut.test)
		if obtainedOutput != ut.expectedOutput {
			t.Fatalf("Desc:%v, Expected:%v, Got:%v", desc, ut.expectedOutput,
				obtainedOutput)
		}
	}
}
