# Code Summary

這個軟體是一個程式碼摘要工具，具有以下主要特點：它能夠根據指定的程式碼檔案，使用聊天GPT模型生成程式碼的摘要報告。軟體提供了彈性的配置選項，包括摘要規則、問題提示和輸出設定等，用戶可以根據需要進行定製。此外，軟體還提供了生成摘要報告的功能，以便使用者可以方便地查看和分享生成的摘要結果。整體而言，這個軟體提供了一個自動化的方式來產生程式碼的摘要，幫助程式開發人員更快速地理解和分析程式碼。

| File | Description |
| --- | --- |
| cmd/codesum/main.go | 這個檔案是程式碼摘要工具的主要進入點。它負責解析命令列參數、讀取配置檔、執行程式碼摘要的相關操作。 |
| pkg/cfgs/cfgs.go | 這個檔案定義了程式碼摘要工具的配置結構和方法。它提供了載入和儲存配置的函式，以及處理配置檔案的相關邏輯。 |
| pkg/cfgs/chatgpt.go | 這個檔案定義了聊天GPT模型的配置結構。它包含了與聊天GPT相關的各種設定，如端點、API金鑰、存取權杖、模型和代理。 |
| pkg/cfgs/prompt.go | 這個檔案定義了摘要提示的配置結構。它包含了系統提示、程式碼摘要、摘要表格和最終摘要的相關設定。 |
| pkg/cfgs/rule.go | 這個檔案定義了摘要規則的配置結構。它包含了要包含和排除的摘要規則、輸出目錄和輸出檔案名稱的設定。 |
| pkg/sumer/question.go | 這個檔案定義了關於摘要問題的函式。它提供了根據檔案名稱和內容生成問題的功能，以及根據檔案列表生成摘要表格問題的功能。 |
| pkg/sumer/report.go | 這個檔案定義了程式碼摘要報告相關的結構和方法。它提供了儲存報告、產生 Markdown 格式報告和操作部分摘要、檔案摘要的相關功能。 |
| pkg/sumer/summary.go | 這個檔案定義了程式碼摘要相關的結構和方法。它包含了程式碼摘要、檔案摘要和問答相關的功能，用於處理程式碼的摘要和生成摘要報告。 |
| pkg/utils/helper.go | 該文件包含一組工具函數，用於幫助其他程式模組進行不同的操作。它包括更新 API 伺服器存取權杖的函數、處理 Markdown 表格摘要的函數，以及判斷特定錯誤類型的函數。 |

## cmd/codesum/main.go

# cmd/codesum/main.go

這個檔案是程式的進入點，主要執行程式的邏輯。

- 匯入相關的套件
- 定義全域變數 BuildName、BuildVersion、BuildHash、BuildTime
- 定義命令列參數的結構 args
- 初始化函式 `init()`：設定命令列參數的預設值，並解析命令列參數。
- 主程式進入點 `main()`：執行 `run()` 函式，如果有錯誤則使用 `log.Fatal()` 輸出錯誤訊息。
- 函式 `run()`：載入並合併設定檔，套用日誌設定，印出程式版本資訊，更新 API 伺服器的存取權杖，建立 Summarizer 物件，開始進行摘要處理。
- 函式 `startSummarize()`：建立並管理 Goroutine，接收系統信號以終止處理。
- 函式 `pushMessage()`：用於輸出訊息。
- 函式 `parseArgs()`：解析命令列參數，檢查參數數量，並設定 `args.codeFolder` 的值。
- 函式 `applyLog()`：套用日誌設定。

這個檔案的主要邏輯為載入設定檔、設定日誌、建立 Summarizer 物件並開始進行摘要處理，並對命令列參數進行解析和驗證。

## pkg/cfgs/cfgs.go

這個檔案實現了與設定檔相關的功能，包括載入、合併、保存和寫入設定檔。

