// Copyright 2025 Google LLC
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

package gnoi_file_test

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"io"
	"testing"
	"time"

	"github.com/openconfig/gnoi/common"
	fpb "github.com/openconfig/gnoi/file"
	"github.com/openconfig/gnoi/types"
	"github.com/openconfig/lemming/internal/binding"
	"github.com/openconfig/ondatra"
)

const (
	// Test file paths
	testConfigFilePath = "/config/running.conf"
	testBackupFilePath = "/config/backup.conf"

	// Test directory
	testConfigDir = "/config"

	// File permissions
	testPermissions uint32 = 0o644

	// Simple configuration for testing
	testConfig = `{
		"openconfig-interfaces:interfaces": {
			"interface": [
				{
					"name": "Ethernet1/1",
					"config": {
						"name": "Ethernet1/1",
						"type": "iana-if-type:ethernetCsmacd",
						"description": "WAN Interface",
						"enabled": true
					},
					"subinterfaces": {
						"subinterface": [
							{
								"index": 0,
								"openconfig-if-ip:ipv4": {
									"addresses": {
										"address": [
											{
												"ip": "192.168.1.1",
												"config": {
													"ip": "192.168.1.1",
													"prefix-length": 24
												}
											}
										]
									}
								}
							}
						]
					}
				}
			]
		},
		"openconfig-system:system": {
			"config": {
				"hostname": "ROUTER-01",
				"domain-name": "example.com"
			}
		}
	}`
)

func TestMain(m *testing.M) {
	ondatra.RunTests(m, binding.KNE(".."))
}

// TestGNIOFileService tests all gNOI file service RPCs
func TestGNIOFileService(t *testing.T) {
	dut := ondatra.DUT(t, "dut")

	// Test cases for different file operations
	t.Run("Put_Get_Stat_Remove_Complete_Workflow", func(t *testing.T) {
		testCompleteFileWorkflow(t, dut)
	})

	t.Run("Put_File_With_Different_Hash_Methods", func(t *testing.T) {
		testPutWithDifferentHashMethods(t, dut)
	})

	t.Run("TransferToRemote_Simulation", func(t *testing.T) {
		testTransferToRemote(t, dut)
	})
}

