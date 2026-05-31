The stabilization release is here. After a multi-month pause and three betas, the radar comes back with the **Radiant Wilds Protocol18 cleanup finished**, **Mists detection restored**, **multi-interface capture** that survives ExitLag and VPN toggles, **LAN access from your phone**, a **coherent logging system** with one-click pcap recording, and an **848-test safety net** (591 frontend + 257 Go) that pinned every known bug before fixing it.

Closes #24, #25, #29, #30, #32, #36, #52, #57, #65, #69, #80, #81, #83, #90, #91, #92, #98 (and more).

---

## 🎯 What you'll notice first

- 🌫️ **Mists are visible again** - portals (Solo, Duo) with rarity, feu follets before the portal opens, wisp cages inside the zone.
- 🛡️ **ExitLag-safe capture** - the radar listens on every selected interface in parallel, so VPN starts, ExitLag activations, and WiFi-to-Ethernet handoffs no longer silence it. ExitLag users on the default WFP redirection should switch to **NDIS (Legacy)** in ExitLag's advanced settings (the README screenshots the exact path).
- 📱 **Open the radar from your phone** - load `http://<your-pc-ip>:5001` over the LAN, the WebSocket connects without configuration. Startup banner now prints both URLs side by side.
- 🚨 **Alerts work everywhere** - Pulsating Border and Screen Flash now render on the radar canvas itself, so they show up in **Picture-in-Picture mode** and on the floating overlay, not just the full DOM. Pulse is also wider, more blurred, and zone-aware (fires on faction-flagged players in Black Zones).
- 🎯 **Pick your network interfaces from the browser** - **Settings → Network** lets you check / uncheck capture interfaces live; the radar no longer keys off a hardcoded IP (`ip.txt` migrates once to a `network.json` with stable interface identifiers).
- 🎨 **Cleaner radar (optional)** - **Resource Color Badges** is a new toggle (off by default) for players who want a tier-first view without the game icons. When on, harvestables render as colored squares with a `T<tier>+<enchant>` overlay. Living variants get a gold border so they stay distinct from static nodes.
- 🏰 **Every Avalonian / Roads dungeon family back on the radar** - Solo and Group alike were silently filtered out, including `T6_MORGANA`, `T6_KEEPER`, `T6_UNDEAD`, `T5_PORTAL_ROYAL_SOLO`, `T6_PORTAL`. Now visible across all enchant tiers.
- 🐟 **Fishing pools no longer phantom** - empty-string spawn type is now accepted, so the 3-of-5 pools that were silently dropped per scenario are back. `settingFishing` toggles instantly.
- ⚖️ **Tier and resource detection fixed across the board** - living mob tiers match the in-game tooltip (validated on 6,469 pcap events with zero outliers), static + DEAD harvestables route through the correct filter at render time, and the `-1` / `0xFFFF` mobileTypeId sentinel correctly routes to the static path.

---

## 🌐 Protocol18 stabilization (Radiant Wilds aftermath)

