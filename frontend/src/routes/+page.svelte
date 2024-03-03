<script lang="ts">
    import {onDestroy, onMount} from "svelte";
    import {fade} from "svelte/transition"
    import Loading from "$lib/components/Loading.svelte";
    import MountError from "$lib/components/MountError.svelte";
    import InputField from "$lib/components/InputField.svelte";
    import Hint from "$lib/components/Hint.svelte";
    import {currentScreen, input, isDisabled, lastInput, loadingState, sessionStatus, start, word} from "$lib/store";
    import terminal from "$lib/logging"
    import {hide, reduce} from "$lib/word";
    import IconButton from "$lib/components/shared/IconButton.svelte";
    import {ArrowPath} from "@steeze-ui/heroicons";
    import Info from "$lib/components/screens/Info.svelte";
    import Header from "$lib/components/Header.svelte";
    import {LogError, LogInfo} from "$lib/wailsjs/runtime";
    import {withTimeout} from "$lib/network";
    import {EnsureDataset, GetDataset, GetDefinitions} from "$lib/wailsjs/go/exponie/App";

    let loaded = false

    const DEFAULT_DEFINITION = 'Looking for the definition...';
    let definition: string = DEFAULT_DEFINITION;

    let hintShown = false;

    let end: number = -1;

    let throttle = false;
    let offline = false;

    let autoSuggestedDetected = false

    let mountErrors: string[] = []

    let dataset: string[] = []
    let definitions: {[key: string]: string}

    $: reduced = reduce($word);
    $: hidden = hide($word);

    let offlineTimer = setInterval(fasterCheckInternetConnection, 2 * 1000)
    let fasterOfflineTimer = setInterval(checkInternetConnection, 15 * 1000)
    onDestroy(() => {
        clearInterval(offlineTimer);
        clearInterval(fasterOfflineTimer);
    });

    onMount(async () => {
        $loadingState = "Checking internet connection..."
        await checkInternetConnection()
        try {
            if (!offline) {
                $loadingState = "Checking for dataset updates..."
                await EnsureDataset()
            }
            $loadingState = "Loading dataset..."
            dataset = await GetDataset()
            $loadingState = "Loading definitions..."
            definitions = await GetDefinitions()

            LogInfo("size of dataset: " + dataset.length + ", offline: " + offline + ", definitions: " + Object.keys(definitions).length)
            if ((dataset.length === 0 || definitions == null) && offline) {
                mountErrors = [
                    "You need internet connection for the first time opening exponie.me, this is because we need to get " +
                    "the dataset needed for the application to work. Please try again later."
                ]
                return
            }

            $loadingState = "Application should have loaded already, this is a bug."
            if (!navigator.onLine) {
                setOffline()
            }

            await reset()
        } catch (error: any) {
            if (error instanceof Error) {
                mountErrors = [error.message]
            } else {
                mountErrors = [error]
            }
        }
    })

    function fasterCheckInternetConnection() {
        if (!navigator.onLine && !offline) {
            setOffline()
        }
    }

    async function checkInternetConnection() {
        if (navigator.onLine) {
            try {
                await withTimeout(1_000, "check_connection", async (signal) => {
                    await fetch('https://exponie.mihou.pw/hello.txt', { signal: signal }).then((response) => {
                        if (response.ok) {
                            offline = false;
                            if (definition === DEFAULT_DEFINITION) {
                                define($word, false);
                            }
                        } else {
                            if (!offline) {
                                setOffline()
                            }
                        }
                    }).catch((e) => {
                        if (!offline) {
                            setOffline()
                        }
                    })
                })
            } catch (e: any) {
                if (e === "timeout" || e.message === "timeout") {
                    setOffline()
                    return
                }
                LogError("Network request for check_connection errored: " + e)
            }
        } else {
            if (!offline) {
                setOffline();
            }
        }
    }

    function setOffline() {
        offline = true;
        hintShown = true;
    }

    function navigate(to: string) {
        $currentScreen = to;
    }

    function random(): string {
        return dataset[Math.floor(Math.random() * dataset.length)].toLowerCase();
    }

    function define(word: string, useCache: boolean = true) {
        if (useCache) {
            let def = definitions[word]
            if (def != null) {
                terminal.event({ ev: "definition", cached: true, word: word, definition: def })
                definition = def
                return
            }
        }

        terminal.network(word)
        return withTimeout(2_000, "definition", async (signal) => {
            await fetch('https://api.dictionaryapi.dev/api/v2/entries/en/' + word, { signal: signal })
                .then((response) => response.ok === true ? response.json() : null)
                .then((data: Array<any>) => {
                    if (data != null) {
                        definition = data.at(0).meanings[0].definitions[0].definition
                    } else {
                        definition = 'No definition found.';
                        throw {error: 'No definition found, this exception is intentional.'}
                    }
                })
        })
    }

    async function reset() {
        if (throttle) {
            terminal.event('throttled')
            return;
        }

        throttle = true;
        document.getElementById('container')?.classList.add('animate-pulse');

        try {
            $start = -1;
            end = -1;

            $input = '';
            $lastInput = '';

            hintShown = true;
            let n = random()

            await define(n);

            $word = n;
            const inputField = document.getElementById('input')
            $isDisabled = false;

            if (inputField) {
                $sessionStatus = 0
                //@ts-ignore
                inputField.value = '';
                inputField.focus()

                // important: Subsequent focus is for mobile purposes since mobile requires two focuses for some reason.
                setTimeout(() => inputField.focus(), 500)
            }

            setTimeout(() => {
                throttle = false;
                document.getElementById('container')?.classList.remove('animate-pulse');
            }, 150);
            terminal.event({ev: 'res', word: $word, def: definition})
        } catch (e) {
            throttle = false;
            await reset()
        }
    }

    function hint() {
        if (!offline) {
            hintShown = !hintShown;
        }

        terminal.event({ev: 'tog', opt: 'hint'})
    }

    function cheated() {
        $input = $lastInput;
        autoSuggestedDetected = true;
    }

    function erase() {
        $input = '';
        $lastInput = '';
    }

    function hideAutoSuggestionWarning() {
        autoSuggestedDetected = false;
    }

    async function complete() {
        if (end !== -1) return

        $isDisabled = true;
        hintShown = true;

        reduced = $word;
        end = Date.now();

        if ($input === $word) {
            $sessionStatus = 1
            terminal.event({ev: 'compl', s: true})
            return;
        }

        $sessionStatus = 2
        terminal.event({ev: 'compl', s: false})
    }

    function handleGlobalKeyDown(event: KeyboardEvent) {
        if (event.key === 'Tab') {
            event.preventDefault();
            document.getElementById('reset')?.focus();

        }

        if (event.key === 'Enter' && end !== -1 && ((end + 100) < Date.now())) {
            event.preventDefault();
            reset();
        }
    }

    async function handleResetKeyDown(event: any) {
        if (event.key === 'Enter') {
            event.preventDefault();
            await reset();

            document.getElementById('input')!!.focus();
        }
    }

    function _head(ev: CustomEvent<string>) {
        navigate(ev.detail)
    }
