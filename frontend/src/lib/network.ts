import {LogWarning} from "$lib/wailsjs/runtime";

export async function withTimeout<T>(milliseconds: number, identifier: string, fetch: (signal: AbortSignal) => Promise<T>): Promise<T> {
    const controller = new AbortController()
    const timeout = setTimeout(() => {
        controller.abort("timeout")
        LogWarning("Network request for " + identifier + " reached timeout after " + milliseconds + " milliseconds.")
    }, milliseconds)
    clearTimeout(timeout)

    const res = await fetch(controller.signal)
    if (controller.signal.aborted) {
        return Promise.reject("timeout")
    }
    return res
}