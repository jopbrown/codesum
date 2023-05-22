<script lang="ts">
	import Markdown from '$lib/components/markdown/Markdown.svelte';
	import { selectedFileStore, statusStore, tileStore, codeFolderStore } from '$lib/stores';
	import { ShowError, ShowSuccess } from '$lib/utils';
	import { CodeSummarize, SelectFolder, Stop } from '$lib/wailsjs/go/app/App';
	import type { openai } from '$lib/wailsjs/go/models';
	import { EventsOff, EventsOn } from '$lib/wailsjs/runtime';
	import { Loader } from 'lucide-svelte';
	import { afterUpdate, beforeUpdate, onMount } from 'svelte/internal';

	const eventPushMessage = 'pushMessage';
	const eventStreamAnswer = 'streamAnswer';

	let elemChat: HTMLElement;

	let messages: openai.ChatCompletionMessage[] = [];
	let autoscroll = false;

	function scrollChatBottom(behavior?: ScrollBehavior): void {
		elemChat.scrollTo({ top: elemChat.scrollHeight, behavior });
	}

	beforeUpdate(() => {
		autoscroll =
			elemChat && elemChat.offsetHeight + elemChat.scrollTop > elemChat.scrollHeight - 50;
	});

	afterUpdate(() => {
		if (autoscroll) scrollChatBottom();
	});

	onMount(() => {
		scrollChatBottom();
	});

	function selectFolder() {
		SelectFolder()
			.then((result) => {
				if (result) {
					$codeFolderStore = result;
				}
			})
			.catch((err) => {
				statusStore.set(err);
			});
	}

	let isRunning = false;
	let streamedAnswer = '';
	function start() {
		if ($codeFolderStore.length == 0) {
			ShowError("Code folder can't not be empty");
			return;
		}
		statusStore.set('Running...');

		EventsOn(eventPushMessage, (data) => {
			const msg = data as openai.ChatCompletionMessage;
			messages = [...messages, msg];
			streamedAnswer = '';
		});
		EventsOn(eventStreamAnswer, (data) => {
			const delta = data as string;
			streamedAnswer = streamedAnswer + delta;
		});

		isRunning = true;
		messages = [];
		streamedAnswer = '';

		CodeSummarize($codeFolderStore)
			.then((result) => {
				ShowSuccess('Summarize success', {
					label: 'Check Report',
					response: () => {
						$selectedFileStore = result;
						$tileStore = 'report';
					},
				});
			})
			.catch(ShowError)
			.finally(() => {
				isRunning = false;
				streamedAnswer = '';
				EventsOff(eventPushMessage);
				EventsOff(eventStreamAnswer);
			});
	}

	function stop() {
		Stop().then();
	}
</script>

<div class="flex-none flex flex-col">
	<section
		class="flex-auto h-[calc(100vh_-_108px)] overflow-y-auto p-2 space-y-4"
		bind:this={elemChat}
	>
		{#if messages.length == 0}
			<section class="h-full flex items-center justify-center">
				<p>no message yet</p>
			</section>
		{:else}
			{#each messages as bubble}
				{#if bubble.role === 'assistant'}
					<div class="flex gap-2">
						<div class="flex-none bg-secondary-400 h-12 w-12 rounded-full text-center pt-3">A</div>
						<div class="flex-auto card p-4 bg-secondary-200 rounded-tl-none space-y-2">
							<Markdown source={bubble.content} />
						</div>
					</div>
				{:else}
					<div class="flex gap-2">
						<div class="flex-auto card p-4 bg-primary-200 rounded-tr-none space-y-2">
							<Markdown source={bubble.content} />
						</div>
						<div class="flex-none bg-primary-400 h-12 w-12 rounded-full text-center pt-3">Q</div>
					</div>
				{/if}
			{/each}
			{#if isRunning && messages.slice(-1)[0].role == 'user'}
				<div class="flex gap-2">
					<div class="flex-none bg-secondary-400 h-12 w-12 rounded-full text-center pt-3">A</div>
					<div class="flex-auto card p-4 bg-secondary-200 rounded-tl-none space-y-2">
						<Markdown source={streamedAnswer} />
						<Loader class="animate-spin w-4" />
					</div>
				</div>
			{/if}
		{/if}
	</section>
	<section class="flex-none h-[60px] flex p-2 gap-4">
		<div
			class="flex-auto input-group input-group-divider grid-cols-[1fr_auto]
			{$codeFolderStore ? '' : 'input-error'}"
		>
			<input
				type="text"
				placeholder="Input code folder..."
				bind:value={$codeFolderStore}
				disabled={isRunning}
			/>
			<button class="variant-filled-tertiary" on:click={selectFolder} disabled={isRunning}>
				Select
			</button>
		</div>
		<button
			class="flex-none btn {isRunning ? 'variant-filled-error' : 'variant-filled-primary'} w-40"
			on:click={() => {
				isRunning ? stop() : start();
			}}
		>
			{isRunning ? 'Stop' : 'Start'}
		</button>
	</section>
</div>
