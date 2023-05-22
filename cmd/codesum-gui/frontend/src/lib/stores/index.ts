import { writable, type Writable } from 'svelte/store';
import type { cfgs } from '$lib/wailsjs/go/models';
import type { Tile } from '$lib/types';
import { localStorageStore } from '@skeletonlabs/skeleton';

export const tileStore: Writable<Tile> = writable('home');
export const statusStore = writable('ready');
export const configStore = writable({} as cfgs.Config);
export const selectedFileStore = writable('');
export const codeFolderStore: Writable<string> = localStorageStore('codeFolderStore', '');
