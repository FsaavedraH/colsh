package integrationtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

const transportistaIDPrueba = "939baa61-34d4-42f2-8b53-6c7fe668db31"

// crearPedidoEnDespachoAux lleva un pedido hasta el estado "En despacho"
func crearPedidoEnDespachoAux(t *testing.T, idProducto string, cantidad int) string {
	t.Helper()
	idPedido := crearPedidoEnEmpaqueAux(t, idProducto, cantidad)

	body := map[string]string{
		"id_pedido":   idPedido,
		"responsable": responsableEmpaquePrueba,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/empaque", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Empaque")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error al confirmar empaque auxiliar: %v", err)
	}
	resp.Body.Close()

	return idPedido
}

// TestGenerarDespacho - RF-19, RF-20, RF-21
func TestGenerarDespacho(t *testing.T) {
	idPedido := crearPedidoEnDespachoAux(t, "86d331a6-7420-4573-869d-f7d902424b43", 1)

	body := map[string]string{
		"id_pedido":     idPedido,
		"transportista": transportistaIDPrueba,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/despacho", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Transportista")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if res["estado"] != "en_despacho" {
		t.Errorf("esperaba estado 'en_despacho', obtuve '%v'", res["estado"])
	}
	if res["codigo_seguimiento"] == "" || res["codigo_seguimiento"] == nil {
		t.Error("esperaba un codigo_seguimiento no vacio")
	}
}

// TestConfirmarEntrega - RF-22, RF-23
func TestConfirmarEntrega(t *testing.T) {
	idPedido := crearPedidoEnDespachoAux(t, "86d331a6-7420-4573-869d-f7d902424b43", 1)

	bodyDespacho := map[string]string{
		"id_pedido":     idPedido,
		"transportista": transportistaIDPrueba,
	}
	bb, _ := json.Marshal(bodyDespacho)
	reqDespacho, _ := http.NewRequest("POST", testServer.URL+"/api/despacho", bytes.NewReader(bb))
	reqDespacho.Header.Set("Content-Type", "application/json")
	reqDespacho.Header.Set("X-User-Role", "Transportista")
	respDespacho, _ := http.DefaultClient.Do(reqDespacho)
	respDespacho.Body.Close()

	bodyEntrega := map[string]string{
		"id_pedido":     idPedido,
		"transportista": transportistaIDPrueba,
	}
	bodyBytes, _ := json.Marshal(bodyEntrega)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/entrega", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Transportista")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	var res map[string]string
	json.NewDecoder(resp.Body).Decode(&res)

	if res["estado"] != "entregado" {
		t.Errorf("esperaba estado 'entregado', obtuve '%s'", res["estado"])
	}
}

// TestDespacho_RolNoAutorizado - RNF-04
func TestDespacho_RolNoAutorizado(t *testing.T) {
	idPedido := crearPedidoEnDespachoAux(t, "86d331a6-7420-4573-869d-f7d902424b43", 1)

	body := map[string]string{
		"id_pedido":     idPedido,
		"transportista": transportistaIDPrueba,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/despacho", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Cliente") // no autorizado

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("esperaba 403, obtuve %d", resp.StatusCode)
	}
}