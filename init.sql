CREATE TABLE IF NOT EXISTS clientes (
    cliente_id SERIAL PRIMARY KEY,
    nome VARCHAR(30),
    tel BIGINT
);

CREATE TABLE IF NOT EXISTS comandas (
    comanda_id SERIAL PRIMARY KEY,
    cliente_id INTEGER REFERENCES clientes(cliente_id),
    data_entrada DATE,
    data_entrega DATE,
    data_retirada DATE,
    servico VARCHAR(100),
    preco_total NUMERIC(10,2),
    sinal NUMERIC(10,2),
    valor_restante NUMERIC(10,2),
    tipo_pag VARCHAR(10),
    status_name VARCHAR(10)
); 

CREATE TABLE IF NOT EXISTS log_comandas (
    log_id SERIAL PRIMARY KEY,
    comanda_id INTEGER REFERENCES comandas(comanda_id),
    cliente_id INTEGER REFERENCES clientes(cliente_id),
    data_evento DATE,
    servico VARCHAR(100),
    preco_total NUMERIC(10,2),
    sinal NUMERIC(10,2),
    valor_restante NUMERIC(10,2),
    status_name VARCHAR(10)
);

CREATE TABLE IF NOT EXISTS registros (
    registro_id SERIAL PRIMARY KEY,
    valor NUMERIC(10,2),
    data_entrada DATE,
    comanda_id INTEGER REFERENCES comandas(comanda_id),
    cliente_id INTEGER REFERENCES clientes(cliente_id),
    status_name VARCHAR(10),
    tipo_pag VARCHAR(10),
    categoria VARCHAR(30),
    descricao VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS fluxo (
    fluxo_id SERIAL PRIMARY KEY,
    valor NUMERIC(10,2),
    data_entrada DATE,
    categoria VARCHAR(30),
    comanda_id INTEGER REFERENCES comandas(comanda_id)
);

CREATE TABLE IF NOT EXISTS user (
    user_id SERIAL PRIMARY KEY,
    whatsapp BIGINT,
    created_at TIMESTAMP,
    active BOOLEAN
)

