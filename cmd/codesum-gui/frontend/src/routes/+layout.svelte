<script lang="ts">
	// The ordering of these imports is critical to your app working properly
	import '@skeletonlabs/skeleton/themes/theme-skeleton.css';
	// If you have source.organizeImports set to true in VSCode, then it will auto change this ordering
	import '@skeletonlabs/skeleton/styles/skeleton.css';
	// Most of your app wide CSS should be put in this file
	import '../app.postcss';

	import Statusbar from '$lib/components/Statusbar.svelte';
	import { configStore, tileStore } from '$lib/stores';
	import { ShowError } from '$lib/utils';
	import { LoadAndApplyConfig } from '$lib/wailsjs/go/app/App';
	import {
		AppRail,
		AppRailTile,
		AppShell,
		LightSwitch,
		Toast,
		autoModeWatcher,
		storeHighlightJs,
	} from '@skeletonlabs/skeleton';
	import hljs from 'highlight.js';
	import 'highlight.js/styles/github-dark.css';
	import { FileCog, FileText, Home, Settings } from 'lucide-svelte';
	import { onMount } from 'svelte';

	storeHighlightJs.set(hljs);

	onMount(() => {
		autoModeWatcher();
		LoadAndApplyConfig()
			.then((result) => {
				$configStore = result;
			})
			.catch(ShowError);
	});
</script>

<Toast position="t" />

<!-- <slot /> -->

<AppShell>
	<svelte:fragment slot="sidebarLeft">
		<AppRail>
			<AppRailTile bind:group={$tileStore} name="home" value="home">
				<svelte:fragment slot="lead"><Home class="inline-block" /></svelte:fragment>
				<span>Home</span>
			</AppRailTile>
			<AppRailTile bind:group={$tileStore} name="report" value="report">
				<svelte:fragment slot="lead"><FileText class="inline-block" /></svelte:fragment>
				<span>Report</span>
			</AppRailTile>
			<AppRailTile bind:group={$tileStore} name="config" value="config">
				<svelte:fragment slot="lead"><FileCog class="inline-block" /></svelte:fragment>
				<span>Config</span>
			</AppRailTile>

			<svelte:fragment slot="trail">
				<AppRailTile bind:group={$tileStore} name="setting" value="setting">
					<svelte:fragment slot="lead"><Settings class="inline-block" /></svelte:fragment>
					<span>Setting</span>
				</AppRailTile>
			</svelte:fragment>
		</AppRail>
	</svelte:fragment>
	<slot />
	<svelte:fragment slot="footer">
		<div class="flex p-2">
			<div class="flex-auto w-full"><Statusbar /></div>
			<div class="flex-none"><LightSwitch /></div>
		</div>
	</svelte:fragment>
</AppShell>
