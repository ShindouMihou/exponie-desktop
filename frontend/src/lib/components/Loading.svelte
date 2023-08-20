<script>
    import {loadingState} from "$lib/store";
</script>

<div class="w-full m-auto" style="--wails-draggable:drag">
    <div class="flex flex-col gap-2 w-full items-center justify-center m-auto animate-pulse">
        <p class="font-bold text-2xl uppercase text-black bg-white p-1 animate-bounce select-none">EX</p>
        <p class="fontlight text-sm uppercase text-[#3a3a3a] p-1">{$loadingState}</p>
    </div>
</div>