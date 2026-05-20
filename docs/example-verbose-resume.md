# Alex Rivera — Verbose Resume

> Fictional example only — not a real person or employer. Use as a structural reference for your own `verboseResume.md`.

## Contact

- **Location:** Portland, OR (open to remote)
- **Phone:** 555-0142
- **Email:** alex.rivera@example.com
- **Website:** https://example.com
- **GitHub:** https://github.com/example/arivera

## Summary

Platform and backend engineer with about twelve years of experience building internal developer platforms, Kubernetes estates, and data-adjacent infrastructure. I tend to join during growth phases: when a product org needs paved roads, when observability is reactive, or when platform work has been split across too many part-time owners.

I am comfortable writing Terraform and Helm, pairing with application teams on service onboarding, and explaining tradeoffs to directors who care about cost, risk, and time-to-ship. I have led small platform squads, mentored contractors, and partnered with security and networking on private connectivity patterns.

This file is intentionally long. Most of it should never appear on a one-page tailored resume.

## Skills

### Languages

- Go
- Python
- Bash
- SQL (PostgreSQL, BigQuery)

### Cloud & Containers

- Kubernetes (EKS, GKE)
- AWS (VPC, IAM, RDS, MSK)
- GCP (GKE, Cloud Run, Pub/Sub)
- Linux administration

### CI / CD & GitOps

- GitHub Actions
- Argo CD
- Helm
- Terraform

### Observability

- Prometheus
- Grafana
- OpenTelemetry
- Loki

### Data & Messaging

- Kafka (MSK, Confluent Cloud)
- Kafka Connect
- Debezium (evaluation)

## Experience

### Meridian Analytics — Staff Platform Engineer

**Dates:** Mar 2022 – Present | **Location:** Remote (US)

#### Overview

Meridian sells analytics dashboards to mid-market retailers. When I joined, forty product microservices ran on a single overstuffed EKS cluster with hand-rolled Helm and no standard labels. Deployments were manual, on-call was noisy, and new hires took weeks to ship their first change.

I was hired onto the **Platform Foundations** group (eight engineers, two managers). My charter was cluster standards, developer onboarding, and observability baselines—not owning every data pipeline in the company.

Within six months the group split into **Runtime** (clusters, networking, GitOps) and **Developer Experience** (templates, docs, local dev). I stayed on Runtime and became the de facto tech lead when the previous lead moved to management.

#### Teams

- Platform Foundations (initial home team)
- Runtime Platform (post-split)
- Incident working group (rotating member)
- Retail Integrations program (dotted-line partner for Kafka onboarding)

#### Notable work

##### EKS cluster standards and fleet rollout

Defined a “golden cluster” checklist: IRSA for controllers, separate node groups for system vs workload, enforced Pod Security standards, and mandatory `app.kubernetes.io/*` labels for service discovery.

Wrote Terraform modules consumed by three product lines. Migrated the legacy cluster in phases (namespace-by-namespace) with rollback runbooks. Cut deployment-related Sev2 pages by roughly half over two quarters (internal metric; not a customer-facing SLA).

##### GitHub Actions reusable workflows for service repos

Partnered with Developer Experience to publish composite actions: lint, test, image build, Trivy scan, and promote-to-Argo. Documented copy-paste examples for Node and Go services.

Reduced bespoke pipeline YAML in the org; most new services adopted the template within one sprint of launch.

##### Observability baseline (Prometheus + Grafana + OTel)

Standardized scrape annotations, recording rules for HTTP latency, and dashboards per tier (ingress, service, datastore). Introduced OpenTelemetry collector Helm chart with tail sampling for high-cardinality paths.

Wrote an onboarding doc “Your first RED dashboard in 30 minutes” that became required reading in engineering onboarding.

##### Kafka onboarding for Retail Integrations

Retail Integrations needed near-real-time inventory feeds. I was lent 20% time to help them land on MSK with TLS and IAM auth.

Delivered Terraform for topics, ACL patterns, and a reference producer in Go. Paired with their team on consumer lag alerts and retry semantics. Did not own their business logic—only the paved road.

##### On-call and incident hygiene

Authored postmortems for two major incidents (cert expiry on internal CA; Argo CD repo-server OOM). Drove action items: cert-manager adoption, repo-server HA, and game days for platform.

Mentored two mid-level engineers on constructive postmortem writing and blameless facilitation.

---

### Lattice Systems — Senior Software Engineer, Infrastructure

**Dates:** Jun 2017 – Feb 2022 | **Location:** Seattle, WA (hybrid)

#### Overview

Lattice built B2B workflow automation. I started on the **Core Services** team maintaining JVM APIs and PostgreSQL. After a reorg I moved to **Cloud Infrastructure** reporting to the same director but with a platform mandate.

#### Teams

- Core Services (years 1–2)
- Cloud Infrastructure (years 3–5)
- Security champions guild (volunteer)

#### Notable work

##### Core Services — API reliability push

First role at Lattice. Owned payment webhook ingestion path: idempotency keys, dead-letter queue, replay tooling.

Reduced duplicate charges in staging environments before a major partner launch. Learned the domain deeply enough to interview new backend hires.

##### Move to Cloud Infrastructure

Promotion followed reorg. Charter: move from bespoke EC2 + Ansible to Kubernetes on AWS.

##### kops to EKS migration

Led technical discovery comparing managed control planes vs self-managed. Wrote ADR adopted by leadership.

Executed migration over four months with dual-write period for critical config services. Training sessions for app teams on Deployments vs their old pet servers.

##### Secrets management (Vault)

Integrated HashiCorp Vault with Kubernetes auth. Rotated database credentials without application restarts using sidecar pattern (later simplified when apps adopted dynamic secrets SDK).

Partnered with Security on policy-as-code for path naming and audit log review.

##### Cost visibility

Built monthly reports tagging AWS spend by team label. Surfaced “orphan ELBs” and oversized RDS instances.

Not glamorous work, but it funded two additional platform headcount the following year.

---

### Northwind Digital — Software Engineer

**Dates:** Aug 2013 – May 2017 | **Location:** Boise, ID

#### Overview

Agency-style consulting: two-week to six-month engagements. Mostly PHP and JavaScript early; later small Go utilities for clients modernizing stacks.

#### Teams

- Delivery Team A (rotating client assignments)
- Internal tooling (10% time)

#### Notable work

##### Client: regional bank — online account opening

Built form flows and backend validation. Learned PCI scope boundaries and why platform teams exist.

##### Client: logistics startup — shipment tracking MVP

Introduced Docker Compose for local dev; helped them hire their first full-time DevOps contractor.

##### Internal — proposal estimation spreadsheet

Scripted historical velocity export from Jira to improve bid accuracy (Python). Pet project; still referenced in all-hands jokingly.

## Education

- B.S. Computer Science | Boise State University | Boise, ID | May 2013

## Certifications

- Certified Kubernetes Administrator (CKA), 2021
- AWS Solutions Architect – Associate, 2019 (lapsed; plan to renew)

## Appendix — facts I might forget without this file

- Meridian re:Invent booth 2024 — demoed GitOps flow; no production changes that week.
- Lattice office move Q3 2019 — one week of flaky VPN; documented split-tunnel workaround.
- Preferred interview stories: EKS migration ADR, Kafka onboarding, postmortem culture.
