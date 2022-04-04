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
	"google.golang.org/grpc"
	// authzpb "github.com/openconfig/gnsi/authz"
	// certpb "github.com/openconfig/gnsi/cert"
	// consolepb "github.com/openconfig/gnsi/console"
	// pathzpb "github.com/openconfig/gnsi/pathz"
	// sshpb "github.com/openconfig/gnsi/ssh"
)

/*
type authz struct {
}

func (a *authz) Rotate(authzpb.AuthzManagement_RotateServer) error {
	return grpc.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

type cert struct {
}

func (c *cert) Install(certpb.CertificateManagement_InstallServer) error {
	return grpc.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

func (c *cert) Rotate(certpb.CertificateManagement_RotateServer) error {
	return grpc.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

type console struct {
}

func (c *console) MutateAccountPassword(consolepb.Console_MutateAccountPasswordServer) error {
	return grpc.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

type pathz struct {
}

func (p *pathz) Install(pathzpb.PathzManagement_InstallServer) error {
	return grpc.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

func (p *pathz) Rotate(pathzpb.PathzManagement_RotateServer) error {
	return grpc.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

type ssh struct {
}

func (s *ssh) MutateAccountCredentials(sshpb.Ssh_MutateAccountCredentialsServer) error {
	return grpc.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

func (s *ssh) MutateHostCredentials(sshpb.Ssh_MutateHostCredentialsServer) error {
	return grpc.Errorf(codes.Unimplemented, "Fake UnImplemented")
}

*/

// Server is a fake gNSI implementation.
type Server struct {
	s *grpc.Server
	//	authz   *authz
	//	cert    *cert
	//	console *console
	//	pathz   *pathz
	//	ssh     *ssh
}

// New returns a new fake gNMI server.
func New(s *grpc.Server) *Server {
	srv := &Server{
		s: s,
		//		authz:   &authz{},
		//		cert:    &cert{},
		//		console: &console{},
		//		pathz:   &pathz{},
		//		ssh:     &ssh{},
	}
	//	authzpb.RegisterAuthzManagementServer(s, srv.authz)
	//	certpb.RegisterCertificateManagementServer(s, srv.cert)
	//	consolepb.RegisterConsoleServer(s, srv.console)
	//	pathzpb.RegisterPathzManagementServer(s, srv.pathz)
	//	sshpb.RegisterSshServer(s, srv.ssh)

	return srv
}
