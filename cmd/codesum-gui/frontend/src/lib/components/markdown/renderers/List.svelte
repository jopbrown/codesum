<script lang="ts">
	export let ordered: boolean;
	export let start: any;
</script>

{#if ordered}
	<ol class="list-decimal list-inside pl-5" {start}><slot /></ol>
{:else}
	<ul class="list-disc list-inside pl-5"><slot /></ul>
{/if}
