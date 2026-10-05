#!/usr/bin/env python3
"""Verify release-sensitive fields on the public Pages HTML."""

from __future__ import annotations

import re
import sys
from pathlib import Path

from catalog_runtime import runtime_variant_ids


ROOT = Path(__file__).resolve().parents[1]
INDEX = ROOT / "docs" / "index.html"
WINDOWS_AMD64_ASSET = "AnimalsDesktop-windows-amd64.zip"
WINDOWS_AMD64_NO_NETWORK_ASSET = "AnimalsDesktop-windows-amd64-no-network.zip"
WINDOWS_386_ASSET = "AnimalsDesktop-windows-386.zip"
MAC_ARM64_ASSET = "AnimalsDesktop-macos-arm64.zip"
MAC_AMD64_ASSET = "AnimalsDesktop-macos-amd64.zip"
CHECKSUM_ASSET = "SHA256SUMS.txt"
CATALOG = ROOT / "internal" / "catalog" / "catalog.go"
EXPECTED_RELEASE = "v2.0.0"
EXPECTED_CURRENT_ANIMALS = 60
EXPECTED_BUILD_DATE = "2026-10-05"
EXPECTED_UPCOMING: list[str] = []


def fail(message: str) -> None:
    print(f"verify_page_release: {message}", file=sys.stderr)
    raise SystemExit(1)


def one(pattern: str, html: str, label: str) -> str:
    matches = re.findall(pattern, html, flags=re.S)
    if not matches:
        fail(f"missing {label}")
    if len(matches) > 1:
        fail(f"expected one {label}, found {len(matches)}")
    match = matches[0]
    if isinstance(match, tuple):
        return next(part for part in match if part)
    return match


def release_tag(asset: str, html: str, label: str) -> str:
    matches = re.findall(
        rf'releases/download/(v[0-9][^"/]*)/{re.escape(asset)}',
        html,
        flags=re.S,
    )
    if not matches:
        fail(f"missing {label}")
    tags = sorted(set(matches))
    if len(tags) > 1:
        fail(f"{label} points to multiple versions: {', '.join(tags)}")
    return tags[0]


def testid_href(testid: str, html: str) -> str:
    return one(
        rf'<a(?=[^>]*data-testid="{re.escape(testid)}")'
        rf'(?=[^>]*href="([^"]+)")[^>]*>',
        html,
        f"{testid} link",
    )


def current_page_variant_ids(html: str) -> list[str]:
    block = one(
        r"const animalCatalog = \[(.*?)\n\s*\];",
        html,
        "JavaScript animal catalog",
    )
    ids = re.findall(r'\["([a-z0-9_]+)",\s*"[^"]+"\]', block)
    if not ids:
        fail("JavaScript animal catalog has no IDs")
    return ids


def english_animal_labels(html: str) -> list[str]:
    block = one(
        r"const animalTextEn = \[(.*?)\n\s*\];",
        html,
        "English animal labels",
    )
    labels = re.findall(r'"([^"]+)"', block)
    if not labels:
        fail("English animal label list is empty")
    return labels


def upcoming_page_ids(html: str) -> list[str]:
    matches = re.findall(
        r'<div class="future-grid">(.*?)</div>\s*</section>',
        html,
        flags=re.S,
    )
    if not matches:
        return []
    if len(matches) > 1:
        fail(f"expected one upcoming animal grid, found {len(matches)}")
    block = matches[0]
    ids = re.findall(r'assets/upcoming-silhouettes/([a-z0-9_]+)\.png', block)
    return ids


def verify_asset_refs(html: str, page_ids: list[str]) -> None:
    refs = set(re.findall(r'(?:src|href)="(assets/[^"]+)"', html))
    refs.update(re.findall(r'url\(["\']?(assets/[^"\')]+)', html))
    refs.update(f"assets/animal-icons/current-{variant_id}.png" for variant_id in page_ids)
    if not refs:
        fail("no local asset references found")
    normalized = sorted({re.split(r"[?#]", ref, maxsplit=1)[0] for ref in refs})
    missing = [ref for ref in normalized if not (ROOT / "docs" / ref).exists()]
    if missing:
        fail(f"missing local page assets: {missing}")


def verify_current_icons(page_ids: list[str], icon_dir: Path) -> None:
    expected = {f"current-{variant_id}.png" for variant_id in page_ids}
    actual = {path.name for path in icon_dir.glob("current-*.png")}
    missing = sorted(expected - actual)
    stale = sorted(actual - expected)
    if missing or stale:
        fail(
            "current animal icons differ from runtime roster; "
            f"missing={missing}, stale={stale}"
        )


