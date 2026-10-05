#!/bin/bash
set -euo pipefail
here="$(cd "$(dirname "$0")" && pwd)"
exec env ANIMALSDESKTOP_MOTIONS=legacy "$here/AnimalsDesktop.app/Contents/MacOS/AnimalsDesktop"
