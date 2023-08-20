import {writable} from "svelte/store";
import type {Writable} from "svelte/store";

export const word = writable('')

export const input = writable('')
export const lastInput = writable('')

export const start = writable(-1)

export const isDisabled = writable(false)
export const sessionStatus: Writable<number> = writable(0)

export const loadingState: Writable<string> = writable("Loading configuration...")