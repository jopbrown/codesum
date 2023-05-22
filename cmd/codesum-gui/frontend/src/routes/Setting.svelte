<script lang="ts">
	import { configStore, tileStore } from '$lib/stores';
	import { ShowError, ShowSuccess } from '$lib/utils';
	import { ApplyAndSaveConfig, GetConfig } from '$lib/wailsjs/go/app/App';
	import type { cfgs } from '$lib/wailsjs/go/models';
	import { Loader } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { SlideToggle } from '@skeletonlabs/skeleton';

	let busy = false;

	let debug_mode = false;
	let log_path = '';
	let rules: cfgs.SummaryRules = {};
	let includeText = '';
	let excludeText = '';

	function reset() {
		busy = true;
		GetConfig()
			.then((result) => {
				$configStore = result;
				debug_mode = $configStore.debug_mode ?? false;
				log_path = $configStore.log_path ?? '';
				rules = { ...$configStore.summary_rules };
				includeText = rules.include!.join('\n');
				excludeText = rules.exclude!.join('\n');
			})
			.catch(ShowError)
			.finally(() => {
				busy = false;
			});
	}

	function apply() {
		busy = true;
		ApplyAndSaveConfig($configStore)
			.then(() => {
				ShowSuccess('Apply Success');
			})
			.catch(ShowError)
			.finally(() => {
				busy = false;
			});
	}

	$: if ($tileStore === 'config') {
		reset();
	}

	$: $configStore.debug_mode = debug_mode;
	$: $configStore.log_path = log_path;
	$: $configStore.summary_rules = { ...rules };
	$: $configStore.summary_rules!.include = includeText.split('\n').filter((val) => val);
	$: $configStore.summary_rules!.exclude = excludeText.split('\n').filter((val) => val);
</script>

<div class="flex flex-col">
	<div class="flex-auto flex flex-col h-[calc(100vh_-_108px)] overflow-y-auto gap-2 m-2">
		<h2 class="h2 w-full">Setting</h2>
		<SlideToggle class="ml-3" name="debug-mode" size="md" bind:checked={debug_mode} disabled={busy}>
			Debug Mode
		</SlideToggle>
		<label class="label">
			<span>Log Path</span>
			<input
				class="input"
				type="text"
				placeholder="{'${appDir}'}/log/codesum.log"
				bind:value={log_path}
				disabled={busy}
			/>
		</label>
		<label class="label">
			<span>Report Folder</span>
			<input
				class="input"
				type="text"
				placeholder="{'${appDir}'}/log/codesum.log"
				bind:value={rules.out_dir}
				disabled={busy}
			/>
		</label>
		<label class="label">
			<span>Report Naming</span>
			<input
				class="input"
				type="text"
				placeholder="{'${appDir}'}/log/codesum.log"
				bind:value={rules.out_file_name}
				disabled={busy}
			/>
		</label>
		<label class="label">
			<span>Include (line by line)</span>
			<textarea class="textarea" rows="3" placeholder="" bind:value={includeText} disabled={busy} />
		</label>
		<label class="label">
			<span>Exclude (line by line)</span>
			<textarea class="textarea" rows="3" placeholder="" bind:value={excludeText} disabled={busy} />
		</label>
	</div>
	<div class="flex-none flex m-1 h-[50px]">
		<div class="flex-auto" />
		<button
			type="button"
			class="btn variant-filled-warning flex-none w-40 mx-2"
			on:click={reset}
			disabled={busy}>Reset</button
		>
		<button
			type="button"
			class="btn variant-filled-primary flex-none w-40 mx-2"
			on:click={apply}
			disabled={busy}
		>
			Apply
			{#if busy}
				<Loader class="animate-spin" />
			{/if}
		</button>
	</div>
</div>

<!-- <div class="flex flex-col">
	<div class="flex-auto flex flex-col h-[calc(100vh_-_108px)] overflow-y-auto">
		<div class="flex flex-col border-2 border-gray-300 m-2 p-2">
			<h2 class="h2 w-full">ChatGPT</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
				<label class="label">
					<span>End Point</span>
					<input
						class="input"
						type="text"
						placeholder="https://api.openai.com/v1"
						bind:value={chat_gpt.end_point}
						disabled={busy}
					/>
				</label>
				<label class="label">
					<span>Protocol</span>
					<select class="select" bind:value={protocol} disabled={busy}>
						<option value="webapp">ChatGPT Web APP</option>
						<option value="api">OpenAI API</option>
					</select>
				</label>
				<label class="label">
					{#if protocol === 'webapp'}
						<span>Access Token</span>
						<input
							class="input"
							type="text"
							placeholder="eyxxxxxxxxxxxxxxxxxxxx"
							bind:value={chat_gpt.access_token}
							disabled={busy}
						/>
					{:else}
						<span>API Key</span>
						<input
							class="input"
							type="text"
							placeholder="sk-xxxxxxxxxxxxxxxxxxx"
							bind:value={chat_gpt.api_key}
							disabled={busy}
						/>
					{/if}
				</label>
				<label class="label">
					<span>Model</span>
					<input
						class="input"
						type="text"
						placeholder="gpt-3.5-turbo"
						bind:value={chat_gpt.model}
						disabled={busy}
					/>
				</label>
				<label class="label">
					<span>HTTP Proxy</span>
					<input
						class="input"
						type="text"
						placeholder="https://xxxxxxxxxxxxxxx:xx"
						bind:value={chat_gpt.proxy}
						disabled={busy}
					/>
				</label>
			</div>
		</div>
		<div class="flex flex-col border-2 border-gray-300 m-2 p-2">
			<h2 class="h2 w-full">Prompt</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2">
				<label class="label">
					<span>System</span>
					<textarea
						class="textarea"
						rows="6"
						placeholder=""
						bind:value={prompt.system}
						disabled={busy}
					/>
				</label>
				<label class="label">
					<span>Code Summaize</span>
					<textarea
						class="textarea"
						rows="6"
						placeholder={'${fileName}\n${fileContent}'}
						bind:value={prompt.code_summary}
						disabled={busy}
					/>
				</label>
				<label class="label">
					<span>Summary Table</span>
					<textarea
						class="textarea"
						rows="6"
						placeholder={'${filesCommaList}'}
						bind:value={prompt.summary_table}
						disabled={busy}
					/>
				</label>
				<label class="label">
					<span>Final Summary</span>
					<textarea
						class="textarea"
						rows="3"
						placeholder=""
						bind:value={prompt.final_summary}
						disabled={busy}
					/>
				</label>
			</div>
		</div>
		<div class="flex-auto" />
	</div>
	<div class="flex-none flex m-1 h-[50px]">
		<div class="flex-auto" />
		<button
			type="button"
			class="btn variant-filled-warning flex-none w-40 mx-2"
			on:click={reset}
			disabled={busy}>Reset</button
		>
		<button
			type="button"
			class="btn variant-filled-primary flex-none w-40 mx-2"
			on:click={apply}
			disabled={busy}
		>
			Apply
			{#if busy}
				<Loader class="animate-spin" />
			{/if}
		</button>
	</div>
</div> -->
