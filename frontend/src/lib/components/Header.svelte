<script>
    import {createEventDispatcher, onMount} from "svelte";
    import {Quit, WindowIsMaximised, WindowMinimise, WindowToggleMaximise} from "$lib/wailsjs/runtime/runtime.js";
    import {Icon} from "@steeze-ui/svelte-icon";
    import {Cross2, EnterFullScreen, ExitFullScreen} from "@steeze-ui/radix-icons"
    import {currentScreen} from "$lib/store";

    const dispatch = createEventDispatcher()
    let maximized = false

    onMount(async () => {
        maximized = await WindowIsMaximised()
    })

    function toggleMaximize() {
        WindowToggleMaximise()
        maximized = !maximized
    }
</script>

<div class="w-full justify-between flex flex-row items-center" style="--wails-draggable:drag">
    <div>
        <button on:click={() => dispatch('show', $currentScreen !== 'PLAY' ? 'PLAY' : 'INFO')}
                class="font-bold uppercase text-white w-fit p-1 px-[0.84rem] hover:opacity-80 duration-300 ease-in-out text-xs">
            Exponie.me
        </button>
    </div>
    <div class="flex flex-row gap-2 items-center">
        <button on:click={WindowMinimise}
                class="font-bold uppercase text-white w-fit p-1 px-[0.84rem] hover:opacity-80 duration-300 ease-in-out text-xs">
            -
        </button>
        <button on:click={toggleMaximize}
                class="font-bold uppercase text-white w-fit p-1 px-[0.84rem] hover:opacity-80 duration-300 ease-in-out text-xs">
            <Icon src={maximized ? ExitFullScreen : EnterFullScreen} size="16"/>
        </button>
        <button on:click={Quit}
                class="font-bold uppercase text-red-500 w-fit p-1 px-[0.84rem] hover:opacity-80 duration-300 ease-in-out text-xs">
            <Icon src={Cross2} size="16"/>
        </button>
    </div>
</div>