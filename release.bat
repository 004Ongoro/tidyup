@echo off
set /p VERSION="Enter new version (e.g., 0.7.0): "

echo Updating version in cmd/root.go...
:: This uses PowerShell to do a find-and-replace in the root.go file
powershell -Command "(gc cmd/root.go) -replace 'Version = \".*\"', 'Version = \"%VERSION%\"' | Out-File -encoding utf8 cmd/root.go"

echo Staging changes...
git add .

echo Committing version bump to v%VERSION%...
git commit -m "chore: bump version to v%VERSION%"

echo Creating tag v%VERSION%...
git tag -a v%VERSION% -m "Release v%VERSION%"

echo Pushing to GitHub...
:: Pushing both the branch and the tag separately to avoid RPC errors
git push origin main
git push origin v%VERSION%

echo Done! Check your GitHub Actions tab for the build status.
pause