// testCompleteFileWorkflow tests the complete workflow: Put -> Get -> Stat -> Remove
func testCompleteFileWorkflow(t *testing.T, dut *ondatra.DUTDevice) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	gnoiClient, err := dut.RawAPIs().BindingDUT().DialGNOI(ctx)
	if err != nil {
		t.Fatalf("Error dialing gNOI: %v", err)
	}

	// Step 1: Put the file
	t.Logf("Step 1: Putting file %s", testConfigFilePath)
	putStream, err := gnoiClient.File().Put(ctx)
	if err != nil {
		t.Fatalf("Failed to create Put stream: %v", err)
	}
	defer putStream.CloseSend()

	// Send Open message
	openReq := &fpb.PutRequest{
		Request: &fpb.PutRequest_Open{
			Open: &fpb.PutRequest_Details{
				RemoteFile:  testConfigFilePath,
				Permissions: testPermissions,
			},
		},
	}
	if err := putStream.Send(openReq); err != nil {
		t.Fatalf("Failed to send Open request: %v", err)
	}

	// Send content
	contentReq := &fpb.PutRequest{
		Request: &fpb.PutRequest_Contents{
			Contents: []byte(testConfig),
		},
	}
	if err := putStream.Send(contentReq); err != nil {
		t.Fatalf("Failed to send content: %v", err)
	}

	// Send hash
	h := md5.New()
	h.Write([]byte(testConfig))
	hashReq := &fpb.PutRequest{
		Request: &fpb.PutRequest_Hash{
			Hash: &types.HashType{
				Method: types.HashType_MD5,
				Hash:   h.Sum(nil),
			},
		},
	}
	if err := putStream.Send(hashReq); err != nil {
		t.Fatalf("Failed to send hash: %v", err)
	}

	// Close and receive response
	putResp, err := putStream.CloseAndRecv()
	if err != nil {
		t.Fatalf("Failed to close Put stream: %v", err)
	}
	t.Logf("Put response: %v", putResp)

	// Step 2: Get the file
	t.Logf("Step 2: Getting file %s", testConfigFilePath)
	getReq := &fpb.GetRequest{
		RemoteFile: testConfigFilePath,
	}
	getStream, err := gnoiClient.File().Get(ctx, getReq)
	if err != nil {
		t.Fatalf("Failed to create Get stream: %v", err)
	}

	var receivedContent []byte
	var receivedHash *types.HashType

	for {
		getResp, err := getStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("Failed to receive Get response: %v", err)
		}

		switch resp := getResp.Response.(type) {
		case *fpb.GetResponse_Contents:
			receivedContent = append(receivedContent, resp.Contents...)
		case *fpb.GetResponse_Hash:
			receivedHash = resp.Hash
		}
	}

	// Verify content
	if string(receivedContent) != testConfig {
		t.Errorf("Content mismatch: got %s, want %s", string(receivedContent), testConfig)
	}

	// Verify hash
	if receivedHash == nil {
		t.Error("No hash received in Get response")
	} else {
		expectedHash := md5.Sum([]byte(testConfig))
		if !bytes.Equal(receivedHash.Hash, expectedHash[:]) {
			t.Errorf("Hash mismatch: got %v, want %v", receivedHash.Hash, expectedHash[:])
		}
	}

	// Step 3: Stat the file
	t.Logf("Step 3: Getting file stats for %s", testConfigDir)
	statReq := &fpb.StatRequest{
		Path: testConfigDir,
	}
	statResp, err := gnoiClient.File().Stat(ctx, statReq)
	if err != nil {
		t.Fatalf("Failed to get file stats: %v", err)
	}

	// Find our file in the stats
	var foundFile *fpb.StatInfo
	for _, stat := range statResp.Stats {
		if stat.Path == testConfigFilePath {
			foundFile = stat
			break
		}
	}

	if foundFile == nil {
		t.Errorf("File %s not found in directory listing", testConfigFilePath)
	} else {
		t.Logf("File stats: %+v", foundFile)
		// Verify file size
		if foundFile.Size != uint64(len(testConfig)) {
			t.Errorf("File size mismatch: got %d, want %d", foundFile.Size, len(testConfig))
		}
		// Verify permissions
		if foundFile.Permissions != testPermissions {
			t.Errorf("File permissions mismatch: got %o, want %o", foundFile.Permissions, testPermissions)
		}
	}

	// Step 4: Remove the file
	t.Logf("Step 4: Removing file %s", testConfigFilePath)
	removeReq := &fpb.RemoveRequest{
		RemoteFile: testConfigFilePath,
	}
	removeResp, err := gnoiClient.File().Remove(ctx, removeReq)
	if err != nil {
		t.Fatalf("Failed to remove file: %v", err)
	}
	t.Logf("Remove response: %v", removeResp)

	// Step 5: Verify file is removed
	t.Logf("Step 5: Verifying file is removed")
	statResp2, err := gnoiClient.File().Stat(ctx, statReq)
	if err != nil {
		t.Fatalf("Failed to get file stats after removal: %v", err)
	}

	// Check that our file is no longer in the listing
	for _, stat := range statResp2.Stats {
		if stat.Path == testConfigFilePath {
			t.Errorf("File %s still exists after removal", testConfigFilePath)
		}
	}
}

