import {LogInfo} from "$lib/wailsjs/runtime";

function network(item: string) {
    info({ type: 'resource', item: item })
}

function event(item: any) {
    info({ type: "event", item: item })
}

function info(item: any) {
    LogInfo(JSON.stringify(item))
}

export default { network, info, event }