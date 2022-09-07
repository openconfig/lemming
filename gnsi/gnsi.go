// Copyright 2022 Google LLC
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

package gnsi

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authzpb "github.com/openconfig/gnsi/authz"
	certpb "github.com/openconfig/gnsi/cert"
	credentialzpb "github.com/openconfig/gnsi/credentialz"
	pathzpb "github.com/openconfig/gnsi/pathz"
)

type authz struct {
	authzpb.UnimplementedAuthzServer
}

func (a *authz) Rotate(authzpb.Authz_RotateServer) error {
	return status.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

func (a *authz) Get(context.Context, *authzpb.GetRequest) (*authzpb.GetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

type cert struct {
	certpb.UnimplementedCertificateManagementServer
}

func (c *cert) Install(certpb.CertificateManagement_InstallServer) error {
	return status.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

func (c *cert) Rotate(certpb.CertificateManagement_RotateServer) error {
	return status.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

type pathz struct {
	pathzpb.UnimplementedPathzServer
}

func (p *pathz) Rotate(pathzpb.Pathz_RotateServer) error {
	return status.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

func (p *pathz) Probe(context.Context, *pathzpb.ProbeRequest) (*pathzpb.ProbeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

type credentialz struct {
	credentialzpb.UnimplementedCredentialzServer
}

func (c *credentialz) MutateAccountCredentials(credentialzpb.Credentialz_RotateAccountCredentialsServer) error {
	return status.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

func (c *credentialz) MutateHostCredentials(credentialzpb.Credentialz_RotateHostCredentialsServer) error {
	return status.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

// Server is a fake gNSI implementation.
type Server struct {
	s     *grpc.Server
	authz *authz
	cert  *cert
	pathz *pathz
	credz *credentialz
}

// New returns a new fake gNMI server.
func New(s *grpc.Server) *Server {
	srv := &Server{
		s:     s,
		authz: &authz{},
		cert:  &cert{},
		pathz: &pathz{},
		credz: &credentialz{},
	}
	authzpb.RegisterAuthzServer(s, srv.authz)
	certpb.RegisterCertificateManagementServer(s, srv.cert)
	credentialzpb.RegisterCredentialzServer(s, srv.credz)
	pathzpb.RegisterPathzServer(s, srv.pathz)

	return srv
}