// testPutWithDifferentHashMethods tests Put with different hash methods
func testPutWithDifferentHashMethods(t *testing.T, dut *ondatra.DUTDevice) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	gnoiClient, err := dut.RawAPIs().BindingDUT().DialGNOI(ctx)
	if err != nil {
		t.Fatalf("Error dialing gNOI: %v", err)
	}

	testCases := []struct {
		name       string
		hashMethod types.HashType_HashMethod
	}{
		{
			name:       "MD5_Hash",
			hashMethod: types.HashType_MD5,
		},
		{
			name:       "SHA256_Hash",
			hashMethod: types.HashType_SHA256,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			putStream, err := gnoiClient.File().Put(ctx)
			if err != nil {
				t.Fatalf("Failed to create Put stream: %v", err)
			}
			defer putStream.CloseSend()

			// Send Open message
			openReq := &fpb.PutRequest{
				Request: &fpb.PutRequest_Open{
					Open: &fpb.PutRequest_Details{
						RemoteFile:  testBackupFilePath,
						Permissions: testPermissions,
					},
				},
			}
			if err := putStream.Send(openReq); err != nil {
				t.Fatalf("Failed to send Open request: %v", err)
			}

			// Send content
			contentReq := &fpb.PutRequest{
				Request: &fpb.PutRequest_Contents{
					Contents: []byte(testConfig),
				},
			}
			if err := putStream.Send(contentReq); err != nil {
				t.Fatalf("Failed to send content: %v", err)
			}

			// Compute hash based on method
			var hash []byte
			switch tc.hashMethod {
			case types.HashType_MD5:
				h := md5.New()
				h.Write([]byte(testConfig))
				hash = h.Sum(nil)
			case types.HashType_SHA256:
				h := sha256.New()
				h.Write([]byte(testConfig))
				hash = h.Sum(nil)
			}

			// Send hash
			hashReq := &fpb.PutRequest{
				Request: &fpb.PutRequest_Hash{
					Hash: &types.HashType{
						Method: tc.hashMethod,
						Hash:   hash,
					},
				},
			}
			if err := putStream.Send(hashReq); err != nil {
				t.Fatalf("Failed to send hash: %v", err)
			}

			// Close and receive response
			putResp, err := putStream.CloseAndRecv()
			if err != nil {
				t.Fatalf("Failed to close Put stream: %v", err)
			}
			t.Logf("Put response for %s: %v", tc.name, putResp)

			// Clean up
			removeReq := &fpb.RemoveRequest{
				RemoteFile: testBackupFilePath,
			}
			_, err = gnoiClient.File().Remove(ctx, removeReq)
			if err != nil {
				t.Logf("Failed to remove test file: %v", err)
			}
		})
	}
}

// testTransferToRemote tests the TransferToRemote RPC
func testTransferToRemote(t *testing.T, dut *ondatra.DUTDevice) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	gnoiClient, err := dut.RawAPIs().BindingDUT().DialGNOI(ctx)
	if err != nil {
		t.Fatalf("Error dialing gNOI: %v", err)
	}

	// First, create a file to transfer
	putStream, err := gnoiClient.File().Put(ctx)
	if err != nil {
		t.Fatalf("Failed to create Put stream: %v", err)
	}

	// Send Open message
	openReq := &fpb.PutRequest{
		Request: &fpb.PutRequest_Open{
			Open: &fpb.PutRequest_Details{
				RemoteFile:  testConfigFilePath,
				Permissions: testPermissions,
			},
		},
	}
	if err := putStream.Send(openReq); err != nil {
		t.Fatalf("Failed to send Open request: %v", err)
	}

	// Send content
	contentReq := &fpb.PutRequest{
		Request: &fpb.PutRequest_Contents{
			Contents: []byte(testConfig),
		},
	}
	if err := putStream.Send(contentReq); err != nil {
		t.Fatalf("Failed to send content: %v", err)
	}

	// Send hash
	h := md5.New()
	h.Write([]byte(testConfig))
	hashReq := &fpb.PutRequest{
		Request: &fpb.PutRequest_Hash{
			Hash: &types.HashType{
				Method: types.HashType_MD5,
				Hash:   h.Sum(nil),
			},
		},
	}
	if err := putStream.Send(hashReq); err != nil {
		t.Fatalf("Failed to send hash: %v", err)
	}

	_, err = putStream.CloseAndRecv()
	if err != nil {
		t.Fatalf("Failed to close Put stream: %v", err)
	}

	// Test TransferToRemote
	transferReq := &fpb.TransferToRemoteRequest{
		LocalPath: testConfigFilePath,
		RemoteDownload: &common.RemoteDownload{
			Path:     "sftp://backup-server.example.com/configs/router-01.json",
			Protocol: common.RemoteDownload_SFTP,
		},
	}

	transferResp, err := gnoiClient.File().TransferToRemote(ctx, transferReq)
	if err != nil {
		t.Fatalf("Failed to transfer file: %v", err)
	}

	t.Logf("Transfer response: %+v", transferResp)

	// Verify response has hash
	if transferResp.Hash == nil {
		t.Error("Transfer response missing hash")
	} else {
		t.Logf("Transfer hash: method=%v, hash=%v", transferResp.Hash.Method, transferResp.Hash.Hash)
	}

	// Clean up
	removeReq := &fpb.RemoveRequest{
		RemoteFile: testConfigFilePath,
	}
	_, err = gnoiClient.File().Remove(ctx, removeReq)
	if err != nil {
		t.Logf("Failed to remove test file: %v", err)
	}
}
