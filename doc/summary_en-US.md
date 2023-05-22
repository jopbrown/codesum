# Code Summary

The analyzed software is a code summarization application. It is designed to generate summaries for code files using the ChatGPT model. The software provides a command-line interface for users to interact with and configure various aspects of the summarization process, including prompts, rules, and output settings. It also includes features for generating code summary questions and producing comprehensive reports. Additionally, utility functions are available to handle common operations and ensure smooth execution. Overall, the software aims to simplify the process of code comprehension and analysis by providing concise and informative summaries.

| File | Description |
| --- | --- |
| cmd/codesum/main.go | The main.go file in the cmd/codesum package is the entry point of the code summarization application. It contains the main function that initializes the application and handles command-line arguments. |
| pkg/cfgs/cfgs.go | The cfgs.go file in the pkg/cfgs package provides configuration structures and methods for managing various settings and options used in the code summarization process. It includes functions for loading and parsing configuration files. |
| pkg/cfgs/chatgpt.go | The chatgpt.go file in the pkg/cfgs package defines the ChatGpt structure, which represents the configuration settings related to the ChatGPT model used for code summarization. It specifies the endpoint, API key, access token, model, and proxy settings. |
| pkg/cfgs/prompt.go | The prompt.go file in the pkg/cfgs package defines the Prompt structure, which represents the configuration settings related to prompts used for code summarization. It includes fields for system prompts, code summary prompts, summary table prompts, and final summary prompts. |
| pkg/cfgs/rule.go | The rule.go file in the pkg/cfgs package defines the SummaryRules structure, which represents the configuration settings related to summary rules for code summarization. It includes fields for inclusion/exclusion rules, output directory, and output file name. |
| pkg/sumer/question.go | The question.go file in the pkg/sumer package provides functions for generating code summary questions. It includes functions for creating questions based on file summaries and summary tables, as well as a function for generating the final summary question. |
| pkg/sumer/report.go | The report.go file in the pkg/sumer package contains data structures and methods for generating code summary reports. It includes structures for code summaries, partial summaries, file summaries, and question-answer pairs. The file provides methods for writing the code summary as a Markdown file, generating the file summary table, and managing partial and final summaries. |
| pkg/sumer/summary.go | The 'summary.go' file in the 'pkg/sumer' package implements a code summarization feature. It defines a 'Summarizer' struct with methods for initializing a summarizer, summarizing code files, and sending requests to the GPT (OpenAI) API for generating summaries. It also includes helper functions for handling file content, saving a markdown report, and checking if the context is canceled. |
| pkg/utils/helper.go | The 'helper.go' file in the 'pkg/utils' package provides utility functions for various operations. It includes functions for updating the API server access token, trimming the header of a Markdown table from a summary string, and checking for specific types of errors related to HTTP status codes and unexpected end-of-file. These functions are used in other parts of the codebase to perform specific tasks. |


## cmd/codesum/main.go

This file contains the main entry point and the logic for the code summarization application.

1. The file imports various packages required for the application's functionality.
2. It declares variables `BuildName`, `BuildVersion`, `BuildHash`, and `BuildTime` representing the application's build information.
3. The `args` struct is declared to hold command-line arguments.
4. The `init` function sets up command-line flags and parses the arguments.
5. The `main` function is the entry point of the application. It calls the `run` function and handles any errors by logging them and terminating the program.
6. The `run` function loads the configuration from a specified file, applies logging configuration, and updates the access token for an API server if provided in the configuration. Then it creates a new `Summarizer` instance and starts the code summarization process.
7. The `startSummarize` function sets up a context with cancellation support and listens for termination signals. It creates an error group and executes the code summarization in a separate goroutine. If a termination signal is received or the context is canceled, the function terminates the code summarization. It waits for the goroutine to finish and returns any encountered errors.
8. The `pushMessage` function is a callback that receives a message from the summarizer and logs it, truncating the content if it exceeds a maximum length.
9. The `parseArgs` function parses the command-line arguments and sets the `args.codeFolder` variable to the provided code folder path. It also prints usage information and terminates the program if the arguments are invalid.
10. The `applyLog` function opens a log file based on the provided configuration, creates a logger, and sets it as the global logger.

