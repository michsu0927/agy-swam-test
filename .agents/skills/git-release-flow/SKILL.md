---
name: git-release-flow
description: 標準化的 Git 分支管理、Conventional Commits 提交規範、推送與 Pull Request 建立流程。用於審查通過後的版本控制與發布作業。
---

# Git Release Flow (`/git-release-flow`)

本技能規範專案從分支開發、合併、撰寫 Commit Message 到推送遠端與建立 Pull Request (PR) 的標準作業流程。

---

## 1. 分支管理規範 (Branching Strategy)

- **主分支 (`main`)**：永遠保持穩定、可部署且測試全綠的狀態。嚴禁直接在 `main` 上進行未經審查的功能開發。
- **功能/任務分支**：一律從最新的 `main` 切出獨立分支，命名格式如下：
  - `feat/<feature-name>`：新功能開發
  - `fix/<issue-name>`：錯誤修復
  - `refactor/<scope>`：程式碼重構
  - `chore/<task-name>`：建置工具、設定或文件維護
  - `test/<scope>`：新增或調整測試

---

## 2. 提交前預檢 (Pre-Commit Verification)

在執行 `git commit` 之前，必須確認：

1. **已通過 `/code-review-gate`**：所有單元測試與靜態檢查皆為綠燈（PASS）。
2. **工作區乾淨度檢查**：執行 `git status` 與 `git diff --staged`，確認：
   - 沒有誤加入臨時檔案、編譯產物（Build artifacts）、環境變數檔（`.env`）或 Worktree 暫存目錄。
   - 每個 Commit 保持原子性（Atomic Commit），只包含單一邏輯變更。

---

## 3. Conventional Commits 規範

Commit Message 必須嚴格遵循 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```text
<type>(<optional-scope>): <subject>

<optional-body>

<optional-footer>
```

### 常用 `<type>` 定義：
- `feat`: 新增功能（Feature）
- `fix`: 修復錯誤（Bug fix）
- `docs`: 僅文件變更（Documentation）
- `style`: 不影響程式邏輯的格式調整（空白、排版、分號等）
- `refactor`: 既非新增功能也非修復錯誤的程式碼重構
- `perf`: 提升效能的變更
- `test`: 新增或修正測試案例
- `chore`: 建置流程、輔助工具、設定檔或相依套件的變更

### 撰寫原則：
- `<subject>` 使用祈使句、小寫開頭、結尾不加句號，長度控制在 50~72 字元內。
- 若有重大不相容變更（Breaking Changes），於 `<type>` 後加 `!`（如 `feat!: ...`）並於 footer 註明 `BREAKING CHANGE: <說明>`。

---

## 4. 推送與 PR 建立流程 (Push & Pull Request)

1. **同步與整合**：
   - 若使用平行 Worktree 開發，先將各 Worker 通過審查的變更乾淨合併至目標功能分支（如 `feat/xxx`），並再次執行完整測試確認無整合衝突。
2. **推送至遠端 (Push)**：
   - 確認已設定遠端倉庫（`git remote -v`）。
   - 執行 `git push -u origin <branch-name>`。
3. **建立 Pull Request (PR) 前人工確認**：
   - **重要規範**：在執行建立 PR（如 `gh pr create`）或直接合併回 `main` 之前，**必須先向使用者報告變更摘要與 Commit 紀錄，取得使用者明確同意後**才可執行。
   - PR 內容應包含：變更摘要（Summary）、測試驗證結果（Test Plan）以及相關 Issue 連結。
