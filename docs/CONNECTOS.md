# ConnectOS setup and personalization

ConnectOS can manage multiple companies in one repository. Create one workspace per company; workspace membership controls access. Operational records belong in PostgreSQL; company knowledge belongs in the workspace's Optimal Engine. Keep private company records and credentials out of this public GitHub repository.

## Start locally

Requirements: Node 22+, pnpm (repository pins 9.15.4), Go, PostgreSQL 17 with pgvector, Redis, Elixir/OTP, CMake, Snappy and Python 3. The optional Electron desktop is bundled; the source launcher currently installs/builds its dependencies even in headless mode.

```sh
BUSINESSOS_HEADLESS=1 BUSINESSOS_PG_BIN=/opt/homebrew/opt/postgresql@17/bin make dev-local
make dev-local-status
BUSINESSOS_HEADLESS=1 make dev-local-verify
```

Open http://localhost:5273/register to create your local account, then http://localhost:5273/os. This account is separate from GitHub and Google; OAuth requires its own setup. Use `make dev-local-stop` when finished. Restart with the same headless and PostgreSQL flags. Do not run `make clean`: that removes Docker volumes.

The launcher writes ignored `.env.dev` (ports), `frontend/.env.development.local`, `desktop/backend-go/.env` (local security keys and connection strings), and `.run/` (data, logs and binaries). Default ports: frontend 5273, backend 8801, PostgreSQL 25432, Redis 26379, engine 4200. It uses an isolated checkout database, not an existing company database.

## Sample data versus company data

`make demo OWNER=<registered-email>` bootstraps HQ, Vantage, Launchpad Cohort and Personal, then loads sample records. It is a demo refresh, not a general-purpose import for your real companies. Do not rerun it against customized sample records you want to preserve. Alex, Vantage and Launchpad are examples, not your business data.

The platform supports arbitrary company workspaces. However, the inherited `/os` operator console still contains hard-coded demo venture mappings. Creating a workspace alone does not remap every `/os` finance/funnel/agent card. Use the platform Projects, Relationships and Pipelines within each selected company workspace; adapting `/os` to arbitrary companies is a separate change.

## Where to customize

| Area | In the app / durable location | Source files |
| --- | --- | --- |
| Companies | Workspace switcher/settings; workspaces and memberships in PostgreSQL | `config/founderos/engine-topology.yaml` is the bootstrap's four sample workspace/memory routes; `desktop/backend-go/internal/services/workspace_service.go` implements company workspaces |
| Projects | `/projects`, select company workspace first | `frontend/src/routes/(app)/projects/+page.svelte`; backend `internal/handlers/projects.go` |
| Contacts and sales | `/clients` and `/pipelines`; `/os/clients` and `/os/funnel` are operator surfaces | `frontend/src/routes/(app)/clients/+page.svelte`, `(app)/pipelines/+page.svelte`; `desktop/backend-go/internal/founderos/pages/clients/`, `pages/funnel/` |
| Finance | `/finance` platform records; `/os/finances` operator overview | `frontend/src/routes/(app)/finance/+page.svelte`; `desktop/backend-go/internal/founderos/pages/finances/`; integrations in `internal/founderos/connectors/` |
| Knowledge | `/knowledge`, `/os/brain`; Settings > Optimal Engine per workspace | `frontend/src/routes/(app)/knowledge/+page.svelte`; `desktop/backend-go/internal/founderos/memory/`; bundled `optimal-engine/` |
| Agent definitions | `/os/agents`; registry/runtime and specialized workers | `desktop/backend-go/internal/founderos/agents/runtime.go`, `agents/roster/{comms,sales,clients,tech}/`; sample agents in `seed/demo/founderos-os.sql` |
| Demo venture identity | Currently Vantage and Launchpad; requires coordinated mapping changes for your companies | `desktop/backend-go/internal/founderos/pages/ventures/ventures.go`, topology YAML and seed/ETL files |
| Branding | Console, login/register, metadata, setup and welcome copy | `frontend/src/lib/founderos/chrome/`, `kit/OsMark.svelte`, `commandModules.ts`, `frontend/src/lib/components/auth/AuthLayout.svelte`, `frontend/src/app.html`, `frontend/src/routes/(founderos)/os/+layout.svelte` |
| Demo records | SQL fixtures are examples, not private data storage | `desktop/backend-go/internal/founderos/seed/demo/{founderos-os,bank,paykit,ledger}.sql`; loader `cmd/founderos-seed`, transformations `internal/founderos/etl/` |

Paths under `internal/` above are relative to `desktop/backend-go/`. Existing `founderos_*` table names, `/api/founderos` endpoints, theme/storage keys and `FOUNDEROS_*` variables remain compatible with upstream. MIT license and author attribution are retained. The existing emblem is retained pending a new ConnectOS logo.

## Environment and external dependencies

- Root `.env.example` documents the Docker configuration; `desktop/backend-go/.env.example` is the full backend/provider/OAuth reference. The local launcher generates its own ignored configuration.
- Required infrastructure: `DATABASE_URL`, `REDIS_URL`, `SECRET_KEY`, `TOKEN_ENCRYPTION_KEY`, `REDIS_KEY_HMAC_SECRET`; engine uses `OPTIMAL_ENGINE_URL` and `OPTIMAL_ENGINE_API_KEY`. Never commit generated values.
- Frontend/backend routing uses `BACKEND_PORT`, `FRONTEND_PORT`, `BACKEND_TARGET`, `VITE_BACKEND_URL`, `PUBLIC_API_URL`/`ORIGIN` where applicable. Treat Docker and local settings separately.
- AI: local Ollama is the default optional provider (`AI_PROVIDER`, `OLLAMA_LOCAL_URL`, `DEFAULT_MODEL`); cloud providers require your own API key. Installing the app does not install a model or connect a paid provider.
- OAuth providers include Google, Microsoft, Slack, Notion, HubSpot, Linear, ClickUp and Airtable; set client ID/secret and callback URI for each chosen provider. Email/password local login does not need Google OAuth.
- The operator connector implementations include email/calendar, Slack, Stripe, PayPal, Wise, PayKit, Typeform, Meta Ads, Beehiiv, Fathom, Plaud, Loom, Miro, Docusign, Robinhood and others. Each requires its own credentials/account identifiers. Example variables include `STRIPE_SECRET_KEY`, `STRIPE_VANTAGE_KEY`, `PAYKIT_LC_KEY`, `PAYKIT_VANTAGE_KEY`, `BEEHIIV_API_KEY`, `BEEHIIV_PUBLICATION_ID`, `META_ADS_ACCESS_TOKEN`. Configure only services you use and review sample venture-specific account mappings first.
- Connector credentials may come from process environment, `~/.founderos/.env`, or the Connections board. This setup does not read other apps' secrets or configure accounts automatically.
- `FOUNDEROS_WRITES` and `FOUNDEROS_CRONS` stay off unless explicitly enabled. Keep them off while exploring.
- Google Fonts supplies JetBrains Mono. Electron, Go modules, npm packages and Elixir dependencies require downloads on first start. Docker is optional; the knowledge engine uses native libraries such as SQLite/RocksDB.

## Next three customization steps

1. Sign in with your local account and use separate workspaces for each company. Add real projects, clients and deals through the platform modules, beginning with a small set of records.
2. Connect each workspace's knowledge engine and add company foundation documents (services, people, processes and goals). Adjust the operator console's venture mapping only after deciding what each company's dashboard should show.
3. Configure selected email/calendar/finance and AI providers, verify read-only results, then separately decide whether to enable outbound actions or scheduled agents.
