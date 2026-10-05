param(
    [string]$OutDir = "docs/art-source/external-ai-trials/google-flow-20260625/hamster-cream-color-pilot",
    [string]$Animal = "golden Syrian hamster recolored to cream coat, same accepted AnimalsDesktop sprite style",
    [string]$Family = "hamster",
    [string]$Sequence = "set00",
    [string]$FrameList = "0,4,12,26,32,40,52,56",
    [string]$AnchorFrames = "0,4,12"
)

$ErrorActionPreference = "Stop"

$repoRoot = Split-Path -Parent $PSScriptRoot
Set-Location $repoRoot

$anchorPaths = @()
foreach ($frame in $AnchorFrames.Split(",")) {
    $frameNumber = [int]$frame.Trim()
    $anchorPaths += "docs/art-source/$Family/motion-source/accepted-frames/$Sequence/frame-$($frameNumber.ToString("00")).png"
}

go run ./cmd/flowgridtemplate `
    -out-dir $OutDir `
    -animal $Animal `
    -sequence $Sequence `
    -frame-list $FrameList `
    -anchors ($anchorPaths -join ",")

Get-ChildItem -Path $OutDir -Filter "._*" -File -ErrorAction SilentlyContinue | Remove-Item -Force
