import { zodValidationError, IsStringDateFormat } from './index';
import { STRING_DATE_FORMATS } from './constants';
import { ZodErrorResult } from './interfaces';

// Global vars
const zodError = {
	success: false,
	error: {
		issues: [
			{
				code: 'too_small',
				minimum: 10000,
				type: 'number',
				inclusive: true,
				path: ['address', 'zipCode'],
				message: 'Value should be greater than or equal to 10000',
			},
			{
				code: 'invalid_type',
				expected: 'string',
				received: 'number',
				path: ['names', 1],
				message: 'Invalid input: expected string, received number',
			},
			{
				code: 'unrecognized_keys',
				keys: ['extra'],
				path: ['address'],
				message: 'Unrecognized key(s) in object: "field"',
			},
		],
	},
} as unknown as ZodErrorResult;

/**
 * Start test functionality
 */
describe('Zod validation', () => {

	it('Success: zodValidation', async () => {
		// Arrange
		const messages = zodError
			.error
			.issues
			.map((issue) => {
				let { message, path } = issue;
				message = String(message).toLowerCase();
				return `${path.join('.')} ${message}`;
			});

		// Act
		const result = await zodValidationError(zodError);

		// Assert
		expect(result).toEqual(messages);
	});

	describe('IsStringDateFormat', () => {

		const dateExamples = [
			'2023-01-01',              // Date
			'2023-12-31',              // Las day of year
			'9999-12-31',              // Date max
			'1609459200',              // Unix timestamp on seconds (2021-01-01T00:00:00Z)
			'1625097600',              // Unix timestamp on seconds (2021-07-01T00:00:00Z)
			'1672531199',              // Unix timestamp on seconds (2022-12-31T23:59:59Z)
			'1609459200000',           // Unix timestamp on milliseconds
			'1625097600000',           // Unix timestamp on milliseconds
			'1672531199000',           // Unix timestamp on milliseconds
			'2023-01-01 12:00:00',     // Data com hora
			'2023-02-14 18:30:25',     // Date time
			'2023-10-01 15:30:45.123',  // Date time
		];

		const expectedResults = [
			new Date('2023-01-01'),
			new Date('2023-12-31'),
			new Date('9999-12-31'),
			new Date('2021-01-01T00:00:00Z'),
			new Date('2021-07-01T00:00:00Z'),
			new Date('2022-12-31T23:59:59Z'),
			new Date('2021-01-01T00:00:00Z'),
			new Date('2021-07-01T00:00:00Z'),
			new Date('2022-12-31T23:59:59Z'),
			new Date('2023-01-01T12:00:00'),
			new Date('2023-02-14T18:30:25'),
			new Date('2023-10-01T15:30:45.123'),
		];

		it.each(dateExamples)('should validate and transform valid date "%s" to Date object', (input) => {
			const result = IsStringDateFormat.parse(input);
			expect(result).toEqual(input);
		});

		it('should throw an error for invalid date formats', () => {
			const invalidDates = ['2023-13-01', 'Hello World', '123', '12-31-2023'];

			invalidDates.forEach((date) => {
				expect(() => IsStringDateFormat.parse(date)).toThrow();
			});
		});

		it('should return an error message indicating valid formats', () => {
			const invalidDate = 'InvalidDate';

			try {
				IsStringDateFormat.parse(invalidDate);
			} catch (error: any) {
				if (typeof error === 'object' && error !== null && 'errors' in error) {
					expect(error.errors[0].message).toBe(`Invalid date. Allowed formats: ${[
						...STRING_DATE_FORMATS,
						'Unix timestamp on seconds (10 digits)',
						'Unix timestamp on milliseconds (13 digits)',
					].join(', ')}`);
				}
			}
		});
	});
});
