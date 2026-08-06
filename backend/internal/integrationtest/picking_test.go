package integrationtest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

const (
	ubicacionCorrectaPrueba = "B-12-03-02"
	responsablePickingPrueba = "68215789-6d0d-441c-ac90-eedf8aca4a61"
)

func crearYValidarPedidoAux(t *testing.T, idProducto string, cantidad int) string {
	t.Helper()
	idPedido := crearPedidoAux(t, idProducto, cantidad)

	body := map[string]string{"id_pedido": idPedido}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/inventario/validar", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Picking")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error al validar pedido auxiliar: %v", err)
	}
	resp.Body.Close()

	return idPedido
}

// TestEscanearUbicacion_Correcta - RF-11
func TestEscanearUbicacion_Correcta(t *testing.T) {
	idPedido := crearYValidarPedidoAux(t, productoIDPrueba, 1)

	body := map[string]string{
		"id_pedido":            idPedido,
		"id_producto":          productoIDPrueba,
		"ubicacion_escaneada": ubicacionCorrectaPrueba,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/picking/escanear-ubicacion", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Picking")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if coincide, ok := res["coincide"].(bool); !ok || !coincide {
		t.Errorf("esperaba coincide=true, obtuve %v", res["coincide"])
	}
}

// TestEscanearUbicacion_Incorrecta - RF-11, RF-26
func TestEscanearUbicacion_Incorrecta(t *testing.T) {
	idPedido := crearYValidarPedidoAux(t, productoIDPrueba, 1)

	body := map[string]string{
		"id_pedido":            idPedido,
		"id_producto":          productoIDPrueba,
		"ubicacion_escaneada": "Z-99-99-99",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/picking/escanear-ubicacion", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Picking")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("esperaba 409, obtuve %d", resp.StatusCode)
	}

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	if coincide, ok := res["coincide"].(bool); !ok || coincide {
		t.Errorf("esperaba coincide=false, obtuve %v", res["coincide"])
	}
}

// TestEscanearProducto_Incorrecto - RF-12, RF-13, RF-26
func TestEscanearProducto_Incorrecto(t *testing.T) {
	idPedido := crearYValidarPedidoAux(t, productoIDPrueba, 1)

	body := map[string]string{
		"id_pedido":              idPedido,
		"id_producto_esperado":  productoIDPrueba,
		"id_producto_escaneado": "86d331a6-7420-4573-869d-f7d902424b43",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/picking/escanear-producto", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Picking")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		t.Errorf("esperaba 409, obtuve %d", resp.StatusCode)
	}
}

// TestConfirmarRecoleccion - RF-14, RF-15
func TestConfirmarRecoleccion(t *testing.T) {
	idPedido := crearYValidarPedidoAux(t, productoIDPrueba, 1)

	body := map[string]interface{}{
		"id_pedido":   idPedido,
		"id_producto": productoIDPrueba,
		"cantidad":    1,
		"responsable": responsablePickingPrueba,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/recoleccion", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Role", "Picking")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	var res map[string]string
	json.NewDecoder(resp.Body).Decode(&res)

	if res["estado"] != "recolectado" {
		t.Errorf("esperaba estado 'recolectado', obtuve '%s'", res["estado"])
	}
}