- 匯入相關的套件
- 初始化函式 `init()`：註冊環境變數的擴展處理器。
- 定義類型 `Expander` 為 `strutil.Expander` 的別名。
- 定義結構類型 `Config`，包含設定檔中的各項設定。
- 定義嵌入資源 `defaultCfgFs`，用於存儲預設設定檔。
- 函式 `DefaultConfig()`：載入並回傳預設設定檔。
- 函式 `LoadConfig(fname string)`：從指定的檔案載入設定檔，並回傳 `Config` 物件。
- 函式 `ReadConfig(r io.Reader)`：從 `io.Reader` 載入設定檔，並回傳 `Config` 物件。
- 函式 `MergeDefault()`：將預設設定檔合併到當前的 `Config` 物件中。
- 函式 `SaveConfig(fname string)`：將 `Config` 物件保存到指定的檔案。
- 函式 `WriteConfig(w io.Writer)`：將 `Config` 物件寫入 `io.Writer`。


## pkg/cfgs/chatgpt.go

這個檔案定義了與 ChatGpt 相關的設定項目。

- 定義結構類型 `ChatGpt`，包含與 ChatGpt 相關的設定項目，如端點（EndPoint）、API 金鑰（APIKey）、存取權杖（AccessToken）、模型（Model）和代理（Proxy）等。

## pkg/cfgs/prompt.go

這個檔案定義了與提示訊息（Prompt）相關的設定項目。

- 定義結構類型 `Prompt`，包含與提示訊息相關的設定項目，如系統提示（System）、程式碼摘要提示（CodeSummary）、摘要表格提示（SummaryTable）和最終摘要提示（FinalSummary）等。

## pkg/cfgs/rule.go

這個檔案定義了摘要規則（SummaryRules）相關的設定項目。

- 定義結構類型 `SummaryRules`，包含摘要規則的設定項目，如包含條件（Include）、排除條件（Exclude）、輸出目錄（OutDir）和輸出檔名（OutFileName）等。這些設定項目用於指定在進行摘要處理時的檔案選取和輸出設定。

## pkg/sumer/question.go

這個檔案定義了摘要相關的問題（Question）函式。

- 函式 `FileSummaryQuestion(prompt strutil.Expander, fileName, fileContent string)`：根據指定的提示擴展器、檔案名稱和檔案內容，生成一個用於檔案摘要的問題。
- 函式 `SummaryTableQuestion(prompt strutil.Expander, fileList []string)`：根據指定的提示擴展器和檔案清單，生成一個用於摘要表格的問題。將檔案清單以逗號分隔後擴展到提示中。
- 函式 `FinalSummaryQuestion(prompt strutil.Expander)`：根據指定的提示擴展器，生成一個用於最終摘要的問題。回傳原始的提示內容。

## pkg/sumer/report.go

這個檔案定義了程式碼摘要報告相關的結構和方法。

- 型別 `Message`：用於表示開放AI的聊天完成訊息。
- 型別 `CodeSummary`：程式碼摘要的結構，包含部分摘要列表和最終摘要問答。
  - 方法 `NewCodeSummary() *CodeSummary`：建立並回傳一個新的程式碼摘要結構。
  - 方法 `SaveMarkdown(fname string) error`：將程式碼摘要以 Markdown 格式儲存到指定的檔案中。
  - 方法 `WriteMarkdown(w io.Writer) error`：將程式碼摘要以 Markdown 格式寫入指定的 io.Writer。
  - 方法 `WriteFileSummaryTable(w io.Writer)`：將檔案摘要表格以 Markdown 格式寫入指定的 io.Writer。
  - 方法 `GetFileSummaryTable() string`：取得檔案摘要表格的字串表示。
  - 方法 `AddPartialSummaries(systemPrompt string) *PartialSummaries`：新增部分摘要並回傳。
  - 方法 `RequestFinalSummaryMessages(question string) []Message`：根據問題回傳最終摘要的訊息列表。
  - 方法 `SetFinalSummaryQA(question, answer string) *QA`：設定最終摘要的問題與回答。
  - 方法 `FinalSummary() string`：取得最終摘要的內容。
