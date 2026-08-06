package integrationtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

const responsableEmpaquePrueba = "e83f286d-7f7c-4f6d-ae41-b3c9bf63df14"

// crearPedidoEnEmpaqueAux crea un pedido y lo lleva hasta el estado "En empaque"
func crearPedidoEnEmpaqueAux(t *testing.T, idProducto string, cantidad int) string {
	t.Helper()
	idPedido := crearYValidarPedidoAux(t, idProducto, cantidad)

	body := map[string]interface{}{
		"id_pedido":   idPedido,
		"id_producto": idProducto,
		"cantidad":    cantidad,
		"responsable": responsablePickingPrueba,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/recoleccion", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Picking")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error al confirmar recoleccion auxiliar: %v", err)
	}
	resp.Body.Close()

	return idPedido
}

// TestRecepcionEmpaque - RF-16
func TestRecepcionEmpaque(t *testing.T) {
	idPedido := crearPedidoEnEmpaqueAux(t, "2994aaa3-5291-42fc-9784-43156f5b6567", 1)

	body := map[string]string{"id_pedido": idPedido}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/empaque/recepcion", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Empaque")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}

// TestEscanearValidacionEmpaque_Incorrecto - RF-17, RF-26
func TestEscanearValidacionEmpaque_Incorrecto(t *testing.T) {
	idPedido := crearPedidoEnEmpaqueAux(t, "2994aaa3-5291-42fc-9784-43156f5b6567", 1)

	body := map[string]string{
		"id_pedido":              idPedido,
		"id_producto_esperado":  "2994aaa3-5291-42fc-9784-43156f5b6567",
		"id_producto_escaneado": "86d331a6-7420-4573-869d-f7d902424b43",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/empaque/escanear", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Empaque")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("esperaba 409, obtuve %d", resp.StatusCode)
	}
}

// TestConfirmarEmpaque - RF-18
func TestConfirmarEmpaque(t *testing.T) {
	idPedido := crearPedidoEnEmpaqueAux(t, "2994aaa3-5291-42fc-9784-43156f5b6567", 1)

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
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	var res map[string]string
	json.NewDecoder(resp.Body).Decode(&res)

	if res["estado"] != "empacado" {
		t.Errorf("esperaba estado 'empacado', obtuve '%s'", res["estado"])
	}
}