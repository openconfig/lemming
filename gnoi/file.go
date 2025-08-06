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

package gnoi

import (
	"context"
	"crypto/md5" //nolint:gosec // MD5 required
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	log "github.com/golang/glog"
	"github.com/openconfig/gnoi/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fpb "github.com/openconfig/gnoi/file"
)

const (
	// 100MB max file size
	maxFileSize = 100 * 1024 * 1024
	// 64KB chunks
	maxChunkSize = 64 * 1024
	defaultUmask = 022
)

// file implements the gNOI file service.
type file struct {
	fpb.UnimplementedFileServer
	mu    sync.RWMutex
	files map[string]*fileInfo
}

// fileInfo represents a file in the simulated file system.
type fileInfo struct {
	path        string
	content     []byte
	permissions uint32
	created     time.Time
	modified    time.Time
}

// newFile creates a new file service instance.
func newFile() *file {
	return &file{
		files: make(map[string]*fileInfo),
	}
}

// Put implements the gNOI file service Put RPC.
func (f *file) Put(stream fpb.File_PutServer) error {
	var (
		filePath     string
		permissions  uint32
		content      []byte
		expectedHash []byte
		hashMethod   types.HashType_HashMethod
	)

	// Process the stream
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed to receive request: %v", err)
		}

		if req == nil || req.Request == nil {
			return status.Error(codes.InvalidArgument, "received nil request")
		}

		switch r := req.Request.(type) {
		case *fpb.PutRequest_Open:
			var err error
			filePath, err = f.validatePath(r.Open.RemoteFile)
			if err != nil {
				return err
			}
			permissions = r.Open.Permissions
			log.Infof("File Put: opening file %s with permissions %o", filePath, permissions)

		case *fpb.PutRequest_Contents:
			if filePath == "" {
				return status.Error(codes.InvalidArgument, "must send Open message before Contents")
			}
			newSize := len(content) + len(r.Contents)
			if err := f.validateFileSize(newSize); err != nil {
				return err
			}
			content = append(content, r.Contents...)
			log.Infof("File Put: received %d bytes for %s (total: %d)", len(r.Contents), filePath, len(content))

		case *fpb.PutRequest_Hash:
			if filePath == "" {
				return status.Error(codes.InvalidArgument, "must send Open message before Hash")
			}
			expectedHash = r.Hash.Hash
			hashMethod = r.Hash.Method
			log.Infof("File Put: received hash for %s using method %v", filePath, hashMethod)

		default:
			return status.Error(codes.InvalidArgument, "unknown request type")
		}
	}

	// Validate required messages
	if filePath == "" {
		return status.Error(codes.InvalidArgument, "no Open message received")
	}
	if expectedHash == nil {
		return status.Error(codes.InvalidArgument, "no Hash message received")
	}

	// Verify hash
	actualHash, err := f.computeHash(content, hashMethod)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to compute hash: %v", err)
	}

	if !f.compareHashes(actualHash, expectedHash) {
		log.Errorf("File Put: hash mismatch for %s", filePath)
		return status.Error(codes.InvalidArgument, "hash verification failed")
	}

	// Store the file
	now := time.Now()
	f.mu.Lock()
	f.files[filePath] = &fileInfo{
		path:        filePath,
		content:     content,
		permissions: permissions,
		created:     now,
		modified:    now,
	}
	f.mu.Unlock()

	log.Infof("File Put: successfully stored file %s (%d bytes)", filePath, len(content))
	return stream.SendAndClose(&fpb.PutResponse{})
}

// Get implements the gNOI file service Get RPC.
func (f *file) Get(req *fpb.GetRequest, stream fpb.File_GetServer) error {
	filePath, err := f.validatePath(req.RemoteFile)
	if err != nil {
		return err
	}

	f.mu.RLock()
	file, exists := f.files[filePath]
	f.mu.RUnlock()

	if !exists {
		return status.Errorf(codes.NotFound, "file %s not found", filePath)
	}

	log.Infof("File Get: streaming file %s (%d bytes)", filePath, len(file.content))

	// Send file content in chunks
	for i := 0; i < len(file.content); i += maxChunkSize {
		end := min(i+maxChunkSize, len(file.content))

		chunk := file.content[i:end]
		if err := stream.Send(&fpb.GetResponse{
			Response: &fpb.GetResponse_Contents{
				Contents: chunk,
			},
		}); err != nil {
			return status.Errorf(codes.Internal, "failed to send content: %v", err)
		}
	}

	// Send hash
	hash, err := f.computeHash(file.content, types.HashType_MD5)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to compute hash: %v", err)
	}

	if err := stream.Send(&fpb.GetResponse{
		Response: &fpb.GetResponse_Hash{
			Hash: &types.HashType{
				// Make it configurable in the future
				Method: types.HashType_MD5,
				Hash:   hash,
			},
		},
	}); err != nil {
		return status.Errorf(codes.Internal, "failed to send hash: %v", err)
	}

	return nil
}

