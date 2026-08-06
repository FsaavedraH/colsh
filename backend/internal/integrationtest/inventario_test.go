package integrationtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func crearPedidoAux(t *testing.T, idProducto string, cantidad int) string {
	t.Helper()
	body := map[string]interface{}{
		"cliente_id":        clienteIDPrueba,
		"direccion_entrega": "Cra 50 #10-25",
		"productos": []map[string]interface{}{
			{"id_producto": idProducto, "cantidad": cantidad},
		},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/pedidos", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Cliente")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error al crear pedido auxiliar: %v", err)
	}
	defer resp.Body.Close()

	var creado map[string]string
	json.NewDecoder(resp.Body).Decode(&creado)
	return creado["id_pedido"]
}

// TestValidarInventario_StockSuficiente - RF-05, RF-09
func TestValidarInventario_StockSuficiente(t *testing.T) {
	idPedido := crearPedidoAux(t, productoIDPrueba, 1)

	body := map[string]string{"id_pedido": idPedido}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/inventario/validar", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Picking")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if disponible, ok := res["disponibilidad"].(bool); !ok || !disponible {
		t.Errorf("esperaba disponibilidad=true, obtuve %v", res["disponibilidad"])
	}
	if res["estado_pedido"] != "En recoleccion" {
		t.Errorf("esperaba estado 'En recoleccion', obtuve '%v'", res["estado_pedido"])
	}
}

// TestValidarInventario_StockInsuficiente - RF-06, RF-08
func TestValidarInventario_StockInsuficiente(t *testing.T) {
	idPedido := crearPedidoAux(t, productoIDPrueba, 999999)

	body := map[string]string{"id_pedido": idPedido}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/inventario/validar", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Picking")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if disponible, ok := res["disponibilidad"].(bool); !ok || disponible {
		t.Errorf("esperaba disponibilidad=false, obtuve %v", res["disponibilidad"])
	}
	if res["estado_pedido"] != "En espera por inventario" {
		t.Errorf("esperaba estado 'En espera por inventario', obtuve '%v'", res["estado_pedido"])
	}
}

// TestValidarInventario_RolNoAutorizado - RNF-04
func TestValidarInventario_RolNoAutorizado(t *testing.T) {
	idPedido := crearPedidoAux(t, productoIDPrueba, 1)

	body := map[string]string{"id_pedido": idPedido}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/inventario/validar", bytes.NewReader(bodyBytes))
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