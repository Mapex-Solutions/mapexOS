import { HttpStatus } from '@src/http';
import { Response } from 'express';

/**
 * Default error is when no has error defined
 */

/**
 * Check the reason for the exception and reject the request
 * @param res: Response - HTTP request
 * @param data: any - Data to request
 * @public
 * @async
 * @example
 *
 * const data = { name: 'Julian', lastName: 'Bastos', age: 28 };
 * await success(res, data);
 * // output http response 200 with data "data"
 */
export const success = async (res: Response, data: any = null) => {
  const httpStatus: number = HttpStatus.OK;
  const error = null;

  return res
    .status(httpStatus)
    .json({ httpStatus, error, data });
}

/**
 * Check the reason for the exception and reject the request
 * @param res: Response - HTTP request
 * @param data: any - Data to request
 * @public
 * @async
 * @example
 *
 * const data = { name: 'Julian', lastName: 'Bastos', age: 28 };
 * await success(res, data);
 * // output http response 201
 */
export const created = async (res: Response, data: any = null) => {
  const httpStatus: number = HttpStatus.CREATED;
  const error = null;

  return res
    .status(httpStatus)
    .json({ httpStatus, error, data });
}
