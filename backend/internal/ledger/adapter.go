package ledger

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type LedgerAdapter struct {
	gateway  *client.Gateway
	network  *client.Network
	contract *client.Contract
	conn     *grpc.ClientConn
}

const (
	mspID         = "Org1MSP"
	peerEndpoint  = "localhost:7051"
	gatewayPeer   = "peer0.org1.example.com"
	channelName   = "colshchannel"
	chaincodeName = "trazabilidad"

	certPath    = "internal/ledger/certs/user-cert.pem"
	keyPath     = "internal/ledger/certs/user-key.pem"
	tlsCertPath = "internal/ledger/certs/peer-tls-ca.crt"
)

// NuevoLedgerAdapter establece la conexion con la red Fabric y devuelve un adapter listo para usar.
func NuevoLedgerAdapter() (*LedgerAdapter, error) {
	conn, err := newGrpcConnection()
	if err != nil {
		return nil, fmt.Errorf("error al crear conexion grpc: %w", err)
	}

	id, err := newIdentity()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error al crear identidad: %w", err)
	}

	sign, err := newSign()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error al crear firmante: %w", err)
	}

	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(conn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("error al conectar al gateway: %w", err)
	}

	network := gw.GetNetwork(channelName)
	contract := network.GetContract(chaincodeName)

	return &LedgerAdapter{
		gateway:  gw,
		network:  network,
		contract: contract,
		conn:     conn,
	}, nil
}

// Cerrar libera la conexion. Llamar al apagar el backend.
func (l *LedgerAdapter) Cerrar() {
	if l.gateway != nil {
		l.gateway.Close()
	}
	if l.conn != nil {
		l.conn.Close()
	}
}

// RegistrarEnLedger invoca RegistrarEvento en el chaincode (transaccion, va al orderer).
func (l *LedgerAdapter) RegistrarEnLedger(ctx context.Context, idEvento, idPedido, estado, fecha, responsable string) error {
	_, err := l.contract.SubmitTransaction("RegistrarEvento", idEvento, idPedido, estado, fecha, responsable)
	if err != nil {
		return fmt.Errorf("error al registrar evento en el ledger: %w", err)
	}
	return nil
}

// ConsultarEvento invoca ConsultarEvento en el chaincode (consulta, no va al orderer).
func (l *LedgerAdapter) ConsultarEvento(ctx context.Context, idEvento string) ([]byte, error) {
	result, err := l.contract.EvaluateTransaction("ConsultarEvento", idEvento)
	if err != nil {
		return nil, fmt.Errorf("error al consultar evento: %w", err)
	}
	return result, nil
}

// ConsultarHistorialPedido invoca ConsultarHistorialPedido en el chaincode.
func (l *LedgerAdapter) ConsultarHistorialPedido(ctx context.Context, idPedido string) ([]byte, error) {
	result, err := l.contract.EvaluateTransaction("ConsultarHistorialPedido", idPedido)
	if err != nil {
		return nil, fmt.Errorf("error al consultar historial: %w", err)
	}
	return result, nil
}

func newGrpcConnection() (*grpc.ClientConn, error) {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer certificado TLS: %w", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	return grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
}

func newIdentity() (*identity.X509Identity, error) {
	certificatePEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer certificado de usuario: %w", err)
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		return nil, err
	}

	return identity.NewX509Identity(mspID, certificate)
}

func newSign() (identity.Sign, error) {
	privateKeyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer llave privada: %w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}

	return identity.NewPrivateKeySign(privateKey)
}
