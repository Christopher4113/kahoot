<script lang="ts">
  export let quiz: { _id: string, name: string };
</script>

<div class='bg-white border-2 border-green-500 p-4 rounded-xl'>
    <h2 class="text-lg font-bold mb-2">{quiz.name}</h2>
</div>