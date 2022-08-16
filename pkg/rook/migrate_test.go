package rook

import (
	"context"
	"encoding/json"
	"github.com/replicatedhq/kurl/pkg/rook/testfiles"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

// test function only, contains panics
func runtimeFromPodlistJson(podListJson []byte) []runtime.Object {
	podList := corev1.PodList{}
	err := json.Unmarshal(podListJson, &podList)
	if err != nil {
		panic(err) // this is only called for unit tests, not at runtime
	}

	runtimeObjects := []runtime.Object{}
	for idx, _ := range podList.Items {
		runtimeObjects = append(runtimeObjects, &podList.Items[idx])
	}
	return runtimeObjects
}

func Test_countRookOSDs(t *testing.T) {
	tests := []struct {
		name      string
		resources []runtime.Object
		hostpath  int
		block     int
	}{
		{
			name:      "6 blockdevice osds",
			resources: runtimeFromPodlistJson(testfiles.SixBlockDevicePods),
			hostpath:  0,
			block:     6,
		},
		{
			name:      "one hostpath osd",
			resources: runtimeFromPodlistJson(testfiles.HostpathPods),
			hostpath:  1,
			block:     0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := require.New(t)
			clientset := fake.NewSimpleClientset(tt.resources...)

			gotHostpath, gotBlock, err := countRookOSDs(context.Background(), clientset)
			req.NoError(err)
			req.Equal(tt.hostpath, gotHostpath)
			req.Equal(tt.block, gotBlock)
		})
	}
}

func Test_getRookOSDs(t *testing.T) {
	tests := []struct {
		name      string
		resources []runtime.Object
		rookOsds  []RookOSD
	}{
		{
			name:      "6 blockdevice osds",
			resources: runtimeFromPodlistJson(testfiles.SixBlockDevicePods),
			rookOsds: []RookOSD{
				{
					Num:        3,
					IsHostpath: false,
				},
				{
					Num:        4,
					IsHostpath: false,
				},
				{
					Num:        5,
					IsHostpath: false,
				},
				{
					Num:        6,
					IsHostpath: false,
				},
				{
					Num:        7,
					IsHostpath: false,
				},
				{
					Num:        8,
					IsHostpath: false,
				},
			},
		},
		{
			name:      "one hostpath osd",
			resources: runtimeFromPodlistJson(testfiles.HostpathPods),
			rookOsds: []RookOSD{
				{
					Num:        0,
					IsHostpath: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := require.New(t)
			clientset := fake.NewSimpleClientset(tt.resources...)

			rookOsds, err := getRookOSDs(context.Background(), clientset)
			req.NoError(err)
			req.Equal(tt.rookOsds, rookOsds)
		})
	}
}