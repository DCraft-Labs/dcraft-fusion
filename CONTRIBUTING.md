# Contributing to DCraft Fusion

Thank you for your interest in contributing to DCraft Fusion. This guide covers how to propose changes, the Developer Certificate of Origin (DCO) requirement, and how to validate your work locally.

## Ways to contribute

- Report bugs and request features via [GitHub Issues](https://github.com/DCraft-Labs/dcraft-fusion/issues).
- Ask questions and share ideas in [GitHub Discussions](https://github.com/DCraft-Labs/dcraft-fusion/discussions).
- Improve docs, tests, Helm charts, or engine adapters with a pull request.

## Fork and pull request workflow

1. Fork [DCraft-Labs/dcraft-fusion](https://github.com/DCraft-Labs/dcraft-fusion).
2. Create a topic branch from `main` (or the current default branch).
3. Make focused commits that are easy to review.
4. Push your branch and open a pull request against the upstream repository.
5. Fill out the PR template and link related issues when applicable.
6. Ensure CI is green and respond to review feedback.

Please keep PRs scoped. Prefer small, reviewable changes over large multi-purpose patches.

## Developer Certificate of Origin (DCO)

All contributions must be signed off under the [Developer Certificate of Origin](https://developercertificate.org/).

Each commit message must include a `Signed-off-by` line:

```text
Signed-off-by: Your Name <your.email@example.com>
```

You can add the trailer automatically:

```bash
git commit -s -m "Your commit message"
```

By signing off, you certify that you have the right to submit the work under the project's Apache-2.0 license.

## Local development and tests

Prerequisites: Node.js 24+, npm 11+, and Docker Desktop (for the local compose stack).

Install dependencies and run workspace checks from the repository root:

```bash
npm install
npm test
npm run typecheck
npm run build
```

Start the local Phase 1 stack:

```powershell
docker compose -f infra\local-dev\docker-compose.yml up --build -d
```

- Web: http://127.0.0.1:5174
- API: http://127.0.0.1:8080

For frontend-only iteration:

```bash
npm run dev:web
```

## Code of conduct and security

- Community standards: [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
- Security reports: [SECURITY.md](SECURITY.md)

## License

Contributions are accepted under the [Apache License 2.0](LICENSE). See [NOTICE](NOTICE) for attribution.
