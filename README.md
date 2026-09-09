# meta-extra

`meta-extra` is a fork of [mautrix/meta](https://github.com/mautrix/meta) that keeps the upstream Messenger and Instagram bridge implementation while adding protocol primitives and capabilities needed by downstream projects.

The Go module path is:

```text
go.mewis.me/meta-extra
```

## Upstream

This fork tracks `mautrix/meta` and keeps upstream history intact. Changes should stay as small and generic as practical so useful protocol improvements can still be contributed upstream.

Existing compatibility identifiers, database ownership values, and bridge behavior remain unchanged unless a change explicitly requires otherwise.

## Bridges

The repository contains two bridges:

- `mautrix-meta` for Facebook Messenger.
- `mautrix-instagram` for Instagram DMs.

Binary names are currently kept compatible with upstream.

## Documentation

Upstream setup and usage documentation remains the primary reference for bridge operation:

- [mautrix-meta documentation](https://docs.mau.fi/bridges/go/meta/index.html)
- [Bridge setup](https://docs.mau.fi/bridges/go/setup.html?bridge=meta)
- [Docker setup](https://docs.mau.fi/bridges/general/docker-setup.html?bridge=meta)
- [Authentication](https://docs.mau.fi/bridges/go/meta/authentication.html)

Fork-specific behavior should be documented in this repository when it differs from upstream.

## Features and roadmap

[ROADMAP.md](ROADMAP.md) contains the upstream feature overview. Additional `meta-extra` capabilities may extend it over time.

## License and attribution

This project is derived from `mautrix/meta`. Upstream copyright and license notices are retained.

See [LICENSE](LICENSE), [LICENSE.exceptions](LICENSE.exceptions), and [NOTICE](NOTICE). The special exceptions in `LICENSE.exceptions` originate from upstream copyright holders and are not reissued or expanded by this fork.

## Upstream discussion

For upstream bridge discussion, see the Matrix room [#meta:maunium.net](https://matrix.to/#/#meta:maunium.net).
