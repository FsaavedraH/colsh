--
-- PostgreSQL database dump
--

\restrict Tp5vISHmeSLkn1Cpjb5EZRSkCeFsZ4blUC3SdPv9VbKo9TdwD49gGcHRdh6aHTm

-- Dumped from database version 17.3
-- Dumped by pg_dump version 18.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: usuario; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.usuario (id_usuario, nombre, email, rol, password_hash) FROM stdin;
68215789-6d0d-441c-ac90-eedf8aca4a61	Juan P‚rez	juan.perez@colsh.com	Picking	placeholder
e83f286d-7f7c-4f6d-ae41-b3c9bf63df14	Mar¡a G¢mez	maria.gomez@colsh.com	Empaque	placeholder
fbada7ec-ee17-4f0a-80db-b5d469adb4d4	Distribuidora del Valle	cliente@colsh.com	Cliente	placeholder
5b0e84c6-4f36-4c29-be9d-3c40fcde8629	Admin ColSh	admin@colsh.com	Administrador	placeholder
\.


--
-- Data for Name: pedido; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.pedido (id_pedido, fecha_creacion, estado, id_cliente, direccion_entrega) FROM stdin;
\.


--
-- Data for Name: producto; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.producto (id_producto, nombre) FROM stdin;
088b064e-fdb3-4032-b360-57174fcc60ee	Pastillas de freno delanteras
2994aaa3-5291-42fc-9784-43156f5b6567	Filtro de aire
86d331a6-7420-4573-869d-f7d902424b43	Buj¡a NGK
\.


--
-- Data for Name: detalle_pedido; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.detalle_pedido (id_detalle, id_pedido, id_producto, cantidad) FROM stdin;
\.


--
-- Data for Name: evento_trazabilidad; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.evento_trazabilidad (id_evento, id_pedido, estado, fecha, responsable, tx_id) FROM stdin;
\.


--
-- Data for Name: inventario; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.inventario (id_inventario, id_producto, stock, ubicacion) FROM stdin;
f6940b72-e9bc-4455-8978-e0c673d1c65b	088b064e-fdb3-4032-b360-57174fcc60ee	10	B-12-03-02
d43ddbf0-083f-4d89-938f-8e8b69b73499	2994aaa3-5291-42fc-9784-43156f5b6567	15	B-05-01-01
\.


--
-- PostgreSQL database dump complete
--

\unrestrict Tp5vISHmeSLkn1Cpjb5EZRSkCeFsZ4blUC3SdPv9VbKo9TdwD49gGcHRdh6aHTm

