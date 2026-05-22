import { SchemaErrorResponse } from '@src/commons/interfaces';

/**
 * Represents an error that occurs during schema validation.
 * Extends the standard Error object to include additional schema error details.
 */
export class SchemaError extends Error {
  public messages: SchemaErrorResponse

  /**
   * Constructs a new SchemaError instance.
   *
   * @param messages - An object containing detailed information about the schema validation errors.
   */
  constructor(messages: SchemaErrorResponse) {
    super('SchemaError')
    this.messages = messages
    Error.captureStackTrace(this, this.constructor);
  }
}