def main() -> None:
    html = INDEX.read_text(encoding="utf-8")

    if "releases/download/v0.2.0" in html:
        fail("blocked public Windows v0.2.0 release link remains")

    windows_tag = release_tag(WINDOWS_AMD64_ASSET, html, "Windows amd64 download version tag")
    windows_no_network_tag = release_tag(
        WINDOWS_AMD64_NO_NETWORK_ASSET,
        html,
        "Windows amd64 no-network download version tag",
    )
    windows_386_tag = release_tag(WINDOWS_386_ASSET, html, "Windows 386 download version tag")
    checksum_tag = release_tag(CHECKSUM_ASSET, html, "SHA256SUMS download version tag")
    mac_arm64_tag = release_tag(MAC_ARM64_ASSET, html, "macOS arm64 download version tag")
    mac_amd64_tag = release_tag(MAC_AMD64_ASSET, html, "macOS amd64 download version tag")
    if windows_tag != EXPECTED_RELEASE:
        fail(f"Windows amd64 download tag {windows_tag} does not match expected {EXPECTED_RELEASE}")

    expected_testid_links = {
        "download-windows-primary": (
            f"https://github.com/UDteach/AnimalsDesktopV2/releases/download/"
            f"{EXPECTED_RELEASE}/{WINDOWS_AMD64_ASSET}"
        ),
        "download-mac-primary": (
            f"https://github.com/UDteach/AnimalsDesktopV2/releases/download/"
            f"{EXPECTED_RELEASE}/{MAC_ARM64_ASSET}"
        ),
        "download-windows-no-network": (
            f"https://github.com/UDteach/AnimalsDesktopV2/releases/download/"
            f"{EXPECTED_RELEASE}/{WINDOWS_AMD64_NO_NETWORK_ASSET}"
        ),
        "download-windows-x86": (
            f"https://github.com/UDteach/AnimalsDesktopV2/releases/download/"
            f"{EXPECTED_RELEASE}/{WINDOWS_386_ASSET}"
        ),
        "download-mac-intel": (
            f"https://github.com/UDteach/AnimalsDesktopV2/releases/download/"
            f"{EXPECTED_RELEASE}/{MAC_AMD64_ASSET}"
        ),
    }
    for testid, expected_href in expected_testid_links.items():
        actual_href = testid_href(testid, html)
        if actual_href != expected_href:
            fail(
                f"{testid} points to {actual_href}, expected {expected_href}"
            )

    windows_badge = one(
        r"<span[^>]*>\s*Windows版\s*<strong>(v[0-9][^<]*)</strong>\s*</span>",
        html,
        "visible Windows version badge",
    )
    release_badge = one(
        r"<strong data-release-version>(v[0-9][^<]*)</strong>",
        html,
        "public release version badge",
    )
    mac_badge = one(
        r"<span[^>]*>\s*Mac版\s*<strong>(v[0-9][^<]*)</strong>\s*</span>",
        html,
        "visible Mac version badge",
    )
    build_date = one(
        r"<strong data-build-date>([^<]+)</strong>",
        html,
        "visible build date",
    )

    for label, got in {
        "Windows 386 download": windows_386_tag,
        "Windows no-network download": windows_no_network_tag,
        "macOS arm64 download": mac_arm64_tag,
        "macOS amd64 download": mac_amd64_tag,
        "SHA256SUMS download": checksum_tag,
        "Windows badge": windows_badge,
        "Mac badge": mac_badge,
        "release badge": release_badge,
    }.items():
        if got != windows_tag:
            fail(f"{label} {got} does not match Windows amd64 tag {windows_tag}")
    if build_date != EXPECTED_BUILD_DATE:
        fail(f"build date {build_date} does not match expected {EXPECTED_BUILD_DATE}")

    try:
        runtime_ids = runtime_variant_ids(CATALOG)
    except ValueError as error:
        fail(str(error))
    page_ids = current_page_variant_ids(html)
    if page_ids != runtime_ids:
        fail(f"current animal grid {page_ids} does not match runtime variants {runtime_ids}")
    if len(page_ids) != EXPECTED_CURRENT_ANIMALS:
        fail(
            f"current animal count {len(page_ids)} does not match expected "
            f"{EXPECTED_CURRENT_ANIMALS}"
        )
    english_labels = english_animal_labels(html)
    if len(english_labels) != len(page_ids):
        fail(
            f"English animal label count {len(english_labels)} does not match "
            f"current animal count {len(page_ids)}"
        )
    verify_current_icons(page_ids, ROOT / "docs" / "assets" / "animal-icons")

    upcoming_ids = upcoming_page_ids(html)
    if upcoming_ids != EXPECTED_UPCOMING:
        fail(f"upcoming animal grid {upcoming_ids} does not match expected priority {EXPECTED_UPCOMING}")

    for required in (
        '<time datetime="2026-07-27">2026-07-27</time>',
        '<time datetime="2026-07-27">July 27, 2026</time>',
        "v0.2.15 / 2026-07-04",
        "v0.2.15 / July 4, 2026",
        "固定表示とランダム表示",
        "種類で絞り込め",
        "fixed animals and random slots",
        "filtered by animal type",
        "v0.2.14 / 2026-07-03",
        "v0.2.14 / July 3, 2026",
        "ハシビロコウが別の動物として復元",
        "stable IDs",
        "groups the Windows and Mac animal pickers by type",
        "v0.2.13 / 2026-07-03",
        "v0.2.13 / July 3, 2026",
        "roster additions do not remap saved choices",
        "v0.2.12 / 2026-07-03",
        "v0.2.12 / July 3, 2026",
        "AnimalsDesktop-windows-amd64-no-network.zip",
        "Smart App Control FAQ",
        "SmartScreen reputation",
        "GlobalSignのSmart App Control事例",
        "GlobalSign Smart App Control case study",
        "v0.2.11 / 2026-07-02",
        "v0.2.11 / July 2, 2026",
        "v0.2.10 / 2026-07-02",
        "v0.2.10 / July 2, 2026",
        "v0.2.9 / 2026-07-02",
        "v0.2.9 / July 2, 2026",
        "v0.2.8 / 2026-07-01",
        "v0.2.8 / July 1, 2026",
        "v0.2.7 / 2026-07-01",
        "v0.2.7 / July 1, 2026",
        "v0.2.6 / 2026-06-30",
        "v0.2.6 / June 30, 2026",
        "v0.2.5 / 2026-06-29",
        "v0.2.5 / June 29, 2026",
        "v0.2.4 / 2026-06-28",
        "v0.2.4 / June 28, 2026",
        "v0.2.3 / 2026-06-28",
        "v0.2.3 / June 28, 2026",
        "v0.2.2 / 2026-06-27",
        "v0.2.2 / June 27, 2026",
        "v0.2.1 / 2026-06-27",
        "v0.2.1 / June 27, 2026",
        "v0.1.5 / 2026-06-21",
        "v0.1.5 / June 21, 2026",
        "v2.0.0は7姿の歩行・待機モーションを更新。残り53姿は従来素材を使用しています。",
        "v2.0.0 updates walking and idle motion for seven animals; the other 53 keep their existing sprites.",
        "公開版 v2.0.0：60種類",
        "Public v2.0.0: 60 animals",
        "Windows版 x64をダウンロード",
        "Mac版 Apple Silicon（macOS 12+）をダウンロード",
        "Download Windows x64",
        "Download Mac Apple Silicon (macOS 12+)",
        "Windows版 x86 / 32-bit",
        "Mac版 Intel / macOS 12+",
        '<h3 class="download-alternatives-title">その他のダウンロード</h3>',
        "Mac版の選び方",
        "M1以降のApple製チップ",
        "Windows x86 / 32-bit",
        "Mac Intel / macOS 12+",
        '<h3 class="download-alternatives-title">Other downloads</h3>',
        "Choose a Mac build",
        "M1 or later Apple chip",
        "フェレット（セーブルパンダ）",
        "フェレット（セーブル）",
        "フェレット（アルビノ）",
        "Sable panda ferret",
        "Sable ferret",
        "Albino ferret",
        "オカメインコ",
        "Cockatiel - normal gray",
        "v2.0.0の詳細",
        "About v2.0.0",
        "最大10匹",
        "Show up to 10 animals",
        "木、種、餌",
        "trees, seeds, and food",
    ):
        if required not in html:
            fail(f"missing version history text: {required}")

    for forbidden in ("次回版", "next-version catalog"):
        if forbidden in html:
            fail(f"stale pre-release copy remains: {forbidden}")

    verify_asset_refs(html, page_ids)

    print(
        f"Pages release links verified: Windows {windows_tag}, macOS {mac_arm64_tag}, "
        f"current animals {len(page_ids)}, upcoming silhouettes {len(upcoming_ids)}"
    )


if __name__ == "__main__":
    main()
