// Copyright (c) Ultraviolet
// SPDX-License-Identifier: Apache-2.0

package current_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net"
	"testing"
	"time"

	ea "github.com/ultravioletrs/cocos/pkg/atls/ea"
	attestation "github.com/ultravioletrs/cocos/pkg/atls/eaattestation"
)

const (
	requestContextSize = 16
	replayContextSize  = 12
	dummyEvidence      = "dummy-attestation-report"
)

var signatureAlgorithmsECDSA = []byte{0x00, 0x02, 0x04, 0x03}

type acceptAllEvidenceVerifier struct{}

func (acceptAllEvidenceVerifier) VerifyEvidence([]byte) error { return nil }

type exactEvidenceVerifier struct {
	want []byte
}

func (v exactEvidenceVerifier) VerifyEvidence(got []byte) error {
	if !bytes.Equal(got, v.want) {
		return fmt.Errorf("unexpected evidence bytes")
	}
	return nil
}

func selfSignedCert(t *testing.T) tls.Certificate {
	t.Helper()

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "ea-regression"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"localhost"},
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatal(err)
	}

	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: priv}
}

func tlsPair(t *testing.T, cert tls.Certificate) (srv, cli *tls.Conn) {
	t.Helper()

	srvConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,
	}
	cliConf := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS13,
		MaxVersion:         tls.VersionTLS13,
	}

	a, b := net.Pipe()
	srv = tls.Server(a, srvConf)
	cli = tls.Client(b, cliConf)

	errCh := make(chan error, 2)
	go func() { errCh <- srv.Handshake() }()
	go func() { errCh <- cli.Handshake() }()
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			t.Fatalf("handshake: %v", err)
		}
	}

	return srv, cli
}

func attestationRequest(t *testing.T, contextSize int) *ea.AuthenticatorRequest {
	t.Helper()

	return authenticatorRequest(t, contextSize, ea.CMWAttestationOfferExtension())
}

func plainRequest(t *testing.T, contextSize int) *ea.AuthenticatorRequest {
	return authenticatorRequest(t, contextSize)
}

func authenticatorRequest(t *testing.T, contextSize int, extraExtensions ...ea.Extension) *ea.AuthenticatorRequest {
	t.Helper()

	ctx, err := ea.NewRandomContext(contextSize)
	if err != nil {
		t.Fatal(err)
	}

	extensions := []ea.Extension{
		{Type: ea.SignatureAlgorithmsExtensionType, Data: signatureAlgorithmsECDSA},
	}
	extensions = append(extensions, extraExtensions...)

	return &ea.AuthenticatorRequest{
		Type:       ea.HandshakeTypeClientCertificateRequest,
		Context:    ctx,
		Extensions: extensions,
	}
}

func attestationDataExtension(t *testing.T, st *tls.ConnectionState, label string, ctx []byte, leaf *x509.Certificate) ea.Extension {
	return attestationDataExtensionWithEvidence(t, st, label, ctx, leaf, []byte(dummyEvidence))
}

