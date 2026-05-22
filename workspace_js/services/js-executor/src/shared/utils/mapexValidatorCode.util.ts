export const mapexValidatorCode = `
				const $mv = (() => {
					// Utility functions
					const isEmail = (str) => /^[^\\s@]+@[^\\s@]+\\.[^\\s@]+$/.test(str);
					const isDate = (value) => value instanceof Date && !isNaN(value.getTime());
					const isObject = (value) => value !== null && typeof value === 'object' && !Array.isArray(value);
					const isEmpty = (value) => value === null || value === undefined || value === '';

					// Create validation error
					const createError = (message, path = '') => ({
						error: {
							message: message,
							path: path,
							details: [{ message, path, type: 'validation' }]
						},
						value: null
					});

					// Create success result
					const createSuccess = (value) => ({ error: null, value });

					// Base validator class
					const createValidator = (type, validators = []) => {
						const validator = {
							_type: type,
							_validators: [...validators],
							_required: false,
							_optional: false,

							validate(value, path = '') {
								// Check required/optional
								if (isEmpty(value)) {
									if (this._required) {
										return createError(\`\${path || 'Value'} is required\`, path);
									}
									if (this._optional) {
										return createSuccess(value);
									}
								}

								// Run all validators
								for (const validator of this._validators) {
									const result = validator(value, path);
									if (result.error) return result;
								}

								return createSuccess(value);
							},

							required() {
								const newValidator = createValidator(this._type, this._validators);
								newValidator._required = true;
								newValidator._optional = false;
								return newValidator;
							},

							optional() {
								const newValidator = createValidator(this._type, this._validators);
								newValidator._optional = true;
								newValidator._required = false;
								return newValidator;
							}
						};

						return validator;
					};

					// String validator
					const string = () => {
						const validator = createValidator('string', [
							(value, path) => {
								if (!isEmpty(value) && typeof value !== 'string') {
									return createError(\`\${path || 'Value'} must be a string\`, path);
								}
								return createSuccess(value);
							}
						]);

						validator.min = (length) => {
							const newValidator = createValidator('string', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value) && value.length < length) {
										return createError(\`\${path || 'Value'} must be at least \${length} characters\`, path);
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						validator.max = (length) => {
							const newValidator = createValidator('string', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value) && value.length > length) {
										return createError(\`\${path || 'Value'} must be at most \${length} characters\`, path);
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						validator.email = () => {
							const newValidator = createValidator('string', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value) && !isEmail(value)) {
										return createError(\`\${path || 'Value'} must be a valid email\`, path);
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						validator.pattern = (regex) => {
							const newValidator = createValidator('string', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value) && !regex.test(value)) {
										return createError(\`\${path || 'Value'} does not match pattern\`, path);
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						return validator;
					};

					// Number validator
					const number = () => {
						const validator = createValidator('number', [
							(value, path) => {
								if (!isEmpty(value) && (typeof value !== 'number' || isNaN(value))) {
									return createError(\`\${path || 'Value'} must be a number\`, path);
								}
								return createSuccess(value);
							}
						]);

						validator.min = (min) => {
							const newValidator = createValidator('number', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value) && value < min) {
										return createError(\`\${path || 'Value'} must be at least \${min}\`, path);
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						validator.max = (max) => {
							const newValidator = createValidator('number', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value) && value > max) {
										return createError(\`\${path || 'Value'} must be at most \${max}\`, path);
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						validator.integer = () => {
							const newValidator = createValidator('number', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value) && !Number.isInteger(value)) {
										return createError(\`\${path || 'Value'} must be an integer\`, path);
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						return validator;
					};

					// Boolean validator
					const boolean = () => {
						return createValidator('boolean', [
							(value, path) => {
								if (!isEmpty(value) && typeof value !== 'boolean') {
									return createError(\`\${path || 'Value'} must be a boolean\`, path);
								}
								return createSuccess(value);
							}
						]);
					};

					// Date validator
					const date = () => {
						const validator = createValidator('date', [
							(value, path) => {
								if (!isEmpty(value)) {
									const dateValue = new Date(value);
									if (!isDate(dateValue)) {
										return createError(\`\${path || 'Value'} must be a valid date\`, path);
									}
								}
								return createSuccess(value);
							}
						]);

						validator.min = (minDate) => {
							const newValidator = createValidator('date', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value)) {
										const dateValue = new Date(value);
										const minDateValue = new Date(minDate);
										if (dateValue < minDateValue) {
											return createError(\`\${path || 'Value'} must be after \${minDate}\`, path);
										}
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						validator.max = (maxDate) => {
							const newValidator = createValidator('date', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value)) {
										const dateValue = new Date(value);
										const maxDateValue = new Date(maxDate);
										if (dateValue > maxDateValue) {
											return createError(\`\${path || 'Value'} must be before \${maxDate}\`, path);
										}
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						return validator;
					};

					// Array validator
					const array = () => {
						const validator = createValidator('array', [
							(value, path) => {
								if (!isEmpty(value) && !Array.isArray(value)) {
									return createError(\`\${path || 'Value'} must be an array\`, path);
								}
								return createSuccess(value);
							}
						]);

						validator.items = (itemSchema) => {
							const newValidator = createValidator('array', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value) && Array.isArray(value)) {
										for (let i = 0; i < value.length; i++) {
											const itemResult = itemSchema.validate(value[i], \`\${path}[\${i}]\`);
											if (itemResult.error) {
												return itemResult;
											}
										}
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						validator.min = (length) => {
							const newValidator = createValidator('array', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value) && value.length < length) {
										return createError(\`\${path || 'Value'} must have at least \${length} items\`, path);
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						validator.max = (length) => {
							const newValidator = createValidator('array', [
								...validator._validators,
								(value, path) => {
									if (!isEmpty(value) && value.length > length) {
										return createError(\`\${path || 'Value'} must have at most \${length} items\`, path);
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						return validator;
					};

					// Object validator
					const object = (schema) => {
						const validator = createValidator('object', [
							(value, path) => {
								if (!isEmpty(value) && !isObject(value)) {
									return createError(\`\${path || 'Value'} must be an object\`, path);
								}
								return createSuccess(value);
							}
						]);

						if (schema) {
							validator._validators.push((value, path) => {
								if (!isEmpty(value) && isObject(value)) {
									for (const key in schema) {
										const fieldPath = path ? \`\${path}.\${key}\` : key;
										const fieldResult = schema[key].validate(value[key], fieldPath);
										if (fieldResult.error) {
											return fieldResult;
										}
									}
								}
								return createSuccess(value);
							});
						}

						validator.keys = (newSchema) => {
							const newValidator = createValidator('object', [
								validator._validators[0], // Keep type check
								(value, path) => {
									if (!isEmpty(value) && isObject(value)) {
										for (const key in newSchema) {
											const fieldPath = path ? \`\${path}.\${key}\` : key;
											const fieldResult = newSchema[key].validate(value[key], fieldPath);
											if (fieldResult.error) {
												return fieldResult;
											}
										}
									}
									return createSuccess(value);
								}
							]);
							newValidator._required = validator._required;
							newValidator._optional = validator._optional;
							return newValidator;
						};

						return validator;
					};

					// Any validator (always passes)
					const any = () => {
						return createValidator('any', [
							(value) => createSuccess(value)
						]);
					};

					// Main API
					return {
						string,
						number,
						boolean,
						date,
						array,
						object,
						any,

						// Utility function
						validate: (value, schema) => schema.validate(value),

						// Version
						version: '1.0.0'
					};
				})();
			`;