</script>

<svelte:window on:keydown={handleGlobalKeyDown}/>

<div class="w-full flex flex-col gap-2">
    <Header on:show={_head}/>
    {#if dataset.length === 0}
        {#if mountErrors.length === 0}
            <Loading/>
        {:else}
            <MountError errors={mountErrors}/>
        {/if}
    {:else}
        {#if $currentScreen === 'PLAY'}
            <div class="w-full m-auto" id="container" in:fade>
                <div class="flex flex-col gap-2 w-full items-center justify-center m-auto">
                    {#if autoSuggestedDetected}
                        <div class="bg-red-500 font-white p-2 text-xs" in:fade out:fade>
                            Auto-correct, or similar was detected and prevented.
                        </div>
                    {/if}
                    <Hint on:click={hint} hidden={hidden} reduced={reduced} hintShown={hintShown}/>
                    {#if !offline}
                        <p class="font-light text-sm lowercase text-center max-w-xl">{definition}</p>
                    {/if}
                    <InputField
                            on:cheated={cheated}
                            on:complete={complete}
                            on:erase={erase}
                            on:input={hideAutoSuggestionWarning}
                    />
                    <div class="flex flex-row gap-4">
                        <IconButton id="reset" icon={ArrowPath} on:click={reset} on:keydown={handleResetKeyDown}/>
                    </div>
                    {#if end !== -1}
                        <p class="font-light text-sm max-w-xl pt-4">{(end - $start) / 1000} seconds</p>
                    {/if}
                </div>
            </div>
            <div class="pt-18 flex flex-col gap-2 text-xs">
                <div class="flex flex-row gap-1 items-center">
                    {#if offline}<p
                            class="bg-red p-1 text-black text-xs font-light bg-red-500 hover:opacity-80 duration-300 ease-in-out"
                            in:fade>OFFLINE</p>{/if}
                </div>
            </div>
        {:else if $currentScreen === 'INFO'}
            <Info on:hide={() =>  navigate('PLAY')}/>
        {/if}
    {/if}
</div>
