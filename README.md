# chaingate

A platform for unified signing, attestation verification, and policy-based admission across AI and software supply chains.

This project continues our position paper, [Attesting LLM Pipelines: Enforcing Verifiable Training and Release Claims](https://doi.org/10.1145/3803437.3805530), by implementing its proposed attestation-aware promotion gate and extending it to unify AI and software supply chain verification.

Under development; the features and stack below describe the planned scope.

## Features

- **Artifact signing and verification:** verify the integrity and signer identity of packages, containers, models, and datasets.
- **Provenance verification:** bind build, training, fine-tuning, and evaluation claims to artifact digests and validate relationships between stages.
- **SBOM and AI/ML-BOM:** track software dependencies, models, adapters, and datasets.
- **Security evidence:** integrate model scanning and evaluation reports with artifact-bound attestations.
- **Agent dependency attribution:** trace agent-introduced dependencies to observed requests and actions, with approval policies.
- **Admission gates:** enforce policies before merge, training, release, deployment, or model loading.
- **Backend and dashboard:** manage artifacts, evidence, trust policies, approvals, and audit records through APIs and a web interface.

## Languages and stack

| Area | Planned technologies |
| --- | --- |
| Backend and CLI | Go |
| AI integrations and workers | Python |
| Web interface | TypeScript, React |
| Storage | PostgreSQL, S3-compatible object storage |
| Provenance and signing | SLSA, in-toto, DSSE, Sigstore/Cosign, model-signing |
| Component inventories | CycloneDX SBOM and AI/ML-BOM |
| AI evidence integrations | ModelAudit or ModelScan, MLflow, lm-evaluation-harness |
| Pipeline integration | GitHub Actions, OCI registries |
