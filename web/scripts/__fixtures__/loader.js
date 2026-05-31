import {readFile} from 'node:fs/promises';
import {fileURLToPath} from 'node:url';
import {dirname, join} from 'node:path';

const here = dirname(fileURLToPath(import.meta.url));

export async function loadFixture(handler, scenario) {
    const path = join(here, 'ws', handler, `${scenario}.json`);
    const body = await readFile(path, 'utf8');
    return JSON.parse(body);
}

export function normalizeParams(stringKeyed) {
    const out = {};
    for (const [k, v] of Object.entries(stringKeyed)) {
        out[Number(k)] = v;
    }
    return out;
}
