use axum::{response::IntoResponse, Json};
use k8s_openapi::apimachinery::pkg::apis::meta::v1::{APIResourceList, APIResource};
use kube::Resource;

pub mod farmpod;
pub mod llama;

pub async fn get_api_operators() -> impl IntoResponse {
    Json(APIResourceList {
        group_version: "app.wetee.tee-worker/v1alpha".to_string(),
        resources: vec![
            APIResource {
                group: Some(llama::Llama::group(&()).into()),
                kind: llama::Llama::kind(&()).into(),
                name: llama::Llama::plural(&()).into(),
                namespaced: true,
                verbs: vec!["list".to_string(), "get".to_string()],
                ..Default::default()
            },
            APIResource {
                group: Some(farmpod::FarmPod::group(&()).into()),
                kind: farmpod::FarmPod::kind(&()).into(),
                name: farmpod::FarmPod::plural(&()).into(),
                namespaced: true,
                verbs: vec!["list".to_string()],
                ..Default::default()
            },
        ],
    })
}
