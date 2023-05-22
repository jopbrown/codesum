<script lang="ts">
	import { tileStore } from '$lib/stores';
	import Config from './Config.svelte';
	import Home from './Home.svelte';
	import Report from './Report.svelte';
	import Setting from './Setting.svelte';
</script>

<div class="flex bg-surface-50-900-token w-full h-full border-2 border-gray-300">
	<div class="flex-auto w-full h-full" hidden={$tileStore !== 'home'}><Home /></div>
	<div class="flex-auto w-full h-full" hidden={$tileStore !== 'report'}><Report /></div>
	<div class="flex-auto w-full h-full" hidden={$tileStore !== 'config'}><Config /></div>
	<div class="flex-auto w-full h-full" hidden={$tileStore !== 'setting'}><Setting /></div>
</div>
