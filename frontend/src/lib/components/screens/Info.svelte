<script lang="ts">
    import {AtSymbol} from "@steeze-ui/heroicons";
    import {createEventDispatcher, onMount} from "svelte";
    import {fade} from "svelte/transition";
    import {Github, Producthunt} from "@steeze-ui/simple-icons";
    import OtherLink from "$lib/components/info/OtherLink.svelte";
    import {GetVersion} from "$lib/wailsjs/go/main/App";

    const dispatch = createEventDispatcher()
    let currentVersion = -1.0

    onMount(async () => {
        currentVersion = await GetVersion()
    })

    function escape(event: KeyboardEvent) {
        if (event.key === 'Escape') {
            dispatch('hide')
        }
    }
</script>
<svelte:window on:keydown={escape}/>
<div class="flex flex-col pt-4" in:fade>
    <div class="flex flex-row gap-2 items-center justify-between">
        <h1 class="text-2xl font-light text-pink-300 w-fit uppercase">Exponentia</h1>
        <p class="font-light text-white w-fit p-1 px-[0.84rem] hover:opacity-80 duration-300 ease-in-out text-xs">
            dataset v{currentVersion}
        </p>
    </div>
    <div class="max-w-[64rem]">
        <p class="text-sm lowercase text-justify">
            There was a lack of proper applications to help spelling contestants. Some were full of ads, lacked
            capabilities or
            weren't fit, and knowing that there might be a time when i will need one... i built one. it's that simple of
            a history.
        </p>
        <h1 class="text-2xl font-light text-pink-300 w-fit uppercase mt-4 mb-2">Privacy Policy</h1>
        <div class="text-sm lowercase text-justify" id="privacy-policy">
            <p>
                exponentia's desktop client does not use any analytical tools, therefore, we do not have any data
                collected, we also
                serve our datasets locally, allowing you to use the application offline.
            </p>
            <p class="my-2">
                other than those, no personal data is collected. all of these are verifiable from the official source
                code, deployments are
                automated to ensure that what is in our official source code is what you are seeing as well (no hidden
                tricks).
            </p>
        </div>
        <div class="flex flex-row gap-2 mt-4">
            <OtherLink href="https://github.com/ShindouMihou" alt="GitHub" icon={Github}/>
            <OtherLink href="https://producthunt.com/posts/exponentia" alt="ProductHunt" color="text-[#DA552F]"
                       icon={Producthunt}/>
            <OtherLink href="mailto:hello@mihou.pw" alt="Email" icon={AtSymbol}
                       color="bg-white text-black rounded-full"/>
        </div>
    </div>
</div>