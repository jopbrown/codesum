import { statusStore } from '$lib/stores';
import { toastStore, type ToastSettings } from '@skeletonlabs/skeleton';

export const ShowError = (err: string) => {
	const t: ToastSettings = {
		message: err,
		background: 'variant-filled-error',
	};
	toastStore.trigger(t);
	statusStore.set('error: ' + err.replace('\n', ' '));
};

export const ShowSuccess = (msg: string, action?: any) => {
	let t: ToastSettings = {
		message: msg,
		background: 'variant-filled-success',
	};
	if (action) {
		t.action = { ...action };
		t.autohide = false;
	}
	toastStore.trigger(t);
	statusStore.set('Ready');
};