func attestationDataExtensionWithEvidence(t *testing.T, st *tls.ConnectionState, label string, ctx []byte, leaf *x509.Certificate, evidence []byte) ea.Extension {
	t.Helper()

	_, aikPubHash, binding, err := attestation.ComputeBinding(st, label, ctx, leaf)
	if err != nil {
		t.Fatal(err)
	}

	payloadBytes, err := attestation.MarshalPayload(attestation.Payload{
		Version:   1,
		Evidence:  evidence,
		MediaType: "application/eat+cwt",
		Binder: attestation.AttestationBinder{
			ExporterLabel: label,
			AIKPubHash:    aikPubHash,
			Binding:       binding,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	ext, err := ea.CMWAttestationDataExtension(payloadBytes)
	if err != nil {
		t.Fatal(err)
	}

	return ext
}

func rootsForCert(t *testing.T, cert tls.Certificate) (*x509.CertPool, *x509.Certificate) {
	t.Helper()

	leaf, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		t.Fatal(err)
	}

	roots := x509.NewCertPool()
	roots.AddCert(leaf)

	return roots, leaf
}

func TestVerifierAcceptsSelfDeclaredNonDefaultExporterLabel(t *testing.T) {
	cert := selfSignedCert(t)
	srv, cli := tlsPair(t, cert)
	defer srv.Close()
	defer cli.Close()

	req := attestationRequest(t, requestContextSize)
	roots, leaf := rootsForCert(t, cert)
	srvState := srv.ConnectionState()

	// Current behavior: the verifier uses the label carried in the payload.
	// This records acceptance of a non-default Cocos profile label.
	ext := attestationDataExtension(t, &srvState, attestation.ExporterLabelAttestationBinding, req.Context, leaf)

	auth, err := ea.CreateAuthenticator(&srvState, ea.RoleServer, req, cert, []ea.Extension{ext})
	if err != nil {
		t.Fatal(err)
	}

	cliState := cli.ConnectionState()
	res, err := ea.ValidateAuthenticatorWithAttestation(&cliState, ea.RoleServer, req, auth, &x509.VerifyOptions{Roots: roots}, attestation.VerificationPolicy{
		EvidenceVerifier: acceptAllEvidenceVerifier{},
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.Attestation == nil {
		t.Fatalf("expected verified attestation result")
	}
	if res.Attestation.UsedExporterLabel != attestation.ExporterLabelAttestationBinding {
		t.Fatalf("got exporter label %q, want %q", res.Attestation.UsedExporterLabel, attestation.ExporterLabelAttestationBinding)
	}
}

func TestVerifierReportsSuccessWithoutAttestationAfterExplicitOffer(t *testing.T) {
	cert := selfSignedCert(t)
	srv, cli := tlsPair(t, cert)
	defer srv.Close()
	defer cli.Close()

	req := attestationRequest(t, requestContextSize)
	srvState := srv.ConnectionState()

	// Current API-level behavior: a syntactically valid authenticator without
	// the offered attestation extension returns no error and no attestation.
	auth, err := ea.CreateAuthenticator(&srvState, ea.RoleServer, req, cert, nil)
	if err != nil {
		t.Fatal(err)
	}

	roots, _ := rootsForCert(t, cert)
	cliState := cli.ConnectionState()
	res, err := ea.ValidateAuthenticatorWithAttestation(&cliState, ea.RoleServer, req, auth, &x509.VerifyOptions{Roots: roots}, attestation.VerificationPolicy{
		EvidenceVerifier: acceptAllEvidenceVerifier{},
	})
	if err != nil {
		t.Fatalf("expected authenticator to be accepted without attestation extension, got %v", err)
	}

	if res.Attestation != nil {
		t.Fatalf("expected no attestation result when attestation extension is omitted")
	}
	if len(res.CMWAttestation) != 0 {
		t.Fatalf("expected no CMW attestation payload")
	}
}

func TestVerifierSeparatesEvidencePolicyFromBinderVerification(t *testing.T) {
	cert := selfSignedCert(t)
	srv, cli := tlsPair(t, cert)
	defer srv.Close()
	defer cli.Close()

	req := attestationRequest(t, requestContextSize)
	roots, leaf := rootsForCert(t, cert)
	srvState := srv.ConnectionState()

	// The binder is correct for this TLS connection, request context, and leaf.
	// The evidence bytes are accepted by policy but are not passed the binder.
	evidence := []byte("policy-accepted-evidence-for-another-binding")
	ext := attestationDataExtensionWithEvidence(t, &srvState, attestation.ExporterLabelAttestation, req.Context, leaf, evidence)

	auth, err := ea.CreateAuthenticator(&srvState, ea.RoleServer, req, cert, []ea.Extension{ext})
	if err != nil {
		t.Fatal(err)
	}

	cliState := cli.ConnectionState()
	res, err := ea.ValidateAuthenticatorWithAttestation(&cliState, ea.RoleServer, req, auth, &x509.VerifyOptions{Roots: roots}, attestation.VerificationPolicy{
		EvidenceVerifier: exactEvidenceVerifier{want: evidence},
	})
	if err != nil {
		t.Fatal(err)
	}

	if res.Attestation == nil {
		t.Fatalf("expected verified attestation result")
	}
	if !res.Attestation.EvidenceVerified || !res.Attestation.BindingVerified {
		t.Fatalf("expected separate evidence and binder verification to succeed")
	}
}

func TestVerifierRejectsLeafKeySubstitution(t *testing.T) {
	bindingCert := selfSignedCert(t)
	authenticatorCert := selfSignedCert(t)

	srv, cli := tlsPair(t, authenticatorCert)
	defer srv.Close()
	defer cli.Close()

	req := attestationRequest(t, requestContextSize)
	_, bindingLeaf := rootsForCert(t, bindingCert)
	roots, _ := rootsForCert(t, authenticatorCert)

	srvState := srv.ConnectionState()
	ext := attestationDataExtension(t, &srvState, attestation.ExporterLabelAttestation, req.Context, bindingLeaf)

	auth, err := ea.CreateAuthenticator(&srvState, ea.RoleServer, req, authenticatorCert, []ea.Extension{ext})
	if err != nil {
		t.Fatal(err)
	}

	cliState := cli.ConnectionState()
	_, err = ea.ValidateAuthenticatorWithAttestation(&cliState, ea.RoleServer, req, auth, &x509.VerifyOptions{Roots: roots}, attestation.VerificationPolicy{
		EvidenceVerifier: acceptAllEvidenceVerifier{},
	})
	if err != attestation.ErrAIKPubHashMismatch && err != attestation.ErrBindingMismatch {
		t.Fatalf("got %v, want %v or %v", err, attestation.ErrAIKPubHashMismatch, attestation.ErrBindingMismatch)
	}
}

func TestSessionRejectsContextReuse(t *testing.T) {
	cert := selfSignedCert(t)
	srv, cli := tlsPair(t, cert)
	defer srv.Close()
	defer cli.Close()

	req := plainRequest(t, replayContextSize)
	createSession := ea.NewSession()
	srvState := srv.ConnectionState()

	// Positive guarantee for the explicit ea.Session path: the same request
	// context cannot be used twice for creation or validation.
	auth, err := createSession.CreateAuthenticator(&srvState, ea.RoleServer, req, cert, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := createSession.CreateAuthenticator(&srvState, ea.RoleServer, req, cert, nil); err != ea.ErrContextReuse {
		t.Fatalf("got %v, want %v", err, ea.ErrContextReuse)
	}

	validateSession := ea.NewSession()
	roots, _ := rootsForCert(t, cert)
	cliState := cli.ConnectionState()

	if _, err := validateSession.ValidateAuthenticator(&cliState, ea.RoleServer, req, auth, &x509.VerifyOptions{Roots: roots}); err != nil {
		t.Fatal(err)
	}
	if _, err := validateSession.ValidateAuthenticator(&cliState, ea.RoleServer, req, auth, &x509.VerifyOptions{Roots: roots}); err != ea.ErrContextReuse {
		t.Fatalf("got %v, want %v", err, ea.ErrContextReuse)
	}
}

func TestAPIWithoutSessionAllowsContextReplay(t *testing.T) {
	cert := selfSignedCert(t)
	srv, cli := tlsPair(t, cert)
	defer srv.Close()
	defer cli.Close()

	req := plainRequest(t, replayContextSize)
	srvState := srv.ConnectionState()
	auth, err := ea.CreateAuthenticator(&srvState, ea.RoleServer, req, cert, nil)
	if err != nil {
		t.Fatal(err)
	}

	roots, _ := rootsForCert(t, cert)
	cliState := cli.ConnectionState()

	// Sessionless validation remains replayable at the API level.
	if _, err := ea.ValidateAuthenticator(&cliState, ea.RoleServer, req, auth, &x509.VerifyOptions{Roots: roots}); err != nil {
		t.Fatal(err)
	}
	if _, err := ea.ValidateAuthenticator(&cliState, ea.RoleServer, req, auth, &x509.VerifyOptions{Roots: roots}); err != nil {
		t.Fatalf("expected replay to remain acceptable without session tracking, got %v", err)
	}
}
