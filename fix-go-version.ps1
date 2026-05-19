# 修复Go版本并提交
cd D:\golang\code\admin

Write-Host "1. 修改go.mod版本..." -ForegroundColor Green
(Get-Content go.mod) -replace 'go 1.25.0', 'go 1.24.0' | Set-Content go.mod

Write-Host "2. 运行go mod tidy..." -ForegroundColor Green
go mod tidy

Write-Host "3. 验证编译..." -ForegroundColor Green
go build ./cmd/admin
if ($LASTEXITCODE -ne 0) {
    Write-Host "编译失败！" -ForegroundColor Red
    exit 1
}

Write-Host "4. 检查go.mod版本..." -ForegroundColor Green
Get-Content go.mod | Select-Object -First 3

Write-Host "5. 提交到Git..." -ForegroundColor Green
git add -A
git status --short
git commit -m "fix: 降级Go版本到1.24.0修复CI(1.25不存在)"

Write-Host "6. 推送到GitHub..." -ForegroundColor Green
git push

Write-Host "✅ 完成！" -ForegroundColor Green
pause
