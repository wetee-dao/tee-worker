#![feature(lazy_cell)]
use std::net::SocketAddr;

use axum::{routing::get, Router};
use axum_server::tls_rustls::RustlsConfig;

mod operators;
use crate::operators::{farmpod, llama, get_api_operators};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let app = Router::new()
        .route("/apis/v0", get(get_api_operators))
        .route(
            "/apis/v0/namespaces/:namespace/llamas",
            get(llama::list_llamas),
        )
        .route(
            "/apis/v0/namespaces/:namespace/llamas/:name",
            get(llama::get_llama),
        )
        .route(
            "/apis/v0/namespaces/:namespace/farmpods",
            get(farmpod::list_farmpods),
        );

    let tls_cert = rcgen::generate_simple_self_signed(vec!["localhost".to_string()])?;
    let tls_config = RustlsConfig::from_der(
        vec![tls_cert.serialize_der()?],
        tls_cert.serialize_private_key_der(),
    )
    .await?;

    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));

    println!("listening on {addr}");

    axum_server::bind_rustls(addr, tls_config)
        .serve(app.into_make_service())
        .await
        .map_err(anyhow::Error::from)
}
