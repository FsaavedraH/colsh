package integrationtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// TestActualizarInventario - RF-14
func TestActualizarInventario(t *testing.T) {
	body := map[string]interface{}{
		"id_producto": productoIDPrueba,
		"cantidad":    1,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/inventario/actualizar", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Picking")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}

// TestListarOrdenesPicking - RF-09, RF-10
func TestListarOrdenesPicking(t *testing.T) {
	req, _ := http.NewRequest("GET", testServer.URL+"/api/picking", nil)
	req.Header.Set("X-User-Role", "Picking")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}

// TestConsultarTrazabilidad - RF-24
func TestConsultarTrazabilidad(t *testing.T) {
	idPedido := crearPedidoAux(t, productoIDPrueba, 1)

	req, _ := http.NewRequest("GET", testServer.URL+"/api/trazabilidad/"+idPedido, nil)
	req.Header.Set("X-User-Role", "Administrador")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	// El ledger esta como nil en pruebas, asi que se espera 503 (Escenario 3)
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("esperaba 503 (ledger no disponible en pruebas), obtuve %d", resp.StatusCode)
	}
}