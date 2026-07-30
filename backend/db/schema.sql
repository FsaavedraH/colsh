--
-- PostgreSQL database dump
--

\restrict iv2FRsNT19QWURGnXrnMqOqVeRl4ee3QpLC9jxAsLqEuPXImpF2SaCEoMaMNRd9

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
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: detalle_pedido; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.detalle_pedido (
    id_detalle uuid DEFAULT gen_random_uuid() NOT NULL,
    id_pedido uuid NOT NULL,
    id_producto uuid NOT NULL,
    cantidad integer NOT NULL,
    CONSTRAINT detalle_pedido_cantidad_check CHECK ((cantidad > 0))
);


ALTER TABLE public.detalle_pedido OWNER TO postgres;

--
-- Name: evento_trazabilidad; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.evento_trazabilidad (
    id_evento uuid DEFAULT gen_random_uuid() NOT NULL,
    id_pedido uuid NOT NULL,
    estado character varying(40) NOT NULL,
    fecha timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    responsable uuid NOT NULL,
    tx_id character varying(100)
);


ALTER TABLE public.evento_trazabilidad OWNER TO postgres;

--
-- Name: inventario; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.inventario (
    id_inventario uuid DEFAULT gen_random_uuid() NOT NULL,
    id_producto uuid NOT NULL,
    stock integer NOT NULL,
    ubicacion character varying(100) NOT NULL,
    CONSTRAINT inventario_stock_check CHECK ((stock >= 0))
);


ALTER TABLE public.inventario OWNER TO postgres;

--
-- Name: pedido; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.pedido (
    id_pedido uuid DEFAULT gen_random_uuid() NOT NULL,
    fecha_creacion timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    estado character varying(40) NOT NULL,
    id_cliente uuid NOT NULL,
    direccion_entrega text NOT NULL,
    CONSTRAINT pedido_estado_check CHECK (((estado)::text = ANY ((ARRAY['Pendiente'::character varying, 'En espera por inventario'::character varying, 'En recolecci¢n'::character varying, 'En empaque'::character varying, 'En despacho'::character varying, 'Entregado'::character varying])::text[])))
);


ALTER TABLE public.pedido OWNER TO postgres;

--
-- Name: producto; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.producto (
    id_producto uuid DEFAULT gen_random_uuid() NOT NULL,
    nombre character varying(200) NOT NULL
);


ALTER TABLE public.producto OWNER TO postgres;

--
-- Name: usuario; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.usuario (
    id_usuario uuid DEFAULT gen_random_uuid() NOT NULL,
    nombre character varying(150) NOT NULL,
    email character varying(255) NOT NULL,
    rol character varying(30) NOT NULL,
    password_hash text NOT NULL,
    CONSTRAINT usuario_rol_check CHECK (((rol)::text = ANY ((ARRAY['Cliente'::character varying, 'Picking'::character varying, 'Empaque'::character varying, 'Transportista'::character varying, 'Administrador'::character varying])::text[])))
);


ALTER TABLE public.usuario OWNER TO postgres;

--
-- Name: detalle_pedido detalle_pedido_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detalle_pedido
    ADD CONSTRAINT detalle_pedido_pkey PRIMARY KEY (id_detalle);


--
-- Name: evento_trazabilidad evento_trazabilidad_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.evento_trazabilidad
    ADD CONSTRAINT evento_trazabilidad_pkey PRIMARY KEY (id_evento);


--
-- Name: inventario inventario_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventario
    ADD CONSTRAINT inventario_pkey PRIMARY KEY (id_inventario);


--
-- Name: pedido pedido_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pedido
    ADD CONSTRAINT pedido_pkey PRIMARY KEY (id_pedido);


--
-- Name: producto producto_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.producto
    ADD CONSTRAINT producto_pkey PRIMARY KEY (id_producto);


--
-- Name: usuario usuario_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.usuario
    ADD CONSTRAINT usuario_email_key UNIQUE (email);


--
-- Name: usuario usuario_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.usuario
    ADD CONSTRAINT usuario_pkey PRIMARY KEY (id_usuario);


--
-- Name: idx_evento_pedido; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_evento_pedido ON public.evento_trazabilidad USING btree (id_pedido);


--
-- Name: idx_pedido_estado; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pedido_estado ON public.pedido USING btree (estado);


--
-- Name: idx_pedido_fecha; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_pedido_fecha ON public.pedido USING btree (fecha_creacion);


--
-- Name: detalle_pedido detalle_pedido_id_pedido_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detalle_pedido
    ADD CONSTRAINT detalle_pedido_id_pedido_fkey FOREIGN KEY (id_pedido) REFERENCES public.pedido(id_pedido) ON DELETE CASCADE;


--
-- Name: detalle_pedido detalle_pedido_id_producto_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.detalle_pedido
    ADD CONSTRAINT detalle_pedido_id_producto_fkey FOREIGN KEY (id_producto) REFERENCES public.producto(id_producto) ON DELETE RESTRICT;


--
-- Name: evento_trazabilidad evento_trazabilidad_id_pedido_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.evento_trazabilidad
    ADD CONSTRAINT evento_trazabilidad_id_pedido_fkey FOREIGN KEY (id_pedido) REFERENCES public.pedido(id_pedido) ON DELETE CASCADE;


--
-- Name: evento_trazabilidad evento_trazabilidad_responsable_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.evento_trazabilidad
    ADD CONSTRAINT evento_trazabilidad_responsable_fkey FOREIGN KEY (responsable) REFERENCES public.usuario(id_usuario) ON DELETE RESTRICT;


--
-- Name: inventario inventario_id_producto_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inventario
    ADD CONSTRAINT inventario_id_producto_fkey FOREIGN KEY (id_producto) REFERENCES public.producto(id_producto) ON DELETE RESTRICT;


--
-- Name: pedido pedido_id_cliente_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.pedido
    ADD CONSTRAINT pedido_id_cliente_fkey FOREIGN KEY (id_cliente) REFERENCES public.usuario(id_usuario) ON DELETE RESTRICT;


--
-- PostgreSQL database dump complete
--

\unrestrict iv2FRsNT19QWURGnXrnMqOqVeRl4ee3QpLC9jxAsLqEuPXImpF2SaCEoMaMNRd9