**v2.1.1** ported the wire parser to Protocol18 so the radar would speak to the Albion server again after the Radiant Wilds patch broke packet parsing (#51). **v2.2.0 finishes the cleanup that follows any protocol revision**:

- **Photon hashtable marshal failure** (#78, closes #80). `Parameters[103]` shifted from a scalar to a hashtable in Protocol18; the deserializer returned `map[interface{}]interface{}`, which `encoding/json` refuses, so every Join response was silently dropped at marshal time. A new `photon.Hashtable` named type with a custom `MarshalJSON` fixes it. Mists instance maps now render.
- **Black Zone parameter went broken in the same revision** (#87, closes #57). Direct `Parameters[103]` parse silently returned the wrong value, leaving `map.isBZ` stuck on `false`. Now derived from `zonesDatabase.getPvpType(mapId)`.
- **Mist clones inheriting the wrong PvP type** (#103, closes #90). A BZ Mist now correctly fires hostile alerts; a Yellow Royal Mist no longer trips the alarm.
- **TypeID OFFSET drift exposed across all critters** (#93, closes #92). OFFSET=16 confirmed across 6,469 pcap NewMob events plus 5,889 session-log events, zero outliers; the legacy `t-1` shift was a compensation that masked the drift.
- **452 stale event-code values synced to upstream** (#70). The positional enum from `StatisticsAnalysisTool` had 61 new names, 15 dropped, and 452 numeric mismatches accumulated since the last sync. Now in a JS source plus a generated Go mirror, refreshable via `make refresh-codes`.
- **Router contract pinned by pcap fixtures** (#64). A table-driven test asserts that `Move (22)`, `JoinFinished (2)`, and `ChangeCluster (41)` keep their parameter shapes. The next protocol revision catches us early.

---

## ✨ New Features

- **Multi-interface capture** (#94) - dynamic interface selection persisted in `network.json` (replaces `ip.txt` with stable `{name, description}` identifiers and migrates the legacy file once on first boot). Pick which interfaces to listen on from **Settings → Network**. Categorization with badges so you can spot Ethernet, WiFi, VPN, and Virtual interfaces at a glance; ExitLag's NDIS filter is also detected and ranked low so you know which one to avoid.
- **ExitLag support** (#99) - documented NDIS Legacy workaround in the README, multi-interface listens covers cases A/B/C of LWF positioning.
- **In-process pcap recording** (#107) - one toggle in **Settings → Logging** records raw network traffic to `logs/captures/capture_<timestamp>_<iface>.pcap`. No more external tcpdump for parser debugging.
- **LAN access** (#88) - dynamic WebSocket URL built from `window.location`, mobile responsive baseline (375x667 portrait, no horizontal scroll, settings forms collapse).
- **Picture-in-Picture alerts** (#85) - both Pulsating Border and Screen Flash now mirror on the radar UI canvas, so users in PiP mode see threats too.
- **Mists detection** (#78) - portals via NewMob, feu follets via the same path, wisp cages via NewCagedObject. Filters by Solo / Duo plus rarity E0-E4.
- **Resource Color Badges and Icon Size slider** (#105) - toggle in Display settings (off by default) for players who want a tier-first view without the game icons: harvestables render as `T<tier>+<enchant>` colored squares per family (Fiber green, Hide tan, Wood brown, Ore blue, Rock purple), gold border on living variants. The Icon Size slider (0.5x-2.0x) scales markers and circles for dense screens.

---

## 🐛 Bug Fixes

- **All Avalonian / Roads dungeons restored** (#78). Enchant source moved from `Parameters[6]` (a dungeon type id outside the 0-4 range, values 2 / 37-39 / 229 / 276 / 310 / 327 / 344) to `Parameters[8]` (the real 0-4 enchant). Five families that were silently filtered out come back: `T6_MORGANA`, `T6_KEEPER`, `T6_UNDEAD`, `T5_PORTAL_ROYAL_SOLO`, `T6_PORTAL` (Solo and Group alike).
- **Phantom fishing spots filtered** (#73, #85, closes #25). Spawns no longer dropped on empty-string type. `settingFishing` gate moved to render-time.
- **Living/DEAD/Static filter at runtime** (#82, closes #32, #30). Render-time gate replaces spawn-time drop, so toggling settings affects the display instantly. DEAD carcasses correctly routed through the Living filter.
- **Settings coherence audit** (#85, closes #81, #65). Orphan keys aligned (`settingShowFish` -> `settingFishing`, etc.), Pulsating Border made zone-aware, Screen Flash mirrored on the canvas, Logger default flipped to `enabled: false`.
- **Shutdown reliability** (#63). `pcap.BlockForever` swapped for a 500 ms timeout so an idle close unblocks instead of hanging on a goroutine; TUI dashboard panic on first `LogMsg` before `WindowSizeMsg` guarded.
- **LAN candidate ranking** (#106). Hyper-V and vEthernet adapters now ranked last so the announced LAN URL is a real one.
- **Living harvestables small fixes** (#71, #74). `mobileTypeId === -1` treated as static; re-gate reads stored `mobileTypeId` instead of hardcoded `isLiving=false`.
- **Chest persistence** (#72, #75). `Parameters[3]` type guard, rarity persisted on the entity (drawing-layer color resolution stays on the v2.3 backlog).

---

## ⚡ Performance and Stability

- **848 tests across the stack** (#68, #93). 591 frontend tests in 22 Vitest suites + 257 Go tests across `internal/photon`, `internal/capture`, `internal/server`, `internal/logger`, `internal/ui`, `cmd/radar`, plus the `anonymize-pcap` and `photon-dump` tooling. Discipline: `@verified`, `@characterization`, `test.fails`. Fixture corpus pcap-derived and PII-scrubbed.
- **Production binary trimmed by 347 KB** (#86). Test files renamed `_*.test.js` so Go embed excludes them. CI guard rejects unprefixed `*.test.js`; `embed_prod_test.go` walks the embed FS to make sure nothing slips through.
- **Logging coherence** (#107). Backend and frontend logs no longer mix: `logs/sessions/`, `logs/debug/`, `logs/errors/` always-on, `logs/captures/`. Configuration in `network.json`, no localStorage race at boot.
- **Test fixture tools** (#68). `tools/anonymize-pcap` (with `--scrub-string` for the local player name), `tools/photon-dump` (extracts WS-level fixtures), `tools/gen-eventcodes`, `tools/offset-validate`.

---

## 🛠️ Under the Hood

<details>
<summary>Click to expand technical details</summary>

### Multi-interface architecture
`internal/capture/manager.go` owns the active capturer set, diffs against the target on `Reconfigure`, opens new handles before closing old ones so the radar never has zero capturers during a swap. `pcap.OpenLive` failures isolate per-interface; the others keep running. Goroutines drain via wait-group with timeout before handles close (libpcap is not safe to close while `Read` is in flight).

### Categorization regex order
First match wins: `virtual` (Hyper-V, vEthernet, docker, lo) -> `exitlag` -> `vpn` (wireguard, tap, tun, ppp) -> `wifi` -> `ethernet` -> `other`. Order keeps "Wi-Fi Direct Virtual Adapter" out of the wifi bucket and a VPN over WiFi tagged as VPN, not WiFi. ExitLag's NDIS filter is detected (description `ExitLag LightWeight Filter`) and ranked second-to-last so it does not silently land at the top of the auto-pick list.

### Logging routing
Source-based, not severity-based.

| Source | Level | sessions/ | debug/ | errors/ |
|---|---|:-:|:-:|:-:|
| Backend Go | DEBUG/INFO/WARN | gated | - | - |
| Backend Go | ERROR/CRITICAL | gated | - | always-on |
| Frontend | DEBUG/INFO/WARN | - | gated | - |
| Frontend | ERROR/CRITICAL | - | gated | gated |

### In-process pcap recording
`Capturer.StartRecording(dir)` writes a `capture_<timestamp>_<sanitized-iface>.pcap` per active interface. `pcapgo.Writer` keeps frame metadata so the output is replayable through `pcap.OpenOffline`. `Manager.StartRecording` propagates and persists the flag so future capturers added by `Reconfigure` start recording too.

### TypeID OFFSET=16
Cross-validated 6,469 pcap NewMob events plus 5,889 session-log events at OFFSET=16, zero outliers. The previous OFFSET=15 was a never-HP-verified deduction; the legacy `t-1` shift compensated the drift on alive non-DYNAMIC critters and exposed it on DEAD/DYNAMIC. `tools/offset-validate` re-anchors at every DB refresh.

### Test discipline
- `@verified YYYY-MM-DD: <reason>` for tests that pass and match an external reference.
- `@characterization YYYY-MM-DD: current code does X` for observed behavior under directional uncertainty.
- `test.fails(...)` for bugs with a known correct value (CI green while broken, red when fixed).

### Frontend test fixtures
Pcap-derived fixtures live in `web/scripts/__fixtures__/ws/<handler>/<scenario>.json`. Synthetic fixtures allowed only for scenarios not observable in the corpus (stale cleanup with `Date.now()` offset, settings injection). Real game data backs every test that touches the database layer (`installRealDatabasesOnWindow()`); mocked DBs hide the class of bugs where the mock lies in sync with a wrong assertion.

</details>

---

## 📦 Migration notes

<details>
<summary>What changes if you are upgrading from 2.1.x</summary>

- `ip.txt` is replaced by `network.json` with `{name, description}` interface identifiers. Migration runs once on first boot; the legacy file is deleted afterwards.
- `localStorage.settingServerLogsEnabled` is overwritten by the value from `network.json` on the first settings page load. The toggle state carries over without user action.
- `/api/settings/server-logs` is replaced by `/api/settings/logging` (GET and POST). Old single-toggle clients break (none in the wild).
- `protocol16.go` is gone. The deserializer is now a cluster: `deserializer.go`, `packet.go`, `events.go`, `readers.go`, `types.go`, `typecodes.go`.

</details>

---

## ⚠️ Known Limitations

- **Player live positions** stay encrypted (XOR with a KeySync XorCode wrapped by Photon AES). Out of scope without a Photon MITM proxy.
- **Some Black Zone map tiles** missing for zone IDs 4000+. Workaround: disable map background in **Settings → Map**.
- **Event 46 unreliability**: `HarvestableChangeState` can skip size values or fire late depending on server batching. The radar reflects what the wire delivers.

---

## 📚 Documentation

Curated in [docs/](https://github.com/Nouuu/Albion-Online-OpenRadar/tree/main/docs). Highlights:

- [docs/releases/RELEASE_2.2.0.md](https://github.com/Nouuu/Albion-Online-OpenRadar/blob/main/docs/releases/RELEASE_2.2.0.md) - the long-form per-PR breakdown
- [docs/technical/MISTS_DETECTION.md](https://github.com/Nouuu/Albion-Online-OpenRadar/blob/main/docs/technical/MISTS_DETECTION.md)
- [docs/technical/CAPTURE_INTERFACES.md](https://github.com/Nouuu/Albion-Online-OpenRadar/blob/main/docs/technical/CAPTURE_INTERFACES.md)
- [docs/technical/LOGGING.md](https://github.com/Nouuu/Albion-Online-OpenRadar/blob/main/docs/technical/LOGGING.md)
- [docs/dev/DEV_GUIDE.md](https://github.com/Nouuu/Albion-Online-OpenRadar/blob/main/docs/dev/DEV_GUIDE.md)

---

## 🙏 Thanks

The breadth of bugs caught in this release is on you. The protocol surface is too wide for one person to test it all, and every issue you opened pointed at something real.

Special thanks to the bug reporters whose tickets shaped the v2.2 cycle:

- @TakeshiKatayama (LAN access #83, sound and flash alert #65, mists report #66)
- @furtulkerim55-a11y (Fiber tier #95, Mist in BZ #90)
- @Saltra2336 (mists detection #69)
- @gregory5993 (Wood living tier #109)
- @edgarsamar090 (T8 Mist living nodes #101)
- @djfaizp (deatheye offset cross-check #76)
- @christiangemesi, @sbibannedaku-tech, @diksco, @Renatoleall, @Sstier134 (Radiant Wilds breakage reports #104, #84, #56, #50, #49)
- @Amebatv91, @aleistershadow4-create, @shindaniel0128-cyber, @tc155426 (earlier cycle reports that informed the stabilization plan)

And to the [Discussions](https://github.com/Nouuu/Albion-Online-OpenRadar/discussions) crew, who turned a one-line "doesn't work" into a 30-message thread that pinned the real problem more than once:

- @S1NGLE-S1 (post-04-14 breakage thread #60, 8 messages of triage)
- @Neiipofer (ExitLag #89 and Mists #102; both threads fed directly into PR #99 and #78)
- @srn3dcom (network adapter switching #96 and icon size #97; both shipped as #94 and #105)
- @declaranet-dev, @0paii, @Gralios, @mexicano1806-oss, @danilka3113-lab (early ideas, language barriers crossed, edge cases reported)

And to everyone who left a kind word on [#44](https://github.com/Nouuu/Albion-Online-OpenRadar/issues/44) during the pause. @pmbstyle, @Charmathan and the others: those messages mattered. They turned a maintenance side-quest back into a project worth shipping.

If your handle is missing here it is on me, not on the value of what you reported. Open a comment on the release or the discussions board and I will edit you in.

---

### Verification

```bash
sha256sum -c checksums-sha256.txt
```

### Requirements

**Windows:** Windows 10/11 (64-bit), [Npcap 1.87+](https://npcap.com/)

**Linux:** libpcap (`apt install libpcap0.8`)

---

**Full Changelog**: https://github.com/Nouuu/Albion-Online-OpenRadar/compare/2.1.1...2.2.0
