@echo off
:: 切换到当前脚本所在目录
cd /d "%~dp0"

:: 添加所有更改
echo 正在添加所有更改...
git add .

:: 检查是否有改动需要提交
git status | find "no changes added to commit"
if %errorlevel% == 0 (
    echo 没有需要提交的更改。
    goto end
)

:: 提交更改
set /p commitMessage="请输入提交信息: "
git commit -m "%commitMessage%"

:: 拉取最新更改（从 origin/main）
echo 正在拉取最新更改...
git pull origin main

:: 检查拉取是否成功
if %errorlevel% neq 0 (
    echo 拉取过程中出现问题，请检查错误信息并手动解决。
    pause
    exit /b %errorlevel%
)

:: 推送更改到远程仓库
echo 正在推送更改...
git push origin main

:end
echo 所有操作完成！
pause