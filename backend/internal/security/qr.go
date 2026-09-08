package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"strings"
)

// obtenerClaveSecreta lee la clave desde la variable de entorno QR_SECRET_KEY.
// Si no esta definida, usa una clave de desarrollo fija (solo para el prototipo;
// en un entorno real esto NUNCA deberia tener un valor por defecto).
func obtenerClaveSecreta() []byte {
	clave := os.Getenv("QR_SECRET_KEY")
	if clave == "" {
		clave = "colsh-clave-desarrollo-cambiar-en-produccion"
	}
	return []byte(clave)
}

// FirmarValor genera un token firmado: "valor|firma_hexadecimal".
// La firma se calcula con HMAC-SHA256 usando una clave secreta que solo
// conoce el backend. Sin esa clave, es matematicamente inviable generar
// una firma valida para un valor arbitrario.
func FirmarValor(valor string) string {
	mac := hmac.New(sha256.New, obtenerClaveSecreta())
	mac.Write([]byte(valor))
	firma := hex.EncodeToString(mac.Sum(nil))
	return valor + "|" + firma
}

// ValidarValorFirmado separa el token en valor+firma, recalcula la firma
// esperada, y compara de forma segura contra timing attacks. Si el token
// fue alterado o generado fuera del sistema (por ejemplo, un QR fabricado
// a mano con un generador externo usando el mismo texto visible), la
// validacion falla.
func ValidarValorFirmado(tokenFirmado string) (string, error) {
	partes := strings.SplitN(tokenFirmado, "|", 2)
	if len(partes) != 2 {
		return "", errors.New("formato de codigo QR invalido")
	}

	valor := partes[0]
	firmaRecibida := partes[1]

	mac := hmac.New(sha256.New, obtenerClaveSecreta())
	mac.Write([]byte(valor))
	firmaEsperada := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(firmaRecibida), []byte(firmaEsperada)) {
		return "", errors.New("codigo QR invalido o no reconocido por el sistema")
	}

	return valor, nil
}