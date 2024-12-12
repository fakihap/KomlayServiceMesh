<script lang="ts">
    import { createQuery } from '@tanstack/svelte-query'

    const query = createQuery({
        queryKey: ['ping'],
        queryFn: async () => 
            await fetch('http://booking-service:8080/ping')
                .then((res) => 
                   res.json(),
                )
        
    })
</script>

<div>
    {#if $query.isLoading}
        <p>Loading...</p>
    {:else if $query.isError}
        <p>Error: {$query.error.message}</p>
    {:else if $query.isSuccess}
        {#each $query.data as todo}
        <p>{ping}</p>
        {/each}
    {/if}

    <h1>SVELTEKIT!!!!!!!</h1>
    <p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>

</div>