Overall, this code file sets up the necessary configuration, initializes the logger, and orchestrates the code summarization process based on the provided command-line arguments and configuration.

## pkg/cfgs/cfgs.go

This file contains the logic for loading, merging, and saving configurations for the code summarization application.

1. The file imports various packages required for the configuration handling.
2. The `init` function registers an environment variable lookup function with the `strutil` package to enable expanding environment variables in configuration values.
3. The `Expander` type is defined as an alias for `strutil.Expander`.
4. The `Config` struct represents the configuration for the application, including fields for debug mode, log path, chat GPT settings, summary rules, and prompts.
5. The `defaultCfgFs` variable is declared as an embedded file system that contains the default configuration files.
6. The `DefaultConfig` function retrieves the default configuration by opening the default configuration file from the embedded file system and reading its contents. It returns the parsed configuration.
7. The `LoadConfig` function loads a configuration from the specified file by opening and reading it. It returns the parsed configuration or an error if the file cannot be opened or the configuration is invalid.
8. The `ReadConfig` function decodes a configuration from an `io.Reader` and returns the parsed configuration or an error if the decoding fails.
9. The `MergeDefault` method merges the current configuration with the default configuration using the `mergo.Merge` function. It returns an error if the merge fails.
10. The `SaveConfig` method saves the configuration to the specified file by opening the file for writing and writing the encoded configuration to it. It returns an error if the file cannot be opened or the writing fails.
11. The `WriteConfig` method encodes the configuration and writes it to the specified `io.Writer`. It returns an error if the encoding or writing fails.

Overall, this code file provides functions and methods for loading, merging, and saving configurations for the code summarization application. It supports default configurations, reading configurations from files, merging configurations with defaults, and writing configurations to files.

## pkg/cfgs/chatgpt.go

This file defines the `ChatGpt` struct, which represents the configuration settings for the ChatGpt component of the code summarization application.

1. The file declares the `ChatGpt` struct.
2. The `ChatGpt` struct has the following fields:
   - `EndPoint`: An `Expander` type that represents the endpoint for the ChatGpt API.
   - `APIKey`: An `Expander` type that represents the API key for accessing the ChatGpt API.
   - `AccessToken`: An `Expander` type that represents the access token for the ChatGpt API.
   - `Model`: An `Expander` type that represents the model to be used by the ChatGpt component.
   - `Proxy`: An `Expander` type that represents the proxy configuration for the ChatGpt component.

Overall, this code file defines the structure of the ChatGpt configuration, specifying the various fields required for the ChatGpt component of the code summarization application.

## pkg/cfgs/prompt.go

This file defines the `Prompt` struct, which represents the configuration settings for the prompts used in the code summarization application.

1. The file declares the `Prompt` struct.
2. The `Prompt` struct has the following fields:
   - `System`: An `Expander` type that represents the prompt for the system-related messages.
   - `CodeSummary`: An `Expander` type that represents the prompt for code summary messages.
   - `SummaryTable`: An `Expander` type that represents the prompt for the summary table.
   - `FinalSummary`: An `Expander` type that represents the prompt for the final summary.

Overall, this code file defines the structure of the Prompt configuration, specifying the various fields required for the prompts used in the code summarization application.

## pkg/cfgs/rule.go

This file defines the `SummaryRules` struct, which represents the configuration settings for the summary rules in the code summarization application.

1. The file declares the `SummaryRules` struct.
2. The `SummaryRules` struct has the following fields:
   - `Include`: A string slice representing the patterns to include for code summarization.
   - `Exclude`: A string slice representing the patterns to exclude from code summarization.
   - `OutDir`: An `Expander` type that represents the output directory for the code summaries.
   - `OutFileName`: An `Expander` type that represents the output file name for the code summaries.

