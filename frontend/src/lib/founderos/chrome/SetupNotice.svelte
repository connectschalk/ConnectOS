<!-- Shown in place of the /os views when the signed-in account doesn't own a
     ConnectOS yet (fresh install, or a teammate's account): one panel with
     the one command that sets it up, instead of every card failing. -->
<script lang="ts">
	import { PageHeader } from '$lib/founderos/kit';

	let { email }: { email: string } = $props();
	const command = $derived(`make demo OWNER=${email || '<your email>'}`);
</script>

<section class="bn-card mx-auto mt-10 max-w-2xl p-6" data-part="setup-notice">
	<PageHeader eyebrow="first run" title="Set up ConnectOS" />
	<p class="bn-muted mt-4 font-mono text-[12px] leading-relaxed">
		This account doesn't have a ConnectOS yet. From the repository root, create its workspaces (HQ, Vantage,
		Launchpad Cohort, Personal) and load the demo data:
	</p>
	<pre class="setup-cmd bn-text mt-4 overflow-x-auto px-4 py-3 font-mono text-[12px]">{command}</pre>
	<p class="bn-dim mt-4 font-mono text-[11px] leading-relaxed">Then reload this page. These are sample companies. Re-running this command refreshes sample records; keep your real companies in separate workspaces.</p>
</section>

<style>
	.setup-cmd {
		background: var(--bn-surface-2);
		border: 1px solid var(--bn-border);
		border-radius: var(--bn-r-chip);
	}
</style>
