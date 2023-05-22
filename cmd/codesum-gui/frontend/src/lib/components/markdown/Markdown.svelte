<script lang="ts">
	import SvelteMarkdown from 'svelte-markdown';
	import * as md from './renderers';

	export let source: any;

	const renderers = {
		text: md.Text,
		tablerow: md.TableRow,
		tablehead: md.TableHead,
		tablecell: md.TableCell,
		tablebody: md.TableBody,
		strong: md.Strong,
		paragraph: md.Paragraph,
		listitem: md.ListItem,
		link: md.Link,
		image: md.Image,
		html: md.Html,
		hr: md.Hr,
		heading: md.Heading,
		em: md.Em,
		del: md.Del,
		codespan: md.Codespan,
		br: md.Br,
		blockquote: md.Blockquote,
		code: md.Code,
		list: md.List,
		table: md.Table
	};
</script>

<div class="flex flex-col gap-4">
	<SvelteMarkdown options={{ breaks: true, gfm: true }} {source} {renderers} />
</div>

<slot />
