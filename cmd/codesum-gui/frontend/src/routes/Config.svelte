<script lang="ts">
	import { configStore, tileStore } from '$lib/stores';
	import { ShowError, ShowSuccess } from '$lib/utils';
	import { ApplyAndSaveConfig, GetConfig } from '$lib/wailsjs/go/app/App';
	import type { cfgs } from '$lib/wailsjs/go/models';
	import { Loader } from 'lucide-svelte';
	import { onMount } from 'svelte';

	type ProtocolType = 'webapp' | 'api';
	let protocol: ProtocolType = 'webapp';

	let busy = false;

	let chat_gpt: cfgs.ChatGpt = {};
	let prompt: cfgs.Prompt = {};

	function reset() {
		busy = true;
		GetConfig()
			.then((result) => {
				$configStore = result;
				chat_gpt = { ...$configStore.chat_gpt };
				prompt = { ...$configStore.prompt };
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

	$: $configStore.chat_gpt = { ...chat_gpt };
	$: $configStore.prompt = { ...prompt };
</script>

<div class="flex flex-col">
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
</div>
