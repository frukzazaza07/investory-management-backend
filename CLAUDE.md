# Project Rules

## Frontend integration docs

Whenever a new API-facing feature is added (new endpoint, new field with constrained values, new response shape, etc.), update [FRONTEND_GUIDE.md](FRONTEND_GUIDE.md) in the same change:

- Add/update the relevant TypeScript type in the "TypeScript Types" section.
- Add/update the React Query hook(s) in the matching "Hooks — ..." section, following the existing `KEYS`/`useQuery`/`useMutation` pattern.
- If the feature changes how a form field should be filled (e.g. a value must come from an allow-list endpoint rather than free text), note that constraint directly under the hook/example.

Do this for backend changes only — skip it for internal refactors, docs, or infra changes with no frontend-visible effect.
