---
name: code-review-gate
description: 執行嚴格的程式碼審查與測試驗證關卡。用於功能開發完成後、合併或提交前的逐行品質審查、型別/Lint 檢查與單元測試全綠驗證。
---

# Code Review Gate (`/code-review-gate`)

本技能定義了專案在合併分支或提交變更前必須通過的品質審查關卡（Quality Gate）。所有 Verifier 或執行審查的 Agent 必須嚴格逐項檢查，**任何一項未通過即視為 REJECT，必須打回給開發者（Worker）修正直到全綠為止**。

---

## 1. 自動化測試與靜態分析關卡 (Automated Checks)

執行審查時，必須先在實際環境中跑過以下檢查並確認 Exit Code 為 `0`：

- [ ] **單元測試與整合測試全綠**：執行專案測試指令，確保所有既有測試與新增測試 100% 通過，無跳過（skipped）或隱藏失敗。
- [ ] **測試覆蓋率與邊界條件**：新功能或修改的邏輯必須包含對應的單元測試，涵蓋：
  - 正常路徑（Happy Path）
  - 邊界數值與空值處理（Edge Cases / Null / Empty）
  - 錯誤與例外處理（Error / Exception Path）
- [ ] **Linter / Formatter / 型別檢查**：確保無編譯錯誤、型別錯誤（Type Errors）、未使用的變數或違反專案排版規範之處。

---

## 2. 逐行程式碼審查清單 (Line-by-Line Review Checklist)

請執行 `git diff` 逐行檢視所有變更檔案，核對以下項目：

### A. 正確性與架構設計 (Correctness & Architecture)
- [ ] **需求符合度**：實作是否完全符合需求規格？有無遺漏或過度設計（Over-engineering）？
- [ ] **單一職責與模組化**：函式與類別是否保持精簡、職責單一？是否遵循既有專案目錄結構與設計模式？
- [ ] **狀態與併發安全**：非同步操作（Async/Await）、競態條件（Race Conditions）、資源釋放（File handles / Connections）是否正確處理？

### B. 安全性防護 (Security)
- [ ] **零機密外洩**：程式碼、設定檔與測試中嚴禁硬編碼（Hardcode）API Keys、Token、密碼或私鑰。
- [ ] **輸入驗證與注入防護**：外部輸入是否有驗證與消毒（Sanitization）？是否防範 SQL Injection、XSS、Command Injection、Path Traversal 等風險？
- [ ] **相依套件安全**：若新增第三方套件，是否為必要且可信賴的來源？

### C. 可讀性與維護性 (Readability & Maintainability)
- [ ] **命名清晰度**：變數、函式、類別與檔案命名是否語意明確且符合語言慣例？
- [ ] **錯誤處理**：嚴禁吞掉例外（Empty catch blocks）。錯誤訊息應具備足夠上下文以便除錯。
- [ ] **無殘留除錯碼**：移除所有臨時的 `console.log`、`print`、註解掉的廢棄程式碼（Dead code）或未解決的 `TODO`/`FIXME`。

---

## 3. 審查結果判定與回饋格式 (Verdict Protocol)

審查完成後，Verifier 必須輸出明確的結構化報告：

```markdown
### [Code Review Gate Report]
- **判定結果**: ✅ PASS | ❌ REJECT
- **測試結果**: [通過數 / 總數] (附上執行指令與簡要輸出)
- **審查發現 (Findings)**:
  1. `[檔案路徑:行號]` - [嚴重度: Critical / Major / Minor] - 問題描述與具體修正建議
- **下一步行動**: (若 REJECT，明確指派對應 Worker 修正指定項目；若 PASS，准許進入合併與 `/git-release-flow` 流程)
```
