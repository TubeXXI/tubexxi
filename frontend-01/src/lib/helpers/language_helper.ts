import { type RequestEvent } from '@sveltejs/kit';
import { translate, type SingleResponse, type BatchResponse } from "@siamf/google-translate";
import { BaseHelper } from './base_helper';

export class LanguageHelper extends BaseHelper {
	constructor(event: RequestEvent) {
		super(event);
	}

	async singleTranslate(text: string, targetLang: string): Promise<SingleResponse> {
		try {
			return await translate.single(text, targetLang);
		} catch (error) {
			console.error('Error translating text:', error);
			const message = error instanceof Error ? error.message : 'Unknown error';
			return await translate.single(message, 'auto');
		}
	}

	async batchTranslate(texts: string, targetLang: string[]): Promise<BatchResponse> {
		try {
			return await translate.batch(texts, targetLang);
		} catch (error) {
			console.error('Error translating texts:', error);
			const message = error instanceof Error ? error.message : 'Unknown error';
			return await translate.batch(message, ['auto']);
		}
	}
}