Overall, this code file defines the structure of the SummaryRules configuration, specifying the various fields required for configuring the summary rules in the code summarization application. It allows specifying patterns to include and exclude, as well as the output directory and file name for the code summaries.

## pkg/sumer/question.go

This file contains functions that generate questions or prompts used by the code summarization component of the application.

1. The file imports the necessary packages for the code.
2. The `FileSummaryQuestion` function takes a prompt (represented by an `strutil.Expander`), a file name, and the content of a file. It expands the prompt by replacing placeholders with the file name and content, and returns the resulting question or prompt.
3. The `SummaryTableQuestion` function takes a prompt (represented by an `strutil.Expander`) and a list of file names. It joins the file names with commas and expands the prompt by replacing a placeholder with the comma-separated list of file names. The resulting question or prompt is then returned.
4. The `FinalSummaryQuestion` function takes a prompt (represented by an `strutil.Expander`). It directly returns the expanded prompt as a string without any additional placeholders or modifications.

Overall, this code file provides utility functions for generating questions or prompts used in the code summarization process. These functions use the `strutil.Expander` to dynamically expand the prompts by replacing specific placeholders with relevant values.

## pkg/sumer/report.go

This file contains various data structures and methods related to generating code summaries and generating reports.

1. The file imports the necessary packages for the code.
2. The `CodeSummary` struct represents a code summary and contains partial summaries and a final summary question and answer.
3. The `NewCodeSummary` function creates a new instance of `CodeSummary` and initializes its partial list.
4. The `SaveMarkdown` method saves the code summary as a Markdown file with the provided filename.
5. The `WriteMarkdown` method writes the code summary in Markdown format to the provided writer.
6. The `WriteFileSummaryTable` method writes the file summary table in Markdown format to the provided writer.
7. The `GetFileSummaryTable` method returns the file summary table as a string.
8. The `AddPartialSummaries` method adds partial summaries to the code summary with the given system prompt and returns the added partial summaries.
9. The `RequestFinalSummaryMessages` method constructs and returns a list of messages (chat completion messages) required for generating the final summary. It includes the system prompt, all partial summaries, and the final summary question.
10. The `SetFinalSummaryQA` method sets the final summary question and answer of the code summary and returns the created `QA` instance.
11. The `FinalSummary` method returns the content of the final summary.
12. The `PartialSummaries` struct represents partial summaries for a code summary and contains the system prompt, file summaries, and a summary question and answer.
13. The `FileSummary` struct represents a summary for a specific file and contains the filename and a question-answer pair.
14. The `QA` struct represents a question and answer pair and contains two messages: the question (user message) and the answer (assistant message).
15. The `NewQA` function creates a new instance of `QA` with the given question and answer content.
16. The `SetSystemPrompt` method sets the system prompt for the partial summaries and returns the created system message.
17. The `RequestFileSummaryMessages` method constructs and returns a list of messages (chat completion messages) required for generating the file summary. It includes the system prompt and all file summary questions and answers.
18. The `AddFileSummary` method adds a file summary to the partial summaries with the given filename, question, and answer, and returns the added file summary.
19. The `PopFileSummary` method removes and returns the last file summary from the partial summaries.
20. The `GetSummary` method returns the content of the summary (answer) for a file or the entire partial summaries.
21. The `RequestSummaryMessages` method constructs and returns a list of messages (chat completion messages) required for generating the summary. It includes the system prompt, all file summary questions and answers, and the summary question.
22. The `SetSummaryQA` method sets the summary question and answer for the partial summaries and returns the created `QA` instance.
23. The `GetSummary` method returns the content of the summary (answer) for the partial summaries.
24. The `FileList` method returns a list of filenames associated with the partial summaries.

Overall, this code file provides data structures and methods for managing code summaries, file summaries, and generating reports in Markdown format. It allows adding partial summaries, setting system prompts, generating questions and answers, and saving the code summary report.

