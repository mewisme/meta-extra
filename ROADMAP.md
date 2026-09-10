# Roadmap

`meta-extra` tracks capabilities by protocol surface instead of duplicating a Matrix-to-platform and platform-to-Matrix matrix for every feature.

- `[x]` means the capability is implemented in the current repository and has deterministic coverage, live coverage where appropriate, or inherited upstream coverage.
- `[ ]` means there is a real implementation gap that is still useful to pursue.
- Items under **Deferred or intentionally unsupported** are not active TODOs. They require new current-protocol evidence before implementation.

Protocol primitives and bridge-level exposure are listed separately. A low-level Messenger primitive can be supported even when the bridge does not yet map it to a Matrix feature.

## Messenger - unencrypted

### Messaging

- [x] Text send/receive
- [x] Formatting
- [x] Replies
- [x] Mentions
- [x] Images
- [x] Videos
- [x] Files
- [x] Voice messages
- [x] GIFs and stickers on receive
- [x] Story, reel and clip shares on receive
- [x] Profile/contact shares on receive
- [x] Message reactions
- [x] Message edits
- [x] Message unsend, including realtime self-recall
- [x] Message history
- [x] Read receipts
- [x] Typing notifications
- [x] Incoming static locations
- [ ] Poll messages in the normal bridge message pipeline
- [ ] Product shares

### Threads and participants

- [x] Create/open conversations
- [x] Add participants
- [x] Remove participants
- [x] Incoming participant leave events
- [ ] Explicit outbound leave operation
- [x] Thread name changes
- [x] Thread avatar changes
- [x] Admin state on receive
- [ ] Complete bridge-level power/admin mutation mapping
- [ ] Bridge-level per-chat nickname mapping

### meta-extra Messenger protocol extensions

These are typed protocol capabilities available to downstream users even when the upstream bridge does not expose a matching Matrix action.

- [x] Per-thread participant nickname mutation
- [x] Custom thread emoji mutation
- [x] Admin approval mode mutation
- [x] Archive/unarchive
- [x] Theme catalog and theme mutation
- [x] Call notification mute/unmute
- [x] Message pin/unpin
- [x] Pinned-message state/query
- [x] Poll state/details query
- [x] Per-thread message search with cursor pagination
- [x] Typed realtime thread system events
  - [x] Nickname changes
  - [x] Emoji changes
  - [x] Approval-mode changes
  - [x] Theme changes
  - [x] Participant add/remove
  - [x] Participant admin changes
  - [x] Pin/unpin changes
  - [x] Poll changes
- [x] Contact-card sharing
- [x] Messenger restrict/unrestrict task
- [x] Messenger message-block/unblock task
- [x] Restriction state decoding
- [x] Message/full-block state decoding
- [x] Current self-recall `cleanUpOnRecall` decoding

### Remaining useful Messenger work

- [ ] Complete high-level poll create/vote/message integration around the existing typed poll primitives
- [ ] Expose per-chat nickname cleanly through the bridge
- [ ] Add an explicit outbound leave-group operation if the current protocol path is revalidated
- [ ] Complete bridge-level participant admin/power mutation support
- [ ] Add product-share decoding only if current payloads can be covered with stable typed models

### Deferred or intentionally unsupported

- **Presence / active status:** current Messenger Web uses PresenceUnified/Bladerunner. Legacy Lightspeed presence rows and app-state tasks are not authoritative enough for a stable public API.
- **Regular Messenger outbound static location:** current Web has no maintained regular Lightspeed send path. Historical task `263` was server-accepted but produced no message.
- **Regular Messenger outbound live location:** current Web no longer exposes maintained start/update/end write operations. Existing inbound decoders remain for compatibility.
- **Outbound message power-ups/effects:** current Web send paths hardcode `MessagePowerUp.NONE`; accepting an old field is not treated as working support.
- **Messenger Notes:** owned by Facebook Web/GraphQL rather than the Messenger transport and intentionally left to downstream Facebook-layer code.
- **Sound bites:** legacy attachment compatibility remains, but no active current catalog/send path was found.
- **Magic words:** legacy schema compatibility remains, but no active current configuration/mutation path was found.
- **Dedicated shared-media/files/links browser API:** current attachment-range behavior does not add content semantics beyond typed history/backfill and is intentionally not duplicated.

## Messenger - encrypted

### Messaging

- [x] Text
- [x] Formatting
- [x] Replies
- [x] Mentions
- [x] Images
- [x] Videos
- [x] Files
- [x] Voice messages
- [x] Stickers
- [x] GIFs
- [x] Static location payloads with coordinate validation
- [x] Message reactions
- [x] Message edits
- [x] Message unsend
- [x] Read receipts
- [x] Typing notifications
- [ ] Polls

### Sync and thread management

- [x] Incoming participant add/remove/leave
- [x] Incoming thread metadata changes
- [x] Initial thread metadata
- [x] Admin state on receive
- [ ] Write chat backup
- [ ] Read message history from chat backup
- [ ] Outbound leave operation
- [ ] Outbound thread name changes
- [ ] Outbound thread avatar changes
- [ ] Per-chat nickname support
- [ ] Complete bridge-level power/admin mutation mapping

## Instagram

### Messaging

- [x] Text
- [x] Formatting
- [x] Replies
- [x] Mentions
- [x] Images
- [x] Videos
- [x] Voice messages
- [x] GIFs and stickers on receive
- [x] Story, reel and clip shares on receive
- [x] Message reactions
- [x] Message edits
- [x] Message unsend
- [x] Message history
- [x] Read receipts
- [x] Typing notifications

### Threads and participants

- [x] Incoming participant add/remove/leave
- [x] Incoming thread title changes
- [x] Incoming thread avatar changes
- [x] Initial thread metadata
- [x] User name/avatar metadata
- [x] Outbound thread name changes
- [x] Outbound thread avatar changes
- [ ] Outbound participant invite
- [ ] Outbound participant kick
- [ ] Outbound leave operation
- [ ] Per-chat nickname support
- [ ] Reliable admin-state mapping where the server does not return it
- [ ] Complete bridge-level power/admin mutation mapping

### Deferred

- **Presence:** do not expose a new stable presence API until the current Instagram transport provides authoritative semantics suitable for the bridge.

## Shared bridge capabilities

- [x] Multi-user support
- [x] Shared group-chat portals
- [x] Messenger encryption
- [x] Matrix encryption
- [x] Automatic portal creation at startup, on membership changes and on incoming messages
- [x] Private-chat creation by inviting the remote-user puppet
- [x] Double-puppeting / using the Matrix account for messages sent from other Messenger or Instagram clients

## Current priorities

1. Complete unencrypted Messenger poll integration using the typed poll state/query primitives already present.
2. Finish bridge-level nickname and participant-admin mappings without duplicating Messenger protocol details outside `pkg/messagix`.
3. Improve encrypted Messenger backup/history and outbound thread-management coverage.
4. Close Instagram outbound membership and per-chat nickname gaps.
5. Revisit deferred protocol surfaces only when current Messenger/Instagram clients provide new evidence that they are active and maintainable.

## Compatibility rule

Do not implement a roadmap item from historical task labels, GraphQL document IDs or old client behavior alone. Revalidate the current protocol first, keep requests and decoded state typed, add deterministic contract tests, and use reversible live verification for mutations where it is safe to do so.
