//
// GENERIC
//
export const IS_REQUIRED = 'is required.';
export const INVALID_TYPE = 'the value specified is a invalid type.';
export const MUST_NOT_BE_EMPTY = 'must not be empty.';

//
// STRING
//
export const NOT_IS_NUMERIC_STRING = 'Value must be a number.';

//
// NUMBER
//
export const INTEGER_NUMBER = 'Number must be an integer';
export const POSITIVE_NUMBER = 'Number must be a positive';
export const VALUE_MUST_BE_GREATER_THAN = 'Value must be greater than';
export const VALUE_MUST_BE_LESS_THAN = 'Value must be less than';

// STRING DATE OR DATETIME

export const STRING_DATE_FORMATS = [
  'YYYY-MM-DD',                       // Only date
  'X',                                // Unix timestamp in seconds
  'x',                                // Unix timestamp in milliseconds
  'YYYY-MM-DD HH:mm:ss',              // Date and time
  'YYYY-MM-DD HH:mm:ss.SSS',          // Date, time with milliseconds
  'YYYY-MM-DDTHH:mm:ss.SSSZ',         // ISO 8601 with timezone (Z)
  'YYYY-MM-DDTHH:mm:ssZ',             // ISO 8601 without milliseconds
  'YYYY-MM-DDTHH:mm:ss.SSS[Z]',       // ISO 8601 with literal Z
  'YYYY-MM-DDTHH:mm:ss[Z]',           // ISO 8601 without ms, literal Z
];