## pkg/sumer/summary.go

This file contains the implementation of a code summarization package called `sumer`. The package provides a `Summarizer` struct with various methods for summarizing code files.

#### Struct: Summarizer
- `cfg`: Configuration object.
- `gptClient`: Client for the OpenAI GPT model.
- `fileFilter`: Matcher for filtering code files based on inclusion and exclusion rules.
- `pushMsgCallback`: Callback function for pushing chat completion messages.

#### Function: NewSummarizer
- Creates a new instance of `Summarizer`.
- Accepts the configuration object (`cfg`) and a callback function (`pushMsgCallback`) as parameters.
- Initializes the `cfg` and `gptClient` fields of the `Summarizer` instance.
- Compiles the inclusion and exclusion rules for file filtering.
- Returns the initialized `Summarizer` instance or an error.

#### Method: Summarizer.Summarize
- Summarizes the code files in a given folder.
- Accepts a context (`ctx`) and the path to the code folder (`codeFolder`) as parameters.
- Retrieves the list of code files using the `fsutil.ListWithMatcher` function and the `fileFilter` matcher.
- Initializes a `CodeSummary` object.
- Adds partial summaries for the system prompt using the `Prompt.System` configuration.
- Calls the `pushMsgCallback` with the system partial summary message.
- Iterates over each code file and performs the following steps:
  - Checks if the context is canceled and breaks the loop if it is.
  - Retrieves the file name and content using the `getCodeFileContent` function.
  - Constructs a question for the file summary using the file name and content.
  - Requests file summary messages using the partial summary object.
  - Sends a request to the OpenAI GPT model using the `sendRequest` method.
  - Adds the file summary to the `CodeSummary` object.
- Iterates over each partial summary and performs the following steps:
  - Constructs a question for the summary table using the file list.
  - Requests summary messages using the partial summary object.
  - Sends a request to the OpenAI GPT model using the `sendRequest` method.
  - Updates the summary for the partial summary object.
- Constructs a question for the final summary using the `Prompt.FinalSummary` configuration.
- Requests final summary messages using the `CodeSummary` object.
- Sends a request to the OpenAI GPT model using the `sendRequest` method.
- Sets the final summary for the `CodeSummary` object.
- Saves the Markdown report using the `saveMarkdownReport` method.
- Returns

## pkg/utils/helper.go

This file contains utility functions for various operations.

#### Function: UpdateApiServerAccessToken
- Updates the access token for the API server.
- Accepts the API server endpoint and the new token as parameters.
- Parses the endpoint URL and constructs the request URL.
- Converts the token to JSON payload.
- Creates an HTTP PATCH request with the request URL and payload.
- Sets the request headers.
- Sends the request using an HTTP client.
- Checks the response status code and returns an error if it's not 200 OK.

#### Function: TrimMDTableHeader
- Trims the header of a Markdown table from the given summary string.
- Accepts the summary string as a parameter.
- Uses a scanner to iterate through each line of the summary.
- Skips empty lines and keeps track of the line number.
- Appends non-empty lines to a string builder after skipping the first two lines (the header).
- Returns the trimmed summary string.

#### Function: IsErrHTTP413
- Checks if the given error indicates an HTTP 413 Request Entity Too Large error.
- Accepts an error as a parameter.
- Checks if the error is of type `*openai.RequestError` and has an HTTP status code of 413.
- Checks if the error is of type `*openai.APIError` and has an HTTP status code of 413.
- Returns `true` if the error indicates a 413 error, `false` otherwise.

#### Function: IsErrUnexpectedEOF
- Checks if the given error indicates an unexpected end-of-file error.
- Accepts an error as a parameter.
- Uses the `errors.Is` function to check if the error is equal to `io.ErrUnexpectedEOF`.
- Returns `true` if the error indicates an unexpected end-of-file, `false` otherwise.
