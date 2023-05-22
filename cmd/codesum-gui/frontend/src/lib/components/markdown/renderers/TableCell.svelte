<script lang="ts">
	export let header: boolean;
	export let align: any;
</script>

{#if header}
	<th class="table-cell-fit" {align}><slot /></th>
{:else}
	<td class="table-cell-fit" {align}><slot /></td>
{/if}

<style>
	td {
		white-space: normal !important;
	}
</style>
