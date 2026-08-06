package integrationtest

import (
	"encoding/json"
	"net/http"
	"testing"
)

// TestListarReportePedidos_Autorizado - RF-28
func TestListarReportePedidos_Autorizado(t *testing.T) {
	req, _ := http.NewRequest("GET", testServer.URL+"/api/reportes/pedidos", nil)
	req.Header.Set("X-User-Role", "Administrador")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}

// TestListarReportePedidos_RolNoAutorizado - RF-27, RNF-04
func TestListarReportePedidos_RolNoAutorizado(t *testing.T) {
	req, _ := http.NewRequest("GET", testServer.URL+"/api/reportes/pedidos", nil)
	req.Header.Set("X-User-Role", "Cliente")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("esperaba 403, obtuve %d", resp.StatusCode)
	}
}

// TestReportes_SinHeaderRol - RNF-04
func TestReportes_SinHeaderRol(t *testing.T) {
	req, _ := http.NewRequest("GET", testServer.URL+"/api/reportes/pedidos", nil)
	// sin X-User-Role

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("esperaba 401, obtuve %d", resp.StatusCode)
	}
}

// TestTiemposPorEtapa - RF-29
func TestTiemposPorEtapa(t *testing.T) {
	req, _ := http.NewRequest("GET", testServer.URL+"/api/reportes/tiempos", nil)
	req.Header.Set("X-User-Role", "Administrador")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperaba 200, obtuve %d", resp.StatusCode)
	}

	var res []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Errorf("respuesta no es un JSON valido: %v", err)
	}
}

// TestIndicadorIncidencias - Indicador operativo
func TestIndicadorIncidencias(t *testing.T) {
	req, _ := http.NewRequest("GET", testServer.URL+"/api/reportes/incidencias", nil)
	req.Header.Set("X-User-Role", "Administrador")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}