- 型別 `PartialSummaries`：部分摘要的結構，包含系統提示、檔案摘要列表和摘要問答。
- 型別 `FileSummary`：檔案摘要的結構，包含檔案名稱和摘要問答。
- 型別 `QA`：摘要問答的結構，包含問題和回答的聊天完成訊息。
- 方法 `NewQA(question, answer string) *QA`：建立並回傳一個新的摘要問答結構。
- 方法 `SetSystemPrompt(prompt string) *Message`：設定系統提示訊息並回傳。
- 方法 `RequestFileSummaryMessages(question string) []Message`：根據問題回傳檔案摘要的訊息列表。
- 方法 `AddFileSummary(fileName, question, answer string) *FileSummary`：新增檔案摘要並回傳。
- 方法 `PopFileSummary() *FileSummary`：彈出最後一個檔案摘要並回傳。
- 方法 `GetSummary() string`：取得摘要問答的回答內容。
- 方法 `RequestSummaryMessages(question string) []Message`：根據問題

## pkg/sumer/summary.go

以下是 "pkg/sumer/summary.go" 這個程式碼檔案的摘要：

- `NewSummarizer` 函式用於建立 `Summarizer` 的實例。該函式會初始化並返回一個 `Summarizer` 物件，並設定相關的屬性和回調函式。
- `Summarize` 方法用於進行程式碼摘要。它接受一個代表程式碼資料夾的路徑，並在該資料夾中處理程式碼檔案的摘要工作。該方法會遍歷程式碼檔案並使用 GPT 模型向 OpenAI 服務發送請求，獲取程式碼的摘要結果。最終，它會生成摘要報告並保存為 Markdown 格式。
- `sendRequest` 方法用於向 OpenAI 服務發送摘要請求。它接受一個 `ChatCompletionRequest` 對象，將請求發送到 OpenAI 並接收回應。在接收回應期間，它會將回應內容逐步累積起來，最終返回完整的摘要結果。
- `getCodeFileContent` 函式用於讀取程式碼檔案的內容。它接受程式碼檔案的路徑，讀取該檔案的內容並返回檔案名稱和內容兩個字符串。
- `saveMarkdownReport` 方法用於將摘要報告保存為 Markdown 檔案。它接受程式碼資料夾的路徑和 `CodeSummary` 物件，根據配置的規則生成摘要報告的路徑，並將報告內容保存到指定的檔案中。
- `isCtxCanceled` 函式用於檢查是否取消了上下文。它接受一個上下文對象，檢查是否已經收到了取消信號。

這些函式和方法共同協作，通過使用 GPT 模型和 OpenAI 服務來實現程式碼的摘要功能。它們處理程式碼資料夾中的檔案，將相關的請求發送到 OpenAI 服務並獲取回應，最終生成和保存摘要報告。

## pkg/utils/helper.go

### UpdateApiServerAccessToken(endpoint, token string) error
此函式接收兩個參數 `endpoint` 和 `token`，用於更新 API 伺服器的存取權杖。它將檢查 `endpoint` 的有效性，然後建立一個 URL 以進行 PATCH 請求。此請求將包含一個 JSON 載荷，其中包含要更新的權杖。請求的標頭將包含 `Content-Type` 和 `Authorization`。然後，它使用 HTTP 客戶端執行請求並檢查回應狀態碼。如果回應的狀態碼不是 200，則返回錯誤；否則返回 nil。

### TrimMDTableHeader(summary string) string
此函式接收一個參數 `summary`，用於處理 Markdown 表格的摘要。它會逐行讀取並過濾 `summary` 中的內容，從第三行開始將所有非空行寫入字符串建構器中。然後返回建構器的內容。

### IsErrHTTP413(err error) bool
此函式接收一個錯誤參數 `err`，用於判斷是否為 HTTP 413 錯誤。它檢查錯誤是否是 `openai.RequestError` 或 `openai.APIError`，並檢查它們的 `HTTPStatusCode` 是否為 413。如果是，則返回 true；否則返回 false。

### IsErrUnexpectedEOF(err error) bool
此函式接收一個錯誤參數 `err`，用於判斷是否為意外的 EOF（文件結尾）錯誤。它使用 `errors.Is` 函式檢查錯誤是否等於 `io.ErrUnexpectedEOF`。如果是，則返回 true；否則返回 false。

