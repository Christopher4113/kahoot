<script lang="ts">
  import svelteLogo from './assets/svelte.svg'
  import viteLogo from '/vite.svg'
  import Button from './lib/Button.svelte';
  import QuizCard from './lib/QuizCard.svelte';

  let quizzes: {_id: string; name: string}[] = [];
  async function getQuizzes() {
    let response = await fetch('http://localhost:3000/api/quizzes')
    if (!response.ok) {
      alert('Failed to fetch quizzes')
      return
    }
    let json = await response.json()
    quizzes = json
    console.log(json)
  }

  let code = "";
  function connect() {
    let websocket = new WebSocket('ws://localhost:3000/ws')
    websocket.onopen = () => {
      console.log('WebSocket connection opened')
      websocket.send(`join: ${code}`)
    }
    websocket.onmessage = (event) => {
      console.log('Received message:', event.data)
    }
  }
  function hostQuiz(quiz) {
    let websocket = new WebSocket('ws://localhost:3000/ws')
    websocket.onopen = () => {
      console.log('WebSocket connection opened')
      websocket.send(`host: ${quiz._id}`)
    }
    websocket.onmessage = (event) => {
      console.log('Received message:', event.data)
    }
  }

</script>

<Button on:click={getQuizzes}>Get Quizzes</Button>

{#if quizzes.length > 0}
  <ul>
    {#each quizzes as quiz}
      <QuizCard quiz = {quiz} />
    {/each}
  </ul>
{/if}

<input bind:value={code} class="border" type="text" placeholder="Game Code"/>
<Button on:click={connect}>Join Game</Button>

