@echo off
@chcp 65001 >nul
git add .
git commit -m "456546mit"
git pull origin main
git push origin main

:end
echo 所有操作完成！
pause