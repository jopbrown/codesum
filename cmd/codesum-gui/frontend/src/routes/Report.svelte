<script lang="ts">
	import Markdown from '$lib/components/markdown/Markdown.svelte';
	import { selectedFileStore, tileStore } from '$lib/stores';
	import { ShowError } from '$lib/utils';
	import { GetFileContent, GetReportList } from '$lib/wailsjs/go/app/App';
	import { Loader, RotateCw } from 'lucide-svelte';

	let reportFiles: string[] = [];
	let loading = false;
	let fileContent: string = '';

	function reloadFileList() {
		GetReportList()
			.then((result) => {
				reportFiles = result ?? [];
			})
			.catch(ShowError);
	}

	function loadFile() {
		loading = true;
		GetFileContent($selectedFileStore)
			.then((result) => {
				fileContent = result;
			})
			.catch(ShowError)
			.finally(() => {
				loading = false;
			});
	}

	$: if ($selectedFileStore !== '') {
		loadFile();
	}

	$: if ($tileStore === 'report') {
		reloadFileList();
	}

	function baseName(fname: string) {
		fname = fname.replaceAll('\\', '/');
		let base = new String(fname).substring(fname.lastIndexOf('/') + 1);
		if (base.lastIndexOf('.') != -1) base = base.substring(0, base.lastIndexOf('.'));
		return base;
	}
</script>

<div class="flex h-full w-full overflow-hidden">
	<div class="flex-none flex flex-col h-full w-60 border-r-2 border-gray-300">
		<div class="flex-none flex border-b-4 p-2 border-gray-300">
			<div class="flex-auto w-full flex items-center justify-center">
				<p>{reportFiles.length} report{reportFiles.length == 1 ? '' : 's'}</p>
			</div>
			<a
				href={'#'}
				class="flex-none btn-icon btn-icon-sm hover:variant-soft-primary"
				on:click={reloadFileList}
			>
				<RotateCw class="w-4" />
			</a>
		</div>
		<div class="flex-auto h-full">
			<select class="select h-full text-sm w-full p-0" size="2" bind:value={$selectedFileStore}>
				{#each reportFiles as filePath}
					<option value={filePath}>{baseName(filePath)}</option>
				{/each}
			</select>
		</div>
	</div>
	<div class="flex-auto w-full flex flex-col">
		<div class="flex-none flex items-center justify-center">
			<p class="text-xl"><strong>{baseName($selectedFileStore)}</strong></p>
		</div>
		<div class="flex-auto w-full">
			{#if $selectedFileStore === ''}
				<section class="h-full flex items-center justify-center">
					<p>not select yet</p>
				</section>
			{:else if loading}
				<section class="h-full flex items-center justify-center">
					<Loader class="animate-spin" />
				</section>
			{:else}
				<div class="w-[calc(100vh_-320px)] h-[calc(100vh_-_75px)] overflow-auto p-2">
					<Markdown source={fileContent} />
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.select option {
		padding: 0.5rem;
		border-radius: 10px;
		text-overflow: ellipsis;
		overflow: hidden;
	}
</style>
