# Get the path to the rdb executable
$rdbPath = Join-Path $PSScriptRoot "..\bin\rdb.exe"

Write-Host "=== RDB Basic Usage Example ===" -ForegroundColor Green
Write-Host ""

# Create example repository
$repoPath = Join-Path $PWD "example-repo"
Write-Host "Repository will be created at: $repoPath" -ForegroundColor Yellow
Write-Host ""

# Change to the repository directory
Set-Location $repoPath

Write-Host "1. Initializing RDB repository..." -ForegroundColor Yellow
& $rdbPath init --layout tree --types "text,audio,texture"

Write-Host "`n2. Checking status after initialization..." -ForegroundColor Yellow
& $rdbPath status

Write-Host "`n3. Creating sample assets..." -ForegroundColor Yellow
# Create a sample text asset
$textAssetDir = "assets\1030002"
New-Item -ItemType Directory -Path $textAssetDir -Force
New-Item -ItemType Directory -Path "$textAssetDir\data" -Force

# Create sample text files
$enContent = @"
Hello, world!
This is a sample English text file.
"@
$enContent | Out-File -FilePath "$textAssetDir\data\en.txt" -Encoding UTF8

$frContent = @"
Bonjour, monde !
Ceci est un exemple de fichier texte français.
"@
$frContent | Out-File -FilePath "$textAssetDir\data\fr.txt" -Encoding UTF8

# Create metadata file
$metadata = @{
    type = "text"
    id = 1030002
    name = "DialogLine_Intro"
    tags = @("localization", "ui")
    version = 1
    attributes = @{
        language = "en-US"
        platform = @("pc", "xbox")
    }
} | ConvertTo-Json -Depth 3

$metadata | Out-File -FilePath "$textAssetDir\meta.json" -Encoding UTF8

Write-Host "`n4. Adding assets to repository..." -ForegroundColor Yellow
& $rdbPath add "assets\1030002" --type text --id 1030002 --name "DialogLine_Intro"

Write-Host "`n5. Checking status after adding assets..." -ForegroundColor Yellow
& $rdbPath status

Write-Host "`n6. Committing changes..." -ForegroundColor Yellow
& $rdbPath commit -m "Add intro dialog line"

Write-Host "`n7. Viewing commit history..." -ForegroundColor Yellow
& $rdbPath log

Write-Host "`n8. Building package..." -ForegroundColor Yellow
& $rdbPath build

Write-Host "`n9. Checking final status..." -ForegroundColor Yellow
& $rdbPath status

Write-Host "`n=== Example completed! ===" -ForegroundColor Green
Write-Host "Repository created at: $repoPath" -ForegroundColor Cyan
Write-Host "Package created at: dist/" -ForegroundColor Cyan

# Return to original directory
Set-Location .. 