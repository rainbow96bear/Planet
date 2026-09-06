<script lang="ts">
	import { page } from '$app/stores';
	import './layout.css';
	import '../app.css';
	let { children } = $props();
</script>

{@render children()}
