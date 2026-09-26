"use client";

import { useEffect, useId, useRef, useState } from "react";
import { Html5Qrcode } from "html5-qrcode";

interface ScanBoxProps {
  onScan: (texto: string) => void;
  activo?: boolean;
  // Cuando pausado=true, la camara sigue visible pero deja de procesar
  // nuevos codigos: evita que un QR erroneo que sigue frente a la camara
  // dispare escaneos en bucle y borre el mensaje de resultado antes de
  // que el operario alcance a leerlo.
  pausado?: boolean;
}

export default function ScanBox({ onScan, activo = true, pausado = false }: ScanBoxProps) {
  const readerId = "reader-" + useId().replace(/:/g, "");
  const scannerRef = useRef<Html5Qrcode | null>(null);
  const corriendoRef = useRef(false);
  const pausadoRef = useRef(pausado);
  const [camaras, setCamaras] = useState<{ id: string; label: string }[]>([]);
  const [indiceCamara, setIndiceCamara] = useState(0);
  const [estado, setEstado] = useState("Iniciando camara...");

  useEffect(() => {
    pausadoRef.current = pausado;
  }, [pausado]);

  useEffect(() => {
    if (!activo) return;

    let cancelado = false;
    const scanner = new Html5Qrcode(readerId);
    scannerRef.current = scanner;

    Html5Qrcode.getCameras()
      .then((devices) => {
        if (cancelado) return;
        if (!devices || devices.length === 0) {
          setEstado("No se detectaron camaras en este dispositivo.");
          return;
        }
        setCamaras(devices);
        const indiceTrasera = devices.findIndex((d) => /back|rear|environment/i.test(d.label));
        const inicial = indiceTrasera >= 0 ? indiceTrasera : 0;
        setIndiceCamara(inicial);
        iniciarCamara(devices[inicial].id, cancelado);
      })
      .catch((err) => {
        if (!cancelado) setEstado("Error al listar camaras: " + err);
      });

    return () => {
      cancelado = true;
      if (corriendoRef.current && scannerRef.current) {
        scannerRef.current
          .stop()
          .then(() => {
            corriendoRef.current = false;
          })
          .catch(() => {
            corriendoRef.current = false;
          });
      }
    };
  }, [activo]);

  function iniciarCamara(id: string, cancelado: boolean) {
    if (!scannerRef.current || cancelado) return;
    scannerRef.current
      .start(id, { fps: 10, qrbox: { width: 250, height: 250 } }, (decodedText) => {
        // Ignora nuevas detecciones mientras se procesa/muestra un resultado
        // anterior, para que un QR (correcto o incorrecto) que sigue frente
        // a la camara no dispare escaneos repetidos en bucle.
        if (pausadoRef.current) return;
        onScan(decodedText);
      }, undefined)
      .then(() => {
        if (!cancelado) {
          corriendoRef.current = true;
          setEstado("Escaneando...");
        }
      })
      .catch((err) => {
        if (!cancelado) setEstado("Error al iniciar camara: " + err);
      });
  }

  function cambiarCamara() {
    if (camaras.length < 2 || !scannerRef.current || !corriendoRef.current) return;
    scannerRef.current
      .stop()
      .then(() => {
        corriendoRef.current = false;
        const siguiente = (indiceCamara + 1) % camaras.length;
        setIndiceCamara(siguiente);
        iniciarCamara(camaras[siguiente].id, false);
      })
      .catch((err) => {
        corriendoRef.current = false;
        setEstado("Error al cambiar camara: " + err);
      });
  }

  return (
    <div className="w-full">
      <p className="text-sm text-gray-500 mb-2">
        {pausado && corriendoRef.current ? "Escaneo en pausa" : estado}
      </p>
      <div id={readerId} className="rounded-xl overflow-hidden border-2 border-gray-300" />
      {camaras.length > 1 && (
        <button
          onClick={cambiarCamara}
          className="mt-3 px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded-lg text-sm font-medium"
        >
          Cambiar cámara
        </button>
      )}
    </div>
  );
}