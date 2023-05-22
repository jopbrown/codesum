<script lang="ts">
	import { CodeBlock } from '@skeletonlabs/skeleton';
	export let lang: string;
	export let text: string;
</script>

<CodeBlock lineNumbers language={lang ? lang : 'plaintext'} code={text} />
