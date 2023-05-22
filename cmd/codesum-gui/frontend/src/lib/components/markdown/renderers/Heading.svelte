<script>
	export let depth = 0;
	export let raw = '';
</script>

{#if depth === 1}
	<h1 class="h1"><slot /></h1>
{:else if depth === 2}
	<h2 class="h2"><slot /></h2>
{:else if depth === 3}
	<h3 class="h3"><slot /></h3>
{:else if depth === 4}
	<h4 class="h4"><slot /></h4>
{:else if depth === 5}
	<h5 class="h5"><slot /></h5>
{:else if depth === 6}
	<h6 class="h6"><slot /></h6>
{:else}
	{raw}
{/if}
