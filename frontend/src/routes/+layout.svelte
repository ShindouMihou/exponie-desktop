<script>
  import "../app.postcss";
</script>

<svelte:head>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="true">
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@100;200;300;400;500;600;700;800;900&display=swap" rel="stylesheet">
    <title>Exponentia</title>
    <meta name="title" content="Exponentia"/>
    <meta name="og:title" content="Exponentia"/>
    <meta name="og:type" content="website"/>
    <meta name="description" content="A minimalist's spelling test website. Challenge your vocabulary."/>
    <meta name="og:description" content="A minimalist's spelling test website. Challenge your vocabulary."/>
    <meta name="theme-color" content="#fff"/>
    <style>
        .inter { font-family: 'Inter', 'sans-serif'; }
        body { background-color: #000; }
    </style>
</svelte:head>

<div class="p-6 inter flex text-white flex-col min-h-screen">
    <main class="flex w-full flex-grow">
        <slot/>
    </main>
</div>
