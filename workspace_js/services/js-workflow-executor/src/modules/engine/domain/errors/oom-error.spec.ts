/**
 * OOMError Unit Tests
 *
 * Tests the custom OOM error class used for V8 isolate out-of-memory signaling.
 */

import { OOMError } from './oom-error';

describe('OOMError', () => {
	it('should be an instance of Error', () => {
		const error = new OOMError('V8 heap exhausted');

		expect(error).toBeInstanceOf(Error);
		expect(error).toBeInstanceOf(OOMError);
	});

	it('should have name set to OOMError', () => {
		const error = new OOMError('test');

		expect(error.name).toBe('OOMError');
	});

	it('should preserve the message', () => {
		const error = new OOMError('Worker V8 OOM: isolate disposed');

		expect(error.message).toBe('Worker V8 OOM: isolate disposed');
	});

	it('should be catchable as OOMError via instanceof', () => {
		let caught = false;

		try {
			throw new OOMError('test');
		} catch (error) {
			if (error instanceof OOMError) {
				caught = true;
			}
		}

		expect(caught).toBe(true);
	});

	it('should be distinguishable from regular Error', () => {
		const oomError = new OOMError('OOM');
		const regularError = new Error('Regular');

		expect(oomError instanceof OOMError).toBe(true);
		expect(regularError instanceof OOMError).toBe(false);
	});
});
