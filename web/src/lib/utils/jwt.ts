// $lib/utils/jwt.ts
import { ACCESS_TOKEN_SECRET } from '$env/static/private';

function base64urlToUint8Array(base64url: string): Uint8Array<ArrayBuffer> {
	const base64 = base64url.replace(/-/g, '+').replace(/_/g, '/');
	const padded = base64.padEnd(base64.length + ((4 - (base64.length % 4)) % 4), '=');
	const binary = atob(padded);
	const bytes = new Uint8Array(binary.length);
	for (let i = 0; i < binary.length; i++) {
		bytes[i] = binary.charCodeAt(i);
	}
	return bytes as Uint8Array<ArrayBuffer>;
}

function base64urlDecodeToString(base64url: string): string {
	const bytes = base64urlToUint8Array(base64url);
	return new TextDecoder().decode(bytes);
}

async function getHmacKey() {
	return crypto.subtle.importKey(
		'raw',
		new TextEncoder().encode(ACCESS_TOKEN_SECRET),
		{ name: 'HMAC', hash: 'SHA-256' },
		false,
		['verify']
	);
}

export async function verifyJwt(token: string) {
	const [headerB64, payloadB64, signatureB64] = token.split('.');
	if (!headerB64 || !payloadB64 || !signatureB64) {
		throw new Error('Malformed token');
	}

	const key = await getHmacKey();
	const signature = base64urlToUint8Array(signatureB64);
	const data = new TextEncoder().encode(`${headerB64}.${payloadB64}`);

	const isValid = await crypto.subtle.verify('HMAC', key, signature, data);
	if (!isValid) {
		throw new Error('Invalid signature');
	}

	const payload = JSON.parse(base64urlDecodeToString(payloadB64));

	if (payload.exp && payload.exp * 1000 < Date.now()) {
		throw new Error('Token expired');
	}

	return payload;
}
