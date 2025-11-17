<script lang="ts">
  import Button from './Button.svelte';
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
  export let quiz: { _id: string, name: string };
  function Host() {
    dispatch('host');
  }
</script>

<div class='flex justify-between items-center bg-white border-2 border-green-500 p-4 rounded-xl'>
    <h2 class="text-lg font-bold mb-2">{quiz.name}</h2>
    <Button on:click={Host}>Host</Button>
</div>
