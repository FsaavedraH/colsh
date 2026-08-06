package integrationtest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

const (
	clienteIDPrueba  = "fbada7ec-ee17-4f0a-80db-b5d469adb4d4"
	productoIDPrueba = "088b064e-fdb3-4032-b360-57174fcc60ee"
)

// TestCrearPedido_Exitoso - RF-01, RF-03 / HU-1 Escenario 1
func TestCrearPedido_Exitoso(t *testing.T) {
	body := map[string]interface{}{
		"cliente_id":        clienteIDPrueba,
		"direccion_entrega": "Cra 50 #10-25",
		"productos": []map[string]interface{}{
			{"id_producto": productoIDPrueba, "cantidad": 1},
		},
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", testServer.URL+"/api/pedidos", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Cliente")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error al hacer la peticion: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("esperaba 201, obtuve %d: %s", resp.StatusCode, string(b))
	}

	var respuesta map[string]string
	json.NewDecoder(resp.Body).Decode(&respuesta)

	if respuesta["estado"] != "Pendiente" {
		t.Errorf("esperaba estado 'Pendiente', obtuve '%s'", respuesta["estado"])
	}
	if respuesta["id_pedido"] == "" {
		t.Error("esperaba un id_pedido no vacio")
	}
}

// TestCrearPedido_DatosIncompletos - RF-02 / HU-1 Escenario 2
func TestCrearPedido_DatosIncompletos(t *testing.T) {
	body := map[string]interface{}{
		"direccion_entrega": "Cra 50 #10-25",
		"productos":         []map[string]interface{}{},
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", testServer.URL+"/api/pedidos", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Cliente")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error al hacer la peticion: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("esperaba 400, obtuve %d", resp.StatusCode)
	}
}

// TestCrearPedido_RolNoAutorizado - RF-27, RNF-04, Escenario Arquitectural 5
func TestCrearPedido_RolNoAutorizado(t *testing.T) {
	body := map[string]interface{}{
		"cliente_id":        clienteIDPrueba,
		"direccion_entrega": "Cra 50 #10-25",
		"productos": []map[string]interface{}{
			{"id_producto": productoIDPrueba, "cantidad": 1},
		},
	}
	bodyBytes, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", testServer.URL+"/api/pedidos", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Picking") // rol incorrecto para crear pedido

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error al hacer la peticion: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("esperaba 403, obtuve %d", resp.StatusCode)
	}
}

// TestConsultarPedido - RF-04
func TestConsultarPedido(t *testing.T) {
	body := map[string]interface{}{
		"cliente_id":        clienteIDPrueba,
		"direccion_entrega": "Cra 50 #10-25",
		"productos": []map[string]interface{}{
			{"id_producto": productoIDPrueba, "cantidad": 1},
		},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/pedidos", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Cliente")
	resp, _ := http.DefaultClient.Do(req)
	var creado map[string]string
	json.NewDecoder(resp.Body).Decode(&creado)
	resp.Body.Close()

	reqGet, _ := http.NewRequest("GET", testServer.URL+"/api/pedidos/"+creado["id_pedido"], nil)
	reqGet.Header.Set("X-User-Role", "Cliente")

	respGet, err := http.DefaultClient.Do(reqGet)
	if err != nil {
		t.Fatalf("error al consultar pedido: %v", err)
	}
	defer respGet.Body.Close()

	if respGet.StatusCode != http.StatusOK {
		t.Errorf("esperaba 200, obtuve %d", respGet.StatusCode)
	}
}