// Stat implements the gNOI file service Stat RPC.
func (f *file) Stat(ctx context.Context, req *fpb.StatRequest) (*fpb.StatResponse, error) {
	path, err := f.validatePath(req.Path)
	if err != nil {
		return nil, err
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	var stats []*fpb.StatInfo
	for filePath, file := range f.files {
		// Support both exact directory match and recursive listing
		fileDir := filepath.Dir(filePath)
		if fileDir == path || (strings.HasPrefix(filePath, path+"/") && path != "/") || (path == "/" && strings.HasPrefix(filePath, "/")) {
			stats = append(stats, &fpb.StatInfo{
				Path:         filePath,
				LastModified: uint64(file.modified.UnixNano()),
				Permissions:  file.permissions,
				Size:         uint64(len(file.content)),
				Umask:        defaultUmask,
			})
		}
	}

	log.Infof("File Stat: found %d files in %s", len(stats), path)
	return &fpb.StatResponse{Stats: stats}, nil
}

// Remove implements the gNOI file service Remove RPC.
func (f *file) Remove(ctx context.Context, req *fpb.RemoveRequest) (*fpb.RemoveResponse, error) {
	filePath, err := f.validatePath(req.RemoteFile)
	if err != nil {
		return nil, err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if _, exists := f.files[filePath]; !exists {
		return nil, status.Errorf(codes.NotFound, "file %s not found", filePath)
	}

	delete(f.files, filePath)
	log.Infof("File Remove: deleted file %s", filePath)

	return &fpb.RemoveResponse{}, nil
}

// TransferToRemote implements the gNOI file service TransferToRemote RPC.
func (f *file) TransferToRemote(ctx context.Context, req *fpb.TransferToRemoteRequest) (*fpb.TransferToRemoteResponse, error) {
	if req.GetLocalPath() == "" {
		return nil, status.Error(codes.InvalidArgument, "local_path cannot be empty")
	}

	if req.GetRemoteDownload() == nil {
		return nil, status.Error(codes.InvalidArgument, "remote_download cannot be nil")
	}

	localPath, err := f.validatePath(req.GetLocalPath())
	if err != nil {
		return nil, err
	}

	remoteDownload := req.GetRemoteDownload()
	if remoteDownload.GetPath() == "" {
		return nil, status.Error(codes.InvalidArgument, "remote_download.path cannot be empty")
	}

	// Check if the local file exists in our simulated file system
	f.mu.RLock()
	file, exists := f.files[localPath]
	f.mu.RUnlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "local file %s not found", localPath)
	}

	// Validate protocol is supported
	switch remoteDownload.GetProtocol() {
	case 0:
		return nil, status.Error(codes.InvalidArgument, "remote_download.protocol cannot be UNKNOWN")
	case 1, 2, 3, 4:
		// All protocols are supported in simulation
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported protocol: %v", remoteDownload.GetProtocol())
	}

	log.Infof("File TransferToRemote: transferring %s to %s via protocol %v",
		localPath, remoteDownload.GetPath(), remoteDownload.GetProtocol())

	// Simulate the transfer by computing hash of the file content
	hash, err := f.computeHash(file.content, types.HashType_MD5)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to compute hash: %v", err)
	}

	log.Infof("File TransferToRemote: successfully transferred %s (%d bytes)", localPath, len(file.content))

	return &fpb.TransferToRemoteResponse{
		Hash: &types.HashType{
			// Make it configurable in the future
			Method: types.HashType_MD5,
			Hash:   hash,
		},
	}, nil
}

// validatePath validates and normalizes a file path.
func (f *file) validatePath(path string) (string, error) {
	if path == "" {
		return "", status.Error(codes.InvalidArgument, "path cannot be empty")
	}

	cleanPath := filepath.Clean(path)
	if !filepath.IsAbs(cleanPath) {
		return "", status.Error(codes.InvalidArgument, "path must be absolute")
	}

	// Prevent access to sensitive system paths
	if strings.HasPrefix(cleanPath, "/proc") || strings.HasPrefix(cleanPath, "/sys") {
		return "", status.Error(codes.PermissionDenied, "access to system paths not allowed")
	}

	return cleanPath, nil
}

// validateFileSize checks if file size is within limits.
func (f *file) validateFileSize(size int) error {
	if size > maxFileSize {
		return status.Errorf(codes.InvalidArgument, "file size %d exceeds maximum allowed size %d", size, maxFileSize)
	}
	return nil
}

// computeHash computes the hash of the given data using the specified method.
func (f *file) computeHash(data []byte, method types.HashType_HashMethod) ([]byte, error) {
	switch method {
	case types.HashType_MD5:
		h := md5.New() //nolint:gosec // MD5 required
		h.Write(data)
		return h.Sum(nil), nil
	case types.HashType_SHA256:
		h := sha256.New()
		h.Write(data)
		return h.Sum(nil), nil
	case types.HashType_SHA512:
		h := sha512.New()
		h.Write(data)
		return h.Sum(nil), nil
	default:
		return nil, fmt.Errorf("unsupported hash method: %v", method)
	}
}

// compareHashes compares two hashes for equality.
func (f *file) compareHashes(hash1, hash2 []byte) bool {
	if len(hash1) != len(hash2) {
		return false
	}
	for i := range hash1 {
		if hash1[i] != hash2[i] {
			return false
		}
	}
	return true
}

// GetFileInfo returns information about a specific file for testing.
func (f *file) GetFileInfo(filePath string) (*fileInfo, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	file, exists := f.files[filePath]
	return file, exists
}

// ListFiles returns all files in the simulated file system for testing.
func (f *file) ListFiles() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	files := make([]string, 0, len(f.files))
	for filePath := range f.files {
		files = append(files, filePath)
	}
	return files
}

// Reset clears all files from the simulated file system for testing.
func (f *file) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.files = make(map[string]*fileInfo)
}
